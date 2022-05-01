package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kolan92/exchange-rate-api/controllers"
	docs "github.com/kolan92/exchange-rate-api/docs"
	"github.com/kolan92/exchange-rate-api/repositories"
	"github.com/shopspring/decimal"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Rate Exchange API
// @version 1.0
// @description Provides basic functionality for checking currency exchange rate.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
// @schemes http

func main() {
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "localhost:8081"
	decimal.MarshalJSONWithoutQuotes = true

	log.Println("Starting exchange rate api...")

	connectionString := getConnectionString()
	repo := repositories.NewPostgresCurrenciesRepository(connectionString)

	controller := controllers.NewExchangeRatesController(repo)

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/check", HealthCheck)
	}
	controller.RegisterRouter(v1)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run(":8081")
}

// @Summary healthcheck
// @Tags	healthcheck
// @Schemes
// @Description basic healthcheck
// @Produce json
// @Success 200 {object} map[string]string
// @Router /check [get]
func HealthCheck(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"check": "ok"})
}

func getConnectionString() string {
	const missingEnvVariableLog = "Missing required env variable: %s"

	const dbUserVar = "DB_USER"
	dbUser := os.Getenv(dbUserVar)
	if dbUser == "" {
		panic(fmt.Sprintf(missingEnvVariableLog, dbUserVar))
	}

	const dbPasswordVar = "DB_PASSWORD"
	dbPassword := os.Getenv(dbPasswordVar)
	if dbPassword == "" {
		panic(fmt.Sprintf(missingEnvVariableLog, dbPasswordVar))
	}

	const dbNameVar = "DB_NAME"
	dbName := os.Getenv(dbNameVar)
	if dbName == "" {
		panic(fmt.Sprintf(missingEnvVariableLog, dbNameVar))
	}

	const dbHostVar = "DB_HOST"
	dbHost := os.Getenv(dbHostVar)
	if dbHost == "" {
		panic(fmt.Sprintf(missingEnvVariableLog, dbHostVar))
	}

	const dbPostVar = "DB_PORT"
	dbPort := os.Getenv(dbPostVar)
	if dbPort == "" {
		panic(fmt.Sprintf(missingEnvVariableLog, dbPostVar))
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}
