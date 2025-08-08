package settings

import (
	"encoding/binary"
	"io"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

const (
	RegistryKeyPath   = `SOFTWARE\Microsoft\Windows\CurrentVersion\Internet settings\Connections`
	RegistryValueName = "DefaultConnectionSettings"
)

// ReadBinary reads the binary data of DefaultConnectionSettings from the registry.
func ReadBinary() ([]byte, error) {
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

// WriteBinary writes the binary data of DefaultConnectionSettings to the registry.
func WriteBinary(data []byte) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, RegistryKeyPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer func() {
		_ = key.Close()
	}()

	if err = key.SetBinaryValue(RegistryValueName, data); err != nil {
		return err
	}

	return nil
}

// Read reads and decodes the binary data of DefaultConnectionSettings from the registry.
func Read() (*DefaultConnectionSettings, error) {
	data, err := ReadBinary()
	if err != nil {
		return nil, err
	}

	settings := new(DefaultConnectionSettings)
	return settings, settings.UnmarshalBinary(data)
}

// Write encodes and writes the binary data of DefaultConnectionSettings to the registry.
func Write(settings *DefaultConnectionSettings) error {
	data, err := settings.MarshalBinary()
	if err != nil {
		return err
	}

	return WriteBinary(data)
}

func readString(r io.Reader) (string, error) {
	var size int32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return "", err
	}

	s := make([]byte, size)
	if err := binary.Read(r, binary.LittleEndian, s); err != nil {
		return "", err
	}

	return unsafe.String(unsafe.SliceData(s), size), nil
}

func writeString(w io.Writer, s string) error {
	if err := binary.Write(w, binary.LittleEndian, int32(len(s))); err != nil {
		return err
	}

	return binary.Write(w, binary.LittleEndian, unsafe.Slice(unsafe.StringData(s), len(s)))
}
