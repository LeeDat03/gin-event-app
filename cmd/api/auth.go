package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/LeeDat03/gin-event-app/internal/database"
	. "github.com/LeeDat03/gin-event-app/internal/helpers"
	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

func (app *application) registerUser(c *gin.Context) {
	var register registerRequest
	if err := c.ShouldBindJSON(&register); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	register.Password = string(hashedPassword)
	user := database.User{
		Email:    register.Email,
		Password: register.Password,
		Name:     register.Name,
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (app *application) getUserOrAbort(c *gin.Context, id int) *database.User {
	user, err := app.models.Users.Get(id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return nil
	}
	if user == nil {
		ErrorResponse(c, http.StatusNotFound, "user not found")
		return nil
	}
	return user

}
