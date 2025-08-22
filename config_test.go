package hazana

import (
	"flag"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c := ConfigFromFile("config_test.json")
	if len(c.Metadata) != 0 {
		t.Error("expected empty metadata")
	}
	if c.RampupTimeSec != 10 {
		t.Error("expected RampupTimeSec 10")
	}
	if c.DoTimeoutSec != 5 {
		t.Error("expected timeout 5")
	}
}

func TestOverrideLoadedConfig(t *testing.T) {
	_ = flag.Set("rps", "31")
	_ = flag.Set("attack", "32")
	_ = flag.Set("ramp", "33")
	_ = flag.Set("max", "34")
	_ = flag.Set("o", "here")
	_ = flag.Set("verbose", "false")
	_ = flag.Set("s", "?")
	_ = flag.Set("timeout", "35")
	c := ConfigFromFile("config_test.json")
	if got, want := c.RPS, 31; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := c.AttackTimeSec, 32; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := c.RampupTimeSec, 33; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := c.MaxAttackers, 34; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := c.OutputFilename, "here"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := c.Verbose, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := c.RampupStrategy, "?"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := c.MaxAttackers, 34; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := c.DoTimeoutSec, 35; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
