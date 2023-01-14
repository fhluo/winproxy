package cmd

import (
	_ "embed"
	"github.com/fhluo/winproxy"
	"golang.org/x/exp/slog"
	"strings"
	"text/template"
)

var (
	//go:embed data/settings.tmpl
	settingsTemplateString string

	settingsTemplate = template.Must(template.New("").Parse(settingsTemplateString))
)

func formatSettings(s *winproxy.Settings) string {
	b := new(strings.Builder)
	if err := settingsTemplate.Execute(b, s); err != nil {
		slog.Error("failed to execute settings template", err)
	}
	return b.String()
}
