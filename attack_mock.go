package hazana

import (
	"context"
	"time"
)

type attackMock struct {
	sleep                     time.Duration
	afterCalled, beforeCalled bool
}

func (m *attackMock) Setup(c Config) error {
	return nil
}

func (m *attackMock) Do(ctx context.Context) DoResult {
	select {
	case <-time.After(m.sleep):
		return DoResult{}
	case <-ctx.Done():
		return DoResult{Error: ctx.Err()}
	}
}

func (m *attackMock) Teardown() error {
	return nil
}

func (m *attackMock) Clone() Attack {
	return m
}

func (m *attackMock) BeforeRun(c Config) error {
	m.beforeCalled = true
	return nil
}

func (m *attackMock) AfterRun(r *RunReport) error {
	m.afterCalled = true
	return nil
}
