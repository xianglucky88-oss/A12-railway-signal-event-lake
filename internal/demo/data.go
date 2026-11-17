package demo

import "example.com/railway-signal-lake/internal/model"

func Events() []model.Event {
	b := int64(1_768_100_000_000)
	return []model.Event{
		{EventID: "E1", DeviceID: "SIG-01", DeviceType: "signal", Sequence: 1, State: "RED", OccurredAtMs: b, ReceivedAtMs: b + 50},
		{EventID: "E2", DeviceID: "SIG-01", DeviceType: "signal", Sequence: 3, State: "DARK", OccurredAtMs: b + 1000, ReceivedAtMs: b + 1050},
		{EventID: "E3", DeviceID: "TC-01", DeviceType: "track_circuit", Sequence: 1, State: "CLEAR", OccurredAtMs: b + 100, ReceivedAtMs: b + 150},
		{EventID: "E4", DeviceID: "TC-01", DeviceType: "track_circuit", Sequence: 2, State: "OCCUPIED", OccurredAtMs: b + 1200, ReceivedAtMs: b + 1250},
		{EventID: "E5", DeviceID: "TC-01", DeviceType: "track_circuit", Sequence: 3, State: "FAULT", OccurredAtMs: b + 1500, ReceivedAtMs: b + 1550},
		{EventID: "E6", DeviceID: "SIG-01", DeviceType: "signal", Sequence: 3, State: "DARK", OccurredAtMs: b + 1000, ReceivedAtMs: b + 3000},
	}
}
func Links() []model.TopologyLink {
	return []model.TopologyLink{{FromDeviceID: "SIG-01", ToDeviceID: "TC-01", Relation: "route_conflict"}}
}
func WorkOrders() []model.WorkOrder {
	return []model.WorkOrder{{WorkOrderID: "WO-7", DeviceID: "TC-01", OpenedAtMs: 1_768_100_000_500, ClosedAtMs: 1_768_100_005_000, Category: "track_circuit_fault"}}
}
