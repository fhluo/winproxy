package winproxy

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/fhluo/winproxy/go/settings"
)

type Settings struct {
	Proxy        bool
	ProxyAddress string

	Script        bool
	ScriptAddress string

	AutoDetect bool

	BypassList []string
}

// ReadSettings reads the settings from the registry.
func ReadSettings() (s Settings, err error) {
	base, err := settings.Read()
	if err != nil {
		return
	}
	s = Settings{
		Proxy:         base.HasFlag(settings.FlagProxy),
		Script:        base.HasFlag(settings.FlagAutoProxyURL),
		AutoDetect:    base.HasFlag(settings.FlagAutoDetect),
		ProxyAddress:  base.ProxyAddress,
		BypassList:    slices.Collect(parseBypassList(base.BypassList)),
		ScriptAddress: base.ScriptAddress,
	}

	return
}

// parseBypassList returns an iterator over the trimmed non-empty items of a semicolon-separated bypass list.
func parseBypassList(value string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for item := range strings.SplitSeq(value, ";") {
			if item = strings.TrimSpace(item); item != "" && !yield(item) {
				return
			}
		}
	}
}

// formatBypassList returns a semicolon-separated bypass list of the trimmed non-empty items of list.
func formatBypassList(list []string) string {
	size := max(len(list)-1, 0)
	for _, item := range list {
		size += len(item)
	}

	var b strings.Builder
	b.Grow(size)
	for _, item := range list {
		if item = strings.TrimSpace(item); item != "" {
			if b.Len() > 0 {
				b.WriteByte(';')
			}
			b.WriteString(item)
		}
	}

	return b.String()
}

// Apply writes the settings to the registry.
func (s Settings) Apply() error {
	base, err := settings.Read()
	if err != nil {
		return fmt.Errorf("failed to read settings: %w", err)
	}

	base.Version++

	base.SetFlag(settings.FlagProxy, s.Proxy)
	base.SetFlag(settings.FlagAutoProxyURL, s.Script)
	base.SetFlag(settings.FlagAutoDetect, s.AutoDetect)

	base.ProxyAddress = s.ProxyAddress
	base.BypassList = formatBypassList(s.BypassList)
	base.ScriptAddress = s.ScriptAddress

	return settings.Write(base)
}
