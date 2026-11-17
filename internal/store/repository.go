package store

import (
	"path/filepath"

	"example.com/railway-signal-lake/internal/model"
)

type RawRepository struct{ Root string }

func NewRawRepository(root string) *RawRepository { return &RawRepository{Root: root} }
func (r *RawRepository) AppendEvents(rows []model.Event) error {
	return AppendJSONL(filepath.Join(r.Root, "raw", "events.jsonl"), rows)
}
func (r *RawRepository) AppendLinks(rows []model.TopologyLink) error {
	return AppendJSONL(filepath.Join(r.Root, "raw", "topology_links.jsonl"), rows)
}
func (r *RawRepository) AppendWorkOrders(rows []model.WorkOrder) error {
	return AppendJSONL(filepath.Join(r.Root, "raw", "work_orders.jsonl"), rows)
}
func (r *RawRepository) ReadEvents() ([]model.Event, error) {
	return ReadJSONL[model.Event](filepath.Join(r.Root, "raw", "events.jsonl"))
}
func (r *RawRepository) ReadLinks() ([]model.TopologyLink, error) {
	return ReadJSONL[model.TopologyLink](filepath.Join(r.Root, "raw", "topology_links.jsonl"))
}
func (r *RawRepository) ReadWorkOrders() ([]model.WorkOrder, error) {
	return ReadJSONL[model.WorkOrder](filepath.Join(r.Root, "raw", "work_orders.jsonl"))
}

type ArtifactRepository struct{ Root string }

func NewArtifactRepository(root string) *ArtifactRepository { return &ArtifactRepository{Root: root} }
func (r *ArtifactRepository) Save(result model.Result) error {
	if err := WriteJSONL(filepath.Join(r.Root, "artifacts", "state_snapshots.jsonl"), result.Snapshots); err != nil {
		return err
	}
	if err := WriteJSONL(filepath.Join(r.Root, "artifacts", "sequence_issues.jsonl"), result.Sequence); err != nil {
		return err
	}
	if err := WriteJSONL(filepath.Join(r.Root, "artifacts", "incidents.jsonl"), result.Incidents); err != nil {
		return err
	}
	return WriteJSONL(filepath.Join(r.Root, "artifacts", "correlations.jsonl"), result.Correlations)
}
