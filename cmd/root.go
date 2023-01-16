package cmd

import (
	"fmt"
	"github.com/fhluo/winproxy"
	"github.com/fhluo/winproxy/i18n"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"os"
)

var (
	settings winproxy.Settings

	useProxy   bool
	useScript  bool
	autoDetect bool

	proxyAddress  string
	bypassList    []string
	scriptAddress string
)

var rootCmd = &cobra.Command{
	Use: "winproxy",
	Run: func(cmd *cobra.Command, args []string) {
		if !flagsChanged(cmd, "use-proxy", "use-script", "auto-detect", "proxy-address", "bypass-list", "script-address") {
			fmt.Println(formatSettings(settings))
			return
		}

		settings.SetUseProxy(useProxy)
		settings.SetUseScript(useScript)
		settings.SetAutoDetect(autoDetect)
		settings.SetProxyAddress(proxyAddress)
		settings.SetBypassList(bypassList)
		settings.SetScriptAddress(scriptAddress)

		if err := settings.Apply(); err != nil {
			slog.Error("failed to apply settings", err)
		}

		fmt.Println(formatSettings(settings))
	},
}

func init() {
	var err error
	settings, err = winproxy.ReadSettings()
	if err != nil {
		slog.Error("failed to read settings", err)
		os.Exit(1)
	}

	rootCmd.Flags().BoolVar(&useProxy, "use-proxy", settings.UseProxy(), i18n.Localize("use-proxy", "use a proxy server"))
	rootCmd.Flags().BoolVar(&useScript, "use-script", settings.UseScript(), i18n.Localize("use-script", "use setup script"))
	rootCmd.Flags().BoolVar(&autoDetect, "auto-detect", settings.AutoDetect(), i18n.Localize("auto-detect", "automatically detect settings"))

	rootCmd.Flags().StringVar(&proxyAddress, "proxy-address", settings.ProxyAddress(), i18n.Localize("proxy-address", "proxy address"))
	rootCmd.Flags().StringSliceVar(&bypassList, "bypass-list", settings.BypassList(), i18n.Localize("bypass-list", "bypass list"))
	rootCmd.Flags().StringVar(&scriptAddress, "script-address", settings.ScriptAddress(), i18n.Localize("script-address", "script address"))
}

func flagsChanged(cmd *cobra.Command, names ...string) bool {
	for _, name := range names {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	return false
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to execute root command", err)
		os.Exit(1)
	}
}
