package i18n

import (
	"embed"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/sys/windows"
	"golang.org/x/text/language"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "i18n: ", 0)

	//go:embed locales
	locales embed.FS

	localize func(lc *i18n.LocalizeConfig) string
)

func init() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	_, err := bundle.LoadMessageFileFS(locales, "locales/en.toml")
	if err != nil {
		logger.Fatalln(err)
	}

	_, err = bundle.LoadMessageFileFS(locales, "locales/zh.toml")
	if err != nil {
		logger.Fatalln(err)
	}

	languages, err := windows.GetUserPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if err != nil {
		languages, err = windows.GetSystemPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
		if err != nil {
			logger.Fatalln(err)
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
