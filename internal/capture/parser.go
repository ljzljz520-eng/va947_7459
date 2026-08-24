package capture

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"example.com/childfitness/internal/domain"
)

type DeviceRow struct {
	Row        int
	ChildID    string
	Height     string
	Weight     string
	Grip       string
	Reaction   string
	ReceivedAt string
}

type Parser struct {
	StrictColumns bool
}

func NewParser() Parser {
	return Parser{StrictColumns: true}
}

func (p Parser) ReadRows(input string) ([]DeviceRow, []domain.ImportIssue, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.FieldsPerRecord = -1
	rows := make([]DeviceRow, 0)
	issues := make([]domain.ImportIssue, 0)
	rowNumber := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowNumber++
		if err != nil {
			issues = append(issues, domain.ImportIssue{Row: rowNumber, Field: "csv", Message: err.Error()})
			continue
		}
		if rowNumber == 1 && len(record) > 0 && strings.EqualFold(strings.TrimSpace(record[0]), "child_id") {
			continue
		}
		if p.StrictColumns && len(record) != 6 {
			childID := ""
			if len(record) > 0 {
				childID = strings.TrimSpace(record[0])
			}
			issues = append(issues, domain.ImportIssue{Row: rowNumber, ChildID: childID, Field: "columns", Message: fmt.Sprintf("expected 6 columns, got %d", len(record))})
			continue
		}
		if len(record) < 6 {
			issues = append(issues, domain.ImportIssue{Row: rowNumber, Field: "columns", Message: "not enough columns"})
			continue
		}
		rows = append(rows, DeviceRow{Row: rowNumber, ChildID: strings.TrimSpace(record[0]), Height: strings.TrimSpace(record[1]), Weight: strings.TrimSpace(record[2]), Grip: strings.TrimSpace(record[3]), Reaction: strings.TrimSpace(record[4]), ReceivedAt: strings.TrimSpace(record[5])})
	}
	return rows, issues, nil
}

func ParseRow(row DeviceRow, batchID string) (domain.MeasurementRecord, *domain.ImportIssue) {
	record := domain.MeasurementRecord{ID: fmt.Sprintf("%s-%d", batchID, row.Row), BatchID: batchID, ChildID: row.ChildID, ReceivedAt: row.ReceivedAt, Status: domain.RecordValid}
	values := []struct {
		name  string
		value string
		apply func(float64)
	}{
		{name: "height", value: row.Height, apply: func(v float64) { record.HeightCm = v }},
		{name: "weight", value: row.Weight, apply: func(v float64) { record.WeightKg = v }},
		{name: "grip", value: row.Grip, apply: func(v float64) { record.GripKg = v }},
	}
	for _, item := range values {
		value, err := strconv.ParseFloat(item.value, 64)
		if err != nil {
			return rejectedRecord(record), &domain.ImportIssue{Row: row.Row, ChildID: row.ChildID, Field: item.name, Message: fmt.Sprintf("%q is not a number", item.value)}
		}
		item.apply(value)
	}
	reaction, err := strconv.Atoi(row.Reaction)
	if err != nil {
		return rejectedRecord(record), &domain.ImportIssue{Row: row.Row, ChildID: row.ChildID, Field: "reaction", Message: fmt.Sprintf("%q is not an integer", row.Reaction)}
	}
	record.ReactionMs = reaction
	if err := domain.ValidateMeasurement(record); err != nil {
		return rejectedRecord(record), &domain.ImportIssue{Row: row.Row, ChildID: row.ChildID, Field: "range", Message: err.Error()}
	}
	return record, nil
}

func rejectedRecord(record domain.MeasurementRecord) domain.MeasurementRecord {
	record.Status = domain.RecordRejected
	return record
}
