package models

import (
	"encoding/json"
	"time"
)

type ISOTime struct {
	time.Time
}

func (t *ISOTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(time.RFC3339))
}

func (t *ISOTime) UnmarshalJSON(data []byte) error {
	return t.Time.UnmarshalJSON(data)
}
