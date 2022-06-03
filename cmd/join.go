package cmd

import (
	"github.com/spf13/cobra"
	"whatsapp-client/pkg/whatsapp"
)

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join group",
	Args:  cobra.ExactArgs(2),
	Run:   newRun(nil),
}

func (ctx *Context) join() {
	phone, link := ctx.args[1], ctx.args[2]
	client, err := whatsapp.GetClient(phone)
	defer func() {
		if err != nil {
			ctx.Write(err)
		}
	}()

	if err != nil {
		return
	}
	_, err = client.JoinGroupWithLink(link)
	if err == nil {
		ctx.Write("加入群组成功")
	}
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
