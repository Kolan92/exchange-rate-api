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
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Rate Exchange API
// @version 1.0
// @description Provides basic functionality for checking currency exchange rate.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
// @schemes http

func main() {
	docs.SwaggerInfo.BasePath = "/api/v1"

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
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		panic("missing dbUser")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		panic("missing dbPassword")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		panic("missing dbName")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		panic("missing dbHost")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		panic("missing dbPort")
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}
