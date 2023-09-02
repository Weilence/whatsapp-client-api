package controller

import (
	"fmt"
	"time"

	"log/slog"

	"github.com/weilence/whatsapp-client/internal/utils"
	"go.mau.fi/whatsmeow/types"
)

type groupListReq struct {
	JID types.JID `query:"jid"`
}
type groupListRes struct {
	Phone        types.JID   `json:"jid"`
	Name         string      `json:"name"`
	CreateTime   time.Time   `json:"createTime"`
	Participants []types.JID `json:"participants"`
}

func GroupList(c *utils.HttpContext, req *groupListReq) (interface{}, error) {
	client, err := utils.GetClient(req.JID)
	if err != nil {
		return nil, err
	}
	groups, err := client.GetJoinedGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get joined groups: %w", err)
	}

	data := make([]groupListRes, len(groups))
	for i, group := range groups {
		participants := make([]types.JID, len(group.Participants))
		for j, participant := range group.Participants {
			participants[j] = participant.JID
		}

		data[i] = groupListRes{
			Phone:        group.JID,
			Name:         group.Name,
			CreateTime:   group.GroupCreated,
			Participants: participants,
		}
	}
	return data, nil
}

type groupGetReq struct {
	JID      types.JID `query:"jid"`
	GroupJID types.JID `query:"groupJID"`
}

func GroupGet(c *utils.HttpContext, req *groupGetReq) (interface{}, error) {
	client, err := utils.GetClient(req.JID)
	if err != nil {
		return nil, err
	}

	info, err := client.GetGroupInfo(req.GroupJID)
	if err != nil {
		return nil, err
	}

	return info, nil
}

type groupJoinReq struct {
	JID  types.JID `query:"jid"`
	Code string    `query:"code"`
}

func GroupJoin(c *utils.HttpContext, req *groupJoinReq) (interface{}, error) {
	client, err := utils.GetClient(req.JID)
	if err != nil {
		return nil, err
	}

	groupJID, err := client.JoinGroupWithLink(req.Code)
	if err != nil {
		slog.Error("faild to join group", "err", err)
		return nil, err
	}
	slog.Info("joined group", "jid", groupJID)
	return nil, nil
}
