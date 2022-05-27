package whatsapp

import (
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
)

func GetDevices() ([]*store.Device, error) {
	return container.GetAllDevices()
}

func GetDevice(jid types.JID) (*store.Device, error) {
	return container.GetDevice(jid)
}

func DeleteDevice(device *store.Device) error {
	return container.DeleteDevice(device)
}
