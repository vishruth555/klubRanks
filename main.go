package main

import (
	"klubRanks/db"
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
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {your JWT token}"
func main() {
	db.InitDB()
	server := gin.Default()

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.GET("/ping", ping)

	routes.RegisterRoutes(server)

	server.Run(":8080")

}

func ping(context *gin.Context) {
	context.JSON(http.StatusTeapot, gin.H{"message": "pong"})
}
