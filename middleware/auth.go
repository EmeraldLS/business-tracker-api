package middleware

import (
	"fmt"
	"net/http"

	"github.com/EmeraldLS/phsps-api/token"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	signedToken, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"response": "token_error",
			"message":  "No authorization token provided",
		})
		c.Abort()
		return
	}
	if err := token.ValidateToken(signedToken); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"response": "token_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.Next()
}
