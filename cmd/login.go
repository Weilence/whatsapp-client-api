/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"go.mau.fi/whatsmeow"
	"gopkg.in/natefinch/npipe.v2"
	"io"
	"log"
	"strings"
	"whatsapp-client/pkg/utils"
	"whatsapp-client/pkg/whatsapp"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := npipe.Dial(`\\.\pipe\whatsapp-client`)
		if err != nil {
			log.Println(err)
		}
		_, err = fmt.Fprintf(conn, "login %v\n", strings.Join(args, " "))

		for {
			line, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if !errors.Is(err, io.EOF) {
					log.Println("登录失败", err)
				}
				break
			}

			t, msg, found := strings.Cut(line, ",")
			if !found {
				log.Println("获取qrcode时连接中断")
				break
			}

			if t == "code" {
				code, err := qrcode.New(msg, qrcode.Low)
				utils.NoError(err)
				fmt.Println(code.ToSmallString(false))
			} else {
				fmt.Println(msg)
			}
		}
	},
}

func (context *Context) login(args []string) {
	var client *whatsapp.Client
	if len(args) > 1 {
		client = whatsapp.NewClient(args[1])
	} else {
		client = whatsapp.NewClient("")
	}
	qrChan := client.Login()

	if qrChan == nil {
		context.Write("success,直接登录成功")
		return
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			context.Write("code," + evt.Code)
		} else if evt == whatsmeow.QRChannelSuccess {
			context.Write("success,扫码登录成功")
			break
		} else if evt == whatsmeow.QRChannelScannedWithoutMultidevice {
			context.Write("error,请开启多设备测试版")
			break
		} else {
			log.Println(evt)
			context.Write("error,扫码登录失败")
			break
		}
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
