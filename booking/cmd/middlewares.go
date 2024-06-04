package main

import "github.com/gin-gonic/gin"

func (app *application) authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Authorization header missing"})
		c.Abort()
		return
	}

	claims, err := app.jwtUtil.ValidateToken(tokenString)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("userID", claims.UserID)
	c.Set("role", claims.Role)
	c.Next()
}
