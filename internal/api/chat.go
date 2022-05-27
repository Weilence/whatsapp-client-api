package api

import (
	"github.com/gin-gonic/gin"
	"time"
	"whatsapp-client/internal/model"
)

type (
	ChatQueryReq struct {
		Pagination
		JID string `form:"jid"`
	}
	ChatQueryRes struct {
		JID          string    `json:"jid"`
		Name         string    `json:"name"`
		LastSendTime time.Time `json:"lastSendTime"`
		ReadOnly     bool      `json:"readOnly"`
		UnreadCount  uint      `json:"unreadCount"`
	}
)

func ChatQuery(c *gin.Context) {
	var req ChatQueryReq
	err := c.Bind(&req)
	if err != nil {
		panic(err)
	}

	var list []ChatQueryRes
	var total int64
	model.DB.Model(&model.WhatsappChat{}).
		Where("device_jid = ?", req.JID).
		Count(&total).
		Limit(req.Limit()).
		Offset(req.Offset()).
		Find(&list)

	c.JSON(0, gin.H{
		"total": total,
		"list":  list,
	})
}
