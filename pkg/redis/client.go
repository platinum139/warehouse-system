package redis

import (
	"time"
	"warehouse-system/pkg/models"
)

type Client interface {
	Close()
	GetBoughtProductsCache() ([]models.BoughtProductsQuantity, error)
	GetBoughtItemsCache() ([]models.BoughtItemsQuantity, error)
	SetBoughtProductsCache(boughtProducts []models.BoughtProductsQuantity, expiresAfter time.Duration) error
	SetBoughtItemsCache(boughtItems []models.BoughtItemsQuantity, expiresAfter time.Duration) error
	PutRequestToQueue(request string) error
	PopRequestFromQueue() (string, error)
	PublishResult(topic, message string) error
	SubscribeForResult(topic string, timeout int) (string, error)
	GetLockCounter(token string) (int, error)
	Lock(token string) error
	Unlock(token string) error
	GetRetryCount(token string) (int, error)
}
