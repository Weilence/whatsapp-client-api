package whatsapp

import (
	"database/sql"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"whatsapp-client/pkg/utils"
)

var container *sqlstore.Container

func init() {
	name := "Windows"
	store.DeviceProps.Os = &name
	store.DeviceProps.PlatformType = waProto.CompanionProps_CHROME.Enum()
}

func Init(db *sql.DB) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	container = sqlstore.NewWithDB(db, "sqlite3", dbLog)
	err := container.Upgrade()
	utils.NoError(err)
}
