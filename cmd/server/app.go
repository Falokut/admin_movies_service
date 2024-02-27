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
	cfg := config.GetConfig()

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Logger.SetLevel(logLevel)

	tracer, closer, err := jaegerTracer.InitJaeger(cfg.JaegerConfig)
	if err != nil {
		logger.Errorf("Shutting down, error while creating tracer %v", err)
		return
	}
	logger.Info("Jaeger connected")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	logger.Info("Metrics initializing")
	metric, err := metrics.CreateMetrics(cfg.PrometheusConfig.Name)
	if err != nil {
		logger.Errorf("Shutting down, error while creating metrics %v", err)
		return
	}

	shutdown := make(chan error, 1)
	go func() {
		logger.Info("Metrics server running")
		if err := metrics.RunMetricServer(cfg.PrometheusConfig.ServerConfig); err != nil {
			logger.Errorf("Shutting down, error while running metrics server %v", err)
			shutdown <- err
			return
		}
	}()

	moviesEvents := events.NewMoviesEvents(events.KafkaConfig{
		Brokers: cfg.KafkaConfig.Brokers,
	}, logger.Logger)
	defer moviesEvents.Shutdown()

	imagesEvents := events.NewImagesEvents(events.KafkaConfig{
		Brokers: cfg.KafkaConfig.Brokers,
	}, logger.Logger)
	defer imagesEvents.Shutdown()

	imagesService, err := service.NewImageService(logger.Logger,
		cfg.ImageStorageService.BasePictureUrl, cfg.ImageStorageService.Addr, cfg.ImageProcessingService.Addr,
		cfg.ImageProcessingService.SecureConfig, cfg.ImageStorageService.SecureConfig, imagesEvents)
	if err != nil {
		logger.Errorf("Shutting down, connection to the images services is not established: %s", err.Error())
		return
	}
	defer imagesService.Shutdown()

	logger.Info("Database initializing")
	moviesDatabase, err := repository.NewPostgreDB(cfg.DBConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the database is not established: %s", err.Error())
		return
	}

	moviesRepo := repository.NewMoviesRepository(moviesDatabase, logger.Logger)
	defer moviesRepo.Shutdown()

	genresDatabase, err := repository.NewPostgreDB(cfg.DBConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the database is not established: %s", err.Error())
		return
	}
	genresRepo := repository.NewGenresRepository(genresDatabase, logger.Logger)
	defer genresRepo.Shutdown()

	countriesDatabase, err := repository.NewPostgreDB(cfg.DBConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the database is not established: %s", err.Error())
		return
	}
	countriesRepo := repository.NewCountriesRepository(countriesDatabase, logger.Logger)
	defer countriesRepo.Shutdown()

	logger.Info("Healthcheck initializing")
	healthcheckManager := healthcheck.NewHealthManager(logger.Logger,
		[]healthcheck.HealthcheckResource{moviesDatabase, genresDatabase, countriesDatabase}, cfg.HealthcheckPort, nil)
	go func() {
		logger.Info("Healthcheck server running")
		if err := healthcheckManager.RunHealthcheckEndpoint(); err != nil {
			logger.Errorf("Shutting down, can't run healthcheck endpoint %s", err.Error())
			shutdown <- err
			return
		}
	}()

	checker, err := newExistanceChecker(countriesRepo, genresRepo)
	if err != nil {
		logger.Errorf("Shutting down, connection to the persons service is not established: %s", err.Error())
		return
	}

	manager := service.NewMoviesRepositoryWrapper(countriesRepo, genresRepo, moviesRepo, imagesService, convertPictureConfig(cfg.PostersConfig),
		convertPictureConfig(cfg.PreviewPostersConfig),
		convertPictureConfig(cfg.BackgroundConfig), checker, logger.Logger)
	service := service.NewMoviesService(logger.Logger, manager, countriesRepo, genresRepo, moviesEvents)

	logger.Info("Server initializing")
	s := server.NewServer(logger.Logger, service)
	go func() {
		if err := s.Run(getListenServerConfig(cfg), metric, nil, nil); err != nil {
			logger.Errorf("Shutting down, error while running server %s", err.Error())
			shutdown <- err
			return
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	select {
	case <-quit:
		break
	case <-shutdown:
		break
	}

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

func newExistanceChecker(contriesRepo repository.CountriesRepository,
	genresRepo repository.GenresRepository) (service.ExistanceChecker, error) {
	return service.NewExistanceChecker(contriesRepo, genresRepo, logging.GetLogger().Logger), nil
}
