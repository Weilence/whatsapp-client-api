package cmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"whatsapp-client/pkg/whatsapp"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "发送Whatsapp消息",
	Args:  cobra.NoArgs,
	Run:   newRun(nil),
}

func (ctx *Context) send() {
	flagSet := flag.NewFlagSet("send", flag.ExitOnError)
	defineFlags(flagSet)
	err := flagSet.Parse(ctx.args)
	if err != nil {
		ctx.Write(err)
		return
	}
	from, err := flagSet.GetString("from")
	if err != nil {
		ctx.Write(err)
		return
	}
	to, err := flagSet.GetString("to")
	if err != nil {
		ctx.Write(err)
		return
	}
	msgType, err := flagSet.GetString("type")
	if err != nil {
		ctx.Write(err)
		return
	}
	text, err := flagSet.GetString("text")
	if err != nil {
		ctx.Write(err)
		return
	}
	switch msgType {
	case "text":
		client, err := whatsapp.GetClient(from)
		if err != nil {
			ctx.Write(err)
			return
		}
		client.SendTextMessage(to, text)
	}
	ctx.Write("消息发送成功")
}

func defineFlags(flagSet *flag.FlagSet) {
	flagSet.String("from", "", "发送人")
	flagSet.String("to", "", "接收人")
	flagSet.String("type", "", "消息类型，text、image、file中的一个")
	flagSet.String("text", "", "文本消息内容")
}

func init() {
	rootCmd.AddCommand(sendCmd)
	defineFlags(sendCmd.Flags())
}
