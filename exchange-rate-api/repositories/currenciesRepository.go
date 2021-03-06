package repositories

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	customerros "github.com/kolan92/exchange-rate-api/custom-erros"
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
	GetAllExchangeRatesFromDate(date time.Time) ([]models.ExchangeRate, error)
	GetRangeExchangeRate(sourceCurrencyId, destinationCurrencyId int, from, till *time.Time) ([]models.ExchangeRate, error)
	InsertExchangeRate(exchangeRate *models.ExchangeRate) error
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

func (r *PostgresCurrenciesRepository) GetAllExchangeRatesFromDate(date time.Time) ([]models.ExchangeRate, error) {
	exchangeRates := []models.ExchangeRate{}

	const query string = `
	SELECT destination_code.code as destination, source_code.code as source, rates.date, rates.rate
		FROM public.exchange_rates rates
		JOIN public.currencies_codes source_code 
		ON rates.source_currency_id = source_code.id
		JOIN public.currencies_codes destination_code 
		ON rates.destination_currency_id = destination_code.id
		WHERE rates.date = ?
		ORDER BY rates.date DESC
	`

	if err := r.db.Raw(query, date).Scan(&exchangeRates).Error; err != nil {
		return nil, err
	}

	return exchangeRates, nil
}

func (r *PostgresCurrenciesRepository) GetRangeExchangeRate(sourceCurrencyId, destinationCurrencyId int, from, till *time.Time) ([]models.ExchangeRate, error) {
	exchangeRates := []models.ExchangeRate{}

	const query string = `
	SELECT destination_code.code as destination, source_code.code as source, rates.date, rates.rate
		FROM public.exchange_rates rates
		JOIN public.currencies_codes source_code 
		ON rates.source_currency_id = source_code.id
		JOIN public.currencies_codes destination_code 
		ON rates.destination_currency_id = destination_code.id
		WHERE rates.source_currency_id = ?
		AND rates.destination_currency_id = ?
		AND rates.date >= ?
		AND rates.date < ?
		ORDER BY rates.date DESC
	`

	if err := r.db.Raw(query, sourceCurrencyId, destinationCurrencyId, from, till).Scan(&exchangeRates).Error; err != nil {
		return nil, err
	}

	return exchangeRates, nil
}

func (r *PostgresCurrenciesRepository) InsertExchangeRate(exchangeRate *models.ExchangeRate) error {
	codesCurrenciesIdsMap := r.GetCurrenciesCodesIdsMap()

	dbExchangeRate := &models.DbExchangeRate{
		Source:      codesCurrenciesIdsMap[exchangeRate.Source],
		Destination: codesCurrenciesIdsMap[exchangeRate.Destination],
		Date:        exchangeRate.Date,
		Rate:        exchangeRate.Rate,
	}
	if err := r.db.Create(dbExchangeRate).Error; err != nil {
		if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
			switch pgError.Code {
			case "23505":
				return customerros.ErrDuplicateKeyViolation
			default:
				return err
			}

		}
	}
	return nil
}
