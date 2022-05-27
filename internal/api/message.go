package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"time"
	"whatsapp-client/internal/model"
	"whatsapp-client/pkg/utils"
	"whatsapp-client/pkg/whatsapp"
)

type (
	MessagesReq struct {
		Pagination
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
	JID   string               `form:"jid"`
	Phone string               `form:"phone"`
	Type  int                  `form:"type"`
	Text  string               `form:"text"`
	File  multipart.FileHeader `form:"file"`
}

func MessageSend(c *gin.Context) {
	var req SendReq
	if err := c.Bind(&req); err != nil {
		return
	}

	jid := whatsapp.NewUserJID(req.Phone)
	client := whatsapp.GetClient(req.JID)

	if req.File.Size == 0 {
		client.SendTextMessage(jid, req.Text)

		model.DB.Save(&model.WhatsappSendMessage{
			From: req.JID,
			To:   jid.String(),
			Type: req.Type,
			Text: req.Text,
		})
	} else {
		bytes := FormFileData(req.File)

		if req.Type == 1 {
			client.SendImageMessage(jid, bytes, req.Text)
		} else if req.Type == 2 {
			client.SendDocumentMessage(jid, bytes, req.Text)
		}
		model.DB.Save(&model.WhatsappSendMessage{
			From:     req.JID,
			To:       jid.String(),
			Type:     req.Type,
			Text:     req.Text,
			FileName: req.File.Filename,
		})
	}
	c.JSON(0, nil)
}

func MessageQuery(c *gin.Context) {
	var req MessagesReq
	if err := c.Bind(&req); err != nil {
		return
	}

	var list []MessagesRes

	var total int64
	model.DB.Model(&model.WhatsappSendMessage{}).
		Count(&total).
		Order("id desc").
		Limit(req.Limit()).
		Offset(req.Offset()).
		Find(&list)

	c.JSON(0, gin.H{
		"total": total,
		"list":  list,
	})
}

func FormFileData(f multipart.FileHeader) []byte {
	file, err := f.Open()
	defer utils.Close(file)
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return bytes
}
