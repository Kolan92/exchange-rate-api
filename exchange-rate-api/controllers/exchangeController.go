package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	customerros "github.com/kolan92/exchange-rate-api/custom-erros"
	"github.com/kolan92/exchange-rate-api/models"
	"github.com/kolan92/exchange-rate-api/repositories"
	"gorm.io/gorm"
)

const dateLayout = "2006-01-02"

type ExchangeRatesController struct {
	repo repositories.CurrenciesRepository
}

func NewExchangeRatesController(repo repositories.CurrenciesRepository) *ExchangeRatesController {
	return &ExchangeRatesController{repo}
}

func (controller *ExchangeRatesController) RegisterRouter(routerGroup *gin.RouterGroup) {

	currencies := routerGroup.Group("/currencies")
	{
		currencies.GET("/", func(c *gin.Context) {
			controller.GetAllCurrencies(c)
		})
	}

	exchangeRate := routerGroup.Group("/exchange-rate")
	{
		exchangeRate.GET("/last", func(c *gin.Context) {
			controller.GetLastExchangeRate(c)
		})

		exchangeRate.GET("/all-from-date/:date", func(c *gin.Context) {
			controller.GetAllExchangeRatesFromDate(c)
		})

		exchangeRate.POST("/", func(c *gin.Context) {
			controller.InsertExchangeRate(c)
		})

		exchangeRate.GET("/range/", func(c *gin.Context) {

			controller.GetRangeExchangeRate(c)
		})
	}
}

// @Summary GetAllCurrencies
// @Tags		currencies
// @Schemes
// @Accept		json
// @Produce		json
// @Description Returns list of all currencies
// @Router		/currencies	[get]
// @Success 	200		{object}	[]string
func (c *ExchangeRatesController) GetAllCurrencies(g *gin.Context) {
	currencies := c.repo.GetCurrenciesCodes()
	g.JSON(http.StatusOK, currencies)
}

// @Summary GetLastExchangeRate
// @Tags		exchange-rate
// @Schemes
// @Accept		json
// @Produce		json
// @Description Returns most recent exchange rate  which is not null in database for source - destinaion currencies
// @Param		destination		query	string	false	"destination currency, default is USD"
// @Param		source	query	string	true	"source, currency"
// @Router		/exchange-rate/last	[get]
// @Success 	200		{object}	models.ExchangeRate
// @Success 	404
func (c *ExchangeRatesController) GetLastExchangeRate(g *gin.Context) {
	currencyCodesMap := c.repo.GetCurrenciesCodesIdsMap()

	sourceCurrencyId, destinationCurrencyId, err := getCurrenciesIds(g, currencyCodesMap)

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchangeRate, err := c.repo.GetLastExchangeRate(sourceCurrencyId, destinationCurrencyId)

	if err != nil {
		g.JSON(errToStatusCode(err), gin.H{"error": err.Error()})
	} else {
		g.JSON(http.StatusOK, exchangeRate)
	}
}

// @Summary GetAllExchangeRatesFromDate
// @Tags		exchange-rate
// @Schemes
// @Accept		json
// @Produce		json
// @Description Returns all exchange rates for the given date
// @Param		date	path	string	true	"Date for which exchange rates should be retrived. Date must be formated in YYYY-MM-DD"
// @Router		/exchange-rate/all-from-date/{date}	[get]
// @Success 	203		{object}	[]models.ExchangeRate
func (c *ExchangeRatesController) GetAllExchangeRatesFromDate(g *gin.Context) {
	const dateParmKey = "date"
	dateParam := g.Param(dateParmKey)
	dateValue, err := time.Parse(dateLayout, dateParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("date %s is in incorrect format", dateParam)})
		return
	}
	exchangeRatesFromDate, err := c.repo.GetAllExchangeRatesFromDate(dateValue)

	if err != nil {
		g.JSON(errToStatusCode(err), gin.H{"error": err.Error()})
	} else {
		g.JSON(http.StatusOK, exchangeRatesFromDate)
	}
}

