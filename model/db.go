package model

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

var DB *WhatsappDB

type WhatsappDB struct {
	*gorm.DB
}

func NewWhatsappDB(db *gorm.DB) *WhatsappDB {
	return &WhatsappDB{
		db,
	}
}

func init() {
	db, err := gorm.Open(sqlite.Open("file:data.db?_foreign_keys=off"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			NameReplacer: strings.NewReplacer("JID", "Jid"),
		},
	})
	if err != nil {
		panic(err)
	}

	DB = NewWhatsappDB(db)
	DB.AutoMigrate(&WhatsappSendMessage{}, &WhatsappQuickReply{}, &WhatsappAutoReply{}, &WhatsappChat{}, &WhatsappChatMessage{})
	if gin.IsDebugging() {
		DB.DB = DB.Debug()
	}
}

func (db *WhatsappDB) Model(value interface{}) *WhatsappDB {
	return &WhatsappDB{
		db.DB.Model(value),
	}
}

func (db *WhatsappDB) WhereIf(condition bool, query interface{}, args ...interface{}) *WhatsappDB {
	if condition {
		return &WhatsappDB{
			db.DB.Where(query, args),
		}
	}
	return db
}
