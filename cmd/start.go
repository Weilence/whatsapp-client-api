package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/natefinch/npipe.v2"
	"log"
	"net"
	"strings"
	"whatsapp-client/config"
	"whatsapp-client/internal/model"
	"whatsapp-client/pkg/utils"
	"whatsapp-client/pkg/whatsapp"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动Whatsapp服务",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		config.Init()
		model.Init()
		whatsapp.Init(model.SqlDB())

		ln, err := npipe.Listen(NamedPipeAddress)
		if err != nil {
			return
		}
		ctx := &Context{ln: ln}

		for {
			conn, err := ln.Accept()
			if err != nil {
				return err
			}
			ctx.conn = conn

			str, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				return err
			}

			ctx.args = strings.Fields(str)
			go handleConnection(ctx)
		}
	},
}

type Context struct {
	ln   *npipe.PipeListener
	conn net.Conn
	args []string
}

func (ctx *Context) Write(a ...any) {
	_, _ = fmt.Fprintln(ctx.conn, a...)
}

func handleConnection(ctx *Context) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(ctx.conn)

	switch ctx.args[0] {
	case "stop":
		err := ctx.ln.Close()
		utils.NoError(err)
	case "login":
		ctx.login()
	case "logout":
		ctx.logout()
	case "send":
		ctx.send()
	default:
		log.Println("未执行命令")
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}
