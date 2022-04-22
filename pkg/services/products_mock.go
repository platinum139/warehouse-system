package services

import (
	"github.com/stretchr/testify/mock"
	"warehouse-system/pkg/models"
)

type ProductServiceMock struct {
	mock.Mock
}

func (m *ProductServiceMock) GetBoughtProducts(token, uid string) ([]models.BoughtProductsQuantity, error) {
	args := m.Called(token, uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BoughtProductsQuantity), args.Error(1)
}

func (m *ProductServiceMock) GetBoughtItems(token, uid string) ([]models.BoughtItemsQuantity, error) {
	args := m.Called(token, uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BoughtItemsQuantity), args.Error(1)
}
