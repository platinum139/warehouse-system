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
}

func (config *AppConfig) SetDefault() {
	config.PostgresHost = "localhost"
	config.PostgresPort = "5432"
	config.PostgresDB = "postgres"
	config.PostgresUser = "postgres"
	config.PostgresSslMode = "disable"
	config.PostgresMigrationsPath = "file://./migrations"
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
		return viper.Unmarshal(config)
	}
	return viper.Unmarshal(config)
}

func NewAppConfig() *AppConfig {
	var config AppConfig
	config.SetDefault()
	return &config
}
