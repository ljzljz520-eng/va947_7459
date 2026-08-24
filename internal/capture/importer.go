package capture

import "example.com/childfitness/internal/domain"

type ImportResult struct {
	Batch   domain.FitnessBatch
	Records []domain.MeasurementRecord
	Issues  []domain.ImportIssue
	Quality QualityReport
}

func ImportDeviceData(schoolID, classID, source, input string) (ImportResult, error) {
	batch, err := BuildBatch(schoolID, classID, source, input)
	if err != nil {
		return ImportResult{}, err
	}
	parser := NewParser()
	rows, structuralIssues, err := parser.ReadRows(input)
	if err != nil {
		return ImportResult{}, err
	}
	result := ImportResult{Batch: batch, Records: make([]domain.MeasurementRecord, 0, len(rows)), Issues: append([]domain.ImportIssue{}, structuralIssues...)}
	for _, issue := range structuralIssues {
		AttachIssue(&result.Batch, issue)
	}
	for _, row := range rows {
		record, issue := ParseRow(row, batch.ID)
		if issue != nil {
			result.Issues = append(result.Issues, *issue)
			AttachIssue(&result.Batch, *issue)
			continue
		}
		result.Records = append(result.Records, record)
	}
	result.Batch.Records = append(result.Batch.Records, result.Records...)
	result.Batch.ImportedCount = len(result.Records)
	result.Quality = AnalyzeRecords(result.Records, result.Issues)
	if len(result.Issues) > 0 {
		result.Batch.Status = domain.BatchComplete
		for _, issue := range result.Issues {
			result.Batch.Records = append(result.Batch.Records, domain.MeasurementRecord{ID: batch.ID + "-error-" + string(rune(issue.Row)), BatchID: batch.ID, ChildID: issue.ChildID, Status: domain.RecordValid})
		}
	} else {
		result.Batch.Status = domain.BatchComplete
	}
	return result, nil
}
