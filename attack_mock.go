package hazana

import (
	"context"
	"time"
)

type attackMock struct {
	sleep time.Duration
}

func (m *attackMock) Setup(c Config) error {
	return nil
}

func (m *attackMock) Do(ctx context.Context) DoResult {
	time.Sleep(m.sleep)
	return DoResult{}
}

func (m *attackMock) Teardown() error {
	return nil
}

func (m *attackMock) Clone() Attack {
	return m
}
