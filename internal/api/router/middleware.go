package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weilence/whatsapp-client/internal/api"
	"net/http"
)

func NewRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.String(http.StatusInternalServerError, fmt.Sprint(err))
	})
}

func NewAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := api.GetUser(c)
		if user == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
