package model

import (
	"time"
)

type Time time.Time

type Project struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt Time   `json:"created_at"`
	UpdatedAt Time   `json:"updated_at"`
}

const ISO8601Layout = "2006-01-02T15:04:05"

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).UTC().Format(ISO8601Layout) + `"`), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	timeTime, err := ParseTime(`"`+ISO8601Layout+`"`, string(data))
	if err != nil {
		return err
	}

	*t = Time(timeTime)
	return err
}

func ParseTime(layout string, value string) (Time, error) {
	result, err := time.ParseInLocation(layout, value, time.UTC)
	return Time(result), err
}
