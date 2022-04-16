package main

import (
	"context"
	"log"
	"os"
	"warehouse-system/config"
	"warehouse-system/pkg/postgres"
	"warehouse-system/pkg/redis"
)

const (
	path = "."
	file = ".env"
)

func main() {
	logger := log.New(os.Stdout, "[main] ", log.Ldate|log.Ltime)
	ctx := context.Background()

	appConfig := config.NewAppConfig()
	if err := appConfig.Load(path, file); err != nil {
		logger.Printf("unable to load app config from %s/%s\n", path, file)
		return
	}

	redisClient := redis.NewRedisClient(ctx, logger, appConfig)
	if redisClient == nil {
		logger.Println("Unable to create new redis client.")
		return
	}
	defer redisClient.Close()

	postgresClient := postgres.NewClient(logger, appConfig)
	if postgresClient == nil {
		logger.Println("Unable to create new postgres client.")
		return
	}
	defer postgresClient.Close()

	queueHandler := NewQueueHandler(logger, redisClient, postgresClient)
	queueHandler.Run()
}
