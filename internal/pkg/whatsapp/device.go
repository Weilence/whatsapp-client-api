package whatsapp

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
	"github.com/samber/lo"
	"github.com/weilence/whatsapp-client/config"
	"go.mau.fi/whatsmeow/proto/waCompanionReg"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var container *sqlstore.Container

func init() {
	store.DeviceProps.Os = lo.ToPtr("Windows")
	store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_CHROME.Enum()
}

func Setup() {
	db, err := sql.Open("sqlite", "data.db")
	if err != nil {
		panic(err)
	}
	var logger waLog.Logger
	if *config.Env == "dev" {
		logger = waLog.Stdout("Database", "DEBUG", true)
	} else {
		logger = waLog.Stdout("Database", "INFO", false)
	}

	container = sqlstore.NewWithDB(db, "sqlite3", logger)
	err = container.Upgrade()
	if err != nil {
		panic(err)
	}
}

func GetDevices() ([]*store.Device, error) {
	return container.GetAllDevices()
}

func GetDevice(jid types.JID) (*store.Device, error) {
	devices, err := container.GetAllDevices()
	if err != nil {
		return nil, err
	}
	device, _ := lo.Find(devices, func(item *store.Device) bool {
		return *item.ID == jid
	})
	return device, nil
}

func DeleteDevice(jid *types.JID) error {
	devices, err := container.GetAllDevices()
	if err != nil {
		return err
	}

	device, b := lo.Find(devices, func(item *store.Device) bool {
		return *item.ID == *jid
	})

	if !b {
		return nil
	}

	return container.DeleteDevice(device)
}
