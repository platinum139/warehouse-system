package services

import (
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

	// if cache is empty, put request to a queue
	if boughtProductsQuantity == nil {
		ps.log.Println("'products:bought' cache is empty.")
		if err := ps.redisClient.PutRequestToQueue("products:bought"); err != nil {
			ps.log.Printf("Unable to put request to queue: %s\n", err)
			return nil, err
		}
		ps.log.Println("Request is put to queue successfully.")
		return nil, nil
	}

	ps.log.Println("'products:bought' cache is found.")
	return boughtProductsQuantity, nil
}

func (ps *ProductService) GetBoughtItems() ([]models.BoughtItemsQuantity, error) {
	boughtItemsQuantity, err := ps.redisClient.GetBoughtItemsCache()
	if err != nil {
		return nil, err
	}

	// if cache is empty, put request to a queue
	if boughtItemsQuantity == nil {
		ps.log.Println("'items:bought' cache is empty.")
		if err := ps.redisClient.PutRequestToQueue("items:bought"); err != nil {
			ps.log.Printf("Unable to put request to queue: %s\n", err)
			return nil, err
		}
		ps.log.Println("Request is put to queue successfully.")
		return nil, nil
	}

	ps.log.Println("'items:bought' cache is found.")
	return boughtItemsQuantity, nil
}

func NewProductService(log *log.Logger, redisClient *redis.Client) *ProductService {
	log.SetPrefix("[product service] ")
	return &ProductService{
		log:         log,
		redisClient: redisClient,
	}
}
