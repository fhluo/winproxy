package settings

import "golang.org/x/sys/windows/registry"

const (
	RegistryKeyPath   = `SOFTWARE\Microsoft\Windows\CurrentVersion\Internet settings\Connections`
	RegistryValueName = "DefaultConnectionSettings"
)

func GetBinary() ([]byte, error) {
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

func SetBinary(value []byte) error {
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
