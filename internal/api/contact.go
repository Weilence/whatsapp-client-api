package api

import (
	"github.com/gin-gonic/gin"
	"sort"
	"strings"
	"whatsapp-client/pkg/whatsapp"
)

type Contact struct {
	Jid          string `json:"jid"`
	Found        bool   `json:"found"`
	Name         string `json:"name"`
	BusinessName string `json:"businessName"`
}

func ContactQuery(c *gin.Context) {
	client, err := whatsapp.GetClient(c.Query("jid"))
	contacts, err := client.Store.Contacts.GetAllContacts()

	if err != nil {
		panic(err)
	}

	var data []Contact
	for jid, item := range contacts {
		contact := Contact{
			Jid:          jid.String(),
			Found:        false,
			BusinessName: item.BusinessName,
		}
		if item.FullName != "" {
			contact.Name = item.FullName
		} else {
			contact.Name = item.PushName
		}
		data = append(data, contact)

	}
	sort.Slice(data, func(i, j int) bool {
		return strings.Compare(data[i].Name, data[j].Name) < 0
	})

	c.JSON(0, data)
}

type (
	VerifyReq struct {
		JID    string   `json:"jid,omitempty"`
		Phones []string `json:"phones,omitempty"`
	}

	VerifyRes struct {
		JID  string `json:"jid,omitempty"`
		IsIn bool   `json:"isIn,omitempty"`
	}
)

func ContactVerify(c *gin.Context) {
	var data VerifyReq
	c.Bind(&data)
	for i := range data.Phones {
		if data.Phones[0] != "+" {
			data.Phones[i] = "+" + data.Phones[i]
		}
	}
	client, err := whatsapp.GetClient(data.JID)
	res, err := client.IsOnWhatsApp(data.Phones)
	if err != nil {
		panic(err)
	}

	var results []VerifyRes
	for _, item := range res {
		results = append(results, VerifyRes{JID: item.JID.String(), IsIn: item.IsIn})
	}
	c.JSON(0, results)
}
