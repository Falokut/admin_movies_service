package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Falokut/admin_movies_service/internal/config"
	"github.com/Falokut/admin_movies_service/internal/events"
	"github.com/Falokut/admin_movies_service/internal/repository"
	"github.com/Falokut/admin_movies_service/internal/service"
	admin_movies_service "github.com/Falokut/admin_movies_service/pkg/admin_movies_service/v1/protos"
	jaegerTracer "github.com/Falokut/admin_movies_service/pkg/jaeger"
	"github.com/Falokut/admin_movies_service/pkg/metrics"
	server "github.com/Falokut/grpc_rest_server"
	"github.com/Falokut/healthcheck"
	image_processing_service "github.com/Falokut/image_processing_service/pkg/image_processing_service/v1/protos"
	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	kb = 8 << 10
	mb = kb << 10
)

func main() {
	logging.NewEntry(logging.ConsoleOutput)
	logger := logging.GetLogger()

	appCfg := config.GetConfig()
	logLevel, err := logrus.ParseLevel(appCfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Logger.SetLevel(logLevel)

	tracer, closer, err := jaegerTracer.InitJaeger(appCfg.JaegerConfig)
	if err != nil {
		logger.Fatal("cannot create tracer", err)
	}
	logger.Info("Jaeger connected")
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	logger.Info("Metrics initializing")
	metric, err := metrics.CreateMetrics(appCfg.PrometheusConfig.Name)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		logger.Info("Metrics server running")
		if err := metrics.RunMetricServer(appCfg.PrometheusConfig.ServerConfig); err != nil {
			logger.Fatal(err)
		}
	}()

	moviesEvents := events.NewMoviesEvents(events.KafkaConfig{
		Brokers: appCfg.KafkaConfig.Brokers,
	}, logger.Logger)
	defer moviesEvents.Shutdown()

	imagesEvents := events.NewImagesEvents(events.KafkaConfig{
		Brokers: appCfg.KafkaConfig.Brokers,
	}, logger.Logger)
	defer imagesEvents.Shutdown()

	imagesService, err := service.NewImageService(logger.Logger,
		appCfg.ImageStorageService.BasePictureUrl, appCfg.ImageStorageService.Addr, appCfg.ImageProcessingService.Addr, imagesEvents)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the images services is not established: %s", err.Error())
	}

	logger.Info("Database initializing")
	moviesDatabase, err := repository.NewPostgreDB(appCfg.DBConfig)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the database is not established: %s", err.Error())
	}

	moviesRepo := repository.NewMoviesRepository(moviesDatabase, logger.Logger)
	defer moviesRepo.Shutdown()

	genresDatabase, err := repository.NewPostgreDB(appCfg.DBConfig)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the database is not established: %s", err.Error())
	}
	genresRepo := repository.NewGenresRepository(genresDatabase, logger.Logger)
	defer genresRepo.Shutdown()

	countriesDatabase, err := repository.NewPostgreDB(appCfg.DBConfig)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the database is not established: %s", err.Error())
	}
	countriesRepo := repository.NewCountriesRepository(countriesDatabase, logger.Logger)
	defer countriesRepo.Shutdown()

	logger.Info("Healthcheck initializing")
	healthcheckManager := healthcheck.NewHealthManager(logger.Logger,
		[]healthcheck.HealthcheckResource{moviesDatabase, genresDatabase, countriesDatabase}, appCfg.HealthcheckPort, nil)
	go func() {
		logger.Info("Healthcheck server running")
		if err := healthcheckManager.RunHealthcheckEndpoint(); err != nil {
			logger.Fatalf("Shutting down, can't run healthcheck endpoint %s", err.Error())
		}
	}()

	checker, err := newExistanceChecker(countriesRepo, genresRepo)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the persons service is not established: %s", err.Error())
	}

	manager := service.NewMoviesRepositoryWrapper(countriesRepo, genresRepo, moviesRepo, imagesService, convertPictureConfig(appCfg.PostersConfig),
		convertPictureConfig(appCfg.PreviewPostersConfig),
		convertPictureConfig(appCfg.BackgroundConfig), checker, logger.Logger)
	service := service.NewMoviesService(logger.Logger, manager, countriesRepo, genresRepo, moviesEvents)

	logger.Info("Server initializing")
	s := server.NewServer(logger.Logger, service)

	s.Run(getListenServerConfig(appCfg), metric, nil, nil)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	s.Shutdown()
}

func ConvertResizeType(resizeType string) image_processing_service.ResampleFilter {
	resizeType = strings.ToTitle(resizeType)
	switch resizeType {
	case "Box":
		return image_processing_service.ResampleFilter_Box
	case "CatmullRom":
		return image_processing_service.ResampleFilter_CatmullRom
	case "Lanczos":
		return image_processing_service.ResampleFilter_Lanczos
	case "Linear":
		return image_processing_service.ResampleFilter_Linear
	case "MitchellNetravali":
		return image_processing_service.ResampleFilter_MitchellNetravali
	default:
		return image_processing_service.ResampleFilter_NearestNeighbor
	}
}

func convertPictureConfig(cfg config.PicturesConfig) service.PictureConfig {
	return service.PictureConfig{
		ValidateImage: cfg.ValidateImage,
		CheckImage: struct {
			MaxImageWidth  int32
			MaxImageHeight int32
			MinImageWidth  int32
			MinImageHeight int32
			AllowedTypes   []string
		}{
			MaxImageWidth:  cfg.CheckImage.MaxImageWidth,
			MaxImageHeight: cfg.CheckImage.MaxImageHeight,
			MinImageWidth:  cfg.CheckImage.MinImageWidth,
			MinImageHeight: cfg.CheckImage.MinImageHeight,
		},
		ResizeImage: cfg.ResizeImage,
		ImageProcessingConfig: struct {
			ImageHeight       int32
			ImageWidth        int32
			ImageResizeMethod image_processing_service.ResampleFilter
		}{
			ImageHeight:       cfg.ImageProcessingConfig.ImageHeight,
			ImageWidth:        cfg.ImageProcessingConfig.ImageWidth,
			ImageResizeMethod: ConvertResizeType(cfg.ImageProcessingConfig.ImageResizeMethod),
		},
		Category: cfg.Category,
	}
}

func getListenServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Mode:        cfg.Listen.Mode,
		Host:        cfg.Listen.Host,
		Port:        cfg.Listen.Port,
		ServiceDesc: &admin_movies_service.MoviesServiceV1_ServiceDesc,
		RegisterRestHandlerServer: func(ctx context.Context, mux *runtime.ServeMux, service any) error {
			serv, ok := service.(admin_movies_service.MoviesServiceV1Server)
			if !ok {
				return errors.New("can't convert")
			}

			return admin_movies_service.RegisterMoviesServiceV1HandlerServer(context.Background(),
				mux, serv)
		},
		MaxRequestSize: cfg.Listen.MaxRequestSize * mb,
	}
}

func newExistanceChecker(contriesRepo repository.CountriesRepository, genresRepo repository.GenresRepository) (service.ExistanceChecker, error) {
	return service.NewExistanceChecker(contriesRepo, genresRepo, logging.GetLogger().Logger), nil
}