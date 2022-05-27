package whatsapp

import (
	"context"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func (c Client) SendTextMessage(jid types.JID, text string) {
	c.SendMessage(jid, "", &waProto.Message{Conversation: proto.String(text)})
}

func (c Client) SendImageMessage(jid types.JID, image []byte, caption string) {
	uploaded, err := c.Upload(context.Background(), image, whatsmeow.MediaImage)
	if err != nil {
		panic(err)
	}
	c.SendMessage(jid, "", &waProto.Message{
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
}

func (c Client) SendDocumentMessage(jid types.JID, file []byte, filename string) {
	uploaded, err := c.Upload(context.Background(), file, whatsmeow.MediaDocument)
	if err != nil {
		panic(err)
	}
	c.SendMessage(jid, "", &waProto.Message{
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
}
