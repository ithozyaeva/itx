package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

const DateFormat = "2006-01-02"

func NewDateOnly(t *time.Time) *DateOnly {
	if t == nil {
		return nil
	}
	d := DateOnly(*t)
	return &d
}

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte(`null`), nil
	}
	return []byte(`"` + t.Format(DateFormat) + `"`), nil
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		return nil
	}

	t, err := time.Parse(DateFormat, str)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) Value() (driver.Value, error) {
	t := time.Time(d)
	if t.IsZero() {
		return nil, nil
	}
	return t.Format(DateFormat), nil
}

func (d *DateOnly) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = DateOnly(v)
		return nil
	case []byte:
		t, err := time.Parse(DateFormat, string(v))
		if err != nil {
			return err
		}
		*d = DateOnly(t)
		return nil
	case string:
		t, err := time.Parse(DateFormat, v)
		if err != nil {
			return err
		}
		*d = DateOnly(t)
		return nil
	default:
		return fmt.Errorf("cannot scan value %v into DateOnly", value)
	}
}
