package hazana

import "testing"

func TestLoadConfig(t *testing.T) {
	c := ConfigFromFile("config_test.json")
	if len(c.Metadata) != 0 {
		t.Error("expected empty metadata")
	}
	if c.RampupTimeSec != 10 {
		t.Error("expected RPS 10")
	}
}
