package capture

import (
	"fmt"
	"strings"

	"example.com/childfitness/internal/domain"
)

func CanonicalBatchID(schoolID, classID, source string) string {
	parts := []string{schoolID, classID, source}
	for i, part := range parts {
		parts[i] = strings.ToLower(strings.TrimSpace(part))
	}
	return strings.Join(parts, "-")
}

func CanonicalSource(source string) string {
	clean := strings.ToLower(strings.TrimSpace(source))
	if clean == "" {
		return "device"
	}
	return clean
}

func BuildBatch(schoolID, classID, source, input string) (domain.FitnessBatch, error) {
	if strings.TrimSpace(schoolID) == "" || strings.TrimSpace(classID) == "" {
		return domain.FitnessBatch{}, fmt.Errorf("school and class are required")
	}
	return domain.FitnessBatch{ID: CanonicalBatchID(schoolID, classID, source), SchoolID: strings.TrimSpace(schoolID), ClassID: domain.NormalizeClassID(classID), Source: CanonicalSource(source), Status: domain.BatchPending, Records: []domain.MeasurementRecord{}, Issues: []domain.ImportIssue{}}, nil
}

func AttachIssue(batch *domain.FitnessBatch, issue domain.ImportIssue) {
	if batch == nil {
		return
	}
	batch.Issues = append(batch.Issues, issue)
	batch.ErrorCount = len(batch.Issues)
}
