package whatsapp

import "go.mau.fi/whatsmeow/types"

func (c *Client) GetJoinedGroups() []*types.GroupInfo {
	if len(c.groups) == 0 {
		groups, _ := c.Client.GetJoinedGroups()
		c.groups = groups
	}
	return c.groups
}
