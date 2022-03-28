package main

import (
	"log"
	"os"
	"warehouse-system/config"
	"warehouse-system/pkg/postgres"
)

const (
	path = "."
	file = ".env"
)

func main() {
	logger := log.New(os.Stdout, "[main]", log.Ldate)

	appConfig := config.NewAppConfig()
	if err := appConfig.Load(path, file); err != nil {
		logger.Printf("unable to load app config from %s/%s\n", path, file)
		return
	}

	postgresClient := postgres.NewClient(logger, appConfig)
	if postgresClient == nil {
		logger.Println("unable to create new postgres client")
		return
	}
	defer postgresClient.Close()

	quantities, err := postgresClient.GetBoughtProductsQuantity()
	if err != nil {
		log.Printf("unable to get bought products quantities: %s\n", err)
		return
	}
	for _, quantity := range quantities {
		log.Printf("%s: %d\n", quantity.Manufacturer, quantity.BoughtProductsQuantity)
	}

	itemsQuantities, err := postgresClient.GetBoughtItemsQuantity()
	if err != nil {
		log.Printf("unable to get bought products quantities: %s\n", err)
		return
	}
	for _, quantity := range itemsQuantities {
		log.Printf("%s: %d\n", quantity.Manufacturer, quantity.BoughtItemsQuantity)
	}
}
