package model

import "fmt"

type Event struct {
	EventID      string `json:"event_id"`
	DeviceID     string `json:"device_id"`
	DeviceType   string `json:"device_type"`
	Sequence     int64  `json:"sequence"`
	State        string `json:"state"`
	OccurredAtMs int64  `json:"occurred_at_ms"`
	ReceivedAtMs int64  `json:"received_at_ms"`
}

func (e Event) Validate() error {
	if e.EventID == "" || e.DeviceID == "" || e.DeviceType == "" || e.State == "" {
		return fmt.Errorf("event_id, device_id, device_type and state are required")
	}
	if e.Sequence < 0 {
		return fmt.Errorf("sequence must be non-negative")
	}
	if e.OccurredAtMs <= 0 || e.ReceivedAtMs <= 0 {
		return fmt.Errorf("timestamps must be positive")
	}
	switch e.DeviceType {
	case "signal", "track_circuit", "point_machine", "interlocking":
	default:
		return fmt.Errorf("unsupported device_type %q", e.DeviceType)
	}
	return nil
}

type TopologyLink struct {
	FromDeviceID string `json:"from_device_id"`
	ToDeviceID   string `json:"to_device_id"`
	Relation     string `json:"relation"`
}

type WorkOrder struct {
	WorkOrderID string `json:"work_order_id"`
	DeviceID    string `json:"device_id"`
	OpenedAtMs  int64  `json:"opened_at_ms"`
	ClosedAtMs  int64  `json:"closed_at_ms,omitempty"`
	Category    string `json:"category"`
}

type SequenceIssue struct {
	DeviceID string `json:"device_id"`
	Kind     string `json:"kind"`
	Expected int64  `json:"expected,omitempty"`
	Actual   int64  `json:"actual,omitempty"`
	EventID  string `json:"event_id"`
}

type StateSnapshot struct {
	DeviceID          string `json:"device_id"`
	DeviceType        string `json:"device_type"`
	State             string `json:"state"`
	Sequence          int64  `json:"sequence"`
	OccurredAtMs      int64  `json:"occurred_at_ms"`
	EventID           string `json:"event_id"`
	InvalidTransition bool   `json:"invalid_transition"`
}

type Incident struct {
	IncidentID  string   `json:"incident_id"`
	Kind        string   `json:"kind"`
	DeviceIDs   []string `json:"device_ids"`
	StartMs     int64    `json:"start_ms"`
	EndMs       int64    `json:"end_ms"`
	Severity    string   `json:"severity"`
	EventIDs    []string `json:"event_ids"`
	WorkOrderID string   `json:"work_order_id,omitempty"`
}

type Correlation struct {
	CorrelationID string   `json:"correlation_id"`
	Rule          string   `json:"rule"`
	PrimaryID     string   `json:"primary_incident_id"`
	RelatedIDs    []string `json:"related_incident_ids"`
	DeviceIDs     []string `json:"device_ids"`
	StartMs       int64    `json:"start_ms"`
	EndMs         int64    `json:"end_ms"`
}

type Result struct {
	Snapshots    []StateSnapshot `json:"state_snapshots"`
	Sequence     []SequenceIssue `json:"sequence_issues"`
	Incidents    []Incident      `json:"incidents"`
	Correlations []Correlation   `json:"correlations"`
	Rejected     int             `json:"rejected_events"`
	Duplicates   int             `json:"duplicate_events"`
}
