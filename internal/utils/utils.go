package utils

import (
	"io"

	"log/slog"

	"github.com/samber/lo"
	"github.com/weilence/whatsapp-client/internal/model"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"
	"go.mau.fi/whatsmeow/types"
)

func Close(closer io.Closer) {
	if closer != nil {
		err := closer.Close()
		if err != nil {
			slog.Error("close", "err", err)
		}
	}
}

func GetClient(jid types.JID) (*whatsapp.Client, error) {
	client, ok := lo.Find(whatsapp.GetClients(), func(item *whatsapp.Client) bool {
		return item.Store.ID != nil && *item.Store.ID == jid
	})
	if !ok {
		return nil, model.ResponseModel{
			Code:    1000,
			Message: "客户端未登录",
			Data:    nil,
		}
	}
	return client, nil
}
