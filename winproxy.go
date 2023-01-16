package winproxy

import (
	"github.com/fhluo/winproxy/settings"
	"strings"
)

type Settings struct {
	base *settings.DefaultConnectionSettings
}

func ReadSettings() (s Settings, err error) {
	data, err := settings.GetBinary()
	if err != nil {
		return
	}

	s = Settings{new(settings.DefaultConnectionSettings)}
	err = s.base.UnmarshalBinary(data)
	return
}

func (s Settings) Apply() error {
	s.base.Version++
	data, err := s.base.MarshalBinary()
	if err != nil {
		return err
	}

	return settings.SetBinary(data)
}

func (s Settings) setFlag(flag int32, v bool) {
	if v {
		s.base.Flags |= flag
	} else {
		s.base.Flags &^= flag
	}
}

func (s Settings) UseProxy() bool {
	return s.base.Flags&settings.FlagUseProxy != 0
}

func (s Settings) SetUseProxy(v bool) {
	s.setFlag(settings.FlagUseProxy, v)
}

func (s Settings) UseScript() bool {
	return s.base.Flags&settings.FlagUseScript != 0
}

func (s Settings) SetUseScript(v bool) {
	s.setFlag(settings.FlagUseScript, v)
}

func (s Settings) AutoDetect() bool {
	return s.base.Flags&settings.FlagAutoDetect != 0
}

func (s Settings) SetAutoDetect(v bool) {
	s.setFlag(settings.FlagAutoDetect, v)
}

func (s Settings) ProxyAddress() string {
	return s.base.ProxyAddress
}

func (s Settings) SetProxyAddress(address string) {
	s.base.ProxyAddress = address
}

func (s Settings) BypassList() []string {
	return strings.Split(s.base.BypassList, ";")
}

func (s Settings) SetBypassList(bypassList []string) {
	s.base.BypassList = strings.Join(bypassList, ";")
}

func (s Settings) ScriptAddress() string {
	return s.base.ScriptAddress
}

func (s Settings) SetScriptAddress(address string) {
	s.base.ScriptAddress = address
}
