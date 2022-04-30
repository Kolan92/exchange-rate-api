package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kolan92/exchange-rate-api/repositories"
	"gorm.io/gorm"
)

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

		exchangeRate.GET("/all-from-date/:date", GetAllExchangeRatesFromDate)

		exchangeRate.POST("/", InsertSingleExchangeRate)

		exchangeRate.GET("/range/", func(c *gin.Context) {
			const fromParamKey = "from"

			fromParam := c.Query(fromParamKey)

			const tillParamKey = "till"

			tillParam := c.Query(tillParamKey)

			c.JSON(http.StatusOK, gin.H{"check": "ok", fromParamKey: fromParam, tillParamKey: tillParam})
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
// @Description Returns most recent exchange rate source - destinaion currencies
// @Param		source		query	string	false	"source currency, default is USD"
// @Param		destination	query	string	true	"destination, currency"
// @Router		/exchange-rate/last	[get]
// @Success 	200		{object}	models.ExchangeRate
// @Success 	404
func (c *ExchangeRatesController) GetLastExchangeRate(g *gin.Context) {
	currencyCodesMap := c.repo.GetCurrenciesCodesIdsMap()

	const sourceCurrencyParamKey = "source"
	sourceCurrencyCode := g.Query(sourceCurrencyParamKey)

	const destinationCurrencyParamKey = "destination"
	destinationCurrencyCode := g.Query(destinationCurrencyParamKey)

	if len(sourceCurrencyCode) == 0 {
		sourceCurrencyCode = "USD"
	}

	if len(destinationCurrencyCode) == 0 {
		g.JSON(http.StatusBadRequest, gin.H{"error": "missing destination currency"})
		return
	}

	if sourceCurrencyCode == destinationCurrencyCode {
		g.JSON(http.StatusBadRequest, gin.H{"error": "source and destination currency are the same"})
		return
	}

	sourceCurrencyId, found := currencyCodesMap[sourceCurrencyCode]
	if !found {
		g.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unknown %s source currency", sourceCurrencyCode)})
		return
	}

	destinationCurrencyId, found := currencyCodesMap[destinationCurrencyCode]
	if !found {
		g.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unknown %s destination currency", destinationCurrencyCode)})
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
// @Param		date	path	string	true	"Date for which exchange rates should be retrived"
// @Router		/exchange-rate/all-from-date/{date}	[get]
// @Success 	200		{object}	map[string]string
func GetAllExchangeRatesFromDate(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"check": "ok"})
}

// @Summary InsertSingleExchangeRate
// @Tags		exchange-rate
// @Schemes
// @Accept		json
// @Produce		json
// @Description Inserts new exchange rate
// @Router		/exchange-rate	[post]
// @Success 	200		{object}	map[string]string
func InsertSingleExchangeRate(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"check": "ok"})
}

// @Summary GetrangeExchangeRate
// @Tags		exchange-rate
// @Schemes
// @Accept		json
// @Produce		json
// @Description Returns exchange rates for currencies in the time period
// @Param		source		query	string	true	"source currency, default is USD"
// @Param		destination	query	string	false	"destination currency"
// @Param		from	query	string	true	"From date, inclusive"
// @Param		till	query	string	true	"Till date, exclusive"
// @Router		/exchange-rate/range [get]
// @Success		200	{object}	map[string]string
func GetrangeExchangeRate(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"check": "ok"})
}

func errToStatusCode(err error) int {
	switch err {
	case gorm.ErrRecordNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
