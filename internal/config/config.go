// Package config contains Config struct that is used for configuring application.
package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config contains all parameter for configuring application.
type Config struct {
	LogLevel    string `mapstructure:"log_level"`
	LogEncoding string `mapstructure:"log_encoding"` // json/console
	HTTPPort    int    `mapstructure:"http_port"`
	DB          DB
	Kafka       Kafka
	JWTKey      string `mapstructure:"jwt_key"`
}

// DB contains parameter for configuring repository.
type DB struct {
	Host           string `mapstructure:"db_host"`
	Port           int    `mapstructure:"db_port"`
	User           string `mapstructure:"db_user"`
	Password       string `mapstructure:"db_password"`
	DBName         string `mapstructure:"db_name"`
	MigrationsPath string `mapstructure:"migrations_path"`
}

// Kafka contains parameter for configuring kafka.
type Kafka struct {
	Addr         string        `mapstructure:"kafka_addr"`
	Topic        string        `mapstructure:"kafka_topic"`
	BatchSize    int           `mapstructure:"kafka_batch_size"`
	BatchTimeout time.Duration `mapstructure:"kafka_batch_timeout"`
}

// NewConfig creates a new Config instance with parameters parsed by viber.
func NewConfig() (*Config, error) {
	config := &Config{}
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("log_level", "debug")
	viper.SetDefault("log_encoding", "console")
	viper.SetDefault("http_port", 8080) //nolint:gomnd

	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", 5432) //nolint:gomnd
	viper.SetDefault("db_user", "postgres")
	viper.SetDefault("db_password", "postgres")
	viper.SetDefault("db_name", "postgres")
	viper.SetDefault("migrations_path", "migrations")

	viper.SetDefault("kafka_addr", "127.0.0.1:9092")
	viper.SetDefault("kafka_topic", "companies-mutations")
	viper.SetDefault("kafka_batch_size", 3) //nolint:gomnd,nolintlint
	viper.SetDefault("kafka_batch_timeout", "10s")

	viper.SetDefault("jwt_key", "supersecretkey")

	_ = viper.ReadInConfig()

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config.DB); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config.Kafka); err != nil {
		return nil, err
	}

	return config, nil
}
