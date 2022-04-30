package testhelpers

import (
	"time"

	"github.com/kolan92/exchange-rate-api/models"
)

type MockRepository struct {
	CodesCurrenciesIdsMap                  map[string]int
	LatestExchangeRate                     *models.ExchangeRate
	LatestExchangeRateError                error
	SourceCurrencyId, DestinaionCurrencyId int
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		CodesCurrenciesIdsMap: make(map[string]int),
	}
}

func (m *MockRepository) GetCurrenciesCodesIdsMap() map[string]int {
	return m.CodesCurrenciesIdsMap
}

func (m *MockRepository) GetCurrenciesCodes() []string {
	currencies := []string{}

	for currencyCode := range m.GetCurrenciesCodesIdsMap() {
		currencies = append(currencies, currencyCode)
	}
	return currencies
}

func (m *MockRepository) GetLastExchangeRate(sourceCurrencyId, destinationCurrencyId int) (*models.ExchangeRate, error) {
	m.SourceCurrencyId = sourceCurrencyId
	m.DestinaionCurrencyId = destinationCurrencyId
	return m.LatestExchangeRate, m.LatestExchangeRateError
}

func (m *MockRepository) GetAllExchangeRatesFromDate(date time.Time) ([]models.ExchangeRate, error) {
	return []models.ExchangeRate{}, nil
}

func (m *MockRepository) GetRangeExchangeRate(sourceCurrencyId, destinationCurrencyId int, from, till *time.Time) ([]models.ExchangeRate, error) {
	return []models.ExchangeRate{}, nil
}
