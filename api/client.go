package api

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow"
	"io"
	"log"
	"strings"
	"whatsapp-client/whatsapp"
)

var version = ""

func ClientLogin(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	id := c.Query("jid")

	client, qrItemChan := whatsapp.NewClient(id)

	if qrItemChan == nil {
		client.Login()

		sse.Encode(c.Writer, sse.Event{
			Event: "success",
			Data:  client.Client.Store.ID.String(),
		})
		return
	}
	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Writer.CloseNotify():
			return false
		case evt := <-qrItemChan:
			if evt.Event == "code" {
				sse.Encode(w, sse.Event{
					Event: "message",
					Data:  evt.Code,
				})
				return true
			} else if evt == whatsmeow.QRChannelSuccess {
				client.Login()
				sse.Encode(c.Writer, sse.Event{
					Event: "success",
					Data:  client.Client.Store.ID.String(),
				})
				return false
			} else if evt == whatsmeow.QRChannelScannedWithoutMultidevice {
				sse.Encode(c.Writer, sse.Event{
					Event: "error",
					Data:  "请开启多设备测试版",
				})
				return false
			} else {
				sse.Encode(c.Writer, sse.Event{
					Event: "error",
					Data:  "扫码登录失败",
				})
				return false
			}
		}
	})
}

func ClientLogout(c *gin.Context) {
	client := whatsapp.GetClient(c.Query("jid"))
	client.Logout()
	c.JSON(0, nil)
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
