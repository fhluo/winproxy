package cmd

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/fhluo/winproxy"
	"github.com/samber/lo"
)

type Settings struct {
	Proxy         string // bool
	ProxyAddress  string
	Script        string // bool
	ScriptAddress string
	AutoDetect    string // bool
	BypassList    []string
}

func (s Settings) BaseInfoRows() [][]string {
	//p := i18n.GetPrinter()
	return [][]string{
		{Localize(&i18n.Message{ID: "Use proxy", Other: "Use proxy"}), s.Proxy},
		{Localize(&i18n.Message{ID: "Proxy address", Other: "Proxy address"}), s.ProxyAddress},
		{Localize(&i18n.Message{ID: "Use script", Other: "Use script"}), s.Script},
		{Localize(&i18n.Message{ID: "Script address", Other: "Script address"}), s.ScriptAddress},
		{Localize(&i18n.Message{ID: "Auto-detect", Other: "Auto-detect"}), s.AutoDetect},
	}
}

func (s Settings) BaseInfoTable() *table.Table {
	rows := s.BaseInfoRows()

	t := table.New()

	t.Border(lipgloss.NormalBorder())
	t.BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#31304D")))

	t.StyleFunc(func(row, col int) lipgloss.Style {
		if col == 0 {
			return lipgloss.NewStyle().Align(lipgloss.Center).Width(16).Foreground(lipgloss.Color("#B99470")).Bold(true)
		} else {
			return lipgloss.NewStyle().Align(lipgloss.Center).Width(16).Foreground(lipgloss.Color("#F0ECE5"))
		}
	})

	t.Rows(rows...)

	return t
}

func (s Settings) BypassListTable() *table.Table {
	//p := i18n.GetPrinter()
	headers := []string{Localize(&i18n.Message{ID: "Bypass list", Other: "Bypass list"})}
	rows := lo.Map(s.BypassList, func(item string, _ int) []string {
		return []string{item}
	})

	t := table.New()

	t.Border(lipgloss.NormalBorder())
	t.BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#31304D")))

	t.StyleFunc(func(row, col int) lipgloss.Style {
		if row == 0 {
			return lipgloss.NewStyle().Align(lipgloss.Center).Width(32).Foreground(lipgloss.Color("#B99470")).Bold(true)
		} else {
			return lipgloss.NewStyle().Align(lipgloss.Center).Width(32).Foreground(lipgloss.Color("#F0ECE5"))
		}
	})

	return t.Headers(headers...).Rows(rows...)
}

func (s Settings) String() string {
	return strings.Join([]string{
		fmt.Sprint(s.BaseInfoTable()),
		fmt.Sprint(s.BypassListTable()),
	}, "\n")
}

func CheckMark(b bool) string {
	if b {
		return "✓"
	} else {
		return "𐄂"
	}
}

func Render(s winproxy.Settings) Settings {
	return Settings{
		Proxy:         CheckMark(s.Proxy),
		ProxyAddress:  s.ProxyAddress,
		Script:        CheckMark(s.Script),
		ScriptAddress: s.ScriptAddress,
		AutoDetect:    CheckMark(s.AutoDetect),
		BypassList:    s.BypassList,
	}
}
