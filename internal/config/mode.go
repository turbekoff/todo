package config

import (
	"errors"
	"strings"
)

var (
	ErrMode = errors.New("mode is already set")
)

const (
	M_LOC  = "local"
	M_DEV  = "development"
	M_PROD = "production"
	M_NULL = "undefined"
)

type Mode string

func (m Mode) String() string {
	return string(m)
}

func (m Mode) RunAt(mode string, f func()) {
	if strings.ToLower(strings.TrimSpace(mode)) == string(m) {
		f()
	}
}

func (m *Mode) SetValue(s string) error {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case M_LOC, M_DEV, M_PROD:
		*m = Mode(s)
	default:
		*m = Mode(M_NULL)
	}

	return nil
}