// @Summary InsertExchangeRate
// @Description Inserts new exchange rate
// @Tags		exchange-rate
// @Schemes
// @Accept		json
// @Produce		json
// @Param		newExchangeRate	body	models.ExchangeRate	true	"New exchange rate to insert. Date has to be in RFC3339 format due to gin limitation. Time part will be ignored"
// @Router		/exchange-rate	[post]
// @Success 	204		{object}	models.ExchangeRate
func (c *ExchangeRatesController) InsertExchangeRate(g *gin.Context) {
	newExchangeRate := &models.ExchangeRate{}

	if err := g.ShouldBindJSON(&newExchangeRate); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": "incorrect exchange rate in body " + err.Error()})
		return
	}

	year, month, day := newExchangeRate.Date.Date()
	newExchangeRate.Date = time.Date(year, month, day, 0, 00, 00, 0, time.UTC)

	currencyCodesMap := c.repo.GetCurrenciesCodesIdsMap()
	if _, isFound := currencyCodesMap[newExchangeRate.Source]; !isFound {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Unknown source currency code"})
		return
	}

	if _, isFound := currencyCodesMap[newExchangeRate.Destination]; !isFound {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Unknown destination currency code"})
		return
	}

	if newExchangeRate.Destination == newExchangeRate.Source {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Source and Destination currencies must be different"})
		return
	}

	if err := c.repo.InsertExchangeRate(newExchangeRate); err != nil {
		statusCode := errToStatusCode(err)
		switch statusCode {
		case http.StatusConflict:
			g.JSON(statusCode, gin.H{"error": "Record exists for given currencies and date"})
		default:
			log.Println(fmt.Sprintf("Error while inserting new exchange rate to database: %s", err.Error()))
			g.JSON(statusCode, gin.H{"error": "Error while inserting new exchange rate to database"})
		}
		return
	}

	g.JSON(http.StatusAccepted, &newExchangeRate)
}

// @Summary GetRangeExchangeRate
// @Tags		exchange-rate
// @Schemes
// @Accept		json
// @Produce		json
// @Description Returns exchange rates for currencies in the time period
// @Param		destination		query	string	false	"destination currency, default is USD"
// @Param		source	query	string	true	"source currency"
// @Param		from	query	string	true	"From date, inclusive, must be formated in YYYY-MM-DD"
// @Param		till	query	string	true	"Till date, exclusive, must be formated in YYYY-MM-DD"
// @Router		/exchange-rate/range [get]
// @Success		200	{object}	[]models.ExchangeRate
func (c *ExchangeRatesController) GetRangeExchangeRate(g *gin.Context) {
	currencyCodesMap := c.repo.GetCurrenciesCodesIdsMap()

	sourceCurrencyId, destinationCurrencyId, err := getCurrenciesIds(g, currencyCodesMap)

	from, till, err := parseFromAndTillDates(g)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exchangeRates, err := c.repo.GetRangeExchangeRate(sourceCurrencyId, destinationCurrencyId, from, till)

	if err != nil {
		g.JSON(errToStatusCode(err), gin.H{"error": err.Error()})
	} else {
		g.JSON(http.StatusOK, exchangeRates)
	}
}

func errToStatusCode(err error) int {
	switch err {
	case customerros.ErrDuplicateKeyViolation:
		return http.StatusConflict
	case gorm.ErrRecordNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func getCurrenciesIds(g *gin.Context, currencyCodesMap map[string]int) (sourceCurrencyId, destinationCurrencyId int, err error) {

	const sourceCurrencyParamKey = "source"
	sourceCurrencyCode := g.Query(sourceCurrencyParamKey)

	const destinationCurrencyParamKey = "destination"
	destinationCurrencyCode := g.Query(destinationCurrencyParamKey)

	if len(sourceCurrencyCode) == 0 {
		return 0, 0, errors.New("missing source currency")
	}

	if len(destinationCurrencyCode) == 0 {
		destinationCurrencyCode = "USD"
	}

	if sourceCurrencyCode == destinationCurrencyCode {
		return 0, 0, errors.New("source and destination currency are the same")
	}

	sourceCurrencyId, isFound := currencyCodesMap[sourceCurrencyCode]
	if !isFound {
		return 0, 0, errors.New(fmt.Sprintf("Unknown %s source currency", sourceCurrencyCode))
	}

	destinationCurrencyId, isFound = currencyCodesMap[destinationCurrencyCode]
	if !isFound {
		return 0, 0, errors.New(fmt.Sprintf("Unknown %s destination currency", destinationCurrencyCode))
	}

	return sourceCurrencyId, destinationCurrencyId, nil
}

func parseFromAndTillDates(g *gin.Context) (from, till *time.Time, err error) {
	const fromParamKey = "from"

	fromParam := g.Query(fromParamKey)
	fromValue, err := time.Parse(dateLayout, fromParam)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("from %s is in incorrect format", fromParam))
	}

	const tillParamKey = "till"

	tillParam := g.Query(tillParamKey)
	tillValue, err := time.Parse(dateLayout, tillParam)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("till %s is in incorrect format", tillParam))
	}

	if fromValue.After(tillValue) {
		return nil, nil, errors.New("from must be before till")
	}

	return &fromValue, &tillValue, nil
}
