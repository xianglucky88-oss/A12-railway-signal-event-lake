package alarm

import (
	"fmt"
	"sort"

	"example.com/railway-signal-lake/internal/model"
)

type Config struct {
	DuplicateWindowMs int64
	ChatterWindowMs   int64
	ChatterThreshold  int
	StormWindowMs     int64
	StormDeviceCount  int
}

func (c Config) normalized() Config {
	if c.DuplicateWindowMs <= 0 {
		c.DuplicateWindowMs = 60_000
	}
	if c.ChatterWindowMs <= 0 {
		c.ChatterWindowMs = 120_000
	}
	if c.ChatterThreshold <= 0 {
		c.ChatterThreshold = 3
	}
	if c.StormWindowMs <= 0 {
		c.StormWindowMs = 30_000
	}
	if c.StormDeviceCount <= 0 {
		c.StormDeviceCount = 2
	}
	return c
}

func isAlarm(event model.Event) bool {
	return event.State == "FAULT" || event.State == "DARK" || event.State == "OUT_OF_POSITION" || event.State == "OFFLINE"
}

// Build emits deduplicated alarm incidents, chatter incidents and alarm storms.
func Build(events []model.Event, config Config) []model.Incident {
	config = config.normalized()
	var incidents []model.Incident
	byDevice := make(map[string][]model.Event)
	for _, event := range events {
		if isAlarm(event) {
			byDevice[event.DeviceID] = append(byDevice[event.DeviceID], event)
		}
	}
	for deviceID, items := range byDevice {
		sort.Slice(items, func(i, j int) bool { return items[i].OccurredAtMs < items[j].OccurredAtMs })
		var cluster []model.Event
		flush := func() {
			if len(cluster) == 0 {
				return
			}
			incidents = append(incidents, makeIncident("alarm", []string{deviceID}, cluster))
			if len(cluster) > 1 {
				incidents = append(incidents, makeIncident("duplicate_alarm", []string{deviceID}, cluster))
			}
			cluster = nil
		}
		for _, event := range items {
			if len(cluster) == 0 || event.OccurredAtMs-cluster[len(cluster)-1].OccurredAtMs <= config.DuplicateWindowMs {
				cluster = append(cluster, event)
			} else {
				flush()
				cluster = append(cluster, event)
			}
		}
		flush()
		if countAlarmTransitions(items) >= config.ChatterThreshold {
			incidents = append(incidents, makeIncident("alarm_chatter", []string{deviceID}, items))
		}
	}
	incidents = append(incidents, detectStorms(events, config)...)
	sort.Slice(incidents, func(i, j int) bool {
		if incidents[i].StartMs != incidents[j].StartMs {
			return incidents[i].StartMs < incidents[j].StartMs
		}
		return incidents[i].IncidentID < incidents[j].IncidentID
	})
	return incidents
}

func countAlarmTransitions(items []model.Event) int {
	if len(items) == 0 {
		return 0
	}
	count := 1
	for i := 1; i < len(items); i++ {
		if items[i].State != items[i-1].State {
			count++
		}
	}
	return count
}

func detectStorms(events []model.Event, config Config) []model.Incident {
	alarms := make([]model.Event, 0)
	for _, event := range events {
		if isAlarm(event) {
			alarms = append(alarms, event)
		}
	}
	sort.Slice(alarms, func(i, j int) bool { return alarms[i].OccurredAtMs < alarms[j].OccurredAtMs })
	var result []model.Incident
	for start := 0; start < len(alarms); start++ {
		deviceSet := make(map[string]bool)
		var window []model.Event
		for end := start; end < len(alarms) && alarms[end].OccurredAtMs-alarms[start].OccurredAtMs <= config.StormWindowMs; end++ {
			window = append(window, alarms[end])
			deviceSet[alarms[end].DeviceID] = true
		}
		if len(deviceSet) >= config.StormDeviceCount {
			devices := make([]string, 0, len(deviceSet))
			for device := range deviceSet {
				devices = append(devices, device)
			}
			sort.Strings(devices)
			result = append(result, makeIncident("alarm_storm", devices, window))
		}
	}
	return result
}

func makeIncident(kind string, devices []string, events []model.Event) model.Incident {
	eventIDs := make([]string, 0, len(events))
	start, end := events[0].OccurredAtMs, events[0].OccurredAtMs
	for _, event := range events {
		eventIDs = append(eventIDs, event.EventID)
		if event.OccurredAtMs < start {
			start = event.OccurredAtMs
		}
		if event.OccurredAtMs > end {
			end = event.OccurredAtMs
		}
	}
	return model.Incident{
		IncidentID: fmt.Sprintf("%s:%s:%d:%d", kind, devices[0], start, end),
		Kind:       kind, DeviceIDs: devices, StartMs: start, EndMs: end,
		Severity: severity(kind), EventIDs: eventIDs,
	}
}

func severity(kind string) string {
	switch kind {
	case "alarm_storm":
		return "critical"
	case "alarm_chatter":
		return "high"
	default:
		return "medium"
	}
}
