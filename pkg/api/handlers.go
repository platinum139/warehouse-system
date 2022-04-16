package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"warehouse-system/pkg/services"
)

type WebServer struct {
	log            *log.Logger
	host           string
	port           string
	productService *services.ProductService
}

func (server *WebServer) Run() {
	http.HandleFunc("/products/bought", server.BoughtProductsQuantityHandler)
	http.HandleFunc("/products/items/bought", server.BoughtItemsQuantityHandler)

	addr := fmt.Sprintf("%s:%s", server.host, server.port)
	server.log.Fatalln(http.ListenAndServe(addr, nil))
}

func (server *WebServer) BoughtProductsQuantityHandler(w http.ResponseWriter, r *http.Request) {
	boughtProductsQuantity, err := server.productService.GetBoughtProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	data, err := json.Marshal(boughtProductsQuantity)
	if err != nil {
		server.log.Println("Unable to marshal boughtProductsQuantity struct.")
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(data)
	if err != nil {
		server.log.Printf("Unable to send response: %s\n", err)
	}
}

func (server *WebServer) BoughtItemsQuantityHandler(w http.ResponseWriter, r *http.Request) {
	boughtItemsQuantity, err := server.productService.GetBoughtItems()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	data, err := json.Marshal(boughtItemsQuantity)
	if err != nil {
		server.log.Println("Unable to marshal boughtItemsQuantity struct.")
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(data)
	if err != nil {
		server.log.Printf("Unable to send response: %s\n", err)
	}
}

func NewWebServer(host, port string, log *log.Logger, service *services.ProductService) *WebServer {
	log.SetPrefix("[web server] ")

	return &WebServer{
		log:            log,
		host:           host,
		port:           port,
		productService: service,
	}
}
