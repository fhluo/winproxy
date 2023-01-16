package settings

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"unsafe"
)

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
