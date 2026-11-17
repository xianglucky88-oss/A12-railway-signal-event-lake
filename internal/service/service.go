package service

import (
	"example.com/railway-signal-lake/internal/analysis"
	"example.com/railway-signal-lake/internal/model"
	"example.com/railway-signal-lake/internal/store"
)

type Service struct {
	raw       *store.RawRepository
	artifacts *store.ArtifactRepository
}

func New(root string) *Service {
	return &Service{raw: store.NewRawRepository(root), artifacts: store.NewArtifactRepository(root)}
}
func (s *Service) IngestEvents(rows []model.Event) error         { return s.raw.AppendEvents(rows) }
func (s *Service) IngestLinks(rows []model.TopologyLink) error   { return s.raw.AppendLinks(rows) }
func (s *Service) IngestWorkOrders(rows []model.WorkOrder) error { return s.raw.AppendWorkOrders(rows) }
func (s *Service) Compute(config analysis.Config) (model.Result, error) {
	events, err := s.raw.ReadEvents()
	if err != nil {
		return model.Result{}, err
	}
	links, err := s.raw.ReadLinks()
	if err != nil {
		return model.Result{}, err
	}
	orders, err := s.raw.ReadWorkOrders()
	if err != nil {
		return model.Result{}, err
	}
	result := analysis.Compute(events, links, orders, config)
	if err := s.artifacts.Save(result); err != nil {
		return model.Result{}, err
	}
	return result, nil
}
