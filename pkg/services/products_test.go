package services

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
	"warehouse-system/config"
	e "warehouse-system/errors"
	"warehouse-system/pkg/models"
	"warehouse-system/pkg/redis"
)

const (
	path = "../../"
	file = ".test.env"
)

type SuiteProductServiceTests struct {
	suite.Suite
	log    *log.Logger
	config *config.AppConfig
}

func (suite *SuiteProductServiceTests) SetupSuite() {
	logger := log.New(os.Stdout, "[main] ", log.Ldate|log.Ltime)
	suite.log = logger

	appConfig := config.NewAppConfig()
	if err := appConfig.Load(path, file); err != nil {
		logger.Printf("Unable to load app config from %s/%s\n", path, file)
		return
	}
	suite.config = appConfig
}

func (suite *SuiteProductServiceTests) TestGetBoughtProducts_GetCacheError() {
	token := "token"
	uid := "uid"

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtProductsCache").Return(nil, errors.New("redis error"))

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	boughtProductsQuantity, err := productService.GetBoughtProducts(token, uid)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), boughtProductsQuantity)
}

func (suite *SuiteProductServiceTests) TestGetBoughtProducts_CacheIsFound() {
	token := "token"
	uid := "uid"

	expectedResult := []models.BoughtProductsQuantity{{
		Manufacturer:           "Manufacturer",
		BoughtProductsQuantity: 1000,
	}}
	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtProductsCache").Return(expectedResult, nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtProducts(token, uid)

	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), expectedResult, res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtProducts_PutToQueueError() {
	token := "token"
	uid := "uid"
	request := fmt.Sprintf("%s:%s:products:bought", token, uid)

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtProductsCache").Return(nil, nil)
	redisClient.On("PutRequestToQueue", request).Return(errors.New("redis error"))

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtProducts(token, uid)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtProducts_SubscribeTimeoutError() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:products:bought", token, uid)
	expectedError := e.SubscribeTimeoutError{}

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtProductsCache").Return(nil, nil)
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("", e.SubscribeTimeoutError{})

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtProducts(token, uid)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError.Error())
	assert.Nil(suite.T(), res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtProducts_ResultInternalError() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:products:bought", token, uid)
	expectedError := e.ProcessQueryFailedError{Message: "failed to process query"}

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtProductsCache").Return(nil, nil)
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("internal_err", nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtProducts(token, uid)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError.Error())
	assert.Nil(suite.T(), res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtProducts_ResultMaxRetryCountError() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:products:bought", token, uid)
	expectedError := e.MaxRetryCountExceededError{}

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtProductsCache").Return(nil, nil)
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("max_retry_count", nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtProducts(token, uid)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError.Error())
	assert.Nil(suite.T(), res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtProducts_ResultSuccessGetCacheError() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:products:bought", token, uid)
	expectedError := errors.New("redis error")

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtProductsCache").Return(nil, nil).Once()
	redisClient.On("GetBoughtProductsCache").Return(nil, expectedError)
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("success", nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtProducts(token, uid)
	assert.EqualError(suite.T(), err, expectedError.Error())
	assert.Nil(suite.T(), res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtProducts_ResultSuccessCacheFound() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:products:bought", token, uid)
	expectedRes := []models.BoughtProductsQuantity{{Manufacturer: "Manufacturer", BoughtProductsQuantity: 1000}}

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtProductsCache").Return(nil, nil).Once()
	redisClient.On("GetBoughtProductsCache").Return(expectedRes, nil).Once()
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("success", nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtProducts(token, uid)
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), expectedRes, res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtProducts_DefaultMessage() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:products:bought", token, uid)

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtProductsCache").Return(nil, nil)
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("default", nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtProducts(token, uid)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtItems_GetCacheError() {
	token := "token"
	uid := "uid"

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtItemsCache").Return(nil, errors.New("redis error"))

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	boughtItemsQuantity, err := productService.GetBoughtItems(token, uid)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), boughtItemsQuantity)
}

func (suite *SuiteProductServiceTests) TestGetBoughtItems_CacheIsFound() {
	token := "token"
	uid := "uid"

	expectedResult := []models.BoughtItemsQuantity{{
		Manufacturer:        "Manufacturer",
		BoughtItemsQuantity: 1000,
	}}
	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtItemsCache").Return(expectedResult, nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtItems(token, uid)

	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), expectedResult, res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtItems_PutToQueueError() {
	token := "token"
	uid := "uid"
	request := fmt.Sprintf("%s:%s:items:bought", token, uid)

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtItemsCache").Return(nil, nil)
	redisClient.On("PutRequestToQueue", request).Return(errors.New("redis error"))

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtItems(token, uid)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtItems_SubscribeTimeoutError() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:items:bought", token, uid)
	expectedError := e.SubscribeTimeoutError{}

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtItemsCache").Return(nil, nil)
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("", e.SubscribeTimeoutError{})

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtItems(token, uid)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError.Error())
	assert.Nil(suite.T(), res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtItems_ResultInternalError() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:items:bought", token, uid)
	expectedError := e.ProcessQueryFailedError{Message: "failed to process query"}

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtItemsCache").Return(nil, nil)
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("internal_err", nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtItems(token, uid)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError.Error())
	assert.Nil(suite.T(), res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtItems_ResultMaxRetryCountError() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:items:bought", token, uid)
	expectedError := e.MaxRetryCountExceededError{}

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtItemsCache").Return(nil, nil)
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("max_retry_count", nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtItems(token, uid)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError.Error())
	assert.Nil(suite.T(), res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtItems_ResultSuccessGetCacheError() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:items:bought", token, uid)
	expectedError := errors.New("redis error")

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtItemsCache").Return(nil, nil).Once()
	redisClient.On("GetBoughtItemsCache").Return(nil, expectedError)
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("success", nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtItems(token, uid)
	assert.EqualError(suite.T(), err, expectedError.Error())
	assert.Nil(suite.T(), res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtItems_ResultSuccessCacheFound() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:items:bought", token, uid)
	expectedRes := []models.BoughtItemsQuantity{{Manufacturer: "Manufacturer", BoughtItemsQuantity: 1000}}

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtItemsCache").Return(nil, nil).Once()
	redisClient.On("GetBoughtItemsCache").Return(expectedRes, nil).Once()
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("success", nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtItems(token, uid)
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), expectedRes, res)
}

func (suite *SuiteProductServiceTests) TestGetBoughtItems_DefaultMessage() {
	token := "token"
	uid := "uid"
	timeout := 1
	topic := fmt.Sprintf("%s:%s", token, uid)
	request := fmt.Sprintf("%s:%s:items:bought", token, uid)

	redisClient := redis.NewClientMock()
	redisClient.On("GetBoughtItemsCache").Return(nil, nil)
	redisClient.On("PutRequestToQueue", request).Return(nil)
	redisClient.On("SubscribeForResult", topic, timeout).Return("default", nil)

	productService := NewProductServiceImpl(suite.log, suite.config, redisClient)

	res, err := productService.GetBoughtItems(token, uid)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), res)
}

func TestSuiteProductServiceTests(t *testing.T) {
	suite.Run(t, new(SuiteProductServiceTests))
}
