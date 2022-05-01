package integrationtests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/kolan92/exchange-rate-api/models"
	"github.com/stretchr/testify/assert"
)

var (
	client  *http.Client
	baseURL *url.URL
)

func init() {
	c := &http.Client{}

	client = c

	url, _ := url.Parse("http://localhost:8081")
	baseURL = url

}

func TestRetrivesCurrenciesForDateRange(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const getRangeExchangeRatesPath = "/api/v1/exchange-rate/range/"

	relativeUrl := &url.URL{Path: getRangeExchangeRatesPath}
	url := baseURL.ResolveReference(relativeUrl)
	query := url.Query()
	query.Add("destination", "USD")
	query.Add("source", "CHF")
	query.Add("from", "2017-05-01")
	query.Add("till", "2017-05-06")
	url.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	assert.NoError(t, err)

	response, err := client.Do(request)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	exchangeRates := []models.ExchangeRate{}

	err = json.Unmarshal(body, &exchangeRates)
	assert.NoError(t, err)

	assert.Len(t, exchangeRates, 5)

	for _, exchangeRate := range exchangeRates {
		assert.Equal(t, "USD", exchangeRate.Destination)
		assert.Equal(t, "CHF", exchangeRate.Source)
	}
}
