package repositories

import (
	"fmt"
	"log"
	"sync"

	"github.com/kolan92/exchange-rate-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// codesCurrenciesIdsMap can be initialized once, as it will be never changed.
// If api would support adding new currencies, it should be cached and invalidated on cache miss
var (
	currenciesCodesOnce   sync.Once
	codesCurrenciesIdsMap map[string]int
)

type CurrenciesRepository interface {
	GetCurrenciesCodesIdsMap() map[string]int
	GetCurrenciesCodes() []string
	GetLastExchangeRate(sourceCurrencyId, destinationCurrencyId int) (*models.ExchangeRate, error)
}

type PostgresCurrenciesRepository struct {
	db *gorm.DB
}

func NewPostgresCurrenciesRepository(connectionString string) CurrenciesRepository {
	db, err := gorm.Open(postgres.Open(connectionString))
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	return &PostgresCurrenciesRepository{db}
}

func (r *PostgresCurrenciesRepository) GetCurrenciesCodesIdsMap() map[string]int {
	currenciesCodesOnce.Do((func() {
		var dbCurrencies []models.Currency
		if err := r.db.Find(&dbCurrencies).Error; err != nil {
			log.Println("Can't find currencies")
		}

		currenciesCodesMap := make(map[string]int)

		for _, currencyCode := range dbCurrencies {
			currenciesCodesMap[currencyCode.Code] = currencyCode.Id
		}
		codesCurrenciesIdsMap = currenciesCodesMap
	}))

	return codesCurrenciesIdsMap
}

func (r *PostgresCurrenciesRepository) GetCurrenciesCodes() []string {

	currencies := []string{}

	for currencyCode := range r.GetCurrenciesCodesIdsMap() {
		currencies = append(currencies, currencyCode)
	}
	return currencies
}

func (r *PostgresCurrenciesRepository) GetLastExchangeRate(sourceCurrencyId, destinaionCurrencyId int) (*models.ExchangeRate, error) {
	var exchangeRate models.ExchangeRate

	const query string = `
	SELECT destination_code.code as destination, source_code.code as source, rates.date, rates.rate
		FROM public.exchange_rates rates
		JOIN public.currencies_codes source_code 
		ON rates.source_currency_id = source_code.id
		JOIN public.currencies_codes destination_code 
		ON rates.destination_currency_id = destination_code.id
		WHERE rates.source_currency_id = ?
		AND rates.destination_currency_id = ?
		AND rates.rate IS NOT NULL
		ORDER BY rates.date DESC
		LIMIT 1
	`

	if err := r.db.Raw(query, sourceCurrencyId, destinaionCurrencyId).First(&exchangeRate).Error; err != nil {
		return nil, err
	}

	return &exchangeRate, nil
}
