package main

import (
	"net/http"

	"github.com/LeeDat03/gin-event-app/internal/database"
	. "github.com/LeeDat03/gin-event-app/internal/helpers"
	"github.com/gin-gonic/gin"
)

func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := app.models.Events.Insert(&event)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to get events")
		return
	}

	c.JSON(http.StatusOK, events)
}

func (app *application) getEventById(c *gin.Context) {
	id, err := GetIDFromParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	event, err := app.models.Events.Get(id)

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, event)
}

func (app *application) updateEvent(c *gin.Context) {
	id, err := GetIDFromParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedEvent := &database.Event{}
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedEvent.Id = id
	if err := app.models.Events.Update(updatedEvent); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, updatedEvent)
}

func (app *application) deleteEvent(c *gin.Context) {
	id, err := GetIDFromParam(c, "id")

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"eventId": id,
	})
}
