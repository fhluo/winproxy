package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "winproxy --use-proxy=true",
	Run: func(cmd *cobra.Command, args []string) {
		settings.Proxy = true
		if err := settings.Apply(); err != nil {
			slog.Error("failed to apply settings", "err", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(onCmd)
}
