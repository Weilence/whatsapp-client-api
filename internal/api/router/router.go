package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/weilence/whatsapp-client/internal/api"

	"github.com/weilence/whatsapp-client/config"
	"github.com/weilence/whatsapp-client/internal/api/controller"
	"github.com/weilence/whatsapp-client/internal/api/model"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"

	"github.com/spf13/viper"
)

func init() {
	config.Init()
	model.Init()
	whatsapp.Init(model.SqlDB())
}

type CustomValidator struct{}

func (cv *CustomValidator) Validate(i interface{}) error {
	v := validate.Struct(i)
	if !v.Validate() {
		return v.Errors
	}
	return nil
}

func initRouter() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	group := e.Group("/api")
	{
		group.GET("/info", Wrap(controller.MachineInfo))

		group.GET("/device", Wrap(controller.DeviceQuery))
		group.GET("/device/login", Wrap(controller.DeviceLogin))
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

	return e
}

func Wrap[TReq any, TRes any](f func(*api.HttpContext, *TReq) (TRes, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req TReq

		if err := c.Bind(&req); err != nil {
			return fmt.Errorf("bind request err: %w", err)
		}

		if err := c.Validate(req); err != nil {
			return fmt.Errorf("validate request err: %w", err)
		}

		ctx := &api.HttpContext{Context: c}
		res, err := f(ctx, &req)
		if err != nil {
			return err
		}

		if c.Response().Committed {
			return nil
		}

		err = c.JSON(http.StatusOK, res)
		if err != nil {
			return fmt.Errorf("response err: %w", err)
		}

		return nil
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
