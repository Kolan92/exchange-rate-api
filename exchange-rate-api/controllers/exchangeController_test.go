package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kolan92/exchange-rate-api/models"
	testhelpers "github.com/kolan92/exchange-rate-api/testHelpers"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	repository *testhelpers.MockRepository
	controller *ExchangeRatesController
	recorder   *httptest.ResponseRecorder
	ginContext *gin.Context
)

func setup() {
	r := testhelpers.NewMockRepository()
	repository = r
	repository.CodesCurrenciesIdsMap["USD"] = 1
	repository.CodesCurrenciesIdsMap["CHF"] = 2

	c := NewExchangeRatesController(repository)
	controller = c

	rec := httptest.NewRecorder()
	recorder = rec

	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(recorder)
	ginContext = context
}

func TestGetLastExchangeRateReturnsValue(t *testing.T) {
	setup()
	setDefaultCurrenciesInParams()

	repository.LatestExchangeRate = &models.ExchangeRate{
		Source:      "USD",
		Destination: "CHF",
		Date:        time.Date(2022, 04, 30, 10, 00, 00, 0, time.UTC),
	}
	controller.GetLastExchangeRate(ginContext)

	var actualLatestExchaneRate models.ExchangeRate
	err := json.Unmarshal(recorder.Body.Bytes(), &actualLatestExchaneRate)
	assert.NoError(t, err)
	assert.Equal(t, repository.LatestExchangeRate, &actualLatestExchaneRate)
}

func TestGetLastExchangeRateReturnsError(t *testing.T) {
	setup()
	setDefaultCurrenciesInParams()
	repository.LatestExchangeRateError = gorm.ErrRecordNotFound

	controller.GetLastExchangeRate(ginContext)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetLastExchangeMissingSourceCurrency(t *testing.T) {
	setup()
	setQueryString("destination=USD&source=PLN")

	controller.GetLastExchangeRate(ginContext)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetLastExchangeSameSourceAndDestinationCurrency(t *testing.T) {
	setup()
	setQueryString("source=USD&destination=USD")

	controller.GetLastExchangeRate(ginContext)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetLastExchangeMissingDestinationMapsToUSDCurrency(t *testing.T) {
	setup()
	setQueryString("source=CHF")
	repository.LatestExchangeRate = &models.ExchangeRate{
		Source:      "CHF",
		Destination: "USD",
		Date:        time.Date(2022, 04, 30, 10, 00, 00, 0, time.UTC),
	}

	controller.GetLastExchangeRate(ginContext)
	assert.Equal(t, repository.SourceCurrencyId, 2)
	assert.Equal(t, repository.DestinaionCurrencyId, 1)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetLastExchangeUnknownCurrencies(t *testing.T) {
	setup()

	controller.GetLastExchangeRate(ginContext)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestReturnsCurrencyCodes(t *testing.T) {
	setup()

	controller.GetAllCurrencies(ginContext)

	var currencyCodes []string
	err := json.Unmarshal(recorder.Body.Bytes(), &currencyCodes)
	assert.NoError(t, err)

	exppectedCurrencyCodes := []string{"USD", "CHF"}
	assert.Equal(t, exppectedCurrencyCodes, currencyCodes)
}

func setDefaultCurrenciesInParams() {
	setQueryString("source=USD&destination=CHF")
}

func setQueryString(queryString string) {
	ginContext.Request = &http.Request{
		URL: &url.URL{
			RawQuery: queryString,
		},
	}
}
