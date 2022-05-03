package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"warehouse-system/config"
	"warehouse-system/errors"
	"warehouse-system/pkg/postgres"
	"warehouse-system/pkg/redis"
	"warehouse-system/utils"
)

type Request struct {
	Query string
	Token string
	Uid   string
}

type QueueHandler struct {
	log            *log.Logger
	config         *config.AppConfig
	redisClient    redis.Client
	postgresClient *postgres.Client
}

var (
	mutex        sync.Mutex
	workersCount int
)

func (handler *QueueHandler) Run() {
	for {
		handler.handleQueue()
	}
}

func (handler *QueueHandler) handleQueue() {
	request, err := handler.getRequest()
	if err != nil {
		handler.log.Printf("Unable to get request: %s\n", err)
		return
	}

	if request == nil {
		handler.log.Printf("No requests in queue.")
		return
	}

	if request.Query == "" || request.Uid == "" || request.Token == "" {
		handler.log.Printf("Invalid query. Token=%s, uid=%s, query=%s.\n\n", request.Token, request.Uid, request.Query)
		return
	}

	for workersCount > handler.config.MaxWorkersCount {
		handler.log.Printf("Unable to start more workers than %d. Sleeping for %d millisec...\n",
			handler.config.MaxWorkersCount, handler.config.WorkersCountCheckTime)
		time.Sleep(time.Duration(handler.config.WorkersCountCheckTime) * time.Millisecond)
	}

	go func() {
		handler.log.Printf("Starting %d worker...\n", workersCount)
		handler.worker(request)
		handler.log.Printf("Finishing %d worker.\n", workersCount)
	}()
}

func (handler *QueueHandler) getRequest() (*Request, error) {
	request, err := handler.redisClient.PopRequestFromQueue()
	if err != nil {
		handler.log.Printf("Unable to pop request from queue: %s\n", err)
		return nil, err
	}

	if request == "" {
		handler.log.Printf("No requests in the queue.")
		return nil, nil
	}

	if strings.Count(request, ":") < 2 {
		handler.log.Printf("Invalid request: %s\n", request)
		return nil, errors.InvalidQueryError{Query: request}
	}

	parts := strings.Split(request, ":")
	token := parts[0]
	uid := parts[1]
	query := strings.Join(parts[2:], ":")
	handler.log.Printf("Received: token %s, uid %s, query %s.\n", token, uid, query)
	return &Request{Query: query, Token: token, Uid: uid}, nil
}

func (handler *QueueHandler) worker(request *Request) {
	mutex.Lock()
	workersCount += 1
	mutex.Unlock()

	requestCount, err := handler.redisClient.GetLockCounter(request.Token)
	if err != nil {
		handler.log.Printf("Unable to get lock counter for token: %s\n", err)
		handler.publishResult(request.Token, request.Token, "internal_err")
		return
	}

	if requestCount >= handler.config.MaxRequestsCount {
		handler.handleMaxRequestCount(request)
		return
	}

	result := handler.processQuery(request.Token, request.Query)
	handler.publishResult(request.Token, request.Uid, result)

	mutex.Lock()
	workersCount -= 1
	mutex.Unlock()
}

func (handler *QueueHandler) handleMaxRequestCount(request *Request) {
	handler.log.Printf("Max request count exceeded for token %s.\n", request.Token)

	retryCount, err := handler.redisClient.GetRetryCount(request.Token)
	if err != nil {
		handler.log.Printf("Unable to get retry count for token: %s\n", err)
		handler.publishResult(request.Token, request.Uid, "internal_err")
		return
	}
	// throw away request if max retry count exceeded
	if retryCount > handler.config.MaxRetryCount {
		handler.log.Printf("Max retry count for token: %s\n", err)
		handler.publishResult(request.Token, request.Uid, "max_retry_count")
		return
	}
	// delaying with Fibonacci strategy
	delay := utils.GetFibonacciNumber(retryCount)
	time.Sleep(time.Duration(delay) * time.Second)

	requestStr := strings.Join([]string{request.Token, request.Uid, request.Query}, ":")
	if err := handler.redisClient.PutRequestToQueue(requestStr); err != nil {
		handler.log.Printf("Unable to put request to queue: %s\n", err)
		handler.publishResult(request.Token, request.Uid, "internal_err")
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
