package routes

import (
	"apiproject/models"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerForEvent(c *gin.Context) {
	userID := c.GetInt64("userID")
	eventID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	event, err := models.GetEventByID(eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve event"})
		return
	}

	if err := event.Register(userID); err != nil {
		if errors.Is(err, models.ErrAlreadyRegistered) {
			c.JSON(http.StatusConflict, gin.H{"error": "already registered for this event"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not register for event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "successfully registered for event"})
}

func cancelRegistration(c *gin.Context) {
	userID := c.GetInt64("userID")
	eventID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	event := models.Event{ID: eventID}
	if err := event.CancelRegistration(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "registration not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not cancel registration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully canceled registration"})
}
