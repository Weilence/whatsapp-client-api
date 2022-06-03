package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"whatsapp-client/pkg/whatsapp"

	"github.com/spf13/cobra"
)

// groupsCmd represents the groups command
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List groups",
	Args:  cobra.ExactArgs(1),
	Run: newRun(func(str string) (stop bool) {
		fmt.Print(str)
		return
	}),
}

func (ctx *Context) groups() {
	phone := ctx.args[1]
	var err error
	defer func() {
		if err != nil {
			ctx.Write(err)
		}
	}()
	client, err := whatsapp.GetClient(phone)
	if err != nil {
		return
	}
	groups := client.GetJoinedGroups()

	t := table.NewWriter()
	t.SetOutputMirror(ctx.conn)
	t.AppendHeader(table.Row{"#", "JID", "Name", "Topic"})
	for i, group := range groups {
		t.AppendRow([]interface{}{i + 1, group.JID, group.Name, group.Topic})
	}
	t.Render()
}

func init() {
	rootCmd.AddCommand(groupsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// groupsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
