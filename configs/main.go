package configs

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppName string `env:"APP_NAME" envDefault:"Go Echo MicroService"`

	Port string `env:"PORT,notEmpty"`

	TimeZone string `env:"TIME_ZONE" envDefault:"Asia/Jakarta"`

	UseDatabase        bool   `env:"USE_DATABASE" envDefault:"false"`
	DatabaseConnection string `env:"DATABASE_CONNECTION"`
	DatabaseHost       string `env:"DATABASE_HOST"`
	DatabasePort       string `env:"DATABASE_PORT"`
	DatabaseUsername   string `env:"DATABASE_USERNAME"`
	DatabasePassword   string `env:"DATABASE_PASSWORD"`
	DatabaseName       string `env:"DATABASE_NAME"`
	DatabaseSSLMode    string `env:"DATABASE_SSL_MODE" envDefault:"disable"`
	DatabaseParseTime  string `env:"DATABASE_PARSE_TIME" envDefault:"True"`
	DatabaseCharset    string `env:"DATABASE_CHARSET" envDefault:"utf8mb4"`
	DatabaseTimezone   string `env:"DATABASE_TIMEZONE" envDefault:"Asia/Jakarta"`

	UseRedis      bool   `env:"USE_REDIS" envDefault:"false"`
	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     string `env:"REDIS_PORT"`
	RedisUsername string `env:"REDIS_USERNAME"`
	RedisPassword string `env:"REDIS_PASSWORD"`

	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:","`

	UseSignature               bool   `env:"USE_SIGNATURE" envDefault:"false"`
	SignatureName              string `env:"SIGNATURE_NAME"`
	SignatureValidationName    string `env:"SIGNATURE_VALIDATION_NAME"`
	SignatureTransactionIDName string `env:"SIGNATURE_TRANSACTION_ID_NAME"`

	RSAOAEPKey string `env:"RSA_OAEP_KEY"`

	SecretKey string `env:"SECRET_KEY"`

	ServiceKey string `env:"SERVICE_KEY"`

	DefaultTimeout int `env:"DEFAULT_TIMEOUT" envDefault:"1"`
}

func New() (*Config, error) {
	var (
		tag string = "Config.Main.New."
		cfg Config
	)

	LoadVersion()

	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := godotenv.Load(); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("failed to load environment file")

		return nil, err
	}

	if err := env.Parse(&cfg); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to parse environment")

		return nil, err
	}

	if err := NewBodyDumpLog(); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to initiate a body dump for log")

		return nil, err
	}

	return &cfg, nil
}
