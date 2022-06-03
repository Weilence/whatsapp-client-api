package whatsapp

import (
	"go.mau.fi/whatsmeow/store"
)

func GetDevices() ([]*store.Device, error) {
	return container.GetAllDevices()
}

func GetDevice(phone string) (device *store.Device, err error) {
	devices, err := GetDevices()
	for _, device := range devices {
		if device.ID.User == phone {
			return device, nil
		}
	}
	return nil, err
}

func DeleteDevice(device *store.Device) error {
	return container.DeleteDevice(device)
}
