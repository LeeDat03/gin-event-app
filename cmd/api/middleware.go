package main

import (
	"fmt"
	"net/http"
	"strings"

	. "github.com/LeeDat03/gin-event-app/internal/helpers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (app *application) AuthMiddleWare() gin.HandlerFunc {
	fmt.Println("heheheeh")
	return func(ctx *gin.Context) {
		fmt.Println("runnn middleware")
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ErrorResponse(ctx, http.StatusUnauthorized, "Authorize header")
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			ErrorResponse(ctx, http.StatusUnauthorized, "Bearer token not set")
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(app.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			ErrorResponse(ctx, http.StatusUnauthorized, "Invalid token")
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ErrorResponse(ctx, http.StatusUnauthorized, "Invalid token")
			ctx.Abort()
			return
		}

		userId := claims["userId"].(float64)

		user := app.getUserOrAbort(ctx, int(userId))
		if user == nil {
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
