package analysis

import (
	"example.com/railway-signal-lake/internal/alarm"
	"example.com/railway-signal-lake/internal/correlation"
	"example.com/railway-signal-lake/internal/model"
	"example.com/railway-signal-lake/internal/sequence"
	"example.com/railway-signal-lake/internal/statemachine"
	"example.com/railway-signal-lake/internal/workorder"
)

type Config struct {
	Alarm                alarm.Config
	CorrelationWindowMs  int64
	WorkOrderToleranceMs int64
}

func Compute(events []model.Event, links []model.TopologyLink, orders []model.WorkOrder, config Config) model.Result {
	normalized := sequence.Normalize(events)
	snapshots := statemachine.Rebuild(normalized.Events)
	incidents := alarm.Build(normalized.Events, config.Alarm)
	incidents = workorder.Link(incidents, orders, config.WorkOrderToleranceMs)
	correlations := correlation.Correlate(incidents, links, config.CorrelationWindowMs)
	return model.Result{
		Snapshots: snapshots, Sequence: normalized.Issues, Incidents: incidents,
		Correlations: correlations, Rejected: normalized.Rejected, Duplicates: normalized.Duplicates,
	}
}
