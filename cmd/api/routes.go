package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()
	v1 := g.Group("/api/v1")
	{
		v1.POST("/events", func(ctx *gin.Context) {})
		v1.GET("/events", func(ctx *gin.Context) {})
		v1.GET("/events/:id", func(ctx *gin.Context) {})
		v1.PUT("/events/:id", func(ctx *gin.Context) {})
		v1.DELETE("/events/:id", func(ctx *gin.Context) {})

	}

	return g
}
