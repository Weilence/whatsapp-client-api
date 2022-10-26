package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weilence/whatsapp-client/internal/api"
)

func NewRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.String(http.StatusInternalServerError, fmt.Sprint(err))
	})
}

func NewError() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.IsAborted() {
			s := ctx.Errors.String()
			_, err := ctx.Writer.WriteString(s)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
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
