package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"warehouse-system/config"
	"warehouse-system/pkg/postgres"
	"warehouse-system/pkg/redis"
)

type QueueHandler struct {
	log            *log.Logger
	config         *config.AppConfig
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

		if request == "" {
			continue
		}

		parts := strings.Split(request, ":")
		token := parts[0]
		uid := parts[1]
		query := strings.Join(parts[2:], ":")
		handler.log.Printf("Received: token %s, uid %s, query %s.\n", token, uid, query)

		var result string
		switch query {
		case "products:bought":
			result = handler.handleGetBoughtProductsQuery()
		case "items:bought":
			result = handler.handleGetBoughtItemsQuery()
		}

		topic := fmt.Sprintf("%s:%s", token, uid)
		if err := handler.redisClient.PublishResult(topic, result); err != nil {
			handler.log.Printf("Failed to publish result: %s\n", err)
		} else {
			handler.log.Printf("Result '%s' is published successfully to topic %s.\n", result, topic)
		}
	}
}

func (handler *QueueHandler) handleGetBoughtProductsQuery() string {
	handler.log.Printf("Incoming request: products:bought.")
	boughtProductsQuantity, err := handler.postgresClient.GetBoughtProductsQuantity()
	if err != nil {
		handler.log.Printf("Failed to get bought products quantity from postgres: %s\n", err)
		return "failed"
	}
	expireAfter := time.Duration(handler.config.CacheExpireDuration) * time.Second
	err = handler.redisClient.SetBoughtProductsCache(boughtProductsQuantity, expireAfter)
	if err != nil {
		handler.log.Printf("Failed to set bought products cache: %s\n", err)
		return "failed"
	}
	handler.log.Println("BoughtProductsQuantity is got from db successfully.")
	return "success"
}

func (handler *QueueHandler) handleGetBoughtItemsQuery() string {
	handler.log.Printf("Incoming request: items:bought.")
	boughtItemsQuantity, err := handler.postgresClient.GetBoughtItemsQuantity()
	if err != nil {
		handler.log.Printf("Failed to get bought items quantity from postgres: %s\n", err)
		return "failed"
	}
	expireAfter := time.Duration(handler.config.CacheExpireDuration) * time.Second
	err = handler.redisClient.SetBoughtItemsCache(boughtItemsQuantity, expireAfter)
	if err != nil {
		handler.log.Printf("Failed to set bought items cache: %s\n", err)
		return "failed"
	}
	handler.log.Println("BoughtItemsQuantity is got from db successfully.")
	return "success"
}

func NewQueueHandler(log *log.Logger, config *config.AppConfig, redis *redis.Client, postgres *postgres.Client) *QueueHandler {
	log.SetPrefix("[queue handler] ")
	return &QueueHandler{
		log:            log,
		config:         config,
		redisClient:    redis,
		postgresClient: postgres,
	}
}
