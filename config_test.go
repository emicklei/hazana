package hazana

import "testing"

func TestLoadConfig(t *testing.T) {
	c := ConfigFromFile("config_test.json")
	if len(c.Metadata) != 0 {
		t.Error("expected empty metadata")
	}
}
