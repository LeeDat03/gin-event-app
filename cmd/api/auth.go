package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/LeeDat03/gin-event-app/internal/database"
	. "github.com/LeeDat03/gin-event-app/internal/helpers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginResponse struct {
	Token string `json:"token"`
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

func (app *application) login(c *gin.Context) {
	var auth loginRequest

	if err := c.ShouldBindJSON(&auth); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	existUser, err := app.models.Users.GetByEmail(auth.Email)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(auth.Password))
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "Invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existUser.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(app.jwtSecret))
	fmt.Println(token, tokenStr)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "error gen token")
	}
	c.JSON(http.StatusOK, loginResponse{
		Token: tokenStr,
	})

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
