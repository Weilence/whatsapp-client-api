package api

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow"
	"io"
	"log"
	"strings"
	"whatsapp-client/pkg/whatsapp"
)

var version = ""

func ClientLogin(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	phone := c.Query("phone")
	client := whatsapp.NewClient(phone)

	qrItemChan := client.Login()

	if qrItemChan == nil {
		c.SSEvent("success", client.Store.ID.String())
		return
	}
	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Writer.CloseNotify():
			return false
		case evt := <-qrItemChan:
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
}

func ClientLogout(c *gin.Context) {
	phone := c.Query("phone")
	client, err := whatsapp.GetClient(phone)
	if err != nil {
		BadRequest(c, err)
		return
	}

	err = client.Logout()
	if err != nil {
		BadRequest(c, err)
		return
	}
	Ok(c, nil)
}

func ClientInfo(c *gin.Context) {
	machineId, err := machineid.ProtectedID("whatsapp-client")
	if err != nil {
		log.Fatal(err)
	}
	machineId = strings.ToUpper(machineId[:16])

	c.JSON(0, gin.H{
		"machineId": machineId,
		"version":   version,
	})
}
