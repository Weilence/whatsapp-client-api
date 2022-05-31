package whatsapp

import (
	"context"
	"errors"
	"github.com/mattn/go-ieproxy"
	"github.com/spf13/cast"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.uber.org/atomic"
	"whatsapp-client/internal/model"
	"whatsapp-client/pkg/utils"
)

type Client struct {
	*whatsmeow.Client
	groups    []*types.GroupInfo
	autoReply AutoReply
	DeviceJID string
}

var logId atomic.Uint32

func getClientLog() waLog.Logger {
	return waLog.Stdout("Client"+cast.ToString(logId.Inc()), "DEBUG", true)
}

func NewClient(id string) *Client {
	jid, err := types.ParseJID(id)
	utils.NoError(err)
	device, err := container.GetDevice(jid)
	utils.NoError(err)
	if device == nil {
		device = container.NewDevice()
	}

	client := &Client{Client: whatsmeow.NewClient(device, getClientLog())}
	client.EnableAutoReply()

	client.SetProxy(ieproxy.GetProxyFunc())
	client.AddEventHandler(func(rawEvt interface{}) {
		switch evt := rawEvt.(type) {
		case *events.Connected:
			client.DeviceJID = client.Store.ID.String()
			onlineClientAdd(client)
		case *events.Disconnected:
			onlineClientRemove(client)
		case *events.ClientOutdated:
		case *events.LoggedOut:
			clearChat(client.DeviceJID)
			onlineClientRemove(client)
		case *events.Message:
			var cm model.WhatsappChatMessage
			model.DB.Find(&cm, "msg_id = ? AND device_jid = ?", evt.Info.ID, client.DeviceJID)
			if cm.ID == 0 {
				model.DB.Save(&model.WhatsappChatMessage{
					MsgID:     evt.Info.ID,
					ChatJID:   evt.Info.Chat.String(),
					SenderJID: evt.Info.Sender.String(),
					Msg: &model.Msg{
						Message: evt.Message,
					},
					FromMe:    evt.Info.IsFromMe,
					SendTime:  evt.Info.Timestamp,
					DeviceJID: client.DeviceJID,
				})
			}
		case *events.HistorySync:
			if *evt.Data.SyncType == proto.HistorySync_INITIAL_BOOTSTRAP || *evt.Data.SyncType == proto.HistorySync_RECENT {
				for _, c := range evt.Data.Conversations {
					var chat model.WhatsappChat
					model.DB.Find(&chat, "jid = ? AND device_jid = ?", c.Id, client.DeviceJID)
					if chat.ID > 0 {
						model.DB.Model(chat).Updates(&model.WhatsappChat{
							JID:          c.GetId(),
							Name:         c.GetName(),
							LastSendTime: cast.ToTime(c.ConversationTimestamp),
							ReadOnly:     false,
							UnreadCount:  cast.ToUint(c.UnreadCount),
							DeviceJID:    client.DeviceJID,
						})
					} else {
						model.DB.Save(&model.WhatsappChat{
							JID:          c.GetId(),
							Name:         c.GetName(),
							LastSendTime: cast.ToTime(c.ConversationTimestamp),
							ReadOnly:     false,
							UnreadCount:  cast.ToUint(c.UnreadCount),
							DeviceJID:    client.DeviceJID,
						})
					}
					for _, m := range c.Messages {
						if m.Message.Message == nil {
							continue
						}

						var cm model.WhatsappChatMessage
						model.DB.Find(&cm, "msg_id = ? AND device_jid = ?", m.Message.Key.GetId(), client.DeviceJID)
						if cm.ID == 0 {
							msg := &model.WhatsappChatMessage{
								MsgID:     m.Message.Key.GetId(),
								ChatJID:   m.Message.Key.GetRemoteJid(),
								SenderJID: m.Message.GetParticipant(),
								Msg: &model.Msg{
									Message: m.Message.GetMessage(),
								},
								FromMe:    m.Message.Key.GetFromMe(),
								SendTime:  cast.ToTime(m.Message.GetMessageTimestamp()),
								DeviceJID: client.DeviceJID,
							}
							if msg.FromMe {
								msg.SenderJID = client.Store.ID.ToNonAD().String()
							}
							model.DB.Save(msg)
						}
					}
				}
			}
		}
	})

	return client
}

func clearChat(jid string) {
	model.DB.Unscoped().Delete(&model.WhatsappChat{}, "device_jid = ?", jid)
	model.DB.Unscoped().Delete(&model.WhatsappChatMessage{}, "device_jid = ?", jid)
}

func GetClients() []*Client {
	return onlineClients
}

func GetClient(id string) (*Client, error) {
	for _, client := range onlineClients {
		if client.Store.ID.String() == id {
			return client, nil
		}
	}
	return nil, errors.New("客户端已离线")
}

func (c *Client) Login() <-chan whatsmeow.QRChannelItem {
	if c.Store.ID != nil {
		ch := make(chan bool)
		handlerID := c.AddEventHandler(func(evt interface{}) {
			switch evt.(type) {
			case *events.Connected:
				ch <- true
			case *events.ClientOutdated:
				ch <- false
			}
		})
		utils.NoError(c.Connect())
		if <-ch {
			return nil
		}
		c.RemoveEventHandler(handlerID)
	}

	qrChan, err := c.GetQRChannel(context.Background())
	utils.NoError(err)
	utils.NoError(c.Connect())
	return qrChan
}

func (c *Client) Logout() error {
	for i, client := range onlineClients {
		if client.Store.ID.String() == client.Store.ID.String() {
			onlineClients = append(onlineClients[:i], onlineClients[i+1:]...)
			break
		}
	}

	err := c.Client.Logout()
	if err == nil {
		clearChat(c.DeviceJID)
	}
	return err
}

func (c *Client) Phone() string {
	return c.Store.ID.User
}

var onlineClients []*Client

func onlineClientAdd(client *Client) {
	onlineClients = append(onlineClients, client)
}

func onlineClientRemove(client *Client) {
	for i, c := range onlineClients {
		if c.Store.ID == client.Store.ID {
			onlineClients = append(onlineClients[:i], onlineClients[i+1:]...)
			return
		}
	}
}
