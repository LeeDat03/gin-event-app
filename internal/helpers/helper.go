package helpers

import (
	"errors"
	"strconv"

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

func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status": "fail",
		"error":  message,
	})
}
