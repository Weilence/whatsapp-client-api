package whatsapp

import (
	"context"
	"log"
	"net/http"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
)

func (c *Client) SendTextMessage(phone, text string) (err error) {
	_, err = c.SendMessage(context.Background(), NewUserJID(phone), &proto.Message{Conversation: &text})
	return
}

func (c *Client) SendImageMessage(phone string, image []byte, caption string) (err error) {
	uploaded, err := c.Upload(context.Background(), image, whatsmeow.MediaImage)
	if err != nil {
		log.Panic(err)
	}
	mimeType := http.DetectContentType(image)
	fileLength := uint64(len(image))
	_, err = c.SendMessage(context.Background(), NewUserJID(phone), &proto.Message{
		ImageMessage: &proto.ImageMessage{
			Caption:       &caption,
			Url:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    &fileLength,
		},
	})
	return
}

func (c *Client) SendDocumentMessage(phone string, file []byte, filename string) (err error) {
	uploaded, err := c.Upload(context.Background(), file, whatsmeow.MediaDocument)
	if err != nil {
		log.Panic(err)
	}
	mimeType := http.DetectContentType(file)
	fileLength := uint64(len(file))
	_, err = c.SendMessage(context.Background(), NewUserJID(phone), &proto.Message{
		DocumentMessage: &proto.DocumentMessage{
			Url:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    &fileLength,
			FileName:      &filename,
		},
	})
	return
}
