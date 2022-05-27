package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "whatsapp-client/docs"
	"whatsapp-client/internal/api"
	"whatsapp-client/internal/middleware"
)

func Setup() *gin.Engine {
	g := gin.Default()

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	group := g.Group("/api", cors.Default(), middleware.NewRecovery(), middleware.NewResponse())
	{
		group.POST("/upload", api.UploadAdd)
		group.GET("/upload", api.UploadGet)
		group.GET("/devices", api.DeviceQuery)
		group.DELETE("/device", api.DeviceDelete)
		group.GET("/connect", api.ClientLogin)
		group.GET("/disconnect", api.ClientLogout)
		group.GET("/info", api.ClientInfo)
		group.GET("/groups", api.GroupQuery)
		group.GET("/group", api.GroupGet)
		group.GET("/group/join", api.GroupJoin)
		group.GET("/contacts", api.ContactQuery)
		group.POST("/verify", api.ContactVerify)
		group.POST("/send", api.MessageSend)
		group.GET("/messages", api.MessageQuery)
		group.GET("/quickreply", api.QuickReplyQuery)
		group.POST("/quickreply", api.QuickReplyAdd)
		group.PUT("/quickreply", api.QuickReplyEdit)
		group.DELETE("/quickreply", api.QuickReplyDelete)
		group.GET("/autoreply", api.AutoReplyQuery)
		group.POST("/autoreply", api.AutoReplyAdd)
		group.PUT("/autoreply", api.AutoReplyEdit)
		group.DELETE("/autoreply", api.AutoReplyDelete)
		group.GET("/chats", api.ChatQuery)
	}

	return g
}
