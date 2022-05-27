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

// ClientLogin PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Param username path string true "username"
// @Param passwd path string true "passwd"
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func ClientLogin(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	id := c.Query("jid")

	client, qrItemChan := whatsapp.NewClient(id)

	if qrItemChan == nil {
		client.Login()
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
				client.Login()
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
