package main

import (
	"log"
	"warehouse-system/pkg/postgres"
	"warehouse-system/pkg/redis"
)

type QueueHandler struct {
	log            *log.Logger
	redisClient    *redis.Client
	postgresClient *postgres.Client
}

func (handler *QueueHandler) Run() {
	for {
		request, err := handler.redisClient.PopRequestFromQueue()
		if err != nil {
			handler.log.Printf("Unable to pop request from queue: %s\n", err)
			continue
		}

		if request != "" {
			message := "failed"

			switch request {
			case "products:bought":
				handler.log.Printf("Incoming request: products:bought.")
				boughtProductsQuantity, err := handler.postgresClient.GetBoughtProductsQuantity()
				if err != nil {
					handler.log.Printf("Failed to get bought products quantity from postgres: %s\n", err)
					break
				}
				err = handler.redisClient.SetBoughtProductsCache(boughtProductsQuantity)
				if err != nil {
					handler.log.Printf("Failed to set bought products cache: %s\n", err)
					break
				}
				handler.log.Println("BoughtProductsQuantity is got from db successfully.")
				message = "success"

			case "items:bought":
				handler.log.Printf("Incoming request: items:bought.")
				boughtItemsQuantity, err := handler.postgresClient.GetBoughtItemsQuantity()
				if err != nil {
					handler.log.Printf("Failed to get bought items quantity from postgres: %s\n", err)
					break
				}
				err = handler.redisClient.SetBoughtItemsCache(boughtItemsQuantity)
				if err != nil {
					handler.log.Printf("Failed to set bought items cache: %s\n", err)
					break
				}
				handler.log.Println("BoughtItemsQuantity is got from db successfully.")
				message = "success"
			}

			if err := handler.redisClient.PublishResult(message); err != nil {
				handler.log.Printf("Failed to publish message: %s\n", err)
			} else {
				handler.log.Printf("Message '%s' is published successfully.\n", message)
			}
		}
	}
}

func NewQueueHandler(log *log.Logger, redis *redis.Client, postgres *postgres.Client) *QueueHandler {
	log.SetPrefix("[queue handler] ")
	return &QueueHandler{
		log:            log,
		redisClient:    redis,
		postgresClient: postgres,
	}
}
