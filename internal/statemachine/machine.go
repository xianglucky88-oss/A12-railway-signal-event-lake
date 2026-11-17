package statemachine

import (
	"sort"

	"example.com/railway-signal-lake/internal/model"
)

var allowed = map[string]map[string]map[string]bool{
	"signal": {
		"RED":    {"GREEN": true, "YELLOW": true, "DARK": true},
		"GREEN":  {"RED": true, "YELLOW": true, "DARK": true},
		"YELLOW": {"RED": true, "GREEN": true, "DARK": true},
		"DARK":   {"RED": true, "GREEN": true, "YELLOW": true},
	},
	"track_circuit": {
		"CLEAR":    {"OCCUPIED": true, "FAULT": true},
		"OCCUPIED": {"CLEAR": true, "FAULT": true},
		"FAULT":    {"CLEAR": true, "OCCUPIED": true},
	},
	"point_machine": {
		"NORMAL":          {"REVERSED": true, "OUT_OF_POSITION": true, "FAULT": true},
		"REVERSED":        {"NORMAL": true, "OUT_OF_POSITION": true, "FAULT": true},
		"OUT_OF_POSITION": {"NORMAL": true, "REVERSED": true, "FAULT": true},
		"FAULT":           {"NORMAL": true, "REVERSED": true},
	},
	"interlocking": {
		"ONLINE":   {"DEGRADED": true, "OFFLINE": true},
		"DEGRADED": {"ONLINE": true, "OFFLINE": true},
		"OFFLINE":  {"ONLINE": true, "DEGRADED": true},
	},
}

func validState(deviceType, state string) bool {
	transitions, ok := allowed[deviceType]
	if !ok {
		return false
	}
	for from := range transitions {
		if from == state {
			return true
		}
	}
	for _, to := range transitions {
		if to[state] {
			return true
		}
	}
	return false
}

// Rebuild uses sequence order, not arrival order, so late uploads cannot
// rewind the reconstructed device state.
func Rebuild(events []model.Event) []model.StateSnapshot {
	byDevice := make(map[string][]model.Event)
	for _, event := range events {
		byDevice[event.DeviceID] = append(byDevice[event.DeviceID], event)
	}
	var snapshots []model.StateSnapshot
	for deviceID, items := range byDevice {
		sort.Slice(items, func(i, j int) bool {
			if items[i].Sequence != items[j].Sequence {
				return items[i].Sequence < items[j].Sequence
			}
			if items[i].OccurredAtMs != items[j].OccurredAtMs {
				return items[i].OccurredAtMs < items[j].OccurredAtMs
			}
			return items[i].EventID < items[j].EventID
		})
		previous := ""
		for _, event := range items {
			invalid := !validState(event.DeviceType, event.State)
			if previous != "" && previous != event.State {
				if transitions, ok := allowed[event.DeviceType]; ok && transitions[previous] != nil && !transitions[previous][event.State] {
					invalid = true
				}
			}
			snapshots = append(snapshots, model.StateSnapshot{
				DeviceID: deviceID, DeviceType: event.DeviceType, State: event.State,
				Sequence: event.Sequence, OccurredAtMs: event.OccurredAtMs, EventID: event.EventID,
				InvalidTransition: invalid,
			})
			previous = event.State
		}
	}
	sort.Slice(snapshots, func(i, j int) bool {
		if snapshots[i].DeviceID != snapshots[j].DeviceID {
			return snapshots[i].DeviceID < snapshots[j].DeviceID
		}
		return snapshots[i].Sequence < snapshots[j].Sequence
	})
	return snapshots
}
