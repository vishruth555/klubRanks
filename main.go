package main

import (
	"klubRanks/config"
	"klubRanks/db"
	"klubRanks/logger"
	"klubRanks/models"
	"klubRanks/routes"
	"net"
	"net/http"
	"time"

	_ "klubRanks/docs"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
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

	godotenv.Load()
	config.Load()

	logger.LogDebug("Loaded configuration:", config.AppConfig)

	net.DefaultResolver.PreferGo = true

	db.InitDB()
	createTables()

	server := gin.Default()
	enableCORS(server)

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

func enableCORS(server *gin.Engine) {
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://club-ranks.vercel.app",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func createTables() {
	db.DB.AutoMigrate(
		&models.User{},
		&models.Club{},
		&models.Member{},
		&models.LeaderboardEntry{},
		&models.Message{},
		&models.ActivityLog{},
	)
}
