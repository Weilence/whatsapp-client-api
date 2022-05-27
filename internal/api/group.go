package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow/types"
	"whatsapp-client/pkg/whatsapp"
)

type Group struct {
	JID          string   `json:"jid"`
	Name         string   `json:"name"`
	CreateTime   int64    `json:"createTime"`
	Participants []string `json:"participants"`
}

func GroupQuery(c *gin.Context) {
	groups := whatsapp.GetClient(c.Query("jid")).GetJoinedGroups()

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
	c.JSON(0, data)
}

func GroupGet(c *gin.Context) {
	jid, err := types.ParseJID(c.Query("gjid"))
	if err != nil {
		panic(err)
	}
	info, err := whatsapp.GetClient(c.Query("jid")).GetGroupInfo(jid)
	if err != nil {
		panic(err)
	}
	c.JSON(0, info)
}

func GroupJoin(c *gin.Context) {
	c.Error(errors.New("测试"))

	client := whatsapp.GetClient(c.Query("jid"))
	link := c.Query("link")

	_, err := client.JoinGroupWithLink(link)

	if err != nil {

	}
}
