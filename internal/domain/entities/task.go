package entities

import (
	"time"
)

type Task struct {
	ID        string
	Owner     string
	Name      string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
