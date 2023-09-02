package controller

import (
	"fmt"
	"io"
	"mime/multipart"

	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"
	"github.com/weilence/whatsapp-client/internal/utils"
	"go.mau.fi/whatsmeow/types"
)

type SendReq struct {
	JID   types.JID `query:"jid"`
	Phone string    `form:"phone"`
	Type  int       `form:"type"`
	Text  string    `form:"text"`
}

func MessageSend(c *utils.HttpContext, req *SendReq) (interface{}, error) {
	client, err := utils.GetClient(req.JID)
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
	default:
		if err = client.SendTextMessage(whatsapp.NewUserJID(req.Phone), req.Text); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func FormFileData(f *multipart.FileHeader) []byte {
	file, err := f.Open()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return bytes
}
