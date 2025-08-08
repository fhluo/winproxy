package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/fatih/color"
	winproxy "github.com/fhluo/winproxy/go"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/cobra"
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

	cobra.AddTemplateFuncs(template.FuncMap{
		"FgHiWhite": color.New(color.FgHiWhite).SprintFunc(),
	})

	localizeHelpCommand()
	localizeHelpFlag()
	localizeCompletionCommand()
	localizeUsageTemplate()
	initFlags()
}

func localizeHelpCommand() {
	rootCmd.InitDefaultHelpCmd()

	i := slices.IndexFunc(rootCmd.Commands(), func(cmd *cobra.Command) bool {
		return cmd.Name() == "help"
	})

	if i != -1 {
		helpCmd := rootCmd.Commands()[i]
		helpCmd.Short = Localize(&i18n.Message{ID: "help-short", Other: "Help about any command"})
		helpCmd.Long = LocalizeConfig(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{ID: "help-long", Other: `Help provides help for any command in the application.
Simply type {{.CommandName}} help [path to command] for full details.`},
			TemplateData: map[string]any{
				"CommandName": rootCmd.Name(),
			},
		})
	}
}

func localizeCompletionCommand() {
	rootCmd.InitDefaultCompletionCmd()

	i := slices.IndexFunc(rootCmd.Commands(), func(cmd *cobra.Command) bool {
		return cmd.Name() == "completion"
	})

	if i != -1 {
		completionCmd := rootCmd.Commands()[i]
		completionCmd.Short = Localize(&i18n.Message{ID: "completion-short", Other: "Generate the autocompletion script for the specified shell"})
		completionCmd.Long = LocalizeConfig(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "completion-long",
				Other: `Generate the autocompletion script for {{.CommandName}} for the specified shell.
See each sub-command's help for details on how to use the generated script.
`,
			},
			TemplateData: map[string]any{
				"CommandName": rootCmd.Name(),
			},
		})
	}
}

func localizeHelpFlag() {
	rootCmd.InitDefaultHelpFlag()
	helpFlag := rootCmd.Flags().Lookup("help")
	if helpFlag != nil {
		commandName := Localize(&i18n.Message{ID: "this-command", Other: "this command"})
		if rootCmd.Name() != "" {
			commandName = rootCmd.Name()
		}

		helpFlag.Usage = LocalizeConfig(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{ID: "help-for", Other: "help for {{.CommandName}}"},
			TemplateData: map[string]any{
				"CommandName": commandName,
			},
		})
	}
}

func localizeUsageTemplate() {
	usageTemplate := rootCmd.UsageTemplate()

	usageTemplate = strings.NewReplacer([]string{
		"Usage:", fmt.Sprintf(
			`{{FgHiWhite "%s"}}`,
			Localize(&i18n.Message{ID: "Usage:", Other: "Usage:"}),
		),
		"Aliases:", fmt.Sprintf(
			`{{FgHiWhite "%s"}}`,
			Localize(&i18n.Message{ID: "Aliases:", Other: "Aliases:"}),
		),
		"Examples:", fmt.Sprintf(
			`{{FgHiWhite "%s"}}`,
			Localize(&i18n.Message{ID: "Examples:", Other: "Examples:"}),
		),
		"Available Commands:", fmt.Sprintf(
			`{{FgHiWhite "%s"}}`,
			Localize(&i18n.Message{ID: "Available Commands:", Other: "Available Commands:"}),
		),
		"Additional Commands:", fmt.Sprintf(
			`{{FgHiWhite "%s"}}`,
			Localize(&i18n.Message{ID: "Additional Commands:", Other: "Additional Commands:"}),
		),
		"Flags:", fmt.Sprintf(
			`{{FgHiWhite "%s"}}`,
			Localize(&i18n.Message{ID: "Flags:", Other: "Flags:"}),
		),
		"Global Flags:", fmt.Sprintf(
			`{{FgHiWhite "%s"}}`,
			Localize(&i18n.Message{ID: "Global Flags:", Other: "Global Flags:"}),
		),
		"Additional help topics:", fmt.Sprintf(
			`{{FgHiWhite "%s"}}`,
			Localize(&i18n.Message{ID: "Additional help topics:", Other: "Additional help topics:"}),
		),
		`Use "{{.CommandPath}} [command] --help" for more information about a command.`,
		LocalizeConfig(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Help",
				Other: `Use {{.Help}} for more information about a command.`,
			},
			TemplateData: map[string]any{
				"Help": `"{{.CommandPath}} [command] --help"`,
			},
		}),
	}...).Replace(usageTemplate)

	rootCmd.SetUsageTemplate(usageTemplate)
}

func initFlags() {
	rootCmd.Flags().SortFlags = false
	rootCmd.Flags().BoolVarP(&settings.Proxy, "use-proxy", "p", settings.Proxy,
		Localize(&i18n.Message{
			ID:    "use-proxy",
			Other: "use a proxy server",
		}),
	)
	rootCmd.Flags().StringVar(&settings.ProxyAddress, "proxy-address", settings.ProxyAddress,
		Localize(&i18n.Message{
			ID:    "proxy-address",
			Other: "proxy address",
		}),
	)
	rootCmd.Flags().BoolVarP(&settings.Script, "use-script", "s", settings.Script,
		Localize(&i18n.Message{
			ID:    "use-script",
			Other: "use setup script",
		}),
	)
	rootCmd.Flags().StringVar(&settings.ScriptAddress, "script-address", settings.ScriptAddress,
		Localize(&i18n.Message{
			ID:    "script-address",
			Other: "script address",
		}),
	)
	rootCmd.Flags().BoolVarP(&settings.AutoDetect, "auto-detect", "a", settings.AutoDetect,
		Localize(&i18n.Message{
			ID:    "auto-detect",
			Other: "automatically detect settings",
		}),
	)
	rootCmd.Flags().StringSliceVar(&settings.BypassList, "bypass-list", settings.BypassList,
		Localize(&i18n.Message{
			ID:    "bypass-list",
			Other: "bypass list",
		}),
	)
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
