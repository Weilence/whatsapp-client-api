package cmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"os"
	"path/filepath"
	"whatsapp-client/pkg/whatsapp"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "发送Whatsapp消息",
	Args:  cobra.NoArgs,
	Run:   newRun(nil),
}

func (ctx *Context) send() (err error) {
	flagSet := flag.NewFlagSet("send", flag.ExitOnError)
	defineFlags(flagSet)
	err = flagSet.Parse(ctx.args)
	defer func() {
		if err != nil {
			ctx.Write(err)
		}
	}()
	if err != nil {
		return
	}
	from, err := flagSet.GetString("from")
	if err != nil {
		return
	}
	to, err := flagSet.GetString("to")
	if err != nil {
		return
	}
	msgType, err := flagSet.GetString("type")
	if err != nil {
		return
	}
	text, err := flagSet.GetString("text")
	if err != nil {
		return
	}
	image, err := flagSet.GetString("image")
	if err != nil {
		return
	}
	file, err := flagSet.GetString("file")
	if err != nil {
		return
	}
	client, err := whatsapp.GetClient(from)
	if err != nil {
		return
	}
	switch msgType {
	case "text":
		if err != nil {
			return
		}
		err = client.SendTextMessage(to, text)
		if err != nil {
			return
		}
	case "image":
		bytes, err := os.ReadFile(image)
		if err != nil {
			return err
		}

		caption := text
		if caption == "" {
			caption = filepath.Base(image)
		}
		err = client.SendImageMessage(to, bytes, caption)
		if err != nil {
			return err
		}
	case "file":
		bytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		caption := filepath.Base(file)
		err = client.SendDocumentMessage(to, bytes, caption)
		if err != nil {
			return err
		}
	default:
		ctx.Write("type类型错误")
		return
	}
	ctx.Write("消息发送成功")
	return
}

func defineFlags(flagSet *flag.FlagSet) {
	flagSet.String("from", "", "发送人")
	flagSet.String("to", "", "接收人")
	flagSet.String("type", "", "消息类型，text、image、file中的一个")
	flagSet.String("text", "", "文本消息内容")
	flagSet.String("image", "", "图片路径")
	flagSet.String("file", "", "文件路径")
}

func init() {
	rootCmd.AddCommand(sendCmd)
	defineFlags(sendCmd.Flags())
}
