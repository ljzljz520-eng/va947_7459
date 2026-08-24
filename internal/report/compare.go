package report

import (
	"fmt"
	"sort"

	"example.com/childfitness/internal/domain"
)

type Comparison struct {
	Metric     domain.MetricName `json:"metric"`
	Left       float64           `json:"left"`
	Right      float64           `json:"right"`
	Difference float64           `json:"difference"`
	Direction  string            `json:"direction"`
}

type ChildComparison struct {
	ChildID     string       `json:"child_id"`
	Comparisons []Comparison `json:"comparisons"`
	Headline    string       `json:"headline"`
}

func CompareRecords(left, right domain.MeasurementRecord) ChildComparison {
	comparisons := []Comparison{
		comparison(domain.MetricHeight, left.HeightCm, right.HeightCm),
		comparison(domain.MetricWeight, left.WeightKg, right.WeightKg),
		comparison(domain.MetricGrip, left.GripKg, right.GripKg),
		comparison(domain.MetricReaction, float64(left.ReactionMs), float64(right.ReactionMs)),
	}
	return ChildComparison{ChildID: left.ChildID, Comparisons: comparisons, Headline: comparisonHeadline(comparisons)}
}

func comparison(metric domain.MetricName, left, right float64) Comparison {
	difference := right - left
	direction := "flat"
	if difference > 0.0001 {
		direction = "up"
	} else if difference < -0.0001 {
		direction = "down"
	}
	return Comparison{Metric: metric, Left: left, Right: right, Difference: difference, Direction: direction}
}

func comparisonHeadline(comparisons []Comparison) string {
	up, down := 0, 0
	for _, item := range comparisons {
		if item.Direction == "up" {
			up++
		} else if item.Direction == "down" {
			down++
		}
	}
	if up == 0 && down == 0 {
		return "unchanged"
	}
	return fmt.Sprintf("%d improving, %d needing review", up, down)
}

func SortComparisons(items []ChildComparison) []ChildComparison {
	ordered := append([]ChildComparison{}, items...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ChildID < ordered[j].ChildID })
	return ordered
}

func CompareClassRecords(before, after []domain.MeasurementRecord) []ChildComparison {
	left := make(map[string]domain.MeasurementRecord)
	for _, record := range before {
		left[record.ChildID] = record
	}
	items := make([]ChildComparison, 0)
	for _, record := range after {
		if previous, ok := left[record.ChildID]; ok {
			items = append(items, CompareRecords(previous, record))
		}
	}
	return SortComparisons(items)
}

func ImprovementCount(items []ChildComparison) int {
	count := 0
	for _, item := range items {
		for _, comparison := range item.Comparisons {
			if comparison.Direction == "up" {
				count++
				break
			}
		}
	}
	return count
}
