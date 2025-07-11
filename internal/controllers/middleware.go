package controllers

import (
	"financing-aggregator/internal/exchange"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") || len(authHeader) <= 7 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, exchange.NewErrorResponse("invalid authorization header"))
			return
		}
		c.Next()
	}
}
