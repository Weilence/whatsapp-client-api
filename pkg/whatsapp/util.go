package whatsapp

import "go.mau.fi/whatsmeow/types"

func NewUserJID(phone string) types.JID {
	jid := types.NewJID(phone, types.DefaultUserServer)
	return jid
}

func NewGroupJID(id string) types.JID {
	jid := types.NewJID(id, types.GroupServer)
	return jid
}
