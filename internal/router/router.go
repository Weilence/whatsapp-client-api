package router

import (
	"errors"
	"fmt"
	"net/http"

	_ "github.com/glebarez/go-sqlite"
	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/weilence/whatsapp-client/config"
	"github.com/weilence/whatsapp-client/internal/controller"
	"github.com/weilence/whatsapp-client/internal/model"
	"github.com/weilence/whatsapp-client/internal/utils"
)

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
	if *config.Env == "dev" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	corsConfig := middleware.DefaultCORSConfig
	e.Use(middleware.CORSWithConfig(corsConfig))

	group := e.Group("/api")
	{
		group.GET("/info", Wrap(controller.MachineInfo))
		group.POST("/proxy", Wrap(controller.SetProxy))
		group.GET("/proxy/test", Wrap(controller.TestProxy))

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

		group.POST("/message", Wrap(controller.MessageSend))
	}

	return e
}

func Wrap[TReq any, TRes any](f func(*utils.HttpContext, *TReq) (TRes, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req TReq

		if err := c.Bind(&req); err != nil {
			return fmt.Errorf("bind request err: %w", err)
		}

		if err := c.Validate(req); err != nil {
			return fmt.Errorf("validate request err: %w", err)
		}

		ctx := &utils.HttpContext{Context: c}
		res, err := f(ctx, &req)
		if err != nil {
			var m model.ResponseModel
			if ok := errors.As(err, &m); ok {
				return c.JSON(m.Code, m)
			}

			return c.JSON(http.StatusBadRequest, model.ResponseModel{Code: -1, Message: err.Error()})
		}

		if c.Response().Committed {
			return nil
		}

		return c.JSON(http.StatusOK, model.ResponseModel{Data: res})
	}
}

func RunServer() {
	handler := initRouter()

	server := http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", *config.Port),
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Errorf("server listen err: %w", err))
	}
}
