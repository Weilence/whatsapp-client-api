package cmd

import (
	"bufio"
	"errors"
	"fmt"
	flag "github.com/spf13/pflag"
	"gopkg.in/natefinch/npipe.v2"
	"io"
	"os"
	"strings"
	"whatsapp-client/pkg/utils"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "whatsapp-client-cli",
	Short: "whatsapp客户端命令行",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var NamedPipeAddress = `\\.\pipe\whatsapp-client`

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.whatsapp-client.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func newRun(callback func(str string) (notStop bool)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		conn, err := npipe.Dial(NamedPipeAddress)

		writeArgs := []string{cmd.Name()}
		writeArgs = append(writeArgs, cmd.Flags().Args()...)
		cmd.Flags().Visit(func(flag *flag.Flag) {
			writeArgs = append(writeArgs, "--"+flag.Name+" "+flag.Value.String())
		})

		_, err = fmt.Fprintln(conn, strings.Join(writeArgs, " "))

		utils.NoError(err)

		for {
			str, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Println(err)
				}
				break
			}

			if callback == nil {
				fmt.Println(str)
				continue
			}

			if !callback(str) {
				break
			}
		}
	}
}
