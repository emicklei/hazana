package hazana

import "time"

type attackMock struct {
	sleep time.Duration
}

func (m *attackMock) Setup(c Config) error {
	return nil
}

func (m *attackMock) Do() (requestIndex int, err error) {
	time.Sleep(m.sleep)
	return 0, nil
}

func (m *attackMock) TearDown() error {
	return nil
}

func (m *attackMock) Clone() Attack {
	return m
}
