package service

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"example.com/childfitness/internal/domain"
	"example.com/childfitness/internal/privacy"
)

type ReviewService struct {
	repo  Repository
	audit privacy.AuditTrail
}

func NewReviewService(repo Repository) *ReviewService {
	return &ReviewService{repo: repo, audit: privacy.AuditTrail{Events: make([]privacy.AuditEvent, 0)}}
}

func (s *ReviewService) ReviewClass(schoolID, classID string) (domain.ClassSummary, error) {
	if strings.TrimSpace(schoolID) == "" || strings.TrimSpace(classID) == "" {
		return domain.ClassSummary{}, errors.New("school and class are required")
	}
	classID = domain.NormalizeClassID(classID)
	profiles, err := s.repo.ListProfiles(schoolID, classID)
	if err != nil {
		return domain.ClassSummary{}, err
	}
	records, err := s.repo.RecordsForClass(schoolID, classID)
	if err != nil {
		return domain.ClassSummary{}, err
	}
	summary := summarize(schoolID, classID, profiles, records)
	if err := s.repo.SaveSummary(summary); err != nil {
		return domain.ClassSummary{}, err
	}
	return summary, nil
}

func summarize(schoolID, classID string, profiles []domain.ChildProfile, records []domain.MeasurementRecord) domain.ClassSummary {
	summary := domain.ClassSummary{SchoolID: schoolID, ClassID: classID, TotalChildren: len(profiles), GeneratedFrom: "stored-measurements"}
	var height, weight, grip, reaction float64
	for _, record := range records {
		if record.Status == domain.RecordValid {
			summary.ValidRecords++
			assessment := domain.AssessRecord(10, record)
			summary.ReviewedMetrics += len(assessment.Metrics)
			height += record.HeightCm
			weight += record.WeightKg
			grip += record.GripKg
			reaction += float64(record.ReactionMs)
		} else {
			summary.RejectedRecords++
		}
	}
	if summary.ValidRecords > 0 {
		count := float64(summary.ValidRecords)
		summary.AverageHeightCm = height / count
		summary.AverageWeightKg = weight / count
		summary.AverageGripKg = grip / count
		summary.AverageReactionMs = reaction / count
	}
	return summary
}

func (s *ReviewService) MaskedRoster(schoolID, classID string, audience privacy.Audience) ([]privacy.DisplayChild, error) {
	profiles, err := s.repo.ListProfiles(schoolID, domain.NormalizeClassID(classID))
	if err != nil {
		return nil, err
	}
	policy := privacy.PolicyFor(audience)
	visible := make([]privacy.DisplayChild, 0, len(profiles))
	for _, profile := range profiles {
		if privacy.CanViewClass(policy, schoolID, domain.NormalizeClassID(classID), profile) {
			visible = append(visible, privacy.MaskChild(profile, policy))
			s.audit.Record(policy, schoolID, classID, profile, privacy.MaskingReason(policy, profile))
		} else {
			s.audit.Record(policy, schoolID, classID, profile, "denied")
		}
	}
	sort.Slice(visible, func(i, j int) bool { return visible[i].ID < visible[j].ID })
	return visible, nil
}

func (s *ReviewService) AuditEvents() []privacy.AuditEvent {
	return append([]privacy.AuditEvent{}, s.audit.Events...)
}

func (s *ReviewService) MaskedMeasurements(schoolID, classID string, audience privacy.Audience) ([]privacy.DisplayRecord, error) {
	records, err := s.repo.RecordsForClass(schoolID, domain.NormalizeClassID(classID))
	if err != nil {
		return nil, err
	}
	policy := privacy.PolicyFor(audience)
	display := make([]privacy.DisplayRecord, 0, len(records))
	for _, record := range records {
		display = append(display, privacy.MaskRecord(record, policy))
	}
	return display, nil
}

func (s *ReviewService) BatchStatus(batchID string) (string, error) {
	batch, err := s.repo.GetBatch(batchID)
	if err != nil {
		return "", err
	}
	if batch.ErrorCount > 0 && batch.ImportedCount > 0 {
		return fmt.Sprintf("%s-with-%d-issues", batch.Status, batch.ErrorCount), nil
	}
	return batch.Status, nil
}
