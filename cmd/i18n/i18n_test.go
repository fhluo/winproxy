//go:build windows

package i18n

import (
	"golang.org/x/sys/windows"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"testing"
)

func TestLocalize(t *testing.T) {
	languages, err := windows.GetUserPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("UserPreferredUILanguages: %v", languages)

	languages, err = windows.GetSystemPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("SystemPreferredUILanguages: %v", languages)

	tag, _ := language.MatchStrings(
		language.NewMatcher([]language.Tag{language.English, language.SimplifiedChinese}),
		GetLanguages()...,
	)
	t.Logf("Matched language tag: %v(%v)", tag, display.Self.Name(tag))
}
