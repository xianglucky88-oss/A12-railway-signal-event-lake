package correlation

import (
	"reflect"
	"testing"

	"example.com/railway-signal-lake/internal/model"
)

func TestCorrelatesLinkedDevicesOnce(t *testing.T) {
	incidents := []model.Incident{
		{IncidentID: "a", DeviceIDs: []string{"S"}, StartMs: 100, EndMs: 100},
		{IncidentID: "b", DeviceIDs: []string{"T"}, StartMs: 200, EndMs: 200},
	}
	links := []model.TopologyLink{{FromDeviceID: "S", ToDeviceID: "T"}}
	got := Correlate(incidents, links, 200)
	if len(got) != 1 {
		t.Fatalf("expected one undirected correlation, got %+v", got)
	}
	if got[0].CorrelationID != "corr:a:b" {
		t.Fatalf("expected canonical correlation id, got %+v", got[0])
	}
	if want := []string{"S", "T"}; !reflect.DeepEqual(got[0].DeviceIDs, want) {
		t.Fatalf("expected deduplicated devices %v, got %v", want, got[0].DeviceIDs)
	}
}

func TestSkipsMultiDeviceStorm(t *testing.T) {
	incidents := []model.Incident{
		{IncidentID: "storm", Kind: "alarm_storm", DeviceIDs: []string{"S", "T"}, StartMs: 100, EndMs: 200},
		{IncidentID: "single", Kind: "alarm", DeviceIDs: []string{"T"}, StartMs: 200, EndMs: 200},
	}
	links := []model.TopologyLink{{FromDeviceID: "S", ToDeviceID: "T"}}
	if got := Correlate(incidents, links, 200); len(got) != 0 {
		t.Fatalf("expected aggregate storm to be skipped, got %+v", got)
	}
}
