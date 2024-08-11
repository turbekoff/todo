package entities

import "time"

type Session struct {
	ID       string
	Owner    string
	Device   string
	Token    string
	ExpireAt time.Time
}
