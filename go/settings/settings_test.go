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

	settings := struct {
		Unknown       int32    `json:"unknown"`
		Version       int32    `json:"version"`
		Flags         int32    `json:"flags"`
		ProxyAddress  string   `json:"proxy_address"`
		BypassList    string   `json:"bypass_list"`
		ScriptAddress string   `json:"script_address"`
		Unknown1      [32]byte `json:"unknown1"`
	}(*s)

	data, err := json.Marshal(settings, jsontext.WithIndent("  "))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", unsafe.String(unsafe.SliceData(data), len(data)))
}
