package main

import (
	"fmt"
	"net/http"

	"github.com/LeeDat03/gin-event-app/internal/database"
	. "github.com/LeeDat03/gin-event-app/internal/helpers"
	"github.com/gin-gonic/gin"
)

// CreateEvent creates a new event
//
//	@Summary		Create a new event
//	@Description	Adds a new event to the database with the authenticated user as the owner.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		database.Event	true	"Event object to be created"
//	@Success		201		{object}	database.Event
//	@Router			/api/v1/events [post]
//	@Security		BearerAuth
func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	user := GetUserFromContext(c)

	if err := c.ShouldBindJSON(&event); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	event.OwnerId = user.ID
	err := app.models.Events.Insert(&event)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetEvents returns all events
//
//	@Summary		Returns all events
//	@Description	Returns all events
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]database.Event
//	@Router			/api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to get events")
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEvent returns a single event
//
//	@Summary		Returns a single event
//	@Description	Returns a single event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	database.Event
//	@Router			/api/v1/events/{id} [get]
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

// UpdateEvent updates an existing event
//
//	@Summary		Updates an existing event
//	@Description	Updates an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Event ID"
//	@Param			event	body		database.Event	true	"Event"
//	@Success		200		{object}	database.Event
//	@Router			/api/v1/events/{id} [put]
//	@Security		BearerAuth
func (app *application) updateEvent(c *gin.Context) {
	id, err := GetIDFromParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := GetUserFromContext(c)

	existingEvent := app.getEventOrAbort(c, id)
	if existingEvent == nil {
		return
	}

	if existingEvent.OwnerId != user.ID {
		ErrorResponse(c, http.StatusForbidden, "Not allowed to update this")
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

// DeleteEvent deletes an existing event
//
//	@Summary		Deletes an existing event
//	@Description	Deletes an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Event ID"
//	@Success		204
//	@Router			/api/v1/events/{id} [delete]
//	@Security		BearerAuth
func (app *application) deleteEvent(c *gin.Context) {
	id, err := GetIDFromParam(c, "id")

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := GetUserFromContext(c)
	existingEvent := app.getEventOrAbort(c, id)
	if existingEvent == nil {
		return
	}

	if existingEvent.OwnerId != user.ID {
		ErrorResponse(c, http.StatusForbidden, "Not allowed to update this")
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

// AddAttendeeToEvent adds an attendee to an event
//
//	@Summary		Adds an attendee to an event
//	@Description	Adds an attendee to an event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"Event ID"
//	@Param			userId	path		int	true	"User ID"
//	@Success		201		{object}	database.Attendee
//	@Router			/api/v1/events/{id}/attendees/{userId} [post]
//	@Security		BearerAuth
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

	user := GetUserFromContext(c)
	if user.ID != event.OwnerId {
		ErrorResponse(c, http.StatusForbidden, "Not allowed to update this")
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

// GetAttendeesForEvent returns all attendees for a given event
//
//	@Summary		Returns all attendees for a given event
//	@Description	Returns all attendees for a given event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	[]database.User
//	@Router			/api/v1/events/{id}/attendees [get]
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

// DeleteAttendeeFromEvent deletes an attendee from an event
//
//	@Summary		Deletes an attendee from an event
//	@Description	Deletes an attendee from an event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int	true	"Event ID"
//	@Param			userId	path	int	true	"User ID"
//	@Success		204
//	@Router			/api/v1/events/{id}/attendees/{userId} [delete]
//	@Security		BearerAuth
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

	event := app.getEventOrAbort(c, eventId)
	if event == nil {
		return
	}
	user := GetUserFromContext(c)
	if user.ID != event.OwnerId {
		ErrorResponse(c, http.StatusForbidden, "Not allowed to update this")
		return
	}

	err = app.models.Attendees.Delete(eventId, userId)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed")
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetEventsByAttendee returns all events for a given attendee
//
//	@Summary		Returns all events for a given attendee
//	@Description	Returns all events for a given attendee
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Attendee ID"
//	@Success		200	{object}	[]database.Event
//	@Router			/api/v1/attendees/{id}/events [get]
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
