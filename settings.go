package winproxy

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/sys/windows/registry"
	"io"
)

const (
	keyPath   = `SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings\Connections`
	valueName = "DefaultConnectionSettings"
)

const (
	flagDirect     = 1 << iota
	flagProxy      // use a proxy server
	flagUseScript  // use setup script
	flagAutoDetect // automatically detect settings
)

func getDefaultConnectionSettings() ([]byte, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = key.Close()
	}()

	data, _, err := key.GetBinaryValue(valueName)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func setDefaultConnectionSettings(value []byte) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer func() {
		_ = key.Close()
	}()

	if err = key.SetBinaryValue(valueName, value); err != nil {
		return err
	}

	return nil
}

func readString(r io.Reader) (string, error) {
	var size int32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return "", err
	}

	data := make([]byte, size)
	if _, err := io.ReadFull(r, data); err != nil {
		return "", err
	}

	return string(data), nil
}

func writeString(w io.Writer, s string) error {
	if err := binary.Write(w, binary.LittleEndian, int32(len(s))); err != nil {
		return err
	}

	if _, err := io.WriteString(w, s); err != nil {
		return err
	}

	return nil
}

type Settings struct {
	size          int32
	version       int32
	flags         int32
	ProxyAddress  string
	BypassList    string
	ScriptAddress string
}

func marshal(settings *Settings) (data []byte, err error) {
	data = make([]byte, settings.size+8)
	buffer := bytes.NewBuffer(data[:0])

	if err = binary.Write(buffer, binary.LittleEndian, settings.size); err != nil {
		return
	}
	if err = binary.Write(buffer, binary.LittleEndian, settings.version); err != nil {
		return
	}
	if err = binary.Write(buffer, binary.LittleEndian, settings.flags); err != nil {
		return
	}

	if err = writeString(buffer, settings.ProxyAddress); err != nil {
		return
	}
	if err = writeString(buffer, settings.BypassList); err != nil {
		return
	}
	if err = writeString(buffer, settings.ScriptAddress); err != nil {
		return
	}

	return
}

func unmarshal(data []byte) (settings *Settings, err error) {
	buffer := bytes.NewBuffer(data)
	settings = new(Settings)

	if err = binary.Read(buffer, binary.LittleEndian, &settings.size); err != nil {
		return
	}
	if err = binary.Read(buffer, binary.LittleEndian, &settings.version); err != nil {
		return
	}
	if err = binary.Read(buffer, binary.LittleEndian, &settings.flags); err != nil {
		return
	}

	if settings.ProxyAddress, err = readString(buffer); err != nil {
		return
	}
	if settings.BypassList, err = readString(buffer); err != nil {
		return
	}
	if settings.ScriptAddress, err = readString(buffer); err != nil {
		return
	}

	return
}

func ReadSettings() (*Settings, error) {
	data, err := getDefaultConnectionSettings()
	if err != nil {
		return nil, err
	}

	settings, err := unmarshal(data)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func WriteSettings(settings *Settings) error {
	data, err := marshal(settings)
	if err != nil {
		return err
	}

	return setDefaultConnectionSettings(data)
}

func (s *Settings) setFlag(flag int32, v bool) {
	if v {
		s.flags |= flag
	} else {
		s.flags &^= flag
	}
}

func (s *Settings) Direct() bool {
	return s.flags&flagDirect != 0
}

func (s *Settings) SetDirect(v bool) {
	s.setFlag(flagDirect, v)
}

func (s *Settings) Proxy() bool {
	return s.flags&flagProxy != 0
}

func (s *Settings) SetProxy(v bool) {
	s.setFlag(flagProxy, v)
}

func (s *Settings) UseScript() bool {
	return s.flags&flagUseScript != 0
}

func (s *Settings) SetUseScript(v bool) {
	s.setFlag(flagUseScript, v)
}

func (s *Settings) AutoDetect() bool {
	return s.flags&flagAutoDetect != 0
}

func (s *Settings) SetAutoDetect(v bool) {
	s.setFlag(flagAutoDetect, v)
}
