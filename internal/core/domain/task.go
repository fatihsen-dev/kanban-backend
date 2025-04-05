package domain

import "time"

type Task struct {
	ID        string
	Title     string
	ColumnID  string
	ProjectID string
	CreatedAt time.Time
}
