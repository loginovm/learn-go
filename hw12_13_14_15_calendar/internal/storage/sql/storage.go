package sqlstorage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib" // driver for sql db
	"github.com/jmoiron/sqlx"
	"github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/pressly/goose/v3"
)

var _ storage.EventRepo = (*Storage)(nil)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	var e storage.Event
	err := s.db.GetContext(ctx, &e, `
		SELECT id, title, description, start_at, end_at, remind_period, user_id FROM events
		WHERE id = $1
	`, id)

	return e, err
}

func (s *Storage) GetEventsInPeriod(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]storage.Event, error) {
	events := []storage.Event{}
	err := s.db.SelectContext(ctx, &events, `
SELECT 
    id, 
    title, 
    description, 
    start_at, 
    end_at, 
    remind_period, 
    user_id 
FROM events
WHERE start_at >= $1 AND start_at < $2 OR
      end_at >= $1 AND end_at <= $2`,
		startDate, endDate)

	return events, err
}

func (s *Storage) GetEvents(ctx context.Context) ([]storage.Event, error) {
	events := []storage.Event{}
	err := s.db.SelectContext(ctx, &events, `
SELECT 
    id, 
    title, 
    description, 
    start_at, 
    end_at, 
    remind_period, 
    user_id 
FROM events`)

	return events, err
}

func (s *Storage) CreateEvent(ctx context.Context, e storage.Event) error {
	query := `INSERT INTO events(title, description, start_at, end_at, remind_period, user_id)
VALUES (:title, :description, :start_at, :end_at, :remind_period, :user_id)`

	_, err := s.db.NamedExecContext(ctx, query, e)
	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, e storage.Event) error {
	query := `UPDATE events
SET title = :title,
    description = :description,
    start_at = :start_at,
    end_at = :end_at,
    remind_period = :remind_period,
    user_id = :user_id
WHERE id = :id
`
	_, err := s.db.NamedExecContext(ctx, query, e)

	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)
	return err
}

func (s *Storage) RunMigration(migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(s.db.DB, migrationsDir); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	s.db = db
	return db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}
