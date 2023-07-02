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

type CustomBinder struct{ echo.DefaultBinder }

// Bind implements echo.Binder.
func (b *CustomBinder) Bind(i interface{}, c echo.Context) error {
	if err := b.BindPathParams(c, i); err != nil {
		return err
	}

	if err := b.BindQueryParams(c, i); err != nil {
		return err
	}

	return b.BindBody(c, i)
}

var _ echo.Binder = (*CustomBinder)(nil)

func initRouter() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{}
	e.Binder = &CustomBinder{}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	corsConfig := middleware.DefaultCORSConfig
	e.Use(middleware.CORSWithConfig(corsConfig))

	group := e.Group("/api")
	{
		group.GET("/info", Wrap(controller.MachineInfo))

		group.GET("/device/list", Wrap(controller.DeviceList))
		group.GET("/device/login", Wrap(controller.DeviceLogin))
		group.DELETE("/device", Wrap(controller.DeviceDelete))
		group.GET("/device/status", Wrap(controller.DeviceStatus))
		group.POST("/device/logout", Wrap(controller.DeviceLogout))

		group.GET("/upload", Wrap(controller.UploadGet))
		group.POST("/upload", Wrap(controller.UploadAdd))

		group.GET("/group/list", Wrap(controller.GroupList))
		group.GET("/group", Wrap(controller.GroupGet))
		group.POST("/group/join", Wrap(controller.GroupJoin))

		group.GET("/contact/list", Wrap(controller.ContactQuery))
		group.PUT("/contact/verify", Wrap(controller.ContactVerify))

		group.GET("/message", Wrap(controller.MessageQuery))
		group.POST("/message", Wrap(controller.MessageSend))
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

		err = c.JSON(http.StatusOK, model.ResponseModel{
			Data: res,
		})
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
