package cmd

import (
	"github.com/fhluo/winproxy"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "winproxy --use-proxy=true",
	Run: func(cmd *cobra.Command, args []string) {
		settings.SetUseProxy(true)
		if err := winproxy.WriteSettings(settings); err != nil {
			slog.Error("failed to write settings", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(onCmd)
}
