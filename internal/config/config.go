package config

import (
	"crypto/tls"
	"crypto/x509"
	"sync"

	"github.com/Falokut/admin_movies_service/internal/repository"
	"github.com/Falokut/admin_movies_service/pkg/jaeger"
	"github.com/Falokut/admin_movies_service/pkg/metrics"
	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/ilyakaznacheev/cleanenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
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
		Addr         string                 `yaml:"addr" env:"IMAGE_PROCESSING_SERVICE_ADDRESS"`
		SecureConfig ConnectionSecureConfig `yaml:"secure_config"`
	} `yaml:"image_processing_service"`

	ImageStorageService struct {
		BasePictureUrl string                 `yaml:"base_picture_url" env:"BASE_PICTURE_URL"`
		Addr           string                 `yaml:"addr" env:"IMAGE_STORAGE_SERVICE_ADDRESS"`
		SecureConfig   ConnectionSecureConfig `yaml:"secure_config"`
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

type DialMethod = string

const (
	Insecure                 DialMethod = "INSECURE"
	NilTlsConfig             DialMethod = "NIL_TLS_CONFIG"
	ClientWithSystemCertPool DialMethod = "CLIENT_WITH_SYSTEM_CERT_POOL"
	Server                   DialMethod = "SERVER"
)

type ConnectionSecureConfig struct {
	Method DialMethod `yaml:"dial_method"`
	// Only for client connection with system pool
	ServerName string `yaml:"server_name"`
	CertName   string `yaml:"cert_name"`
	KeyName    string `yaml:"key_name"`
}

func (c ConnectionSecureConfig) GetGrpcTransportCredentials() (grpc.DialOption, error) {
	if c.Method == Insecure {
		return grpc.WithTransportCredentials(insecure.NewCredentials()), nil
	}

	if c.Method == NilTlsConfig {
		return grpc.WithTransportCredentials(credentials.NewTLS(nil)), nil
	}

	if c.Method == ClientWithSystemCertPool {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return grpc.EmptyDialOption{}, err
		}
		return grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(certPool, c.ServerName)), nil
	}

	cert, err := tls.LoadX509KeyPair(c.CertName, c.KeyName)
	if err != nil {
		return grpc.EmptyDialOption{}, err
	}
	return grpc.WithTransportCredentials(credentials.NewServerTLSFromCert(&cert)), nil
}
