package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"warehouse-system/config"
	"warehouse-system/pkg/postgres"
	"warehouse-system/pkg/redis"
	"warehouse-system/utils"
)

type QueueHandler struct {
	log            *log.Logger
	config         *config.AppConfig
	redisClient    redis.Client
	postgresClient *postgres.Client
}

func (handler *QueueHandler) Run() {
	for {
		handler.handleQueue()
	}
}

func (handler *QueueHandler) handleQueue() {
	token, uid, query := handler.getRequest()
	if token == "" || uid == "" || query == "" {
		log.Printf("Invalid query. Token=%s, uid=%s, query=%s.\n\n", token, uid, query)
		return
	}

	requestCount, err := handler.redisClient.GetLockCounter(token)
	if err != nil {
		handler.log.Printf("Unable to get lock counter for token: %s\n", err)
		handler.publishResult(token, uid, "internal_err")
		return
	}

	if requestCount >= handler.config.MaxRequestsCount {
		handler.handleMaxRequestCount(token, uid, query)
		return
	}

	result := handler.processQuery(token, query)
	handler.publishResult(token, uid, result)
}

func (handler *QueueHandler) getRequest() (string, string, string) {
	request, err := handler.redisClient.PopRequestFromQueue()
	if err != nil {
		handler.log.Printf("Unable to pop request from queue: %s\n", err)
		return "", "", ""
	}

	if strings.Count(request, ":") < 2 {
		handler.log.Printf("Invalid request: %s\n", request)
		return "", "", ""
	}

	parts := strings.Split(request, ":")
	token := parts[0]
	uid := parts[1]
	query := strings.Join(parts[2:], ":")
	handler.log.Printf("Received: token %s, uid %s, query %s.\n", token, uid, query)
	return token, uid, query
}

func (handler *QueueHandler) handleMaxRequestCount(token, uid, query string) {
	handler.log.Printf("Max request count exceeded for token %s.\n", token)

	retryCount, err := handler.redisClient.GetRetryCount(token)
	if err != nil {
		handler.log.Printf("Unable to get retry count for token: %s\n", err)
		handler.publishResult(token, uid, "internal_err")
		return
	}
	// throw away request if max retry count exceeded
	if retryCount > handler.config.MaxRetryCount {
		handler.log.Printf("Max retry count for token: %s\n", err)
		handler.publishResult(token, uid, "max_retry_count")
		return
	}
	// delaying with Fibonacci strategy
	delay := utils.GetFibonacciNumber(retryCount)
	time.Sleep(time.Duration(delay) * time.Second)

	request := strings.Join([]string{token, uid, query}, ":")
	if err := handler.redisClient.PutRequestToQueue(request); err != nil {
		handler.log.Printf("Unable to put request to queue: %s\n", err)
		handler.publishResult(token, uid, "internal_err")
		return
	}
	handler.log.Println("Request is put to queue successfully.")
}

func (handler *QueueHandler) processQuery(token, query string) string {
	if err := handler.redisClient.Lock(token); err != nil {
		handler.log.Printf("Unable to lock request: %s\n", err)
		return "internal_err"
	}

	var result string
	switch query {
	case "products:bought":
		result = handler.processGetBoughtProductsQuery()
	case "items:bought":
		result = handler.processGetBoughtItemsQuery()
	}

	if err := handler.redisClient.Unlock(token); err != nil {
		handler.log.Printf("Unable to unlock request: %s\n", err)
		return "internal_err"
	}

	return result
}

func (handler *QueueHandler) processGetBoughtProductsQuery() string {
	handler.log.Printf("Incoming request: products:bought.")
	boughtProductsQuantity, err := handler.postgresClient.GetBoughtProductsQuantity()
	if err != nil {
		handler.log.Printf("Failed to get bought products quantity from postgres: %s\n", err)
		return "internal_err"
	}
	expireAfter := time.Duration(handler.config.CacheExpireDuration) * time.Second
	err = handler.redisClient.SetBoughtProductsCache(boughtProductsQuantity, expireAfter)
	if err != nil {
		handler.log.Printf("Failed to set bought products cache: %s\n", err)
		return "internal_err"
	}
	handler.log.Println("BoughtProductsQuantity is got from db successfully.")
	return "success"
}

func (handler *QueueHandler) processGetBoughtItemsQuery() string {
	handler.log.Printf("Incoming request: items:bought.")
	boughtItemsQuantity, err := handler.postgresClient.GetBoughtItemsQuantity()
	if err != nil {
		handler.log.Printf("Failed to get bought items quantity from postgres: %s\n", err)
		return "internal_err"
	}
	expireAfter := time.Duration(handler.config.CacheExpireDuration) * time.Second
	err = handler.redisClient.SetBoughtItemsCache(boughtItemsQuantity, expireAfter)
	if err != nil {
		handler.log.Printf("Failed to set bought items cache: %s\n", err)
		return "internal_err"
	}
	handler.log.Println("BoughtItemsQuantity is got from db successfully.")
	return "success"
}

func (handler *QueueHandler) publishResult(token, uid, result string) {
	topic := fmt.Sprintf("%s:%s", token, uid)
	if err := handler.redisClient.PublishResult(topic, result); err != nil {
		handler.log.Printf("Failed to publish result: %s\n", err)
	} else {
		handler.log.Printf("Result '%s' is published successfully to topic %s.\n", result, topic)
	}
}

func NewQueueHandler(log *log.Logger, config *config.AppConfig, redis redis.Client, postgres *postgres.Client) *QueueHandler {
	log.SetPrefix("[queue handler] ")
	return &QueueHandler{
		log:            log,
		config:         config,
		redisClient:    redis,
		postgresClient: postgres,
	}
}
