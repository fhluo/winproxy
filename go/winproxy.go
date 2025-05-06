package winproxy

import (
	"github.com/fhluo/winproxy/go/settings"
	"strings"
)

type Settings struct {
	base *settings.DefaultConnectionSettings

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
		base:          base,
		Proxy:         base.Flags&settings.FlagProxy != 0,
		Script:        base.Flags&settings.FlagAutoProxyURL != 0,
		AutoDetect:    base.Flags&settings.FlagAutoDetect != 0,
		ProxyAddress:  base.ProxyAddress,
		BypassList:    strings.Split(strings.TrimSpace(base.BypassList), ";"),
		ScriptAddress: base.ScriptAddress,
	}

	return
}

func (s Settings) setFlag(flag int32, v bool) {
	if v {
		s.base.Flags |= flag
	} else {
		s.base.Flags &^= flag
	}
}

// Apply writes the settings to the registry.
func (s Settings) Apply() error {
	if s.base == nil {
		s.base = settings.New()
	}

	s.base.Version++
	s.setFlag(settings.FlagProxy, s.Proxy)
	s.setFlag(settings.FlagAutoProxyURL, s.Script)
	s.setFlag(settings.FlagAutoDetect, s.AutoDetect)
	s.base.ProxyAddress = s.ProxyAddress

	for i := range s.BypassList {
		s.BypassList[i] = strings.TrimSpace(s.BypassList[i])
	}
	s.base.BypassList = strings.Join(s.BypassList, ";")

	s.base.ScriptAddress = s.ScriptAddress

	return settings.Write(s.base)
}
