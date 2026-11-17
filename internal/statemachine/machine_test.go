package statemachine

import (
	"testing"

	"example.com/railway-signal-lake/internal/model"
)

func TestLateArrivalDoesNotRewindState(t *testing.T) {
	events := []model.Event{
		{EventID: "e2", DeviceID: "D", DeviceType: "signal", Sequence: 2, State: "GREEN", OccurredAtMs: 200, ReceivedAtMs: 900},
		{EventID: "e1", DeviceID: "D", DeviceType: "signal", Sequence: 1, State: "RED", OccurredAtMs: 100, ReceivedAtMs: 1000},
	}
	got := Rebuild(events)
	if len(got) != 2 || got[1].State != "GREEN" {
		t.Fatalf("unexpected snapshots: %+v", got)
	}
}
