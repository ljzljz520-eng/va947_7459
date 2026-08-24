package report

import (
	"sort"

	"example.com/childfitness/internal/domain"
)

type ClassReport struct {
	Summary  domain.ClassSummary  `json:"summary"`
	Children []ChildLine          `json:"children"`
	Records  []RecordLine         `json:"records"`
	Issues   []domain.ImportIssue `json:"issues"`
}

type ChildLine struct {
	ChildID string `json:"child_id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

type RecordLine struct {
	ChildID    string  `json:"child_id"`
	HeightCm   float64 `json:"height_cm"`
	WeightKg   float64 `json:"weight_kg"`
	GripKg     float64 `json:"grip_kg"`
	ReactionMs int     `json:"reaction_ms"`
	Status     string  `json:"status"`
}

func BuildReport(summary domain.ClassSummary, profiles []domain.ChildProfile, records []domain.MeasurementRecord, issues []domain.ImportIssue) ClassReport {
	report := ClassReport{Summary: summary, Children: make([]ChildLine, 0, len(profiles)), Records: make([]RecordLine, 0, len(records)), Issues: append([]domain.ImportIssue{}, issues...)}
	for _, profile := range profiles {
		report.Children = append(report.Children, ChildLine{ChildID: profile.ID, Name: profile.Name, Status: "enrolled"})
	}
	for _, record := range records {
		report.Records = append(report.Records, RecordLine{ChildID: record.ChildID, HeightCm: record.HeightCm, WeightKg: record.WeightKg, GripKg: record.GripKg, ReactionMs: record.ReactionMs, Status: record.Status})
	}
	sort.Slice(report.Children, func(i, j int) bool { return report.Children[i].ChildID < report.Children[j].ChildID })
	sort.Slice(report.Records, func(i, j int) bool { return report.Records[i].ChildID < report.Records[j].ChildID })
	return report
}

func ReportHasWarnings(report ClassReport) bool {
	return len(report.Issues) > 0 || report.Summary.RejectedRecords > 0
}

func RecordCount(report ClassReport) int {
	return len(report.Records)
}
