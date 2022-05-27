package whatsapp

import (
	"context"
	"github.com/mattn/go-ieproxy"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	log "go.mau.fi/whatsmeow/util/log"
	"time"
	"whatsapp-client/internal/model"
)

type Client struct {
	*whatsmeow.Client
	groups    []*types.GroupInfo
	autoReply AutoReply
	DeviceJID string
}

var onlineClients []*Client

func NewClient(id string) (*Client, <-chan whatsmeow.QRChannelItem) {
	var device *store.Device
	if id == "" {
		device = container.NewDevice()
	} else {
		jid, err := types.ParseJID(id)
		if err != nil {
			panic(err)
		}
		device, err = container.GetDevice(jid)

		if err != nil {
			panic(err)
		}
	}

	clientLog := log.Stdout("Client", "DEBUG", true)
	client := &Client{Client: whatsmeow.NewClient(device, clientLog)}
	client.SetProxy(ieproxy.GetProxyFunc())
	client.EnableAutoReply()
	client.AddEventHandler(func(rawEvt interface{}) {
		switch evt := rawEvt.(type) {
		case *events.Connected:
			client.DeviceJID = client.Store.ID.String()
		case *events.ClientOutdated:
		case *events.LoggedOut:
			model.DB.Unscoped().Delete(&model.WhatsappChat{}, "device_jid = ?", client.DeviceJID)
			model.DB.Unscoped().Delete(&model.WhatsappChatMessage{}, "device_jid = ?", client.DeviceJID)
			for i, c := range onlineClients {
				if c.Store.ID == client.Store.ID {
					onlineClients = append(onlineClients[:i], onlineClients[i+1:]...)
					return
				}
			}
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
							LastSendTime: time.Unix(int64(*c.ConversationTimestamp), 0),
							ReadOnly:     false,
							UnreadCount:  uint(*c.UnreadCount),
							DeviceJID:    client.DeviceJID,
						})
					} else {
						model.DB.Save(&model.WhatsappChat{
							JID:          c.GetId(),
							Name:         c.GetName(),
							LastSendTime: time.Unix(int64(*c.ConversationTimestamp), 0),
							ReadOnly:     false,
							UnreadCount:  uint(*c.UnreadCount),
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
								SendTime:  time.Unix(int64(m.Message.GetMessageTimestamp()), 0),
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

	if client.Store.ID == nil {
		c, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			panic(err)
		}
		return client, c
	} else {
		err := client.Connect()
		if err != nil {
			panic(err)
		}
		return client, nil
	}
}

func GetClients() []*Client {
	return onlineClients
}

func GetClient(id string) *Client {
	for _, client := range onlineClients {
		if client.Store.ID.String() == id {
			return client
		}
	}
	return nil
}

func (c *Client) Login() {
	onlineClients = append(onlineClients, c)
}

func (c *Client) Logout() {
	c.Disconnect()
	for i, client := range onlineClients {
		if client.Store.ID.String() == c.Store.ID.String() {
			onlineClients = append(onlineClients[:i], onlineClients[i+1:]...)
			return
		}
	}
}

func (c *Client) Phone() string {
	return c.Store.ID.User
}
