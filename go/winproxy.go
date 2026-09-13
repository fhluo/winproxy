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
		Proxy:         base.Flags&settings.FlagProxy != 0,
		Script:        base.Flags&settings.FlagAutoProxyURL != 0,
		AutoDetect:    base.Flags&settings.FlagAutoDetect != 0,
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

func setFlag(base *settings.DefaultConnectionSettings, flag int32, v bool) {
	if v {
		base.Flags |= flag
	} else {
		base.Flags &^= flag
	}
}

// Apply writes the settings to the registry.
func (s Settings) Apply() error {
	base, err := settings.Read()
	if err != nil {
		return fmt.Errorf("failed to read settings: %w", err)
	}

	base.Version++
	setFlag(base, settings.FlagProxy, s.Proxy)
	setFlag(base, settings.FlagAutoProxyURL, s.Script)
	setFlag(base, settings.FlagAutoDetect, s.AutoDetect)
	base.ProxyAddress = s.ProxyAddress

	for i := range s.BypassList {
		s.BypassList[i] = strings.TrimSpace(s.BypassList[i])
	}
	base.BypassList = strings.Join(s.BypassList, ";")

	base.ScriptAddress = s.ScriptAddress

	return settings.Write(base)
}
