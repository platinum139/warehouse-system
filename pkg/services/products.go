package services

import (
	"fmt"
	"log"
	"warehouse-system/config"
	e "warehouse-system/errors"
	"warehouse-system/pkg/models"
	"warehouse-system/pkg/redis"
)

type ProductService struct {
	log         *log.Logger
	config      *config.AppConfig
	redisClient *redis.Client
}

func (ps *ProductService) GetBoughtProducts(token, uid string) ([]models.BoughtProductsQuantity, error) {
	boughtProductsQuantity, err := ps.redisClient.GetBoughtProductsCache()
	if err != nil {
		return nil, err
	}

	if boughtProductsQuantity != nil {
		ps.log.Println("'products:bought' cache is found.")
		return boughtProductsQuantity, nil
	}

	// if cache is empty, put request to a queue
	ps.log.Println("'products:bought' cache is empty.")

	request := fmt.Sprintf("%s:%s:products:bought", token, uid)
	if err := ps.redisClient.PutRequestToQueue(request); err != nil {
		ps.log.Printf("Unable to put request to queue: %s\n", err)
		return nil, err
	}
	ps.log.Println("Request is put to queue successfully.")

	topic := fmt.Sprintf("%s:%s", token, uid)
	message, err := ps.redisClient.SubscribeForResult(topic, ps.config.SubscribeTimeout)
	if err != nil {
		ps.log.Printf("Subscribing for result failed: %s\n", err)
		return nil, err
	}
	ps.log.Printf("Query worker has processed the request: %s\n", message)

	if message == "failed" {
		return nil, e.ProcessQueryFailedError{Message: "failed to process query"}
	}

	if message == "success" {
		boughtProductsQuantity, err = ps.redisClient.GetBoughtProductsCache()
		if err != nil {
			return nil, err
		}
		return boughtProductsQuantity, nil
	}

	return nil, nil
}

func (ps *ProductService) GetBoughtItems(token, uid string) ([]models.BoughtItemsQuantity, error) {
	boughtItemsQuantity, err := ps.redisClient.GetBoughtItemsCache()
	if err != nil {
		return nil, err
	}

	if boughtItemsQuantity != nil {
		ps.log.Println("'items:bought' cache is found.")
		return boughtItemsQuantity, nil
	}

	// if cache is empty, put request to a queue
	ps.log.Println("'items:bought' cache is empty.")

	request := fmt.Sprintf("%s:%s:items:bought", token, uid)
	if err := ps.redisClient.PutRequestToQueue(request); err != nil {
		ps.log.Printf("Unable to put request to queue: %s\n", err)
		return nil, err
	}
	ps.log.Println("Request is put to queue successfully.")

	topic := fmt.Sprintf("%s:%s", token, uid)
	message, err := ps.redisClient.SubscribeForResult(topic, ps.config.SubscribeTimeout)
	if err != nil {
		return nil, err
	}
	ps.log.Printf("Query worker has processed the request: %s\n", message)

	if message == "failed" {
		return nil, e.ProcessQueryFailedError{Message: "failed to process query"}
	}

	if message == "success" {
		boughtItemsQuantity, err = ps.redisClient.GetBoughtItemsCache()
		if err != nil {
			return nil, err
		}
		return boughtItemsQuantity, nil
	}

	return nil, nil
}

func NewProductService(log *log.Logger, config *config.AppConfig, redisClient *redis.Client) *ProductService {
	log.SetPrefix("[product service] ")
	return &ProductService{
		log:         log,
		config:      config,
		redisClient: redisClient,
	}
}
