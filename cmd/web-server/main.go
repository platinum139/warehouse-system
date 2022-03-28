package main

import (
	"context"
	"log"
	"os"
	"warehouse-system/config"
	"warehouse-system/pkg/redis"
)

const (
	path = "."
	file = ".env"
)

func main() {
	ctx := context.Background()
	logger := log.New(os.Stdout, "[main] ", log.Ldate)
	logger.Println("Web server started.")

	appConfig := config.NewAppConfig()
	if err := appConfig.Load(path, file); err != nil {
		logger.Printf("Unable to load app config from %s/%s\n", path, file)
		return
	}

	redisClient := redis.NewRedisClient(ctx, logger, appConfig)
	if redisClient == nil {
		logger.Println("Unable to create new redis client")
		return
	}
	defer redisClient.Close()
}
