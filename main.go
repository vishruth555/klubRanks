package main

import (
	"klubRanks/config"
	"klubRanks/db"
	"klubRanks/logger"
	"klubRanks/routes"
	"net/http"
	"time"

	_ "klubRanks/docs"

	"github.com/gin-contrib/cors"
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
	// server.Use(cors.Default())
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://192.168.88.18:3000", // if frontend runs here
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// CORS configuration
	//zp: Updated to AllowAllOrigins for easier development integration
	server.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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