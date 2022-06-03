package cmd

import (
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止Whatsapp服务",
	Args:  cobra.NoArgs,
	Run:   newRun(nil),
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
