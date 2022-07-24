package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"
	"go.mau.fi/whatsmeow/types"
	"log"
)

type Group struct {
	JID          string   `json:"jid"`
	Name         string   `json:"name"`
	CreateTime   int64    `json:"createTime"`
	Participants []string `json:"participants"`
}

type groupQueryReq struct {
	JID *types.JID `uri:"jid"`
}

func GroupQuery(c *gin.Context, req *groupQueryReq) (interface{}, error) {
	client, err := whatsapp.GetClient(req.JID)
	if err != nil {
		return nil, err
	}
	groups := client.GetGroups()

	var data = make([]Group, len(groups))
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

func GroupGet(c *gin.Context, req *groupGetReq) (interface{}, error) {
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

func GroupJoin(c *gin.Context, req *groupJoinReq) (interface{}, error) {
	client, err := whatsapp.GetClient(req.JID)
	if err != nil {
		return nil, err
	}

	groupJID, err := client.JoinGroupWithLink(req.Path)
	log.Print("joined " + groupJID.String())
	if err != nil {
		return nil, err
	}
	return nil, nil
}
