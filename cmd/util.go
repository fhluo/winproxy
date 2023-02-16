package cmd

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/fhluo/winproxy"
	"github.com/fhluo/winproxy/cmd/i18n"
	"golang.org/x/exp/slog"
	"os"
	"text/template"
)

func PrintSettings(s winproxy.Settings) {
	w := bufio.NewWriter(os.Stdout)
	p := i18n.GetPrinter()
	templateStr := fmt.Sprintf(`{{red "%s"}}: {{blue "%%v" .Proxy}}
{{red "%s"}}: {{yellow .ProxyAddress}}
{{red "%s"}}: {{blue "%%v" .Script}}
{{red "%s"}}: {{yellow .ScriptAddress}}
{{red "%s"}}: {{blue "%%v" .AutoDetect}}
{{red "%s"}}:
  {{range .BypassList -}}
    - {{yellow .}}
  {{end -}}
`,
		p.Sprintf("Use proxy"), p.Sprintf("Proxy address"),
		p.Sprintf("Use script"), p.Sprintf("Script address"),
		p.Sprintf("Auto-detect"), p.Sprintf("Bypass list"),
	)

	settingsTemplate := template.Must(template.New("").Funcs(
		template.FuncMap{
			"red":    color.RedString,
			"blue":   color.BlueString,
			"yellow": color.YellowString,
		},
	).Parse(templateStr))

	if err := settingsTemplate.Execute(w, s); err != nil {
		slog.Error("failed to execute settings template", err)
		os.Exit(1)
	}

	if err := w.Flush(); err != nil {
		slog.Error("failed to flush", err)
		os.Exit(1)
	}
}
