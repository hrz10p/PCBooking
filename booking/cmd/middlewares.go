package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func (app *application) authMiddleware(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Authorization header missing"})
		c.Abort()
		return
	}

	claims, err := app.jwtUtil.ValidateToken(strings.Split(tokenString, " ")[1])
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	log.Print(claims.UserID)

	c.Set("userID", claims.UserID)
	c.Set("role", claims.Role)
	c.Next()
}

func (app *application) adminMiddleware(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
		return
	}
	c.Next()
}
