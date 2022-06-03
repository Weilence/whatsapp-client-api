package api

import (
	"github.com/gin-gonic/gin"
	"whatsapp-client/internal/model"
	"whatsapp-client/pkg/utils"
	"whatsapp-client/pkg/whatsapp"
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
	utils.NoError(err)

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
	phone := c.Query("phone")
	device, err := whatsapp.GetDevice(phone)
	if err != nil {
		BadRequest(c, err)
		return
	}

	err = whatsapp.DeleteDevice(device)
	if err != nil {
		BadRequest(c, err)
		return
	}

	jid := whatsapp.NewUserJID(phone)
	model.DB.Delete(&model.WhatsappChat{}, "device_jid = ?", jid)
	model.DB.Delete(&model.WhatsappChatMessage{}, "device_jid = ?", jid)
	Ok(c, nil)
}
