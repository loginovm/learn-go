package storage

import "time"

type Event struct {
	ID           string    `db:"id"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	StartAt      time.Time `db:"start_at"`
	EndAt        time.Time `db:"end_at"`
	RemindPeriod Duration  `db:"remind_period"`
	UserID       int       `db:"user_id"`
}
