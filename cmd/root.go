package cmd

import (
	"fmt"
	"github.com/fhluo/winproxy"
	"github.com/fhluo/winproxy/i18n"
	"github.com/spf13/cobra"
	"log"
)

var (
	settings *winproxy.Settings

	direct     bool
	useProxy   bool
	useScript  bool
	autoDetect bool
)

var rootCmd = &cobra.Command{
	Use: "winproxy",
	Run: func(cmd *cobra.Command, args []string) {
		if !flagsChanged(cmd, "use-proxy", "use-script", "auto-detect", "proxy-address", "bypass-list", "script-address") {
			fmt.Println(settings)
			return
		}

		settings.SetUseProxy(useProxy)
		settings.SetUseScript(useScript)
		settings.SetAutoDetect(autoDetect)

		if err := winproxy.WriteSettings(settings); err != nil {
			log.Println(err)
		}

		fmt.Println(settings)
	},
}

func init() {
	log.SetFlags(0)

	var err error
	settings, err = winproxy.ReadSettings()
	if err != nil {
		log.Fatalln(err)
	}

	rootCmd.Flags().BoolVar(&useProxy, "use-proxy", settings.UseProxy(), i18n.Localize("use-proxy", "use a proxy server"))
	rootCmd.Flags().BoolVar(&useScript, "use-script", settings.UseScript(), i18n.Localize("use-script", "use setup script"))
	rootCmd.Flags().BoolVar(&autoDetect, "auto-detect", settings.AutoDetect(), i18n.Localize("auto-detect", "automatically detect settings"))

	rootCmd.Flags().StringVar(&settings.ProxyAddress, "proxy-address", settings.ProxyAddress, i18n.Localize("proxy-address", "proxy address"))
	rootCmd.Flags().StringVar(&settings.BypassList, "bypass-list", settings.BypassList, i18n.Localize("bypass-list", "bypass list"))
	rootCmd.Flags().StringVar(&settings.ScriptAddress, "script-address", settings.ScriptAddress, i18n.Localize("script-address", "script address"))
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
		log.Fatalln(err)
	}
}
