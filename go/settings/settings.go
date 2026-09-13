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
	Unknown       int32
	Version       int32
	Flags         int32
	ProxyAddress  string
	BypassList    string
	ScriptAddress string
	Unknown1      [32]byte
}

// New returns a new DefaultConnectionSettings with Unknown set to 70 and Flags set to FlagDirect.
func New() *DefaultConnectionSettings {
	return &DefaultConnectionSettings{
		Unknown: 70,
		Flags:   FlagDirect,
	}
}

// SetFlag enables or disables the given flag.
func (settings *DefaultConnectionSettings) SetFlag(flag int32, enabled bool) {
	if enabled {
		settings.Flags |= flag
	} else {
		settings.Flags &^= flag
	}
}

// HasFlag reports whether the given flag is set.
func (settings *DefaultConnectionSettings) HasFlag(flag int32) bool {
	return settings.Flags&flag != 0
}

func (settings *DefaultConnectionSettings) fields() iter.Seq[reflect.Value] {
	return func(yield func(reflect.Value) bool) {
		value := reflect.ValueOf(settings).Elem()
		for _, field := range value.Fields() {
			if !yield(field) {
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
