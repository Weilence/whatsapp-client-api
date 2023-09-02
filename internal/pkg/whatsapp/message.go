package whatsapp

import (
	"context"
	"fmt"
	"net/http"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func (c *Client) SendTextMessage(jid types.JID, text string) error {
	_, err := c.SendMessage(context.Background(), jid.ToNonAD(), &waE2E.Message{Conversation: &text})
	if err != nil {
		return fmt.Errorf("send text message failed: %w", err)
	}
	return nil
}

func (c *Client) SendImageMessage(jid types.JID, image []byte, text string) error {
	uploaded, err := c.Upload(context.Background(), image, whatsmeow.MediaImage)
	if err != nil {
		return fmt.Errorf("upload image failed: %w", err)
	}
	mimeType := http.DetectContentType(image)
	fileLength := uint64(len(image))
	_, err = c.SendMessage(context.Background(), jid.ToNonAD(), &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption:       &text,
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &fileLength,
		},
	})
	if err != nil {
		return fmt.Errorf("send image message failed: %w", err)
	}

	return nil
}

func (c *Client) SendDocumentMessage(jid types.JID, file []byte, filename, text string) error {
	uploaded, err := c.Upload(context.Background(), file, whatsmeow.MediaDocument)
	if err != nil {
		return fmt.Errorf("upload document failed: %w", err)
	}
	mimeType := http.DetectContentType(file)
	fileLength := uint64(len(file))
	_, err = c.SendMessage(context.Background(), jid.ToNonAD(), &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			Caption:       &text,
			FileName:      &filename,
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &fileLength,
		},
	})
	if err != nil {
		return fmt.Errorf("send document message failed: %w", err)
	}

	return nil
}
