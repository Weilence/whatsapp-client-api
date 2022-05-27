package model

import (
	"database/sql"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"whatsapp-client/pkg/orm"
	"whatsapp-client/pkg/utils"
)

var DB *orm.DB

func Init() {
	db, err := gorm.Open(sqlite.Open("file:data.db?_foreign_keys=off"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			NameReplacer: strings.NewReplacer("JID", "Jid"),
		},
	})
	utils.NoError(err)

	if viper.GetBool("debug") {
		db = db.Debug()
	}
	err = AutoMigrate(db)
	utils.NoError(err)

	DB = &orm.DB{
		DB: db,
	}
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&WhatsappSendMessage{}, &WhatsappQuickReply{}, &WhatsappAutoReply{}, &WhatsappChat{}, &WhatsappChatMessage{})
	return err
}

func SqlDB() *sql.DB {
	db, err := DB.DB.DB()
	utils.NoError(err)
	return db
}
