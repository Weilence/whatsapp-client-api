package model

import (
	"gorm.io/gorm"
	"time"
)

type WhatsappChat struct {
	gorm.Model
	JID          string
	Name         string
	LastSendTime time.Time
	ReadOnly     bool
	UnreadCount  uint
	DeviceJID    string
}
