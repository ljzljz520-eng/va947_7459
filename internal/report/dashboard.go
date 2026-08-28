package report

import (
	"encoding/csv"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"example.com/childfitness/internal/domain"
)

type Dashboard struct {
	Title       string         `json:"title"`
	Class       string         `json:"class"`
	Headline    string         `json:"headline"`
	MetricCards []MetricCard   `json:"metric_cards"`
	Warnings    []string       `json:"warnings"`
	Rows        []DashboardRow `json:"rows"`
}

type MetricCard struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type DashboardRow struct {
	ChildID    string  `json:"child_id"`
	HeightCm   float64 `json:"height_cm"`
	WeightKg   float64 `json:"weight_kg"`
	GripKg     float64 `json:"grip_kg"`
	ReactionMs int     `json:"reaction_ms"`
	BMI        float64 `json:"bmi"`
	BMIClass   string  `json:"bmi_class"`
}

func BuildDashboard(report ClassReport) Dashboard {
	dashboard := Dashboard{Title: "Child fitness class dashboard", Class: report.Summary.SchoolID + "/" + report.Summary.ClassID, Headline: fmt.Sprintf("%d of %d children have valid results", report.Summary.ValidRecords, report.Summary.TotalChildren), MetricCards: []MetricCard{{Label: "Average height", Value: report.Summary.AverageHeightCm, Unit: "cm"}, {Label: "Average weight", Value: report.Summary.AverageWeightKg, Unit: "kg"}, {Label: "Average grip", Value: report.Summary.AverageGripKg, Unit: "kg"}, {Label: "Average reaction", Value: report.Summary.AverageReactionMs, Unit: "ms"}}, Warnings: make([]string, 0)}
	for _, issue := range report.Issues {
		dashboard.Warnings = append(dashboard.Warnings, fmt.Sprintf("row %d: %s", issue.Row, issue.Message))
	}
	for _, record := range report.Records {
		if record.Status != domain.RecordValid {
			continue
		}
		bmi := domain.ComputeBMI(record.HeightCm, record.WeightKg)
		dashboard.Rows = append(dashboard.Rows, DashboardRow{ChildID: record.ChildID, HeightCm: record.HeightCm, WeightKg: record.WeightKg, GripKg: record.GripKg, ReactionMs: record.ReactionMs, BMI: bmi, BMIClass: domain.ClassifyBMI(bmi)})
	}
	sort.Slice(dashboard.Rows, func(i, j int) bool { return dashboard.Rows[i].ChildID < dashboard.Rows[j].ChildID })
	return dashboard
}

func DashboardCSV(dashboard Dashboard) ([]byte, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write([]string{"child_id", "height_cm", "weight_kg", "grip_kg", "reaction_ms", "bmi", "bmi_class"}); err != nil {
		return nil, err
	}
	for _, row := range dashboard.Rows {
		if err := writer.Write([]string{row.ChildID, strconv.FormatFloat(row.HeightCm, 'f', 1, 64), strconv.FormatFloat(row.WeightKg, 'f', 1, 64), strconv.FormatFloat(row.GripKg, 'f', 1, 64), strconv.Itoa(row.ReactionMs), strconv.FormatFloat(row.BMI, 'f', 2, 64), row.BMIClass}); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return []byte(builder.String()), nil
}

func WarningCount(dashboard Dashboard) int {
	return len(dashboard.Warnings)
}

func MetricCardLabels(dashboard Dashboard) []string {
	labels := make([]string, 0, len(dashboard.MetricCards))
	for _, card := range dashboard.MetricCards {
		labels = append(labels, card.Label)
	}
	return labels
}
