package config

import (
	"sync"

	"github.com/Falokut/admin_movies_service/internal/repository"
	"github.com/Falokut/admin_movies_service/pkg/jaeger"
	"github.com/Falokut/admin_movies_service/pkg/metrics"
	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/ilyakaznacheev/cleanenv"
)

type PicturesConfig struct {
	ValidateImage bool `yaml:"validate_image"`
	CheckImage    struct {
		MaxImageWidth  int32    `yaml:"max_image_height"`
		MaxImageHeight int32    `yaml:"max_image_width"`
		MinImageWidth  int32    `yaml:"min_image_width"`
		MinImageHeight int32    `yaml:"min_image_height"`
		AllowedTypes   []string `yaml:"allowed_types"`
	} `yaml:"check_image"`
	ResizeImage           bool `yaml:"resize_image"`
	ImageProcessingConfig struct {
		ImageHeight       int32  `yaml:"image_height"`
		ImageWidth        int32  `yaml:"image_width"`
		ImageResizeMethod string `yaml:"resize_method"`
	} `yaml:"image_processing"`
	Category string `yaml:"category"`
}

type Config struct {
	LogLevel        string `yaml:"log_level" env:"LOG_LEVEL"`
	HealthcheckPort string `yaml:"healthcheck_port" env:"HEALTHCHECK_PORT"`
	Listen          struct {
		Host           string `yaml:"host" env:"HOST"`
		Port           string `yaml:"port" env:"PORT"`
		Mode           string `yaml:"server_mode" env:"SERVER_MODE"`           // support GRPC, REST, BOTH
		MaxRequestSize int    `yaml:"max_request_size" env:"MAX_REQUEST_SIZE"` // in mb, default 4mb
	} `yaml:"listen"`

	PrometheusConfig struct {
		Name         string                      `yaml:"service_name" ENV:"PROMETHEUS_SERVICE_NAME"`
		ServerConfig metrics.MetricsServerConfig `yaml:"server_config"`
	} `yaml:"prometheus"`

	PostersConfig        PicturesConfig `yaml:"posters_config"`
	PreviewPostersConfig PicturesConfig `yaml:"preview_posters_config"`
	BackgroundConfig     PicturesConfig `yaml:"background_config"`
	KafkaConfig          struct {
		Brokers []string `yaml:"brokers"`
	} `yaml:"kafka"`
	ImageProcessingService struct {
		Addr string `yaml:"addr" env:"IMAGE_PROCESSING_SERVICE_ADDRESS"`
	} `yaml:"image_processing_service"`

	ImageStorageService struct {
		BasePictureUrl string `yaml:"base_picture_url" env:"BASE_PICTURE_URL"`
		Addr           string `yaml:"addr" env:"IMAGE_STORAGE_SERVICE_ADDRESS"`
	} `yaml:"image_storage_service"`

	DBConfig     repository.DBConfig `yaml:"db_config"`
	JaegerConfig jaeger.Config       `yaml:"jaeger"`
}

var instance *Config
var once sync.Once

const configsPath = "configs/"

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		instance = &Config{}

		if err := cleanenv.ReadConfig(configsPath+"config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}
	})

	return instance
}
