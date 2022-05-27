package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"whatsapp-client/internal/model"
	"whatsapp-client/pkg/whatsapp"
)

type AutoReplyQueryReq struct {
	Pagination
}

type AutoReplyQueryRes struct {
	ID   uint   `json:"id,omitempty"`
	JID  string `json:"jid,omitempty"`
	Key  string `json:"key,omitempty"`
	Type int    `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
	File string `json:"file,omitempty"`
}

func AutoReplyQuery(c *gin.Context) {
	var req AutoReplyQueryReq
	c.Bind(&req)

	var total int64
	var list []AutoReplyQueryRes

	model.DB.
		Model(&model.WhatsappAutoReply{}).
		Limit(req.Limit()).
		Offset(req.Offset()).
		Find(&list).
		Count(&total)

	c.JSON(0, gin.H{
		"total": total,
		"list":  list,
	})
}

type AutoReplyAddReq struct {
	JID  string `json:"jid,omitempty"`
	Key  string `json:"key,omitempty"`
	Text string `json:"text,omitempty"`
	File string `json:"file,omitempty"`
}

func AutoReplyAdd(c *gin.Context) {
	var req AutoReplyAddReq
	c.Bind(&req)

	model.DB.Save(&model.WhatsappAutoReply{
		JID:  req.JID,
		Key:  req.Key,
		Text: req.Text,
		File: req.File,
	})

	whatsapp.GetClient(req.JID).RefreshAutoReplay()
	c.JSON(0, nil)
}

type AutoReplyEditReq struct {
	ID   uint   `json:"id,omitempty"`
	JID  string `json:"jid,omitempty"`
	Key  string `json:"key,omitempty"`
	Text string `json:"text,omitempty"`
	File string `json:"file,omitempty"`
}

func AutoReplyEdit(c *gin.Context) {
	var req AutoReplyEditReq
	c.Bind(&req)

	model.DB.Save(&model.WhatsappAutoReply{
		Model: gorm.Model{ID: req.ID},
		Key:   req.Key,
		JID:   req.JID,
		Text:  req.Text,
		File:  req.File,
	})

	whatsapp.GetClient(req.JID).RefreshAutoReplay()
	c.JSON(0, nil)
}

func AutoReplyDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		panic(err)
	}
	model.DB.Unscoped().Delete(&model.WhatsappAutoReply{}, id)

	jid := c.Query("jid")
	whatsapp.GetClient(jid).RefreshAutoReplay()
	c.JSON(0, nil)
}
