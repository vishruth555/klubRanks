package routes

import (
	"klubRanks/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	// server.GET("/events", getEvents)
	// server.GET("/events/:id", getEvent)
	// server.GET("/registrations", getRegistrations)

	// auth := server.Group("/")
	// auth.Use(middlewares.Aunthenticate)
	// auth.POST("/events", createEvent)
	// auth.PUT("/events/:id", updateEvent)
	// auth.DELETE("/events/:id", deleteEvent)
	// auth.POST("/events/:id/register", registerForEvent)
	// auth.DELETE("/events/:id/register", cancelRegistration)

	server.POST("/signup", signup)
	server.POST("/login", login)

	auth := server.Group("/")
	auth.Use(middlewares.Aunthenticate)

	clubs := auth.Group("/clubs")
	{
		clubs.POST("", CreateClub)
		clubs.GET("/", GetMyClubs)
		clubs.GET("/:clubId/members", GetClubMembers)
		clubs.POST("/:clubId/members", AddMember)
	}

}
