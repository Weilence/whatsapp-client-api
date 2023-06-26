package controller

import (
	"time"

	"github.com/weilence/whatsapp-client/internal/api"

	"github.com/weilence/whatsapp-client/internal/api/model"
)

type (
	ChatQueryReq struct {
		model.Pagination
		JID string `query:"jid"`
	}
	ChatQueryRes struct {
		JID          string    `json:"jid"`
		Name         string    `json:"name"`
		LastSendTime time.Time `json:"lastSendTime"`
		ReadOnly     bool      `json:"readOnly"`
		UnreadCount  uint      `json:"unreadCount"`
	}
)

func ChatQuery(c *api.HttpContext, req *ChatQueryReq) (interface{}, error) {
	var list []*ChatQueryRes
	var total int64

	model.DB.Model(&model.WhatsappChat{}).
		Where("device_jid = ?", req.JID).
		Count(&total).
		Scopes(model.Paginate(req.Pagination)).
		Find(&list)

	return model.ResponseList{
		Total: total,
		List:  list,
	}, nil
}
