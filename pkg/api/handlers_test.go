package api

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"warehouse-system/config"
	e "warehouse-system/errors"
	"warehouse-system/pkg/models"
	"warehouse-system/pkg/services"
)

const (
	path = "."
	file = ".env"
)

type SuiteHandlersTests struct {
	suite.Suite
	log    *log.Logger
	config *config.AppConfig
}

func (suite *SuiteHandlersTests) SetupSuite() {
	logger := log.New(os.Stdout, "[main] ", log.Ldate|log.Ltime)
	suite.log = logger

	appConfig := config.NewAppConfig()
	if err := appConfig.Load(path, file); err != nil {
		logger.Printf("Unable to load app config from %s/%s\n", path, file)
		return
	}
	suite.config = appConfig
}

func (suite *SuiteHandlersTests) TestBoughtProductsQuantityHandler_NoToken() {
	host := suite.config.WebServerHost
	port := suite.config.WebServerPort
	productService := services.ProductServiceMock{}
	webServer := NewWebServer(host, port, suite.log, &productService)

	req, err := http.NewRequest(http.MethodGet, "/products/bought", nil)
	assert.NoError(suite.T(), err)

	rr := httptest.NewRecorder()
	webServer.BoughtProductsQuantityHandler(rr, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func (suite *SuiteHandlersTests) TestBoughtProductsQuantityHandler_MaxRetryCountExceededError() {
	token := "token"
	host := suite.config.WebServerHost
	port := suite.config.WebServerPort

	productService := services.ProductServiceMock{}
	productService.On("GetBoughtProducts", token, mock.Anything).
		Return(nil, e.MaxRetryCountExceededError{})

	webServer := NewWebServer(host, port, suite.log, &productService)

	req, err := http.NewRequest(http.MethodGet, "/products/bought", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Token", token)

	rr := httptest.NewRecorder()
	webServer.BoughtProductsQuantityHandler(rr, req)

	assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)
}

func (suite *SuiteHandlersTests) TestBoughtProductsQuantityHandler_Error() {
	token := "token"
	host := suite.config.WebServerHost
	port := suite.config.WebServerPort

	productService := services.ProductServiceMock{}
	productService.On("GetBoughtProducts", token, mock.Anything).
		Return(nil, errors.New("product service error"))

	webServer := NewWebServer(host, port, suite.log, &productService)

	req, err := http.NewRequest(http.MethodGet, "/products/bought", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Token", token)

	rr := httptest.NewRecorder()
	webServer.BoughtProductsQuantityHandler(rr, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)
}

func (suite *SuiteHandlersTests) TestBoughtProductsQuantityHandler_WriteError() {
	token := "token"
	host := suite.config.WebServerHost
	port := suite.config.WebServerPort

	productService := services.ProductServiceMock{}
	productService.On("GetBoughtProducts", token, mock.Anything).
		Return([]models.BoughtProductsQuantity{{
			Manufacturer:           "Manufacturer",
			BoughtProductsQuantity: 1000,
		}}, nil)

	webServer := NewWebServer(host, port, suite.log, &productService)

	req, err := http.NewRequest(http.MethodGet, "/products/bought", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Token", token)

	rr := NewRecorderMock()
	webServer.BoughtProductsQuantityHandler(rr, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)
}

func (suite *SuiteHandlersTests) TestBoughtProductsQuantityHandler_HappyPath() {
	token := "token"
	host := suite.config.WebServerHost
	port := suite.config.WebServerPort

	productService := services.ProductServiceMock{}
	productService.On("GetBoughtProducts", token, mock.Anything).
		Return([]models.BoughtProductsQuantity{{
			Manufacturer:           "Manufacturer",
			BoughtProductsQuantity: 1000,
		}}, nil)

	webServer := NewWebServer(host, port, suite.log, &productService)

	req, err := http.NewRequest(http.MethodGet, "/products/bought", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Token", token)

	rr := httptest.NewRecorder()
	webServer.BoughtProductsQuantityHandler(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
}

func (suite *SuiteHandlersTests) TestBoughtItemsQuantityHandler_NoToken() {
	host := suite.config.WebServerHost
	port := suite.config.WebServerPort
	productService := services.ProductServiceMock{}
	webServer := NewWebServer(host, port, suite.log, &productService)

	req, err := http.NewRequest(http.MethodGet, "/products/items/bought", nil)
	assert.NoError(suite.T(), err)

	rr := httptest.NewRecorder()
	webServer.BoughtItemsQuantityHandler(rr, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func (suite *SuiteHandlersTests) TestBoughtItemsQuantityHandler_MaxRetryCountExceededError() {
	token := "token"
	host := suite.config.WebServerHost
	port := suite.config.WebServerPort

	productService := services.ProductServiceMock{}
	productService.On("GetBoughtItems", token, mock.Anything).
		Return(nil, e.MaxRetryCountExceededError{})

	webServer := NewWebServer(host, port, suite.log, &productService)

	req, err := http.NewRequest(http.MethodGet, "/products/items/bought", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Token", token)

	rr := httptest.NewRecorder()
	webServer.BoughtItemsQuantityHandler(rr, req)

	assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)
}

func (suite *SuiteHandlersTests) TestBoughtItemsQuantityHandler_Error() {
	token := "token"
	host := suite.config.WebServerHost
	port := suite.config.WebServerPort

	productService := services.ProductServiceMock{}
	productService.On("GetBoughtItems", token, mock.Anything).
		Return(nil, errors.New("product service error"))

	webServer := NewWebServer(host, port, suite.log, &productService)

	req, err := http.NewRequest(http.MethodGet, "/products/items/bought", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Token", token)

	rr := httptest.NewRecorder()
	webServer.BoughtItemsQuantityHandler(rr, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)
}

func (suite *SuiteHandlersTests) TestBoughtItemsQuantityHandler_WriteError() {
	token := "token"
	host := suite.config.WebServerHost
	port := suite.config.WebServerPort

	productService := services.ProductServiceMock{}
	productService.On("GetBoughtItems", token, mock.Anything).
		Return([]models.BoughtItemsQuantity{{
			Manufacturer:        "Manufacturer",
			BoughtItemsQuantity: 1000,
		}}, nil)

	webServer := NewWebServer(host, port, suite.log, &productService)

	req, err := http.NewRequest(http.MethodGet, "/products/items/bought", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Token", token)

	rr := NewRecorderMock()
	webServer.BoughtItemsQuantityHandler(rr, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)
}

func (suite *SuiteHandlersTests) TestBoughtItemsQuantityHandler_HappyPath() {
	token := "token"
	host := suite.config.WebServerHost
	port := suite.config.WebServerPort

	productService := services.ProductServiceMock{}
	productService.On("GetBoughtItems", token, mock.Anything).
		Return([]models.BoughtItemsQuantity{{
			Manufacturer:        "Manufacturer",
			BoughtItemsQuantity: 1000,
		}}, nil)

	webServer := NewWebServer(host, port, suite.log, &productService)

	req, err := http.NewRequest(http.MethodGet, "/products/items/bought", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Token", token)

	rr := httptest.NewRecorder()
	webServer.BoughtItemsQuantityHandler(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
}

func TestSuiteHandlersTests(t *testing.T) {
	suite.Run(t, new(SuiteHandlersTests))
}
