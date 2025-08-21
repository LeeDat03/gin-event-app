package helpers

import (
	"errors"
	"strconv"

	"github.com/LeeDat03/gin-event-app/internal/database"
	"github.com/gin-gonic/gin"
)

func GetIDFromParam(c *gin.Context, paramName string) (int, error) {
	idStr := c.Param(paramName)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, errors.New("Invalid id ")
	}
	return id, nil
}
func GetUserFromContext(c *gin.Context) *database.User {
	contextUser, exists := c.Get("user")
	if !exists {
		return &database.User{}
	}

	user, ok := contextUser.(*database.User)
	if !ok {
		return &database.User{}
	}

	return user
}

func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status": "fail",
		"error":  message,
	})
}

func JSONResponse(c *gin.Context, status int, payload any) {
	c.JSON(status, payload)
}
