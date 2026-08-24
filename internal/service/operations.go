package service

import (
	"errors"
	"sort"
	"strings"

	"example.com/childfitness/internal/capture"
	"example.com/childfitness/internal/domain"
	"example.com/childfitness/internal/report"
)

type ClassContext struct {
	SchoolID string
	ClassID  string
	Profiles []domain.ChildProfile
	Batches  []domain.FitnessBatch
	Records  []domain.MeasurementRecord
	Summary  domain.ClassSummary
}

type BatchLedger struct {
	BatchID       string `json:"batch_id"`
	Status        string `json:"status"`
	ImportedCount int    `json:"imported_count"`
	ErrorCount    int    `json:"error_count"`
	HasWarnings   bool   `json:"has_warnings"`
}

func (s *ReviewService) AssembleClassContext(schoolID, classID string) (ClassContext, error) {
	if strings.TrimSpace(schoolID) == "" || strings.TrimSpace(classID) == "" {
		return ClassContext{}, errors.New("school and class are required")
	}
	classID = domain.NormalizeClassID(classID)
	profiles, err := s.repo.ListProfiles(schoolID, classID)
	if err != nil {
		return ClassContext{}, err
	}
	batches, err := s.repo.ListBatches(schoolID, classID)
	if err != nil {
		return ClassContext{}, err
	}
	records, err := s.repo.RecordsForClass(schoolID, classID)
	if err != nil {
		return ClassContext{}, err
	}
	summary := summarize(schoolID, classID, profiles, records)
	return ClassContext{SchoolID: schoolID, ClassID: classID, Profiles: profiles, Batches: batches, Records: records, Summary: summary}, nil
}

func BuildClassReport(context ClassContext) report.ClassReport {
	issues := make([]domain.ImportIssue, 0)
	for _, batch := range context.Batches {
		issues = append(issues, batch.Issues...)
	}
	return report.BuildReport(context.Summary, context.Profiles, context.Records, issues)
}

func Ledger(batch domain.FitnessBatch) BatchLedger {
	return BatchLedger{BatchID: batch.ID, Status: batch.Status, ImportedCount: batch.ImportedCount, ErrorCount: batch.ErrorCount, HasWarnings: len(batch.Issues) > 0}
}

func SortLedgers(batches []domain.FitnessBatch) []BatchLedger {
	ledgers := make([]BatchLedger, 0, len(batches))
	for _, batch := range batches {
		ledgers = append(ledgers, Ledger(batch))
	}
	sort.Slice(ledgers, func(i, j int) bool { return ledgers[i].BatchID < ledgers[j].BatchID })
	return ledgers
}

func (s *IntakeService) RepairBatch(batchID string) (domain.FitnessBatch, capture.QualityReport, error) {
	batch, err := s.repo.GetBatch(batchID)
	if err != nil {
		return domain.FitnessBatch{}, capture.QualityReport{}, err
	}
	clean, quality := capture.ReconcileBatch(batch)
	if err := s.repo.SaveBatch(clean); err != nil {
		return domain.FitnessBatch{}, quality, err
	}
	return clean, quality, nil
}

func (s *IntakeService) IssuesForBatch(batchID string) ([]domain.ImportIssue, error) {
	batch, err := s.repo.GetBatch(batchID)
	if err != nil {
		return nil, err
	}
	issues := append([]domain.ImportIssue{}, batch.Issues...)
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Row == issues[j].Row {
			return issues[i].Field < issues[j].Field
		}
		return issues[i].Row < issues[j].Row
	})
	return issues, nil
}

func (s *ReviewService) LatestBatch(schoolID, classID string) (domain.FitnessBatch, error) {
	batches, err := s.repo.ListBatches(schoolID, domain.NormalizeClassID(classID))
	if err != nil {
		return domain.FitnessBatch{}, err
	}
	if len(batches) == 0 {
		return domain.FitnessBatch{}, errors.New("no batch found")
	}
	sort.Slice(batches, func(i, j int) bool { return batches[i].ID > batches[j].ID })
	return batches[0], nil
}
