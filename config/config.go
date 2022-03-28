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
		return viper.Unmarshal(config)
	}
	return viper.Unmarshal(config)
}

func NewAppConfig() *AppConfig {
	var config AppConfig
	config.SetDefault()
	return &config
}
