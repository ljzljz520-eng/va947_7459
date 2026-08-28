package report

import (
	"testing"

	"example.com/childfitness/internal/domain"
)

func TestBuildReportAndFormat(t *testing.T) {
	summary := domain.ClassSummary{SchoolID: "S1", ClassID: "5A", TotalChildren: 1, ValidRecords: 1}
	report := BuildReport(summary, []domain.ChildProfile{{ID: "C1", Name: "Child"}}, []domain.MeasurementRecord{{ChildID: "C1", HeightCm: 130, Status: domain.RecordValid}}, []domain.ImportIssue{{Row: 2, Field: "grip", Message: "bad"}})
	if !ReportHasWarnings(report) || RecordCount(report) != 1 {
		t.Fatal("report warning or count missing")
	}
	if len(Lines(report)) < 2 {
		t.Fatal("text report too short")
	}
	if _, err := JSON(report); err != nil {
		t.Fatal(err)
	}
}
