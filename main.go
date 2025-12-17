package main

import (
	"klubRanks/config"
	"klubRanks/db"
	"klubRanks/logger"
	"klubRanks/routes"
	"net/http"

	_ "klubRanks/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// @title KlubRanks API
// @version 1.0
// @description API for KlubRanks leaderboard system
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {your JWT token}"
func main() {
	logger.LogInfo("Starting KlubRanks Service...")
	db.InitDB()
	server := gin.Default()

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.GET("/ping", ping)
	server.GET("/health", health)

	routes.RegisterRoutes(server)

	server.Run(":" + config.AppConfig.Server.Port)

}

func ping(context *gin.Context) {
	context.JSON(http.StatusTeapot, gin.H{"message": "pong"})
}
func health(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
