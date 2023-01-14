package i18n

import (
	"embed"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/exp/slog"
	"golang.org/x/sys/windows"
	"golang.org/x/text/language"
	"os"
)

var (
	//go:embed locales
	locales embed.FS

	localize func(lc *i18n.LocalizeConfig) string
)

func init() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	_, err := bundle.LoadMessageFileFS(locales, "locales/en.toml")
	if err != nil {
		slog.Error("failed to load message file", err, "path", "locales/en.toml")
		os.Exit(1)
	}

	_, err = bundle.LoadMessageFileFS(locales, "locales/zh.toml")
	if err != nil {
		slog.Error("failed to load message file", err, "path", "locales/zh.toml")
		os.Exit(1)
	}

	languages, err := windows.GetUserPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if err != nil {
		languages, err = windows.GetSystemPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
		if err != nil {
			slog.Error("failed to get system preferred UI languages", err)
			os.Exit(1)
		}
	}

	localize = i18n.NewLocalizer(bundle, languages...).MustLocalize
}

func Localize(id string, s string) string {
	return localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    id,
			Other: s,
		},
	})
}
