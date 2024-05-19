package settings

import (
	"encoding/json"
	"testing"
	"unsafe"
)

func TestDefaultConnectionSettings(t *testing.T) {
	s, err := Read()
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", unsafe.String(unsafe.SliceData(data), len(data)))
}
