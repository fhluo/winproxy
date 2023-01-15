package cmd

import (
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "winproxy --use-proxy=true",
	Run: func(cmd *cobra.Command, args []string) {
		settings.SetUseProxy(true)
		if err := settings.Apply(); err != nil {
			slog.Error("failed to apply settings", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(onCmd)
}
