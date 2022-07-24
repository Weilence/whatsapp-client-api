package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/weilence/whatsapp-client/internal/api/model"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"io"
)

var version = ""

type deviceLoginReq struct {
	Phone string `json:"phone"`
}

func DeviceLogin(c *gin.Context, req *deviceLoginReq) (_ struct{}, err error) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	client, qrChanItem, err := whatsapp.Login(c, req.Phone)
	if err != nil {
		c.SSEvent("error", err.Error())
		return
	}

	if qrChanItem == nil {
		c.SSEvent("success", client.Store.ID.String())
		return
	}
	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Writer.CloseNotify():
			return false
		case evt := <-qrChanItem:
			if evt.Event == "code" {
				c.SSEvent("message", evt.Code)
				return true
			} else if evt == whatsmeow.QRChannelSuccess {
				c.SSEvent("success", client.Store.ID.String())
				return false
			} else if evt == whatsmeow.QRChannelScannedWithoutMultidevice {
				c.SSEvent("error", "请开启多设备测试版")
				return false
			} else {
				c.SSEvent("error", "扫码登录失败")
				return false
			}
		}
	})
	return
}

type deviceLogoutReq struct {
	JID *types.JID `uri:"id"`
}

func DeviceLogout(c *gin.Context, req *deviceLogoutReq) (interface{}, error) {
	client, err := whatsapp.GetClient(req.JID)
	if err != nil {
		return nil, err
	}

	err = client.Logout()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

type Device struct {
	PushName     string `json:"pushName"`
	Platform     string `json:"platform"`
	Phone        string `json:"phone"`
	Jid          string `json:"jid"`
	BusinessName string `json:"businessName"`
	Online       bool   `json:"online"`
}

func DeviceQuery(c *gin.Context, _ *struct{}) (interface{}, error) {
	devices, err := whatsapp.GetDevices()
	if err != nil {
		return nil, err
	}

	data := make([]Device, len(devices))
	onlineClients := whatsapp.GetOnlineClients()

	for i, device := range devices {
		data[i] = Device{
			PushName:     device.PushName,
			Platform:     device.Platform,
			Phone:        device.ID.User,
			Jid:          device.ID.String(),
			BusinessName: device.BusinessName,
		}

		for _, client := range onlineClients {
			if client.Phone() == data[i].Phone {
				data[i].Online = true
				break
			}
		}
	}

	return data, nil
}

type deviceDeleteReq struct {
	JID *types.JID `uri:"jid"`
}

func DeviceDelete(c *gin.Context, req *deviceDeleteReq) (interface{}, error) {
	err := whatsapp.DeleteDevice(req.JID)
	if err != nil {
		return nil, err
	}

	model.DB.Delete(&model.WhatsappChat{}, "device_jid = ?", req.JID)
	model.DB.Delete(&model.WhatsappChatMessage{}, "device_jid = ?", req.JID)
	return nil, nil
}
