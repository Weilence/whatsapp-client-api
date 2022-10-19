package controller

import (
	"github.com/weilence/whatsapp-client/internal/api"
	"go.mau.fi/whatsmeow/types"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weilence/whatsapp-client/internal/api/model"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"
	"github.com/weilence/whatsapp-client/pkg/utils"
)

type (
	MessagesReq struct {
		model.Pagination
	}
	MessagesRes struct {
		ID        uint      `json:"id,omitempty"`
		From      string    `json:"from,omitempty"`
		To        string    `json:"to,omitempty"`
		Type      int       `json:"type,omitempty"`
		Text      string    `json:"text,omitempty"`
		FileName  string    `json:"fileName,omitempty"`
		CreatedAt time.Time `json:"createdAt,omitempty"`
	}
)

type SendReq struct {
	JID   *types.JID           `form:"jid"`
	Phone string               `form:"phone"`
	Type  int                  `form:"type"`
	Text  string               `form:"text"`
	File  multipart.FileHeader `form:"file"`
}

func MessageSend(c *api.HttpContext, req *SendReq) (interface{}, error) {
	client, err := whatsapp.GetClient(req.JID)
	if err != nil {
		return nil, err
	}

	if req.File.Size == 0 {
		if err = client.SendTextMessage(req.Phone, req.Text); err != nil {
			return nil, err
		}

		model.DB.Save(&model.WhatsappSendMessage{
			From: *req.JID,
			To:   req.Phone,
			Type: req.Type,
			Text: req.Text,
		})
	} else {
		bytes := FormFileData(req.File)

		if req.Type == 1 {
			if err = client.SendImageMessage(req.Phone, bytes, req.Text); err != nil {
				return nil, err
			}
		} else if req.Type == 2 {
			if err = client.SendDocumentMessage(req.Phone, bytes, req.Text); err != nil {
				return nil, err
			}
		}
		model.DB.Save(&model.WhatsappSendMessage{
			From:     *req.JID,
			To:       req.Phone,
			Type:     req.Type,
			Text:     req.Text,
			FileName: req.File.Filename,
		})
	}
	return nil, nil
}

func MessageQuery(c *api.HttpContext, req *MessagesReq) (interface{}, error) {
	var list []MessagesRes
	var total int64
	model.DB.Model(&model.WhatsappSendMessage{}).
		Scopes(model.Paginate(req.Pagination)).
		Count(&total).
		Order("id desc").
		Find(&list)

	return gin.H{
		"total": total,
		"list":  list,
	}, nil
}

func FormFileData(f multipart.FileHeader) []byte {
	file, err := f.Open()
	defer utils.Close(file)
	if err != nil {
		log.Panic(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}

	return bytes
}
