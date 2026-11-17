package workorder

import "example.com/railway-signal-lake/internal/model"

// Link attaches the first open/overlapping work order for each incident device.
func Link(incidents []model.Incident, orders []model.WorkOrder, toleranceMs int64) []model.Incident {
	result := append([]model.Incident(nil), incidents...)
	for i := range result {
		for _, order := range orders {
			if !contains(result[i].DeviceIDs, order.DeviceID) {
				continue
			}
			closed := order.ClosedAtMs
			if closed == 0 {
				closed = order.OpenedAtMs + toleranceMs
			}
			if order.OpenedAtMs <= result[i].EndMs+toleranceMs && closed >= result[i].StartMs-toleranceMs {
				result[i].WorkOrderID = order.WorkOrderID
				break
			}
		}
	}
	return result
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
