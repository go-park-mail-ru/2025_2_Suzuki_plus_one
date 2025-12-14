package dto

import (
	"encoding/json"
	"time"
)

type JSONDate struct {
	time.Time
}

func (d *JSONDate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d JSONDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format("2006-01-02"))
}

func (d JSONDate) GoString() string {
	return d.Format("2006-01-02")
}

type JSONDateTime struct {
	time.Time
}

func NewJSONDateTime(t time.Time) JSONDateTime {
	return JSONDateTime{Time: t}
}

func (dt *JSONDateTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse("15:04 02.01.2006", s)
	if err != nil {
		return err
	}
	dt.Time = t
	return nil
}

func (dt JSONDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.Format("15:04 02.01.2006"))
}

func (dt JSONDateTime) GoString() string {
	return dt.Format("15:04 02.01.2006")
}
