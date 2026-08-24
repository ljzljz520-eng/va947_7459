package domain

import (
	"fmt"
	"math"
	"sort"
)

type MetricName string

const (
	MetricHeight   MetricName = "height"
	MetricWeight   MetricName = "weight"
	MetricGrip     MetricName = "grip"
	MetricReaction MetricName = "reaction"
)

type MetricResult struct {
	Name       MetricName `json:"name"`
	Value      float64    `json:"value"`
	Unit       string     `json:"unit"`
	Band       string     `json:"band"`
	Acceptable bool       `json:"acceptable"`
	Message    string     `json:"message"`
}

type RecordAssessment struct {
	RecordID string         `json:"record_id"`
	Metrics  []MetricResult `json:"metrics"`
	Overall  string         `json:"overall"`
}

type Threshold struct {
	Min float64
	Max float64
}

func AgeForYear(birthYear, measurementYear int) int {
	age := measurementYear - birthYear
	if age < 0 {
		return 0
	}
	return age
}

func BenchmarkFor(age int, metric MetricName) Threshold {
	base := map[MetricName]Threshold{
		MetricHeight:   {Min: 80, Max: 220},
		MetricWeight:   {Min: 10, Max: 160},
		MetricGrip:     {Min: 2, Max: 100},
		MetricReaction: {Min: 100, Max: 3000},
	}[metric]
	if age <= 6 {
		base.Min *= 0.8
		base.Max *= 0.8
	} else if age >= 13 {
		base.Min *= 1.1
		base.Max *= 1.1
	}
	return base
}

func EvaluateMetric(age int, metric MetricName, value float64) MetricResult {
	threshold := BenchmarkFor(age, metric)
	result := MetricResult{Name: metric, Value: value, Acceptable: value >= threshold.Min && value <= threshold.Max}
	switch metric {
	case MetricHeight:
		result.Unit = "cm"
	case MetricWeight, MetricGrip:
		result.Unit = "kg"
	case MetricReaction:
		result.Unit = "ms"
	default:
		result.Unit = "unknown"
	}
	if result.Acceptable {
		result.Band = "in-range"
		result.Message = fmt.Sprintf("%s is within %.1f-%.1f %s", metric, threshold.Min, threshold.Max, result.Unit)
	} else if value < threshold.Min {
		result.Band = "below-range"
		result.Message = fmt.Sprintf("%s is below %.1f %s", metric, threshold.Min, result.Unit)
	} else {
		result.Band = "above-range"
		result.Message = fmt.Sprintf("%s is above %.1f %s", metric, threshold.Max, result.Unit)
	}
	return result
}

func AssessRecord(age int, record MeasurementRecord) RecordAssessment {
	metrics := []MetricResult{
		EvaluateMetric(age, MetricHeight, record.HeightCm),
		EvaluateMetric(age, MetricWeight, record.WeightKg),
		EvaluateMetric(age, MetricGrip, record.GripKg),
		EvaluateMetric(age, MetricReaction, float64(record.ReactionMs)),
	}
	assessment := RecordAssessment{RecordID: record.ID, Metrics: metrics, Overall: "pass"}
	for _, metric := range metrics {
		if !metric.Acceptable {
			assessment.Overall = "review"
			break
		}
	}
	return assessment
}

func ComputeBMI(heightCm, weightKg float64) float64 {
	if heightCm <= 0 {
		return 0
	}
	meters := heightCm / 100
	return weightKg / (meters * meters)
}

func ClassifyBMI(bmi float64) string {
	switch {
	case bmi <= 0:
		return "unknown"
	case bmi < 14:
		return "low"
	case bmi < 22:
		return "typical"
	case bmi < 28:
		return "elevated"
	default:
		return "high"
	}
}

func CompareRecords(left, right MeasurementRecord) int {
	leftScore := left.HeightCm + left.WeightKg + left.GripKg - float64(left.ReactionMs)/100
	rightScore := right.HeightCm + right.WeightKg + right.GripKg - float64(right.ReactionMs)/100
	if math.Abs(leftScore-rightScore) < 0.0001 {
		if left.ID < right.ID {
			return -1
		}
		if left.ID > right.ID {
			return 1
		}
		return 0
	}
	if leftScore > rightScore {
		return -1
	}
	return 1
}

func RankRecords(records []MeasurementRecord) []MeasurementRecord {
	ordered := append([]MeasurementRecord{}, records...)
	sort.SliceStable(ordered, func(i, j int) bool { return CompareRecords(ordered[i], ordered[j]) < 0 })
	return ordered
}
