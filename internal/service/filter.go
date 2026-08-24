package service

import (
	"sort"
	"strings"

	"example.com/childfitness/internal/domain"
)

type ClassFilter struct {
	SchoolID        string
	ClassID         string
	OnlyValid       bool
	ChildQuery      string
	MinimumHeight   float64
	MaximumHeight   float64
	MinimumGrip     float64
	MaximumReaction int
}

type FilterResult struct {
	Records []domain.MeasurementRecord
	Matched int
	Skipped int
}

func FilterRecords(records []domain.MeasurementRecord, filter ClassFilter) FilterResult {
	result := FilterResult{Records: make([]domain.MeasurementRecord, 0)}
	query := strings.ToLower(strings.TrimSpace(filter.ChildQuery))
	for _, record := range records {
		if filter.OnlyValid && record.Status != domain.RecordValid {
			result.Skipped++
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(record.ChildID), query) {
			result.Skipped++
			continue
		}
		if filter.MinimumHeight > 0 && record.HeightCm < filter.MinimumHeight {
			result.Skipped++
			continue
		}
		if filter.MaximumHeight > 0 && record.HeightCm > filter.MaximumHeight {
			result.Skipped++
			continue
		}
		if filter.MinimumGrip > 0 && record.GripKg < filter.MinimumGrip {
			result.Skipped++
			continue
		}
		if filter.MaximumReaction > 0 && record.ReactionMs > filter.MaximumReaction {
			result.Skipped++
			continue
		}
		result.Records = append(result.Records, record)
	}
	result.Matched = len(result.Records)
	return result
}

func GroupRecordsByChild(records []domain.MeasurementRecord) map[string][]domain.MeasurementRecord {
	groups := make(map[string][]domain.MeasurementRecord)
	for _, record := range records {
		groups[record.ChildID] = append(groups[record.ChildID], record)
	}
	for childID := range groups {
		sort.Slice(groups[childID], func(i, j int) bool { return groups[childID][i].ID < groups[childID][j].ID })
	}
	return groups
}

func LatestRecord(records []domain.MeasurementRecord) (domain.MeasurementRecord, bool) {
	if len(records) == 0 {
		return domain.MeasurementRecord{}, false
	}
	ordered := append([]domain.MeasurementRecord{}, records...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ID > ordered[j].ID })
	return ordered[0], true
}

func LatestRecordsByChild(records []domain.MeasurementRecord) []domain.MeasurementRecord {
	groups := GroupRecordsByChild(records)
	latest := make([]domain.MeasurementRecord, 0, len(groups))
	for _, group := range groups {
		if record, ok := LatestRecord(group); ok {
			latest = append(latest, record)
		}
	}
	sort.Slice(latest, func(i, j int) bool { return latest[i].ChildID < latest[j].ChildID })
	return latest
}

func (s *ReviewService) FilterClassRecords(schoolID, classID string, filter ClassFilter) (FilterResult, error) {
	records, err := s.repo.RecordsForClass(schoolID, domain.NormalizeClassID(classID))
	if err != nil {
		return FilterResult{}, err
	}
	filter.SchoolID = schoolID
	filter.ClassID = classID
	return FilterRecords(records, filter), nil
}
