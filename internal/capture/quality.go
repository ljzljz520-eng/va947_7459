package capture

import (
	"fmt"
	"sort"
	"strings"

	"example.com/childfitness/internal/domain"
)

type QualityReport struct {
	TotalRows     int            `json:"total_rows"`
	AcceptedRows  int            `json:"accepted_rows"`
	RejectedRows  int            `json:"rejected_rows"`
	DuplicateRows int            `json:"duplicate_rows"`
	IssueByField  map[string]int `json:"issue_by_field"`
	Warnings      []string       `json:"warnings"`
}

func AnalyzeRecords(records []domain.MeasurementRecord, issues []domain.ImportIssue) QualityReport {
	report := QualityReport{TotalRows: len(records) + len(issues), AcceptedRows: len(records), RejectedRows: len(issues), IssueByField: make(map[string]int), Warnings: []string{}}
	for _, issue := range issues {
		report.IssueByField[issue.Field]++
	}
	if report.RejectedRows > 0 {
		report.Warnings = append(report.Warnings, fmt.Sprintf("%d rows need correction", report.RejectedRows))
	}
	if report.AcceptedRows == 0 && report.TotalRows > 0 {
		report.Warnings = append(report.Warnings, "no valid measurements were received")
	}
	return report
}

func StableRecordID(batchID string, row int) string {
	return fmt.Sprintf("%s-row-%04d", strings.TrimSpace(batchID), row)
}

func DeduplicateRecords(records []domain.MeasurementRecord) ([]domain.MeasurementRecord, int) {
	seen := make(map[string]domain.MeasurementRecord, len(records))
	duplicates := 0
	for _, record := range records {
		key := record.ChildID
		if key == "" {
			key = record.ID
		}
		if _, exists := seen[key]; exists {
			duplicates++
		}
		seen[key] = record
	}
	result := make([]domain.MeasurementRecord, 0, len(seen))
	for _, record := range seen {
		result = append(result, record)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ChildID < result[j].ChildID })
	return result, duplicates
}

func ReconcileBatch(batch domain.FitnessBatch) (domain.FitnessBatch, QualityReport) {
	clean, duplicates := DeduplicateRecords(batch.Records)
	batch.Records = clean
	batch.ImportedCount = len(clean)
	quality := AnalyzeRecords(clean, batch.Issues)
	quality.DuplicateRows = duplicates
	if duplicates > 0 {
		quality.Warnings = append(quality.Warnings, fmt.Sprintf("%d duplicate child rows were collapsed", duplicates))
	}
	if len(batch.Issues) > 0 && len(clean) > 0 {
		batch.Status = domain.BatchPartial
	} else if len(batch.Issues) > 0 {
		batch.Status = domain.BatchPartial
	} else {
		batch.Status = domain.BatchComplete
	}
	return batch, quality
}

func IssueFields(issues []domain.ImportIssue) []string {
	seen := make(map[string]bool)
	fields := make([]string, 0, len(issues))
	for _, issue := range issues {
		if !seen[issue.Field] {
			seen[issue.Field] = true
			fields = append(fields, issue.Field)
		}
	}
	sort.Strings(fields)
	return fields
}
