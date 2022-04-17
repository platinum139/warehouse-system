package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
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
				client.log.Println("Unable to parse quantity to int.")
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
				client.log.Println("Unable to parse quantity to int.")
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
	result, err := client.rds.BLPop(client.ctx, 60*time.Second, "requests").Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return result[1], nil
}

func (client *Client) PublishResult(message string) error {
	return client.rds.Publish(client.ctx, "response", message).Err()
}

func (client *Client) SubscribeForResult() (string, error) {
	pubsub := client.rds.Subscribe(client.ctx, "response")

	message := <-pubsub.Channel()
	client.log.Printf("Received message: %s\n", message.Payload)

	if err := pubsub.Close(); err != nil {
		return "", err
	}

	return message.Payload, nil
}

func NewRedisClient(ctx context.Context, log *log.Logger, config *config.AppConfig) *Client {
	log.SetPrefix("[redis client] ")

	rds := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
	})
	log.Println("Redis client is created successfully.")

	_, err := rds.Ping(ctx).Result()
	if err != nil {
		log.Printf("Unable to ping redis: %s\n", err)
		return nil
	}
	log.Println("Redis Ping-Pong is successful.")

	return &Client{
		ctx: ctx,
		log: log,
		rds: rds,
	}
}
