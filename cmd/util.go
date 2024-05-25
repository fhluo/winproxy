package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/fhluo/winproxy"
	"github.com/fhluo/winproxy/cmd/i18n"
	"strings"
)

type Settings struct {
	Proxy         string // bool
	ProxyAddress  string
	Script        string // bool
	ScriptAddress string
	AutoDetect    string // bool
	BypassList    []string
}

func (s Settings) BaseInfoTable() *table.Table {
	p := i18n.GetPrinter()

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

	t.Row(p.Sprintf("Use proxy"), s.Proxy)
	t.Row(p.Sprintf("Proxy address"), s.ProxyAddress)
	t.Row(p.Sprintf("Use script"), s.Script)
	t.Row(p.Sprintf("Script address"), s.ScriptAddress)
	t.Row(p.Sprintf("Auto-detect"), s.AutoDetect)

	return t
}

func (s Settings) BypassListTable() *table.Table {
	p := i18n.GetPrinter()

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

	t.Headers(p.Sprintf("Bypass list"))

	for _, address := range s.BypassList {
		t.Row(address)
	}

	return t
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

func PrintSettings(s winproxy.Settings) {
	fmt.Println(Render(s))
}
