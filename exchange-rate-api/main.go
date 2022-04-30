package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kolan92/exchange-rate-api/controllers"
	docs "github.com/kolan92/exchange-rate-api/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Rate Exchange API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
// @schemes http

func main() {
	docs.SwaggerInfo.BasePath = "/api/v1"

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/check", HealthCheck)
	}

	server := controllers.NewServer()

	server.RegisterRouter(v1)

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
