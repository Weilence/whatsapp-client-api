package controller

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/weilence/whatsapp-client/internal/api"
	"go.mau.fi/whatsmeow/types"

	"github.com/weilence/whatsapp-client/internal/api/model"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"
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
	JID   types.JID `query:"jid"`
	Phone string    `form:"phone"`
	Type  int       `form:"type"`
	Text  string    `form:"text"`
}

func MessageSend(c *api.HttpContext, req *SendReq) (interface{}, error) {
	client, err := whatsapp.GetClient(req.JID)
	if err != nil {
		return nil, err
	}

	var filename string
	switch req.Type {
	case 1:
		image, err := c.FormFile("image")
		if err != nil {
			return nil, fmt.Errorf("get image error: %w", err)
		}

		bytes := FormFileData(image)
		if err = client.SendImageMessage(whatsapp.NewUserJID(req.Phone), bytes, req.Text); err != nil {
			return nil, err
		}
	case 2:
		file, err := c.FormFile("file")
		if err != nil {
			return nil, fmt.Errorf("get image error: %w", err)
		}

		filename = file.Filename
		bytes := FormFileData(file)
		if err = client.SendDocumentMessage(whatsapp.NewUserJID(req.Phone), bytes, filename, req.Text); err != nil {
			return nil, err
		}
	}

	db := model.DB.Save(&model.WhatsappSendMessage{
		From:     req.JID,
		To:       req.Phone,
		Type:     req.Type,
		Text:     req.Text,
		FileName: filename,
	})
	if db.Error != nil {
		return nil, fmt.Errorf("save message error: %w", db.Error)
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

	return model.ResponseList{
		Total: total,
		List:  list,
	}, nil
}

func FormFileData(f *multipart.FileHeader) []byte {
	file, err := f.Open()
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}

	return bytes
}
