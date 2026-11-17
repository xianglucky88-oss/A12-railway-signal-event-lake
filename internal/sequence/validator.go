package sequence

import (
	"sort"

	"example.com/railway-signal-lake/internal/model"
)

type Normalized struct {
	Events     []model.Event
	Issues     []model.SequenceIssue
	Rejected   int
	Duplicates int
}

// Normalize removes repeated event IDs, retains the latest received copy and
// reports sequence gaps or repeated sequence numbers per device.
func Normalize(events []model.Event) Normalized {
	byID := make(map[string]model.Event)
	rejected := 0
	duplicates := 0
	for _, event := range events {
		if err := event.Validate(); err != nil {
			rejected++
			continue
		}
		if previous, ok := byID[event.EventID]; ok {
			duplicates++
			if event.ReceivedAtMs < previous.ReceivedAtMs {
				continue
			}
		}
		byID[event.EventID] = event
	}
	out := make([]model.Event, 0, len(byID))
	for _, event := range byID {
		out = append(out, event)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].DeviceID != out[j].DeviceID {
			return out[i].DeviceID < out[j].DeviceID
		}
		if out[i].Sequence != out[j].Sequence {
			return out[i].Sequence < out[j].Sequence
		}
		if out[i].OccurredAtMs != out[j].OccurredAtMs {
			return out[i].OccurredAtMs < out[j].OccurredAtMs
		}
		return out[i].EventID < out[j].EventID
	})
	var issues []model.SequenceIssue
	last := make(map[string]int64)
	seen := make(map[string]map[int64]string)
	for _, event := range out {
		if seen[event.DeviceID] == nil {
			seen[event.DeviceID] = make(map[int64]string)
		}
		if firstID, ok := seen[event.DeviceID][event.Sequence]; ok {
			issues = append(issues, model.SequenceIssue{DeviceID: event.DeviceID, Kind: "duplicate_sequence", Actual: event.Sequence, EventID: event.EventID})
			_ = firstID
		} else {
			seen[event.DeviceID][event.Sequence] = event.EventID
		}
		if previous, ok := last[event.DeviceID]; ok && event.Sequence > previous+1 {
			issues = append(issues, model.SequenceIssue{DeviceID: event.DeviceID, Kind: "gap", Expected: previous + 1, Actual: event.Sequence, EventID: event.EventID})
		}
		if current, ok := last[event.DeviceID]; !ok || event.Sequence > current {
			last[event.DeviceID] = event.Sequence
		}
	}
	return Normalized{Events: out, Issues: issues, Rejected: rejected, Duplicates: duplicates}
}
