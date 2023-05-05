package cmd

import (
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "winproxy --use-proxy=false",
	Run: func(cmd *cobra.Command, args []string) {
		settings.Proxy = false
		if err := settings.Apply(); err != nil {
			slog.Error("failed to apply settings", "err", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(offCmd)
}
