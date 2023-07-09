package whatsapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/weilence/whatsapp-client/config"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Client struct {
	*whatsmeow.Client
}

func NewClient(jid types.JID) (*Client, error) {
	device, err := GetDevice(jid)
	if err != nil {
		return nil, fmt.Errorf("device get error: %w", err)
	}

	if device == nil {
		device = container.NewDevice()
	}

	var logger waLog.Logger
	if *config.Env == "dev" {
		logger = waLog.Stdout("Client_"+jid.String(), "DEBUG", true)
	} else {
		logger = waLog.Stdout("Client_"+jid.String(), "INFO", true)
	}

	client := &Client{Client: whatsmeow.NewClient(device, logger)}
	addClient(client)
	return client, nil
}

func (c *Client) Login(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error) {
	if c.Store.ID != nil {
		if err := c.autoLogin(); err == nil {
			return nil, nil
		}
	}

	ch, err := c.GetQRChannel(ctx)
	if err != nil {
		return nil, err
	}

	if err = c.Connect(); err != nil {
		return nil, err
	}
	return ch, nil
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

	removeClient(c)
	return nil
}

func (c *Client) DeviceJID() string {
	return c.Store.ID.String()
}

var clients []*Client = make([]*Client, 0)

func addClient(c *Client) {
	for i, client := range clients {
		if client.Store.ID != nil && c.Store.ID != nil && *client.Store.ID == *c.Store.ID {
			clients[i] = c
			return
		}
	}

	clients = append(clients, c)
}

func removeClient(c *Client) {
	for i, client := range clients {
		if client.Store.ID != nil && c.Store.ID != nil && *client.Store.ID == *c.Store.ID {
			clients = append(clients[:i], clients[i+1:]...)
			return
		}
	}
}

func GetClients() []*Client {
	return clients
}

func GetClient(jid types.JID) (*Client, error) {
	client, ok := lo.Find(clients, func(item *Client) bool {
		return item.Store.ID != nil && *item.Store.ID == jid
	})
	if !ok {
		return nil, errors.New("客户端未登录")
	}
	return client, nil
}
