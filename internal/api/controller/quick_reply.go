package controller

import (
	"github.com/weilence/whatsapp-client/internal/api"
	"github.com/weilence/whatsapp-client/internal/api/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	QueryQuickReplyReq struct {
		model.Pagination
		Text  string `form:"text"`
		Group string `form:"group"`
	}
	QueryQuickReplyRes struct {
		ID    uint   `json:"id"`
		Text  string `json:"text"`
		Group string `json:"group"`
	}
)

func QuickReplyQuery(c *api.HttpContext, req *QueryQuickReplyReq) (interface{}, error) {
	var list []QueryQuickReplyRes
	var total int64
	model.DB.Model(&model.WhatsappQuickReply{}).
		Scopes(
			model.Paginate(req.Pagination),
			model.WhereIf(len(req.Text) > 0, "`text` like ?", "%"+req.Text+"%"),
			model.WhereIf(len(req.Group) > 0, "`group` like ?", "%"+req.Group+"%"),
		).
		Count(&total).
		Find(&list)

	return gin.H{
		"total": total,
		"list":  list,
	}, nil
}

type ReplyAddReq struct {
	Text  string `json:"text"`
	Group string `json:"group"`
}

func QuickReplyAdd(c *api.HttpContext, req *QueryQuickReplyReq) (*struct{}, error) {
	model.DB.Save(&model.WhatsappQuickReply{
		Text:  req.Text,
		Group: req.Group,
	})

	return nil, nil
}

type ReplyEditReq struct {
	ID    uint   `uri:"id"`
	Text  string `json:"text"`
	Group string `json:"group"`
}

func QuickReplyEdit(c *api.HttpContext, req *ReplyEditReq) (interface{}, error) {
	model.DB.Save(&model.WhatsappQuickReply{
		Model: gorm.Model{
			ID: req.ID,
		},
		Text:  req.Text,
		Group: req.Group,
	})

	return nil, nil
}

type ReplyDeleteReq struct {
	ID uint `uri:"id"`
}

func QuickReplyDelete(c *api.HttpContext, req *ReplyEditReq) (interface{}, error) {
	model.DB.Unscoped().Delete(&model.WhatsappQuickReply{}, req.ID)
	return nil, nil
}
