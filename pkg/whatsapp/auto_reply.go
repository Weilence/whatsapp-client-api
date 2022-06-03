package whatsapp

import (
	"go.mau.fi/whatsmeow/types/events"
	"io/ioutil"
	"strings"
	"whatsapp-client/internal/model"
)

type AutoReply struct {
	handlerID uint32
	cache     map[string]model.WhatsappAutoReply
}

func (c *Client) EnableAutoReply() {
	c.RefreshAutoReplay()

	if c.autoReply.handlerID != 0 {
		return
	}
	c.autoReply.handlerID = c.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			msg := v.Message.ExtendedTextMessage
			if msg == nil {
				break
			}
			if v.Info.IsGroup && !contains(msg.ContextInfo.MentionedJid, c.Store.ID.ToNonAD().String()) {
				break
			}

			for key := range c.autoReply.cache {
				if !strings.Contains(*msg.Text, key) {
					continue
				}

				reply := c.autoReply.cache[key]
				if reply.Type == 1 {
					if reply.File == "" {
						c.SendTextMessage(v.Info.Chat.String(), reply.Text)
					} else {
						bytes, err := ioutil.ReadFile(reply.File)
						if err != nil {
							panic(err)
						}

						c.SendImageMessage(v.Info.Chat.String(), bytes, reply.Text)
					}
				} else {
					bytes, err := ioutil.ReadFile(reply.File)
					if err != nil {
						panic(err)
					}

					c.SendDocumentMessage(v.Info.Chat.String(), bytes, reply.Text)
				}
				break
			}
		}
	})
}

func contains(strs []string, str string) bool {
	if len(strs) == 0 {
		return false
	}
	for _, mJid := range strs {
		if mJid == str {
			return true
		}
	}
	return false
}

func (c *Client) RefreshAutoReplay() {
	var autoReplies []model.WhatsappAutoReply
	model.DB.Find(&autoReplies)
	replayCache := make(map[string]model.WhatsappAutoReply)
	for _, reply := range autoReplies {
		replayCache[reply.Key] = reply
	}

	c.autoReply.cache = replayCache
}
