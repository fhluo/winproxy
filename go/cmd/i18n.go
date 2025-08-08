package cmd

import (
	"embed"
	"log/slog"
	"os"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/sys/windows"
	"golang.org/x/text/language"
)

//go:embed locales
var localesFS embed.FS

var Bundle = sync.OnceValue(func() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	_, err := bundle.LoadMessageFileFS(localesFS, "locales/zh-Hans.toml")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return bundle
})

func Languages() []string {
	languages, err := windows.GetUserPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if err != nil {
		languages, err = windows.GetSystemPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
		if err != nil {
			slog.Error("failed to get system preferred UI languages", "err", err)
			os.Exit(1)
		}
	}

	return languages
}

var DefaultLocalizer = sync.OnceValue(func() *i18n.Localizer {
	return i18n.NewLocalizer(Bundle(), Languages()...)
})

func Localize(message *i18n.Message) string {
	r, err := DefaultLocalizer().LocalizeMessage(message)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return r
}

func LocalizeConfig(config *i18n.LocalizeConfig) string {
	r, err := DefaultLocalizer().Localize(config)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return r
}
