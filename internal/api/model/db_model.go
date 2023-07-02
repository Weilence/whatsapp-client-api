package model

import (
	"database/sql/driver"
	"encoding/json"
	"log"
	"time"

	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"gorm.io/gorm"
)

type WhatsappAutoReply struct {
	gorm.Model
	JID  types.JID
	Key  string
	Type int
	Text string
	File string
}

type WhatsappChat struct {
	gorm.Model
	JID          string
	Name         string
	LastSendTime time.Time
	ReadOnly     bool
	UnreadCount  uint
	DeviceJID    string
}

type WhatsappChatMessage struct {
	gorm.Model
	MsgID     string
	ChatJID   string
	SenderJID string
	Msg       *Msg
	FromMe    bool
	SendTime  time.Time
	DeviceJID string
}

type Msg struct {
	*proto.Message
}

func (t *Msg) Scan(v interface{}) error {
	if v != nil {
		var msg Msg
		err := json.Unmarshal([]byte(v.(string)), &msg)
		if err != nil {
			return err
		}
		*t = msg
	}
	return nil
}

func (t *Msg) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
	}
	return string(b), err
}

type WhatsappQuickReply struct {
	gorm.Model
	Group string
	Text  string
}

type WhatsappSendMessage struct {
	gorm.Model
	From     types.JID
	To       string
	Type     int
	Text     string
	FileName string
}
