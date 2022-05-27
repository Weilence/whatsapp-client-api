package model

import "gorm.io/gorm"

type WhatsappAutoReply struct {
	gorm.Model
	JID  string
	Key  string
	Type int
	Text string
	File string
}
