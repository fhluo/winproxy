package cmd

import (
	"github.com/fhluo/winproxy"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "winproxy --use-proxy=false",
	Run: func(cmd *cobra.Command, args []string) {
		settings.SetUseProxy(false)
		if err := winproxy.WriteSettings(settings); err != nil {
			slog.Error("failed to write settings", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(offCmd)
}
