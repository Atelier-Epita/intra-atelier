package api

import (
	"net/http"
	"time"

	"github.com/Atelier-Epita/intra-atelier/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleEvent() {
	events := router.Group("/events")

	events.GET("", GetEventsHandler)
	events.GET("/upcoming", GetUpcomingEventsHandler)
	events.GET("/current", GetCurrentEventsHandler)
	events.GET("/past", GetPastEventsHandler)
	events.POST("/:email", CreateEventHandler)
}

// @Summary Get all events
// @Tags events
// @Produce json
// @Success 200 {array} models.Event
// @Failure 500 "Couldn't get events"
// @Router /events [GET]
func GetEventsHandler(c *gin.Context) {
	zap.S().Info("Getting all events...")

	events, err := models.GetEvents()
	if err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't get events")
		return
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Get upcoming events
// @Tags events
// @Produce json
// @Success 200 {array} models.Event
// @Failure 500 "Couldn't get upcoming events"
// @Router /events/upcoming [GET]
func GetUpcomingEventsHandler(c *gin.Context) {
	events, err := models.GetUpcomingEvents()
	if err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't get upcoming events")
		return
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Get current events
// @Tags events
// @Produce json
// @Success 200 {array} models.Event
// @Failure 500 "Couldn't get current events"
// @Router /events/current [GET]
func GetCurrentEventsHandler(c *gin.Context) {
	events, err := models.GetCurrentEvents()
	if err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't get current events")
		return
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Get past events
// @Tags events
// @Produce json
// @Success 200 {array} models.Event
// @Failure 500 "Couldn't get past events"
// @Router /events/past [GET]
func GetPastEventsHandler(c *gin.Context) {
	events, err := models.GetPastEvents()
	if err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't get past events")
		return
	}

	c.JSON(http.StatusOK, events)
}

type EventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Start_date  time.Time `json:"startDate" binding:"required"`
	End_date    time.Time `json:"endDate" binding:"required"`
	ImageID     uint64    `json:"imageID" binding:"required"`
}

// @Summary Create event
// @Tags events
// @Accept json
// @Param request body EventRequest true "EventRequest"
// @Success 200 "OK"
// @Failure 400 "Bad request"
// @Failure 404 "User not found"
// @Failure 500 "Couldn't create event"
// @Router /events/{email} [POST]
// @Param email path string true "email"
func CreateEventHandler(c *gin.Context) {
	mail := c.Param("email")
	user, err := models.GetUserByMail(mail)
	if err != nil {
		Abort(c, err, http.StatusNotFound, "User using "+mail+" not found")
		return
	}
	var eventRequest EventRequest
	if err := c.BindJSON(&eventRequest); err != nil {
		Abort(c, err, http.StatusBadRequest, "Invalid event request")
		return
	}

	e := models.Event{
		Title:       eventRequest.Title,
		Description: eventRequest.Description,
		Start_date:  eventRequest.Start_date,
		End_date:    eventRequest.End_date,
		OwnerID:     user.Id,
		ImageID:     eventRequest.ImageID,
	}

	if err := e.Insert(); err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't create event")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully created event",
	})
}
