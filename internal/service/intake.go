package service

import (
	"errors"
	"strings"

	"example.com/childfitness/internal/capture"
	"example.com/childfitness/internal/domain"
)

type IntakeService struct {
	repo Repository
}

func NewIntakeService(repo Repository) *IntakeService {
	return &IntakeService{repo: repo}
}

func (s *IntakeService) ImportDeviceBatch(schoolID, classID, source, input string) (domain.FitnessBatch, []domain.ImportIssue, error) {
	if strings.TrimSpace(input) == "" {
		return domain.FitnessBatch{}, nil, errors.New("device input is empty")
	}
	result, err := capture.ImportDeviceData(schoolID, classID, source, input)
	if err != nil {
		return domain.FitnessBatch{}, nil, err
	}
	if err := s.repo.SaveBatch(result.Batch); err != nil {
		return domain.FitnessBatch{}, nil, err
	}
	return result.Batch, result.Issues, nil
}

func (s *IntakeService) RecordManual(batchID string, record domain.MeasurementRecord) error {
	if strings.TrimSpace(batchID) == "" {
		return errors.New("batch id is required")
	}
	if record.BatchID == "" {
		record.BatchID = batchID
	}
	if record.Status == "" {
		record.Status = domain.RecordValid
	}
	if err := domain.ValidateMeasurement(record); err != nil {
		return err
	}
	batch, err := s.repo.GetBatch(batchID)
	if err != nil {
		return err
	}
	batch.Records = append(batch.Records, record)
	batch.ImportedCount++
	return s.repo.SaveBatch(batch)
}

func (s *IntakeService) LoadBatch(batchID string) (domain.FitnessBatch, error) {
	if strings.TrimSpace(batchID) == "" {
		return domain.FitnessBatch{}, errors.New("batch id is required")
	}
	return s.repo.GetBatch(batchID)
}
