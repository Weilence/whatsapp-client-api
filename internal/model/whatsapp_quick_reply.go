package model

import "gorm.io/gorm"

type WhatsappQuickReply struct {
	gorm.Model
	Group string
	Text  string
}
