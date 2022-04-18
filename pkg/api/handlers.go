package api

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"log"
	"net/http"
	"warehouse-system/errors"
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
	token := r.Header.Get("Token")
	if token == "" {
		err := errors.BadRequestError{Message: "token must be provided"}
		server.writeError(w, err, http.StatusBadRequest)
		return
	}
	server.log.Printf("Token: %s\n", token)

	uid := gofakeit.LetterN(16)
	server.log.Printf("Uid: %s\n", uid)

	boughtProductsQuantity, err := server.productService.GetBoughtProducts(token, uid)
	if err != nil {
		server.log.Printf("Get bought products failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(boughtProductsQuantity)
	if err != nil {
		server.log.Println("Unable to marshal boughtProductsQuantity struct.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		server.log.Printf("Unable to send response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *WebServer) BoughtItemsQuantityHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	if token == "" {
		err := errors.BadRequestError{Message: "token must be provided"}
		server.writeError(w, err, http.StatusBadRequest)
		return
	}
	server.log.Printf("Token: %s\n", token)

	uid := gofakeit.LetterN(16)
	server.log.Printf("Uid: %s\n", uid)

	boughtItemsQuantity, err := server.productService.GetBoughtItems(token, uid)
	if err != nil {
		server.log.Printf("Get bought items failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(boughtItemsQuantity)
	if err != nil {
		server.log.Println("Unable to marshal boughtItemsQuantity struct.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		server.log.Printf("Unable to send response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *WebServer) writeError(w http.ResponseWriter, err error, statusCode int) {
	data, err := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(data)
	if err != nil {
		server.log.Printf("Unable to send response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
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
