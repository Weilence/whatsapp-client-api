package middleware

import "github.com/gin-gonic/gin"

func NewRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.JSON(500, err)
	})
}
