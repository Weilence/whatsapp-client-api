package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
	"whatsapp-client/internal/model"
)

type (
	QueryQuickReplyReq struct {
		Pagination
		Text  string `form:"text,omitempty"`
		Group string `form:"group,omitempty"`
	}
	QueryQuickReplyRes struct {
		ID    uint   `json:"id,omitempty"`
		Text  string `json:"text,omitempty"`
		Group string `json:"group,omitempty"`
	}
)

func QuickReplyQuery(c *gin.Context) {
	var req QueryQuickReplyReq
	err := c.Bind(&req)
	if err != nil {
		log.Println(err)
	}

	var list []QueryQuickReplyRes
	var total int64
	model.DB.Model(&model.WhatsappQuickReply{}).
		WhereIf(len(req.Text) > 0, "`text` like ?", "%"+req.Text+"%").
		WhereIf(len(req.Group) > 0, "`group` like ?", "%"+req.Group+"%").
		Count(&total).
		Limit(req.Limit()).
		Offset(req.Offset()).
		Find(&list)

	c.JSON(0, gin.H{
		"total": total,
		"list":  list,
	})
}

type ReplyAddReq struct {
	Text  string `json:"text,omitempty"`
	Group string `json:"group,omitempty"`
}

func QuickReplyAdd(c *gin.Context) {
	var req ReplyAddReq
	c.Bind(&req)

	model.DB.Save(&model.WhatsappQuickReply{
		Text:  req.Text,
		Group: req.Group,
	})

	c.JSON(0, nil)
}

type ReplyEditReq struct {
	ID    uint   `json:"id,omitempty"`
	Text  string `json:"text,omitempty"`
	Group string `json:"group,omitempty"`
}

func QuickReplyEdit(c *gin.Context) {
	var req ReplyEditReq
	c.Bind(&req)

	model.DB.Save(&model.WhatsappQuickReply{
		Model: gorm.Model{
			ID: req.ID,
		},
		Text:  req.Text,
		Group: req.Group,
	})

	c.JSON(0, nil)
}

func QuickReplyDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		panic(id)
	}
	model.DB.Unscoped().Delete(&model.WhatsappQuickReply{}, id)

	c.JSON(0, nil)
}
