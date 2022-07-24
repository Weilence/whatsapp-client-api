package model

import (
	"database/sql"
	"log"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init() {
	db, err := gorm.Open(sqlite.Open("file:data.db?_foreign_keys=off"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			NameReplacer: strings.NewReplacer("JID", "Jid"),
		},
	})
	if err != nil {
		log.Panic(err)
	}

	db = db.Debug()
	err = AutoMigrate(db)
	if err != nil {
		log.Panic(err)
	}

	DB = db
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&WhatsappSendMessage{},
		&WhatsappQuickReply{},
		&WhatsappAutoReply{},
		&WhatsappChat{},
		&WhatsappChatMessage{},
	)
	return err
}

func SqlDB() *sql.DB {
	db, err := DB.DB()
	if err != nil {
		log.Panic(err)
	}
	return db
}
