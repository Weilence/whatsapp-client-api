package whatsapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/mattn/go-ieproxy"
	"github.com/samber/lo"
	"github.com/weilence/whatsapp-client/internal/api/model"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Client struct {
	*whatsmeow.Client
	groups    []*types.GroupInfo
	autoReply AutoReply
}

func GetDevice(phone string) (*store.Device, error) {
	devices, err := container.GetAllDevices()
	if err != nil {
		return nil, err
	}
	device, _ := lo.Find(devices, func(item *store.Device) bool {
		return item.ID.User == phone
	})
	return device, nil
}

func newClient(phone string) (*Client, error) {
	device, err := GetDevice(phone)
	if err != nil {
		return nil, fmt.Errorf("device get error: %w", err)
	}

	if device == nil {
		device = container.NewDevice()
	}

	logger := waLog.Stdout("Client_"+phone, "DEBUG", true)
	client := &Client{Client: whatsmeow.NewClient(device, logger)}
	client.EnableAutoReply()
	client.SetProxy(ieproxy.GetProxyFunc())
	client.AddEventHandler(func(rawEvt interface{}) {
		switch evt := rawEvt.(type) {
		case *events.Connected:
			onlineClientAdd(client)
		case *events.LoggedOut:
			// todo: 清理本地数据
			onlineClientRemove(client)
		case *events.Disconnected:
		case *events.ClientOutdated:
			onlineClientRemove(client)
		case *events.Message:
			var cm model.WhatsappChatMessage
			model.DB.Find(&cm, "msg_id = ? AND device_jid = ?", evt.Info.ID, client.DeviceJID())
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
					DeviceJID: client.Store.ID.String(),
				})
			}
			//case *events.HistorySync:
			//	if *evt.Data.SyncType == proto.HistorySync_INITIAL_BOOTSTRAP || *evt.Data.SyncType == proto.HistorySync_RECENT {
			//		for _, c := range evt.Data.Conversations {
			//			var chat model.WhatsappChat
			//			model.DB.Find(&chat, "jid = ? AND device_jid = ?", c.Id, client.DeviceJID())
			//			if chat.ID > 0 {
			//				model.DB.Model(chat).Updates(&model.WhatsappChat{
			//					JID:          c.GetId(),
			//					Name:         c.GetName(),
			//					LastSendTime: cast.ToTime(c.ConversationTimestamp),
			//					ReadOnly:     false,
			//					UnreadCount:  cast.ToUint(c.UnreadCount),
			//					DeviceJID:    client.DeviceJID(),
			//				})
			//			} else {
			//				model.DB.Save(&model.WhatsappChat{
			//					JID:          c.GetId(),
			//					Name:         c.GetName(),
			//					LastSendTime: cast.ToTime(c.ConversationTimestamp),
			//					ReadOnly:     false,
			//					UnreadCount:  cast.ToUint(c.UnreadCount),
			//					DeviceJID:    client.DeviceJID(),
			//				})
			//			}
			//			for _, m := range c.Messages {
			//				if m.Message.Message == nil {
			//					continue
			//				}
			//
			//				var cm model.WhatsappChatMessage
			//				model.DB.Find(&cm, "msg_id = ? AND device_jid = ?", m.Message.Key.GetId(), client.DeviceJID())
			//				if cm.ID == 0 {
			//					msg := &model.WhatsappChatMessage{
			//						MsgID:     m.Message.Key.GetId(),
			//						ChatJID:   m.Message.Key.GetRemoteJid(),
			//						SenderJID: m.Message.GetParticipant(),
			//						Msg: &model.Msg{
			//							Message: m.Message.GetMessage(),
			//						},
			//						FromMe:    m.Message.Key.GetFromMe(),
			//						SendTime:  cast.ToTime(m.Message.GetMessageTimestamp()),
			//						DeviceJID: client.DeviceJID(),
			//					}
			//					if msg.FromMe {
			//						msg.SenderJID = client.Store.ID.ToNonAD().String()
			//					}
			//					model.DB.Save(msg)
			//				}
			//			}
			//		}
			//	}
		}
	})

	return client, nil
}

func Login(ctx context.Context, phone string) (*Client, <-chan whatsmeow.QRChannelItem, error) {
	c, err := newClient(phone)
	if err != nil {
		return nil, nil, err
	}

	if c.Store.ID != nil {
		err = c.autoLogin()
		if err == nil {
			return c, nil, nil
		}
	}

	var ch <-chan whatsmeow.QRChannelItem
	if ch, err = c.GetQRChannel(ctx); err != nil {
		return nil, nil, err
	}

	if err = c.Connect(); err != nil {
		return nil, nil, err
	}
	return c, ch, nil
}

func (c *Client) autoLogin() error {
	ch := make(chan bool)
	defer close(ch)

	handlerID := c.AddEventHandler(func(evt interface{}) {
		switch evt.(type) {
		case *events.Connected:
			ch <- true
		case *events.ClientOutdated:
			ch <- false
		}
	})
	defer c.RemoveEventHandler(handlerID)

	err := c.Connect()
	if err != nil {
		return err
	}

	if <-ch {
		return nil
	}

	return errors.New("client outdated")
}

func (c *Client) Logout() error {
	err := c.Client.Logout()
	if err != nil {
		return err
	}
	onlineClientRemove(c)
	return nil
}

func (c *Client) Phone() string {
	return c.Store.ID.User
}

func (c *Client) DeviceJID() string {
	return c.Store.ID.String()
}

var onlineClients []*Client

func GetOnlineClients() []*Client {
	return onlineClients
}

func GetClient(jid *types.JID) (*Client, error) {
	client, ok := lo.Find(onlineClients, func(item *Client) bool {
		return *item.Store.ID == *jid
	})
	if !ok {
		return nil, errors.New("客户端未登录")
	}
	return client, nil
}

func onlineClientAdd(client *Client) {
	onlineClients = append(onlineClients, client)
}

func onlineClientRemove(client *Client) {
	onlineClients = lo.DropWhile(onlineClients, func(item *Client) bool {
		return item.Store.ID == client.Store.ID
	})
}

func GetDevices() ([]*store.Device, error) {
	return container.GetAllDevices()
}

func DeleteDevice(jid *types.JID) error {
	devices, err := container.GetAllDevices()
	if err != nil {
		return err
	}

	device, b := lo.Find(devices, func(item *store.Device) bool {
		return *item.ID == *jid
	})

	if !b {
		return nil
	}

	return container.DeleteDevice(device)
}
