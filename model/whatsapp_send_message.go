package model

import "gorm.io/gorm"

type WhatsappSendMessage struct {
	gorm.Model
	From     string
	To       string
	Type     int
	Text     string
	FileName string
}
