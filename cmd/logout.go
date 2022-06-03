package cmd

import (
	"whatsapp-client/pkg/whatsapp"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout by phone",
	Args:  cobra.ExactArgs(1),
	Run:   newRun(nil),
}

func (ctx *Context) logout() {
	client, err := whatsapp.GetClient(ctx.args[1])
	if err != nil {
		ctx.Write(err)
		return
	}
	err = client.Logout()
	if err != nil {
		ctx.Write(err)
		return
	}
	ctx.Write("登出成功")
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
