package api

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"go.mau.fi/whatsmeow/types"
)

type User struct {
	JID types.JID
}

type Context interface {
	context.Context
	CurrentUser() *User
}

type HttpContext struct {
	echo.Context
}

// Deadline implements Context.
func (c *HttpContext) Deadline() (deadline time.Time, ok bool) {
	return c.Request().Context().Deadline()
}

// Done implements Context.
func (c *HttpContext) Done() <-chan struct{} {
	return c.Request().Context().Done()
}

// Err implements Context.
func (c *HttpContext) Err() error {
	return c.Request().Context().Err()
}

// Value implements Context.
func (c *HttpContext) Value(key any) any {
	return c.Request().Context().Value(key)
}

func (c *HttpContext) SSEvent(event string, data string) error {
	_, err := c.Response().Writer.Write([]byte("event: " + event + "\n"))
	if err != nil {
		return fmt.Errorf("write event: %v, err: %w", event, err)
	}
	_, err = c.Response().Writer.Write([]byte("data: " + data + "\n\n"))
	if err != nil {
		return fmt.Errorf("write data: %v, err: %w", data, err)
	}

	c.Response().Flush()
	return nil
}

var _ Context = (*HttpContext)(nil)

func (h *HttpContext) CurrentUser() *User {
	return GetUser(h)
}

func GetUser(c echo.Context) *User {
	jidStr := c.Request().Header.Get("jid")
	if jidStr == "" {
		return nil
	}
	jid, err := types.ParseJID(jidStr)
	if err != nil {
		return nil
	}
	return &User{JID: jid}
}
