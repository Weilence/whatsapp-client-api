package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/weilence/whatsapp-client/internal/api"
	"github.com/weilence/whatsapp-client/internal/api/model"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"
	"go.mau.fi/whatsmeow/types"
	"gorm.io/gorm"
)

type AutoReplyQueryReq struct {
	model.Pagination
}

type AutoReplyQueryRes struct {
	ID   uint   `json:"id"`
	JID  string `json:"jid"`
	Key  string `json:"key"`
	Type int    `json:"type"`
	Text string `json:"text"`
	File string `json:"file"`
}

func AutoReplyQuery(c *api.HttpContext, req *AutoReplyQueryReq) (interface{}, error) {
	var total int64
	var list []AutoReplyQueryRes

	model.DB.
		Model(&model.WhatsappAutoReply{}).
		Find(&list).
		Scopes(model.Paginate(req.Pagination)).
		Count(&total)

	return gin.H{
		"total": total,
		"list":  list,
	}, nil
}

type AutoReplyAddReq struct {
	JID  *types.JID `json:"jid"`
	Key  string     `json:"key"`
	Text string     `json:"text"`
	File string     `json:"file"`
}

func AutoReplyAdd(c *api.HttpContext, req *AutoReplyAddReq) (interface{}, error) {
	model.DB.Save(&model.WhatsappAutoReply{
		JID:  *req.JID,
		Key:  req.Key,
		Text: req.Text,
		File: req.File,
	})

	client, _ := whatsapp.GetClient(req.JID)
	client.RefreshAutoReplay()

	return nil, nil
}

type AutoReplyEditReq struct {
	ID   uint       `uri:"id"`
	JID  *types.JID `json:"jid"`
	Key  string     `json:"key"`
	Text string     `json:"text"`
	File string     `json:"file"`
}

func AutoReplyEdit(c *api.HttpContext, req *AutoReplyEditReq) (interface{}, error) {
	model.DB.Save(&model.WhatsappAutoReply{
		Model: gorm.Model{ID: req.ID},
		Key:   req.Key,
		JID:   *req.JID,
		Text:  req.Text,
		File:  req.File,
	})

	client, _ := whatsapp.GetClient(req.JID)
	client.RefreshAutoReplay()

	return nil, nil
}

type AutoReplyDeleteReq struct {
	ID uint `uri:"id"`
}

func AutoReplyDelete(c *api.HttpContext, req *AutoReplyEditReq) (interface{}, error) {
	var m model.WhatsappAutoReply
	model.DB.Find(m, req.ID)
	model.DB.Unscoped().Delete(&model.WhatsappAutoReply{}, req.ID)

	client, _ := whatsapp.GetClient(&m.JID)
	client.RefreshAutoReplay()

	return nil, nil
}
