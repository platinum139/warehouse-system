package config

import (
	"github.com/spf13/viper"
	"sync"
)

type AppConfig struct {
	PostgresHost           string `mapstructure:"POSTGRES_HOST"`
	PostgresPort           string `mapstructure:"POSTGRES_PORT"`
	PostgresDB             string `mapstructure:"POSTGRES_DB"`
	PostgresUser           string `mapstructure:"POSTGRES_USER"`
	PostgresPassword       string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresSslMode        string `mapstructure:"POSTGRES_SSLMODE"`
	PostgresMigrationsPath string `mapstructure:"POSTGRES_MIGRATIONS_PATH"`
	RedisHost              string `mapstructure:"REDIS_HOST"`
	RedisPort              string `mapstructure:"REDIS_PORT"`
	RedisPassword          string `mapstructure:"REDIS_PASSWORD"`
	WebServerHost          string `mapstructure:"WEB_SERVER_HOST"`
	WebServerPort          string `mapstructure:"WEB_SERVER_PORT"`
	CacheExpireDuration    int    `mapstructure:"CACHE_EXPIRE_DURATION"`
	SubscribeTimeout       int    `mapstructure:"SUBSCRIBE_TIMEOUT"`
	MaxRequestsCount       int    `mapstructure:"MAX_REQUESTS_COUNT"`
	MaxRetryCount          int    `mapstructure:"MAX_RETRY_COUNT"`
}

func (config *AppConfig) SetDefault() {
	config.PostgresHost = "localhost"
	config.PostgresPort = "5432"
	config.PostgresDB = "postgres"
	config.PostgresUser = "postgres"
	config.PostgresSslMode = "disable"
	config.PostgresMigrationsPath = "file://./migrations"
	config.RedisHost = "localhost"
	config.RedisPort = "6379"
	config.WebServerHost = "localhost"
	config.WebServerPort = "80"
	config.CacheExpireDuration = 30
	config.SubscribeTimeout = 5
	config.MaxRequestsCount = 10
	config.MaxRetryCount = 10
}

func (config *AppConfig) Load(path, name string) (err error) {
	once := sync.Once{}
	once.Do(func() {
		viper.AddConfigPath(path)
		viper.SetConfigName(name)
		viper.SetConfigType("env")
		err = viper.ReadInConfig()
	})
	if err != nil {
		viper.AutomaticEnv()
		viper.BindEnv("POSTGRES_HOST")
		viper.BindEnv("POSTGRES_PORT")
		viper.BindEnv("POSTGRES_DB")
		viper.BindEnv("POSTGRES_USER")
		viper.BindEnv("POSTGRES_PASSWORD")
		viper.BindEnv("POSTGRES_SSLMODE")
		viper.BindEnv("POSTGRES_MIGRATIONS_PATH")
		viper.BindEnv("REDIS_HOST")
		viper.BindEnv("REDIS_PORT")
		viper.BindEnv("REDIS_PASSWORD")
		viper.BindEnv("WEB_SERVER_HOST")
		viper.BindEnv("WEB_SERVER_PORT")
		viper.BindEnv("CACHE_EXPIRE_DURATION")
		viper.BindEnv("SUBSCRIBE_TIMEOUT")
		viper.BindEnv("MAX_REQUESTS_COUNT")
		viper.BindEnv("MAX_RETRY_COUNT")
		return viper.Unmarshal(config)
	}
	return viper.Unmarshal(config)
}

func NewAppConfig() *AppConfig {
	var config AppConfig
	config.SetDefault()
	return &config
}
