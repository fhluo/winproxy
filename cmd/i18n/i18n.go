package i18n

import (
	"golang.org/x/exp/slog"
	"golang.org/x/sys/windows"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
)

func GetPrinter(languages ...string) *message.Printer {
	tag, _ := language.MatchStrings(
		language.NewMatcher([]language.Tag{
			language.English, language.SimplifiedChinese,
		}),
		append(languages, GetLanguages()...)...,
	)
	return message.NewPrinter(tag)
}

func GetLanguages() []string {
	languages, err := windows.GetUserPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if err != nil {
		languages, err = windows.GetSystemPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
		if err != nil {
			slog.Error("failed to get system preferred UI languages", err)
			os.Exit(1)
		}
	}

	return languages
}
