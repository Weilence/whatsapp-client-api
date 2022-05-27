package model

import (
	"database/sql/driver"
	"encoding/json"
	"go.mau.fi/whatsmeow/binary/proto"
	"gorm.io/gorm"
	"log"
	"time"
)

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
