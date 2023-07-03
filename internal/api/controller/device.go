package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/weilence/whatsapp-client/internal/api"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

var version = ""

type deviceLoginReq struct {
	JID types.JID `query:"jid"`
}

func DeviceLogin(c *api.HttpContext, req *deviceLoginReq) (_ struct{}, _ error) {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	client, err := whatsapp.NewClient(req.JID)
	if err != nil {
		c.SSEvent("error", err.Error())
		return
	}
	qrChanItem, err := client.Login(c)
	if err != nil {
		c.SSEvent("error", err.Error())
		return
	}

	var jid *types.JID
	if qrChanItem == nil {
		jid = client.Store.ID
	} else {
		for evt := range qrChanItem {
			if evt.Event == "code" {
				c.SSEvent("message", evt.Code)
			} else if evt == whatsmeow.QRChannelSuccess {
				jid = client.Store.ID
			} else if evt == whatsmeow.QRChannelScannedWithoutMultidevice {
				c.SSEvent("error", "请开启多设备测试版")
				return
			} else {
				c.SSEvent("error", "扫码登录失败")
				return
			}
		}
	}

	ticker := time.NewTicker(time.Second)
	timeout := time.After(time.Minute)
	for {
		select {
		case <-ticker.C:
			if client.IsLoggedIn() {
				c.SSEvent("success", jid.String())
				return
			}
		case <-c.Request().Context().Done():
			log.Println("context done")
			return
		case <-timeout:
			c.SSEvent("error", "连接超时")
			return
		}
	}
}

type deviceLogoutReq struct {
	JID types.JID `query:"jid"`
}

func DeviceLogout(c *api.HttpContext, req *deviceLogoutReq) (interface{}, error) {
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

type DeviceListRes struct {
	PushName     string    `json:"pushName"`
	Platform     string    `json:"platform"`
	Phone        string    `json:"phone"`
	Jid          types.JID `json:"jid"`
	BusinessName string    `json:"businessName"`
}

func DeviceList(c *api.HttpContext, _ *struct{}) (interface{}, error) {
	devices, err := whatsapp.GetDevices()
	if err != nil {
		return nil, err
	}

	data := make([]DeviceListRes, len(devices))
	for i, device := range devices {
		data[i] = DeviceListRes{
			PushName:     device.PushName,
			Platform:     device.Platform,
			Jid:          *device.ID,
			BusinessName: device.BusinessName,
		}
	}

	return data, nil
}

type deviceDeleteReq struct {
	JID *types.JID `query:"phone"`
}

func DeviceDelete(c *api.HttpContext, req *deviceDeleteReq) (interface{}, error) {
	err := whatsapp.DeleteDevice(req.JID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

type deviceStatusReq struct {
	JID types.JID `query:"jid"`
}

type deviceStatusRes struct {
	PushName     string `json:"pushName"`
	BusinessName string `json:"businessName"`
	Phone        string `json:"phone"`
	Status       string `json:"status"`
}

func DeviceStatus(c *api.HttpContext, req *deviceStatusReq) (*deviceStatusRes, error) {
	client, err := whatsapp.GetClient(req.JID)
	if err != nil {
		log.Println(fmt.Errorf("get client err: %w", err))
	}

	if client == nil {
		return &deviceStatusRes{Status: "disconnected"}, nil
	}

	res := &deviceStatusRes{
		Phone:        client.Store.ID.User,
		PushName:     client.Store.PushName,
		BusinessName: client.Store.BusinessName,
	}

	if !client.IsLoggedIn() {
		res.Status = "offline"
	} else {
		res.Status = "online"
	}

	return res, nil
}
