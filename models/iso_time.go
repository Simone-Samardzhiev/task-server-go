package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

func (t *ISOTime) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t *ISOTime) Scan(value interface{}) error {
	v, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot scan type %T into ISOTime", value)
	}
	t.Time = v
	return nil
}
