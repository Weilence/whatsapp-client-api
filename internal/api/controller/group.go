package controller

import (
	"log"

	"github.com/weilence/whatsapp-client/internal/api"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"
	"go.mau.fi/whatsmeow/types"
)

type Group struct {
	JID          string   `json:"jid"`
	Name         string   `json:"name"`
	CreateTime   int64    `json:"createTime"`
	Participants []string `json:"participants"`
}

type groupQueryReq struct {
	JID *types.JID `query:"jid" validate:"required"`
}

func GroupQuery(c *api.HttpContext, req *groupQueryReq) (interface{}, error) {
	client, err := whatsapp.GetClient(req.JID)
	if err != nil {
		return nil, err
	}
	groups := client.GetGroups()

	data := make([]Group, len(groups))
	for i, group := range groups {

		participants := make([]string, len(group.Participants))
		for j, participant := range group.Participants {
			participants[j] = participant.JID.String()
		}

		data[i] = Group{
			JID:          group.JID.String(),
			Name:         group.Name,
			CreateTime:   group.GroupCreated.UnixMilli(),
			Participants: participants,
		}
	}
	return data, nil
}

type groupGetReq struct {
	JID *types.JID `uri:"jid"`
}

func GroupGet(c *api.HttpContext, req *groupGetReq) (interface{}, error) {
	client, err := whatsapp.GetClient(req.JID)
	if err != nil {
		return nil, err
	}
	info, err := client.GetGroupInfo(req.JID)
	if err != nil {
		return nil, err
	}
	return info, nil
}

type groupJoinReq struct {
	JID  *types.JID `json:"jid"`
	Path string     `json:"path"`
}

func GroupJoin(c *api.HttpContext, req *groupJoinReq) (interface{}, error) {
	client, err := whatsapp.GetClient(req.JID)
	if err != nil {
		return nil, err
	}

	groupJID, err := client.JoinGroupWithLink(req.Path)
	log.Println("joined " + groupJID.String())
	if err != nil {
		return nil, err
	}
	return nil, nil
}
