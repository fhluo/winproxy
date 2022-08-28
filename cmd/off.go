package cmd

import (
	"github.com/fhluo/winproxy"
	"github.com/spf13/cobra"
	"log"
)

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "winproxy --use-proxy=false",
	Run: func(cmd *cobra.Command, args []string) {
		settings.SetUseProxy(false)
		if err := winproxy.WriteSettings(settings); err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(offCmd)
}
