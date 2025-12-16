package main

import (
	"klubRanks/db"
	"klubRanks/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	server.GET("/ping", ping)

	routes.RegisterRoutes(server)

	server.Run(":8080")

}

func ping(context *gin.Context) {
	context.JSON(http.StatusTeapot, gin.H{"message": "pong"})
}
