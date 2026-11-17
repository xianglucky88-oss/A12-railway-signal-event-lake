package workorder

import (
	"testing"

	"example.com/railway-signal-lake/internal/model"
)

func TestLinksOverlappingOrder(t *testing.T) {
	got := Link([]model.Incident{{IncidentID: "i", DeviceIDs: []string{"D"}, StartMs: 100, EndMs: 200}}, []model.WorkOrder{{WorkOrderID: "W", DeviceID: "D", OpenedAtMs: 50, ClosedAtMs: 300}}, 0)
	if got[0].WorkOrderID != "W" {
		t.Fatalf("unexpected link: %+v", got)
	}
}
