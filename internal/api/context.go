package api

import (
	"context"
	"github.com/gin-gonic/gin"
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
	*gin.Context
}

var _ Context = (*HttpContext)(nil)

func GetUser(c *gin.Context) *User {
	jidStr := c.GetHeader("jid")
	if jidStr == "" {
		return nil
	}
	jid, err := types.ParseJID(jidStr)
	if err != nil {
		return nil
	}
	return &User{JID: jid}
}

func (h *HttpContext) CurrentUser() *User {
	return GetUser(h.Context)
}
