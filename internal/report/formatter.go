package report

import (
	"encoding/json"
	"fmt"
	"strings"
)

func JSON(report ClassReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

func Text(report ClassReport) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Class %s/%s\n", report.Summary.SchoolID, report.Summary.ClassID)
	fmt.Fprintf(&builder, "Children: %d  Valid: %d  Rejected: %d\n", report.Summary.TotalChildren, report.Summary.ValidRecords, report.Summary.RejectedRecords)
	fmt.Fprintf(&builder, "Averages: height %.1f cm, weight %.1f kg, grip %.1f kg, reaction %.1f ms\n", report.Summary.AverageHeightCm, report.Summary.AverageWeightKg, report.Summary.AverageGripKg, report.Summary.AverageReactionMs)
	for _, child := range report.Children {
		fmt.Fprintf(&builder, "- %s %s [%s]\n", child.ChildID, child.Name, child.Status)
	}
	if len(report.Issues) > 0 {
		builder.WriteString("Issues:\n")
		for _, issue := range report.Issues {
			fmt.Fprintf(&builder, "  row %d %s: %s\n", issue.Row, issue.Field, issue.Message)
		}
	}
	return builder.String()
}

func Lines(report ClassReport) []string {
	text := Text(report)
	return strings.Split(strings.TrimSuffix(text, "\n"), "\n")
}
