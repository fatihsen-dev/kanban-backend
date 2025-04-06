package domain

import (
	"time"
)

type Project struct {
	ID        string
	Name      string
	CreatedAt time.Time
}
