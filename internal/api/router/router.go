package router

import (
	"log"
	"net/http"

	"github.com/weilence/whatsapp-client/config"
	"github.com/weilence/whatsapp-client/internal/api/controller"
	"github.com/weilence/whatsapp-client/internal/api/model"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	config.Init()
	model.Init()
	whatsapp.Init(model.SqlDB())
}

func initRouter() *gin.Engine {
	g := gin.New()

	g.Use(gin.Logger())
	g.Use(cors.Default())
	g.Use(NewRecovery())

	group := g.Group("/api")
	{
		group.GET("/info", Wrap(controller.MachineInfo))

		group.GET("/device", Wrap(controller.DeviceQuery))
		group.POST("/device/login", Wrap(controller.DeviceLogin))
		group.POST("/device/:jid/logout", Wrap(controller.DeviceLogout))
		group.DELETE("/device/:jid", Wrap(controller.DeviceDelete))

		group.GET("/upload", Wrap(controller.UploadGet))
		group.POST("/upload", Wrap(controller.UploadAdd))

		group.GET("/group", Wrap(controller.GroupQuery))
		group.GET("/group/:jid", Wrap(controller.GroupGet))
		group.POST("/group/join", Wrap(controller.GroupJoin))

		group.GET("/contact", Wrap(controller.ContactQuery))
		group.PUT("/contact/verify", Wrap(controller.ContactVerify))

		group.GET("/chat", Wrap(controller.ChatQuery))

		group.GET("/message", Wrap(controller.MessageQuery))
		group.POST("/message", Wrap(controller.MessageSend))

		group.GET("/quick-reply", Wrap(controller.QuickReplyQuery))
		group.POST("/quick-reply", Wrap(controller.QuickReplyAdd))
		group.PUT("/quick-reply/:id", Wrap(controller.QuickReplyEdit))
		group.DELETE("/quick-reply/:id", Wrap(controller.QuickReplyDelete))

		group.GET("/auto-reply", Wrap(controller.AutoReplyQuery))
		group.POST("/auto-reply", Wrap(controller.AutoReplyAdd))
		group.PUT("/auto-reply/:id", Wrap(controller.AutoReplyEdit))
		group.DELETE("/auto-reply/:id", Wrap(controller.AutoReplyDelete))
	}

	return g
}

func Wrap[TReq any, TRes any](f func(*gin.Context, *TReq) (TRes, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req TReq
		if err := c.ShouldBindUri(req); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if c.Request.Method != http.MethodGet && c.Request.ContentLength > 0 {
			if err := c.Bind(req); err != nil {
				_ = c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		}

		res, err := f(c, &req)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if c.Request.Response.ContentLength == 0 {
			c.JSON(http.StatusOK, res)
		}
	}
}

func RunServer() {
	handler := initRouter()

	server := http.Server{
		Addr:    viper.GetString("web.host") + ":" + viper.GetString("web.port"),
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server listen err:%s", err)
	}
}
