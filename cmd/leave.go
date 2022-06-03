package cmd

import (
	"github.com/spf13/cobra"
	"whatsapp-client/pkg/whatsapp"
)

// leaveCmd represents the leave command
var leaveCmd = &cobra.Command{
	Use:   "leave",
	Short: "leave group",
	Args:  cobra.ExactArgs(2),
	Run:   newRun(nil),
}

func (ctx *Context) leave() {
	phone, groupJID := ctx.args[1], ctx.args[2]
	var err error
	defer func() {
		if err != nil {
			ctx.Write(err)
		}
	}()
	client, err := whatsapp.GetClient(phone)
	err = client.LeaveGroup(groupJID)
	if err == nil {
		ctx.Write("退出群组成功")
	}
}

func init() {
	rootCmd.AddCommand(leaveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// leaveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// leaveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
