package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		tokens, claims, err := ParseToken(token)
		if err != nil || !tokens.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": 1,
				"status_msg":  "权限不足",
			})
			c.Set("uid", 0)
			c.Abort()
			return
		}
		c.Set("uid", claims.UserId)
		c.Next()
	}
}
