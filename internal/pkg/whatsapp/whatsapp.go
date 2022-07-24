package whatsapp

import (
	"database/sql"
	"log"

	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var container *sqlstore.Container

func init() {
	name := "Windows"
	store.DeviceProps.Os = &name
	store.DeviceProps.PlatformType = proto.DeviceProps_CHROME.Enum()
}

func Init(db *sql.DB) {
	logger := waLog.Stdout("Database", "DEBUG", true)
	container = sqlstore.NewWithDB(db, "sqlite3", logger)
	err := container.Upgrade()
	if err != nil {
		log.Panic(err)
	}
}
