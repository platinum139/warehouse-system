package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
	"warehouse-system/config"
	e "warehouse-system/errors"
	"warehouse-system/pkg/models"
)

type Client struct {
	ctx context.Context
	log *log.Logger
	rds *redis.Client
}

func (client *Client) Close() {
	if err := client.rds.Close(); err != nil {
		client.log.Printf("Unable to close redis client: %s\n", err)
	}
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

func (client *Client) SetBoughtProductsCache(boughtProducts []models.BoughtProductsQuantity, expiresAfter time.Duration) error {
	var input []string
	for _, item := range boughtProducts {
		input = append(input, item.Manufacturer)
		input = append(input, strconv.Itoa(item.BoughtProductsQuantity))
	}
	if err := client.rds.HSet(client.ctx, "products:bought", input).Err(); err != nil {
		return err
	}
	return client.rds.Expire(client.ctx, "products:bought", expiresAfter).Err()
}

func (client *Client) SetBoughtItemsCache(boughtItems []models.BoughtItemsQuantity, expiresAfter time.Duration) error {
	var input []string
	for _, item := range boughtItems {
		input = append(input, item.Manufacturer)
		input = append(input, strconv.Itoa(item.BoughtItemsQuantity))
	}
	if err := client.rds.HSet(client.ctx, "items:bought", input).Err(); err != nil {
		return err
	}
	return client.rds.Expire(client.ctx, "items:bought", expiresAfter).Err()
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

func (client *Client) PublishResult(topic, message string) error {
	return client.rds.Publish(client.ctx, topic, message).Err()
}

func (client *Client) SubscribeForResult(topic string, timeout int) (string, error) {
	pubsub := client.rds.Subscribe(client.ctx, topic)

	var message *redis.Message
	select {
	case message = <-pubsub.Channel():
		client.log.Printf("Received message: %s\n", message.Payload)
	case <-time.After(time.Duration(timeout) * time.Second):
		client.log.Printf("Receiving message time out.")
		return "", e.SubscribeTimeoutError{}
	}

	if err := pubsub.Close(); err != nil {
		return "", err
	}

	return message.Payload, nil
}

func (client *Client) GetLockCounter(token string) (int, error) {
	res, err := client.rds.Get(client.ctx, "lock:"+token).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (client *Client) Lock(token string) error {
	return client.rds.IncrBy(client.ctx, "lock:"+token, 1).Err()
}

func (client *Client) Unlock(token string) error {
	return client.rds.DecrBy(client.ctx, "lock:"+token, 1).Err()
}

func (client *Client) GetRetryCount(token string) (int, error) {
	res, err := client.rds.IncrBy(client.ctx, "retry:"+token, 1).Result()
	if err != nil {
		return 0, err
	}
	return int(res), nil
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
