package redis

import (
	"github.com/stretchr/testify/mock"
	"time"
	"warehouse-system/pkg/models"
)

type ClientMock struct {
	mock.Mock
}

func (client *ClientMock) GetBoughtProductsCache() ([]models.BoughtProductsQuantity, error) {
	args := client.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BoughtProductsQuantity), args.Error(1)
}

func (client *ClientMock) GetBoughtItemsCache() ([]models.BoughtItemsQuantity, error) {
	args := client.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BoughtItemsQuantity), args.Error(1)
}

func (client *ClientMock) SetBoughtProductsCache(boughtProducts []models.BoughtProductsQuantity, expiresAfter time.Duration) error {
	args := client.Called(boughtProducts, expiresAfter)
	return args.Error(0)
}

func (client *ClientMock) SetBoughtItemsCache(boughtItems []models.BoughtItemsQuantity, expiresAfter time.Duration) error {
	args := client.Called(boughtItems, expiresAfter)
	return args.Error(0)
}

func (client *ClientMock) PutRequestToQueue(request string) error {
	args := client.Called(request)
	return args.Error(0)
}

func (client *ClientMock) PopRequestFromQueue() (string, error) {
	args := client.Called()
	return args.Get(0).(string), args.Error(1)
}

func (client *ClientMock) PublishResult(topic, message string) error {
	args := client.Called(topic, message)
	return args.Error(0)
}

func (client *ClientMock) SubscribeForResult(topic string, timeout int) (string, error) {
	args := client.Called(topic, timeout)
	return args.Get(0).(string), args.Error(1)
}

func (client *ClientMock) GetLockCounter(token string) (int, error) {
	args := client.Called(token)
	return args.Get(0).(int), args.Error(1)
}

func (client *ClientMock) Lock(token string) error {
	args := client.Called(token)
	return args.Error(0)
}

func (client *ClientMock) Unlock(token string) error {
	args := client.Called(token)
	return args.Error(0)
}

func (client *ClientMock) GetRetryCount(token string) (int, error) {
	args := client.Called(token)
	return args.Get(0).(int), args.Error(1)
}

func (client *ClientMock) Close() {
	client.Called()
}

func NewClientMock() *ClientMock {
	return &ClientMock{}
}
