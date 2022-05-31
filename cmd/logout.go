/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"gopkg.in/natefinch/npipe.v2"
	"strings"
	"whatsapp-client/pkg/utils"
	"whatsapp-client/pkg/whatsapp"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := npipe.Dial(`\\.\pipe\whatsapp-client`)
		utils.NoError(err)
		_, err = fmt.Fprintf(conn, "logout %v\n", strings.Join(args, " "))
		utils.NoError(err)
	},
}

func (context *Context) logout(args []string) {
	client, err := whatsapp.GetClient(args[1])
	if err != nil {
		context.Write(err)
		return
	}
	err = client.Logout()
	if err != nil {
		context.Write(err)
		return
	}
	context.Write("登出成功")
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
