package service

import (
	"example.com/childfitness/internal/domain"
	"example.com/childfitness/internal/store"
)

type Repository interface {
	SaveProfile(domain.ChildProfile) error
	GetProfile(string, string, string) (domain.ChildProfile, error)
	ListProfiles(string, string) ([]domain.ChildProfile, error)
	SaveBatch(domain.FitnessBatch) error
	GetBatch(string) (domain.FitnessBatch, error)
	ListBatches(string, string) ([]domain.FitnessBatch, error)
	RecordsForClass(string, string) ([]domain.MeasurementRecord, error)
	SaveSummary(domain.ClassSummary) error
	GetSummary(string, string) (domain.ClassSummary, error)
}

var _ Repository = (*store.Store)(nil)
