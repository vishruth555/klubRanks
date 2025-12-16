package routes

import (
	"klubRanks/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getEvent(context *gin.Context) {
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
	context.JSON(http.StatusOK, event)
}

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events"})
		return
	}
	context.JSON(http.StatusOK, events)
}

func createEvent(context *gin.Context) {

	var event models.Event
	err := context.ShouldBindBodyWithJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse data"})
		return
	}

	event.UserID = context.GetInt64("userId")
	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create event"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": event})
}

func updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	userId := context.GetInt64("userId")
	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "event not found"})
		return
	}

	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized to update event"})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindBodyWithJSON(&updatedEvent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse data"})
		return
	}

	updatedEvent.ID = eventId
	err = updatedEvent.UpdateEvent()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "error updating event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "event updated"})

}

func deleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}
	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "event not found"})
		return
	}

	userId := context.GetInt64("userId")
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized to delete event"})
		return
	}

	err = event.DeleteEvent()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "error deleting event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})

}
