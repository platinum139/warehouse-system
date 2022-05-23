package models

import "time"

type BoughtProductsQuantity struct {
	Manufacturer           string
	BoughtProductsQuantity int
}

type BoughtItemsQuantity struct {
	Manufacturer        string
	BoughtItemsQuantity int
}

type User struct {
	Id         int
	ExternalId string
	Username   string
	Phone      string
	Email      string
}

type Manufacturer struct {
	Id         int
	ExternalId string
	Name       string
	Code       string
}

type Product struct {
	Id             int
	ExternalId     string
	Name           string
	ExpiresAt      time.Time
	ManufacturerId int
}

type Order struct {
	Id             int
	ExternalId     string
	Quantity       int
	ManufacturerId int
	ClientId       int
}
