package winproxy

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/sys/windows/registry"
	"reflect"
	"unsafe"
)

const (
	RegistryKeyPath   = `SOFTWARE\Microsoft\Windows\CurrentVersion\Internet settings\Connections`
	RegistryValueName = "DefaultConnectionSettings"
)

func GetRegistryValue() ([]byte, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, RegistryKeyPath, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = key.Close()
	}()

	data, _, err := key.GetBinaryValue(RegistryValueName)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func SetRegistryValue(value []byte) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, RegistryKeyPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer func() {
		_ = key.Close()
	}()

	if err = key.SetBinaryValue(RegistryValueName, value); err != nil {
		return err
	}

	return nil
}

const (
	FlagDirect     = 1 << iota
	FlagUseProxy   // use a proxy server
	FlagUseScript  // use setup script
	FlagAutoDetect // automatically detect settings
)

type DefaultConnectionSettings struct {
	Unknown       int32
	Version       int32
	Flags         int32
	ProxyAddress  string
	BypassList    string
	ScriptAddress string
	UnKnown2      [32]byte
}

func (settings *DefaultConnectionSettings) MarshalBinary() (data []byte, err error) {
	buffer := new(bytes.Buffer)
	value := reflect.ValueOf(settings).Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		switch field.Kind() {
		case reflect.String:
			if err = binary.Write(buffer, binary.LittleEndian, int32(len(field.String()))); err != nil {
				return
			}

			if err = binary.Write(
				buffer, binary.LittleEndian,
				unsafe.Slice(unsafe.StringData(field.String()), len(field.String())),
			); err != nil {
				return
			}
		default:
			if err = binary.Write(buffer, binary.LittleEndian, field.Interface()); err != nil {
				return
			}
		}
	}

	data = buffer.Bytes()
	return
}

func (settings *DefaultConnectionSettings) UnmarshalBinary(data []byte) (err error) {
	buffer := bytes.NewBuffer(data)
	value := reflect.ValueOf(settings).Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		switch field.Kind() {
		case reflect.String:
			var size int32
			if err = binary.Read(buffer, binary.LittleEndian, &size); err != nil {
				return
			}

			s := make([]byte, size)
			if err = binary.Read(buffer, binary.LittleEndian, s); err != nil {
				return
			}

			field.SetString(unsafe.String(unsafe.SliceData(s), len(s)))
		default:
			if err = binary.Read(buffer, binary.LittleEndian, field.Addr().Interface()); err != nil {
				return
			}
		}
	}

	return
}

type Settings struct {
	raw *DefaultConnectionSettings
}

func (s Settings) Apply() error {
	s.raw.Version++
	data, err := s.raw.MarshalBinary()
	if err != nil {
		return err
	}

	return SetRegistryValue(data)
}

func ReadSettings() (s Settings, err error) {
	data, err := GetRegistryValue()
	if err != nil {
		return
	}

	s = Settings{new(DefaultConnectionSettings)}
	err = s.raw.UnmarshalBinary(data)
	return
}

func (s Settings) setFlag(flag int32, v bool) {
	if v {
		s.raw.Flags |= flag
	} else {
		s.raw.Flags &^= flag
	}
}

func (s Settings) UseProxy() bool {
	return s.raw.Flags&FlagUseProxy != 0
}

func (s Settings) SetUseProxy(v bool) {
	s.setFlag(FlagUseProxy, v)
}

func (s Settings) UseScript() bool {
	return s.raw.Flags&FlagUseScript != 0
}

func (s Settings) SetUseScript(v bool) {
	s.setFlag(FlagUseScript, v)
}

func (s Settings) AutoDetect() bool {
	return s.raw.Flags&FlagAutoDetect != 0
}

func (s Settings) SetAutoDetect(v bool) {
	s.setFlag(FlagAutoDetect, v)
}

func (s Settings) ProxyAddress() string {
	return s.raw.ProxyAddress
}

func (s Settings) SetProxyAddress(address string) {
	s.raw.ProxyAddress = address
}

func (s Settings) BypassList() string {
	return s.raw.BypassList
}

func (s Settings) SetBypassList(list string) {
	s.raw.BypassList = list
}

func (s Settings) ScriptAddress() string {
	return s.raw.ScriptAddress
}

func (s Settings) SetScriptAddress(address string) {
	s.raw.ScriptAddress = address
}
