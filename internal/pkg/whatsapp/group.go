package whatsapp

import "go.mau.fi/whatsmeow/types"

func (c *Client) GetGroups() []*types.GroupInfo {
	groups, _ := c.Client.GetJoinedGroups()
	c.groups = groups
	return c.groups
}

func (c *Client) GetGroupInfo(jid *types.JID) (*types.GroupInfo, error) {
	return c.Client.GetGroupInfo(*jid)
}

func (c *Client) LeaveGroup(id string) error {
	err := c.Client.LeaveGroup(NewGroupJID(id))
	return err
}
