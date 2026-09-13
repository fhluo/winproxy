package settings

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"testing"
	"unsafe"
)

func TestDefaultConnectionSettings(t *testing.T) {
	s, err := Read()
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(s, jsontext.WithIndent("  "))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", unsafe.String(unsafe.SliceData(data), len(data)))
}
