package whatsapp

import "go.mau.fi/whatsmeow/types"

func NewUserJID(phone string) types.JID {
	jid := types.NewJID(phone, types.DefaultUserServer)
	return jid
}
