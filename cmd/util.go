package cmd

import (
	_ "embed"
	"github.com/fatih/color"
	"github.com/fhluo/winproxy"
	"golang.org/x/exp/slog"
	"os"
	"text/template"
)

var (
	//go:embed data/settings.tmpl
	settingsTemplateString string

	settingsTemplate = template.Must(template.New("").Funcs(template.FuncMap{
		"red":    color.RedString,
		"blue":   color.BlueString,
		"yellow": color.YellowString,
	}).Parse(settingsTemplateString))
)

func PrintSettings(s winproxy.Settings) {
	if err := settingsTemplate.Execute(os.Stdout, s); err != nil {
		slog.Error("failed to execute settings template", err)
		os.Exit(1)
	}
}
