package correlation

import (
	"sort"

	"example.com/railway-signal-lake/internal/model"
)

// Correlate links incidents that each belong to one device. Multi-device
// incidents are aggregate results and must not be expanded into pairwise links.
func Correlate(incidents []model.Incident, links []model.TopologyLink, windowMs int64) []model.Correlation {
	if windowMs < 0 {
		windowMs = 0
	}
	adjacent := buildAdjacency(links)
	result := make([]model.Correlation, 0)
	for i, primary := range incidents {
		primaryDevices := uniqueDeviceIDs(primary.DeviceIDs)
		if len(primaryDevices) != 1 {
			continue
		}
		for _, other := range incidents[i+1:] {
			otherDevices := uniqueDeviceIDs(other.DeviceIDs)
			if len(otherDevices) != 1 || !overlaps(primary, other, windowMs) {
				continue
			}
			if !adjacent[primaryDevices[0]][otherDevices[0]] {
				continue
			}
			result = append(result, makeCorrelation(primary, other, primaryDevices, otherDevices))
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CorrelationID < result[j].CorrelationID })
	return result
}

func buildAdjacency(links []model.TopologyLink) map[string]map[string]bool {
	adjacent := make(map[string]map[string]bool)
	for _, link := range links {
		if link.FromDeviceID == "" || link.ToDeviceID == "" || link.FromDeviceID == link.ToDeviceID {
			continue
		}
		addNeighbor(adjacent, link.FromDeviceID, link.ToDeviceID)
		addNeighbor(adjacent, link.ToDeviceID, link.FromDeviceID)
	}
	return adjacent
}

func addNeighbor(adjacent map[string]map[string]bool, from, to string) {
	if adjacent[from] == nil {
		adjacent[from] = make(map[string]bool)
	}
	adjacent[from][to] = true
}

func overlaps(left, right model.Incident, windowMs int64) bool {
	return right.StartMs <= left.EndMs+windowMs && left.StartMs <= right.EndMs+windowMs
}

func makeCorrelation(left, right model.Incident, leftDevices, rightDevices []string) model.Correlation {
	first, second := left, right
	if first.IncidentID > second.IncidentID {
		first, second = second, first
	}
	return model.Correlation{
		CorrelationID: "corr:" + first.IncidentID + ":" + second.IncidentID,
		Rule:          "topology_alarm_correlation",
		PrimaryID:     first.IncidentID,
		RelatedIDs:    []string{second.IncidentID},
		DeviceIDs:     mergeDeviceIDs(leftDevices, rightDevices),
		StartMs:       min(left.StartMs, right.StartMs),
		EndMs:         max(left.EndMs, right.EndMs),
	}
}

func uniqueDeviceIDs(deviceIDs []string) []string {
	seen := make(map[string]bool, len(deviceIDs))
	result := make([]string, 0, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		if deviceID == "" || seen[deviceID] {
			continue
		}
		seen[deviceID] = true
		result = append(result, deviceID)
	}
	sort.Strings(result)
	return result
}

func mergeDeviceIDs(left, right []string) []string {
	merged := append(append([]string(nil), left...), right...)
	return uniqueDeviceIDs(merged)
}
