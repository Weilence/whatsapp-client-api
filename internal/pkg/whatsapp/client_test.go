package whatsapp

import (
	"context"
	"github.com/mdp/qrterminal/v3"
	"github.com/weilence/whatsapp-client/internal/api/model"
	"go.mau.fi/whatsmeow"
	"os"
	"testing"
	"time"
)

const phone = ""

func TestLogin(t *testing.T) {
	model.Init()
	Init(model.SqlDB())

	ctx := context.Background()
	client, ch, err := Login(ctx, phone)
	if err != nil {
		t.Error(err)
		return
	}

	if ch != nil {
		for item := range ch {
			if item.Event == "code" {
				qrterminal.GenerateHalfBlock(item.Code, qrterminal.L, os.Stdout)
			} else if item == whatsmeow.QRChannelSuccess {
				t.Logf("登录成功, JID: %s", client.Store.ID.String())
				break
			} else if item == whatsmeow.QRChannelScannedWithoutMultidevice {
				t.Error("请开启多设备测试版")
			} else {
				t.Error("扫码登录失败")
			}
		}

		<-time.After(time.Second * 20)
	}
}
