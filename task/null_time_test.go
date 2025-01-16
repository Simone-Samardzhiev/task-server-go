package task

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"
)

func TestNullTime_MarshalJSON(t *testing.T) {
	nullTime := NullTime{
		sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	data, err := json.Marshal(&nullTime)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

func TestNullTime_UnmarshalJSON(t *testing.T) {
	nullTime := NullTime{
		sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	data, err := json.Marshal(&nullTime)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))

	nullTime.Valid = false
	err = json.Unmarshal(data, &nullTime)
	if err != nil {
		t.Fatal(err)
	}

	if !nullTime.Valid {
		nullTime.Valid = false
	}
}
