package cmd

import (
	"fmt"
	"github.com/fhluo/winproxy"
	"github.com/fhluo/winproxy/cmd/i18n"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

var (
	settings winproxy.Settings
)

var rootCmd = &cobra.Command{
	Use: "winproxy",
	Run: func(cmd *cobra.Command, args []string) {
		if !flagsChanged(cmd, "use-proxy", "use-script", "auto-detect", "proxy-address", "bypass-list", "script-address") {
			fmt.Println(Render(settings))
			return
		}

		if err := settings.Apply(); err != nil {
			slog.Error("failed to apply settings", "err", err)
		}

		fmt.Println(Render(settings))
	},
}

//go:generate go run github.com/fhluo/i18n/tools/gotext@latest -l en-US,zh-Hans -d ./i18n/locales/ -p ./i18n

func init() {
	var err error
	settings, err = winproxy.ReadSettings()
	if err != nil {
		slog.Error("failed to read settings", "err", err)
		os.Exit(1)
	}

	p := i18n.GetPrinter()

	rootCmd.InitDefaultHelpCmd()
	helpCmd, ok := lo.Find(rootCmd.Commands(), func(cmd *cobra.Command) bool {
		return cmd.Name() == "help"
	})
	if ok {
		helpCmd.Short = p.Sprintf("Help about any command")
		helpCmd.Long = p.Sprintf(`Help provides help for any command in the application.
Simply type %s help [path to command] for full details.`, rootCmd.Name())
	}

	rootCmd.InitDefaultCompletionCmd()
	completionCmd, ok := lo.Find(rootCmd.Commands(), func(cmd *cobra.Command) bool {
		return cmd.Name() == "completion"
	})
	if ok {
		completionCmd.Short = p.Sprintf("Generate the autocompletion script for the specified shell")
		completionCmd.Long = p.Sprintf(`Generate the autocompletion script for %[1]s for the specified shell.
See each sub-command's help for details on how to use the generated script.
`, rootCmd.Name())
	}

	rootCmd.InitDefaultHelpFlag()
	helpFlag := rootCmd.Flags().Lookup("help")
	if helpFlag != nil {
		commandName := p.Sprintf("this command")
		if rootCmd.Name() != "" {
			commandName = rootCmd.Name()
		}

		helpFlag.Usage = p.Sprintf("help for %s", commandName)
	}

	rootCmd.SetUsageTemplate(fmt.Sprintf(`%s:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

%s:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

%s:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

%s:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

%s{{end}}
`,
		p.Sprintf("Usage"),
		p.Sprintf("Aliases"),
		p.Sprintf("Examples"),
		p.Sprintf("Available Commands"),
		p.Sprintf("Additional Commands"),
		p.Sprintf("Flags"),
		p.Sprintf("Global Flags"),
		p.Sprintf("Additional help topics"),
		p.Sprintf(`Use "%s" for more information about a command.`, "{{.CommandPath}} [command] --help"),
	))

	rootCmd.Flags().SortFlags = false
	rootCmd.Flags().BoolVarP(&settings.Proxy, "use-proxy", "p", settings.Proxy, p.Sprintf("use a proxy server"))
	rootCmd.Flags().StringVar(&settings.ProxyAddress, "proxy-address", settings.ProxyAddress, p.Sprintf("proxy address"))

	rootCmd.Flags().BoolVarP(&settings.Script, "use-script", "s", settings.Script, p.Sprintf("use setup script"))
	rootCmd.Flags().StringVar(&settings.ScriptAddress, "script-address", settings.ScriptAddress, p.Sprintf("script address"))

	rootCmd.Flags().BoolVarP(&settings.AutoDetect, "auto-detect", "a", settings.AutoDetect, p.Sprintf("automatically detect settings"))
	rootCmd.Flags().StringSliceVar(&settings.BypassList, "bypass-list", settings.BypassList, p.Sprintf("bypass list"))
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
		slog.Error(err.Error())
		os.Exit(1)
	}
}
