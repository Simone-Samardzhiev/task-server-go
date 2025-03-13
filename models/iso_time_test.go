package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestISOTimeMarshalJSON(t *testing.T) {
	isoTime := ISOTime{time.Now()}
	encodedIsoTime, err := isoTime.MarshalJSON()
	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}

	isoTimeString := string(encodedIsoTime)
	isoTimeString = strings.ReplaceAll(isoTimeString, `"`, "")
	if isoTimeString != isoTime.Format(time.RFC3339) {
		t.Errorf("Times doesn't match after marshalling %v != %v", isoTimeString, isoTime.Format(time.RFC3339))
	}
}

func TestISOTimeUnmarshalJSON(t *testing.T) {
	isoTime := ISOTime{time.Now()}
	encodedIsoTime, err := isoTime.MarshalJSON()
	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}

	err = json.Unmarshal(encodedIsoTime, &isoTime)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}
}
