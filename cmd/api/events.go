package main

import (
	"fmt"
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

	event := app.getEventOrAbort(c, id)
	if event == nil {
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

func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := GetIDFromParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "not valid eventId")
		return
	}

	userId, err := GetIDFromParam(c, "userId")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "not valid userId")
		return
	}

	event := app.getEventOrAbort(c, eventId)
	if event == nil {
		return
	}
	userToAdd := app.getUserOrAbort(c, userId)
	if userToAdd == nil {
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(eventId, userId)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve")
		return
	}
	if existingAttendee != nil {
		ErrorResponse(c, http.StatusConflict, "Attendee exists")
		return
	}

	attendee := database.Attendee{
		EventId: eventId,
		UserId:  userId,
	}

	if err := app.models.Attendees.Insert(&attendee); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to insert")
		return
	}

	c.JSON(http.StatusOK, attendee)
}

func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := GetIDFromParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "not valid eventId")
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	eventId, err := GetIDFromParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "not valid eventId")
		return
	}
	userId, err := GetIDFromParam(c, "userId")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "not valid userId")
		return
	}

	err = app.models.Attendees.Delete(eventId, userId)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed")
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (app *application) getEventsByAttendee(c *gin.Context) {
	fmt.Println("run")
	id, err := GetIDFromParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "not valid attendee")
		return
	}

	events, err := app.models.Events.GetByAttendee(id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, events)

}

func (app *application) getEventOrAbort(c *gin.Context, id int) *database.Event {
	event, err := app.models.Events.Get(id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return nil
	}
	if event == nil {
		ErrorResponse(c, http.StatusNotFound, "event not found")
		return nil
	}
	return event
}
