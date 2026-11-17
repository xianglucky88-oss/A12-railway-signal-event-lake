package sequence

import (
	"testing"

	"example.com/railway-signal-lake/internal/model"
)

func TestReportsGapAndDuplicateSequence(t *testing.T) {
	events := []model.Event{
		{EventID: "e1", DeviceID: "D", DeviceType: "signal", Sequence: 1, State: "RED", OccurredAtMs: 100, ReceivedAtMs: 100},
		{EventID: "e2", DeviceID: "D", DeviceType: "signal", Sequence: 3, State: "GREEN", OccurredAtMs: 200, ReceivedAtMs: 200},
		{EventID: "e3", DeviceID: "D", DeviceType: "signal", Sequence: 3, State: "DARK", OccurredAtMs: 300, ReceivedAtMs: 300},
	}
	got := Normalize(events)
	if len(got.Issues) != 2 {
		t.Fatalf("expected gap and duplicate, got %+v", got.Issues)
	}
}
