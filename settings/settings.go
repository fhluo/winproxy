package settings

import (
	"bytes"
	"encoding/binary"
	"iter"
	"reflect"
)

const (
	FlagDirect       = 1 << iota
	FlagProxy        // use an explicitly set proxy server
	FlagAutoProxyURL // use an automatic configuration script downloaded from a specified URL
	FlagAutoDetect   // automatically detect settings
)

// DefaultConnectionSettings is the struct representation of its registry value.
type DefaultConnectionSettings struct {
	Unknown       int32    `json:"unknown"`
	Version       int32    `json:"version"`
	Flags         int32    `json:"flags"`
	ProxyAddress  string   `json:"proxy_address"`
	BypassList    string   `json:"bypass_list"`
	ScriptAddress string   `json:"script_address"`
	UnKnown1      [32]byte `json:"unknown1"`
}

// New returns a new DefaultConnectionSettings with Unknown set to 70 and Flags set to FlagDirect.
func New() *DefaultConnectionSettings {
	return &DefaultConnectionSettings{
		Unknown: 70,
		Flags:   FlagDirect,
	}
}

func (settings *DefaultConnectionSettings) fields() iter.Seq[reflect.Value] {
	return func(yield func(reflect.Value) bool) {
		value := reflect.ValueOf(settings).Elem()
		for i := 0; i < value.NumField(); i++ {
			if !yield(value.Field(i)) {
				break
			}
		}
	}
}

// MarshalBinary encodes itself into a binary form and returns the result.
func (settings *DefaultConnectionSettings) MarshalBinary() (data []byte, err error) {
	buffer := new(bytes.Buffer)

	for field := range settings.fields() {
		switch field.Kind() {
		case reflect.String:
			if err = writeString(buffer, field.String()); err != nil {
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

// UnmarshalBinary decodes the binary data from the registry.
func (settings *DefaultConnectionSettings) UnmarshalBinary(data []byte) (err error) {
	buffer := bytes.NewBuffer(data)

	for field := range settings.fields() {
		switch field.Kind() {
		case reflect.String:
			var s string
			s, err = readString(buffer)
			if err != nil {
				return
			}

			field.SetString(s)
		default:
			if err = binary.Read(buffer, binary.LittleEndian, field.Addr().Interface()); err != nil {
				return
			}
		}
	}

	return
}
