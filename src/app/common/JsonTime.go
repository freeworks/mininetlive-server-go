package common

import (
	"database/sql/driver"
	"time"
)

type JsonTime struct {
	Time  time.Time
	Valid bool
}

func (j JsonTime) format() string {
	return time.Time(j.Time).Format("2006-01-02 15:04")
}

func (j JsonTime) MarshalText() ([]byte, error) {
	return []byte(j.format()), nil
}

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.format() + `"`), nil
}

func (j *JsonTime) Scan(value interface{}) error {
	j.Time, j.Valid = value.(time.Time)
	return nil
}

func (j JsonTime) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return j.Time, nil
}
