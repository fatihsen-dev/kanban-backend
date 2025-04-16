package domain

import (
	"time"
)

type Project struct {
	ID        string
	Name      string
	OwnerID   string
	CreatedAt time.Time
}
