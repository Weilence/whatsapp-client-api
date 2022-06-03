package whatsapp

import (
	"context"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func (c Client) SendTextMessage(phone, text string) (err error) {
	_, err = c.SendMessage(NewUserJID(phone), "", &waProto.Message{Conversation: proto.String(text)})
	return
}

func (c Client) SendImageMessage(phone string, image []byte, caption string) (err error) {
	uploaded, err := c.Upload(context.Background(), image, whatsmeow.MediaImage)
	if err != nil {
		panic(err)
	}
	_, err = c.SendMessage(NewUserJID(phone), "", &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Caption:       proto.String(caption),
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(image)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(image))),
		}})
	return
}

func (c Client) SendDocumentMessage(phone string, file []byte, filename string) (err error) {
	uploaded, err := c.Upload(context.Background(), file, whatsmeow.MediaDocument)
	if err != nil {
		panic(err)
	}
	_, err = c.SendMessage(NewUserJID(phone), "", &waProto.Message{
		DocumentMessage: &waProto.DocumentMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(file)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(file))),
			FileName:      &filename,
		},
	})
	return
}
