-- +goose Up
CREATE TABLE events (
    id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title   text not null,
    description     text,
    start_at    timestamptz not null,
    end_at      timestamptz not null,
    remind_period      bigint not null,
    user_id integer not null
);

CREATE INDEX idx_events_user_id on events (user_id);
CREATE INDEX idx_events_start_at on events (start_at);

INSERT INTO events (title, description, start_at, end_at, remind_period, user_id)
VALUES
    ('daily meet', 'test description daily', '2024-08-01 11:00:00+03', '2024-08-01 11:30:00+03', 900000000000, 1),
    ('weekly meet', 'test description weekly', '2024-08-06 12:00:00+03', '2024-08-06 12:30:00+03', 900000000000, 1);

-- +goose Down
drop table events;
