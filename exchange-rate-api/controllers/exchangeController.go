package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ExchangeRate struct {
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	StartTime   time.Time `json:"date"`
	Rate        *float64  `json:"rate"`
}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) RegisterRouter(routerGroup *gin.RouterGroup) {

	exchangeRate := routerGroup.Group("/exchange-rate")
	{
		exchangeRate.GET("/last", GetLastExchangeRate)

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

// @Summary GetLastExchangeRate
// @Tags		exchange-rates
// @Schemes
// @Accept		json
// @Produce		json
// @Description Returns most recent exchange rate source - destinaion currencies
// @Param		source		query	string	true	"source currency, default is USD"
// @Param		destination	query	string	false	"destination, currency"
// @Router		/last	[get]
// @Success 	200		{object}	map[string]string
func GetLastExchangeRate(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"check": "ok"})
}

// @Summary GetAllExchangeRatesFromDate
// @Tags		exchange-rates
// @Schemes
// @Accept		json
// @Produce		json
// @Description Returns all exchange rates for the given date
// @Param		date	path	string	true	"Date for which exchange rates should be retrived"
// @Router		/all-from-date/{date}	[get]
// @Success 	200		{object}	map[string]string
func GetAllExchangeRatesFromDate(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"check": "ok"})
}

// @Summary InsertSingleExchangeRate
// @Tags		exchange-rates
// @Schemes
// @Accept		json
// @Produce		json
// @Description Inserts new exchange rate
// @Router		/	[post]
// @Success 	200		{object}	map[string]string
func InsertSingleExchangeRate(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"check": "ok"})
}

// @Summary GetrangeExchangeRate
// @Tags		exchange-rates
// @Schemes
// @Accept		json
// @Produce		json
// @Description Returns exchange rates for currencies in the time period
// @Param		source		query	string	true	"source currency, default is USD"
// @Param		destination	query	string	false	"destination currency"
// @Param		from	query	string	true	"From date, inclusive"
// @Param		till	query	string	true	"Till date, exclusive"
// @Router		/range [get]
// @Success		200	{object}	map[string]string
func GetrangeExchangeRate(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"check": "ok"})
}
