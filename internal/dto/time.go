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
	return json.Marshal(d.Time.Format("2006-01-02"))
}

func (d JSONDate) GoString() string {
	return d.Time.Format("2006-01-02")
}
