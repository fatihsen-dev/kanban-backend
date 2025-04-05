package domain

import "time"

type Project struct {
	ID        string
	Title     string
	CreatedAt time.Time
}
