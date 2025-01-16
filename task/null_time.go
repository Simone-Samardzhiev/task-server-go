package task

import (
	"database/sql"
	"encoding/json"
	"time"
)

// NullTime type struct implements [sql.NullTime], with custom json encoding and decoding.
type NullTime struct {
	sql.NullTime
}

func (n *NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

func (n *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		n.Time = time.Time{}
	}

	if err := json.Unmarshal(data, &n.Time); err != nil {
		return err
	}

	n.Valid = true
	return nil
}
