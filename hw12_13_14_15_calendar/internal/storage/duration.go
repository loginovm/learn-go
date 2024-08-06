package storage

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Duration time.Duration

func (d Duration) Value() (driver.Value, error) {
	return driver.Value(int64(d)), nil
}

func (d *Duration) Scan(raw interface{}) error {
	switch v := raw.(type) {
	case int64:
		*d = Duration(v)
	case nil:
		*d = Duration(0)
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Duration from: %#v", v)
	}
	return nil
}
