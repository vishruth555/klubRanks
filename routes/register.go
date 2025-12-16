package routes

import (
	"klubRanks/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "event not found"})
		return
	}

	err = event.Register(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not register user to event " + err.Error()})
		return
	}

	context.JSON(http.StatusAccepted, gin.H{"message": "Registered!"})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}
	var event models.Event
	event.ID = eventId
	event.UserID = userId

	err = event.CancelRegistration(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not cancel registration " + err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Cancelled!"})
}

func getRegistrations(context *gin.Context) {
	data, err := models.GetAllRegistrations()
	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	context.JSON(http.StatusOK, data)
}
