/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/natefinch/npipe.v2"
	"net"
	"strings"
	"whatsapp-client/config"
	"whatsapp-client/internal/model"
	"whatsapp-client/pkg/utils"
	"whatsapp-client/pkg/whatsapp"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		model.Init()
		whatsapp.Init(model.SqlDB())

		var err error
		ln, err = npipe.Listen(`\\.\pipe\whatsapp-client`)
		if err != nil {
			fmt.Printf("启动失败，%v \n", err)
			return
		}
		for {
			conn, err := ln.Accept()
			if err != nil {
				break
			}
			str, err := bufio.NewReader(conn).ReadString('\n')
			utils.NoError(err)

			args := strings.Fields(str)

			go handleConnection(args, &Context{
				conn: conn,
				args: args,
			})
		}
	},
}

var ln *npipe.PipeListener

type Context struct {
	conn net.Conn
	args []string
}

func (context *Context) Write(a ...any) {
	_, _ = fmt.Fprintln(context.conn, a...)
}

func handleConnection(args []string, context *Context) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(context.conn)

	switch args[0] {
	case "stop":
		err := ln.Close()
		utils.NoError(err)
	case "login":
		context.login(args)
	case "logout":
		context.logout(args)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
