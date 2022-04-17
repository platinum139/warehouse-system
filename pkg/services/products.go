package services

import (
	"errors"
	"log"
	"warehouse-system/pkg/models"
	"warehouse-system/pkg/redis"
)

type ProductService struct {
	log         *log.Logger
	redisClient *redis.Client
}

func (ps *ProductService) GetBoughtProducts() ([]models.BoughtProductsQuantity, error) {
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

	if err := ps.redisClient.PutRequestToQueue("products:bought"); err != nil {
		ps.log.Printf("Unable to put request to queue: %s\n", err)
		return nil, err
	}
	ps.log.Println("Request is put to queue successfully.")

	message, err := ps.redisClient.SubscribeForResult()
	if err != nil {
		return nil, err
	}
	ps.log.Printf("Query worker has processed the request: %s\n", message)

	if message == "failed" {
		return nil, errors.New("failed to process query")
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

func (ps *ProductService) GetBoughtItems() ([]models.BoughtItemsQuantity, error) {
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

	if err := ps.redisClient.PutRequestToQueue("items:bought"); err != nil {
		ps.log.Printf("Unable to put request to queue: %s\n", err)
		return nil, err
	}
	ps.log.Println("Request is put to queue successfully.")

	message, err := ps.redisClient.SubscribeForResult()
	if err != nil {
		return nil, err
	}
	ps.log.Printf("Query worker has processed the request: %s\n", message)

	if message == "failed" {
		return nil, errors.New("failed to process query")
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

func NewProductService(log *log.Logger, redisClient *redis.Client) *ProductService {
	log.SetPrefix("[product service] ")
	return &ProductService{
		log:         log,
		redisClient: redisClient,
	}
}
