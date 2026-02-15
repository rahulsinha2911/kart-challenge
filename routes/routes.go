package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kart-challenge/handler"
)

func InitializeRoutes(router *gin.Engine) {
	router.GET("/health", handler.HealthCheck)

	router.GET("/product", handler.ListProducts)

}
