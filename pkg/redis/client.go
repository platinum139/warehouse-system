package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"warehouse-system/config"
	"warehouse-system/pkg/models"
)

type Client struct {
	ctx context.Context
	log *log.Logger
	rds *redis.Client
}

func (client *Client) Close() {
	client.rds.Close()
}

func (client *Client) GetBoughtProductsCache() ([]models.BoughtProductsQuantity, error) {
	result, err := client.rds.HGetAll(client.ctx, "products:bought").Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		var boughtProductsRes []models.BoughtProductsQuantity
		for manufacturer, quantity := range result {
			q, err := strconv.Atoi(quantity)
			if err != nil {
				client.log.Println("unable to parse quantity to int")
			}
			boughtProductsRes = append(boughtProductsRes, models.BoughtProductsQuantity{
				Manufacturer:           manufacturer,
				BoughtProductsQuantity: q,
			})
		}
		return boughtProductsRes, nil
	}
}

func (client *Client) GetBoughtItemsCache() ([]models.BoughtItemsQuantity, error) {
	result, err := client.rds.HGetAll(client.ctx, "items:bought").Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		var boughtItemsRes []models.BoughtItemsQuantity
		for manufacturer, quantity := range result {
			q, err := strconv.Atoi(quantity)
			if err != nil {
				client.log.Println("unable to parse quantity to int")
			}
			boughtItemsRes = append(boughtItemsRes, models.BoughtItemsQuantity{
				Manufacturer:        manufacturer,
				BoughtItemsQuantity: q,
			})
		}
		return boughtItemsRes, nil
	}
}

func (client *Client) SetBoughtProductsCache(boughtProducts []models.BoughtProductsQuantity) error {
	var input []string
	for _, item := range boughtProducts {
		input = append(input, item.Manufacturer)
		input = append(input, strconv.Itoa(item.BoughtProductsQuantity))
	}
	return client.rds.HSet(client.ctx, "products:bought", input).Err()
}

func (client *Client) SetBoughtItemsCache(boughtItems []models.BoughtItemsQuantity) error {
	var input []string
	for _, item := range boughtItems {
		input = append(input, item.Manufacturer)
		input = append(input, strconv.Itoa(item.BoughtItemsQuantity))
	}
	return client.rds.HSet(client.ctx, "items:bought", input).Err()
}

func (client *Client) PutRequestToQueue(request string) error {
	return client.rds.RPush(client.ctx, "requests", request).Err()
}

func (client *Client) PopRequestFromQueue() (string, error) {
	request, err := client.rds.LPop(client.ctx, "requests").Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return request, nil
}

func NewRedisClient(ctx context.Context, log *log.Logger, config *config.AppConfig) *Client {
	log.SetPrefix("[redis.NewClient] ")

	rds := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
	})
	log.Println("Redis client is created successfully.")

	_, err := rds.Ping(ctx).Result()
	if err != nil {
		log.Printf("unable to ping redis: %s\n", err)
		return nil
	}
	log.Println("Redis Ping-Pong is successful.")

	return &Client{
		ctx: ctx,
		log: log,
		rds: rds,
	}
}
