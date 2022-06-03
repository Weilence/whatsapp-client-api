package whatsapp

import "go.mau.fi/whatsmeow/types"

func (c *Client) GetJoinedGroups() []*types.GroupInfo {
	groups, _ := c.Client.GetJoinedGroups()
	c.groups = groups
	return c.groups
}

func (c *Client) LeaveGroup(id string) error {
	err := c.Client.LeaveGroup(NewGroupJID(id))
	return err
}
