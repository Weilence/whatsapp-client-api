package api

import (
	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow/types"
	"whatsapp-client/model"
	"whatsapp-client/whatsapp"
)

type Device struct {
	PushName     string `json:"pushName"`
	Platform     string `json:"platform"`
	Phone        string `json:"phone"`
	Jid          string `json:"jid"`
	BusinessName string `json:"businessName"`
	Online       bool   `json:"online"`
}

func DeviceQuery(c *gin.Context) {
	devices, err := whatsapp.GetDevices()
	if err != nil {
		panic(err)
	}

	data := make([]Device, len(devices))
	clients := whatsapp.GetClients()

	for i, device := range devices {
		data[i] = Device{
			PushName:     device.PushName,
			Platform:     device.Platform,
			Phone:        device.ID.User,
			Jid:          device.ID.String(),
			BusinessName: device.BusinessName,
		}

		for _, client := range clients {
			if client.Phone() == data[i].Phone {
				data[i].Online = true
				break
			}
		}
	}

	c.JSON(0, data)
}

func DeviceDelete(c *gin.Context) {
	jid, err := types.ParseJID(c.Query("jid"))
	if err != nil {
		panic(err)
	}
	device, err := whatsapp.GetDevice(jid)
	if err != nil {
		panic(err)
	}
	err = whatsapp.DeleteDevice(device)
	c.JSON(0, nil)
	model.DB.Unscoped().Delete(&model.WhatsappChat{}, "device_jid = ?", jid)
	model.DB.Unscoped().Delete(&model.WhatsappChatMessage{}, "device_jid = ?", jid)
}
