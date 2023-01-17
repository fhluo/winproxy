package settings

import "golang.org/x/sys/windows/registry"

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
