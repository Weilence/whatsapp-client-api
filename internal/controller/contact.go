package controller

import (
	"sort"
	"strings"

	"github.com/weilence/whatsapp-client/internal/utils"
	"go.mau.fi/whatsmeow/types"
)

type Contact struct {
	Jid          string `json:"jid"`
	Found        bool   `json:"found"`
	Name         string `json:"name"`
	BusinessName string `json:"businessName"`
}

type contactQueryReq struct {
	JID types.JID `query:"jid"`
}

func ContactQuery(c *utils.HttpContext, req *contactQueryReq) (interface{}, error) {
	client, err := utils.GetClient(req.JID)
	if err != nil {
		return nil, err
	}

	contacts, err := client.Store.Contacts.GetAllContacts()
	if err != nil {
		return nil, err
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

	return data, nil
}

type (
	VerifyReq struct {
		JID    types.JID `query:"jid"`
		Phones []string  `json:"phones,omitempty"`
	}

	VerifyRes struct {
		JID  types.JID `json:"jid"`
		IsIn bool      `json:"isIn"`
	}
)

func ContactVerify(c *utils.HttpContext, req *VerifyReq) ([]*VerifyRes, error) {
	for i := range req.Phones {
		if req.Phones[0] != "+" {
			req.Phones[i] = "+" + req.Phones[i]
		}
	}
	client, err := utils.GetClient(req.JID)
	if err != nil {
		return nil, err
	}
	r, err := client.IsOnWhatsApp(req.Phones)
	if err != nil {
		return nil, err
	}

	var res []*VerifyRes
	for _, item := range r {
		res = append(res, &VerifyRes{JID: item.JID, IsIn: item.IsIn})
	}
	return res, nil
}
