package alarm

import (
	"testing"

	"example.com/railway-signal-lake/internal/model"
)

func TestDuplicateAlarmAndStorm(t *testing.T) {
	events := []model.Event{
		{EventID: "a", DeviceID: "D1", DeviceType: "track_circuit", State: "FAULT", OccurredAtMs: 100},
		{EventID: "b", DeviceID: "D1", DeviceType: "track_circuit", State: "FAULT", OccurredAtMs: 200},
		{EventID: "c", DeviceID: "D2", DeviceType: "signal", State: "DARK", OccurredAtMs: 250},
	}
	got := Build(events, Config{DuplicateWindowMs: 1000, StormWindowMs: 1000, StormDeviceCount: 2})
	hasDuplicate, hasStorm := false, false
	for _, item := range got {
		if item.Kind == "duplicate_alarm" {
			hasDuplicate = true
		}
		if item.Kind == "alarm_storm" {
			hasStorm = true
		}
	}
	if !hasDuplicate || !hasStorm {
		t.Fatalf("missing incident: %+v", got)
	}
}
