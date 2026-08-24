package domain

import (
	"math"
	"sort"
)

type MetricSeries struct {
	Metric MetricName `json:"metric"`
	Points []float64  `json:"points"`
}

type MetricRange struct {
	Metric MetricName `json:"metric"`
	Min    float64    `json:"min"`
	Max    float64    `json:"max"`
	Mean   float64    `json:"mean"`
	Count  int        `json:"count"`
}

type ClassMetricProfile struct {
	SchoolID string        `json:"school_id"`
	ClassID  string        `json:"class_id"`
	Ranges   []MetricRange `json:"ranges"`
	Records  int           `json:"records"`
}

func ValuesForMetric(records []MeasurementRecord, metric MetricName) []float64 {
	values := make([]float64, 0, len(records))
	for _, record := range records {
		if record.Status != RecordValid {
			continue
		}
		switch metric {
		case MetricHeight:
			values = append(values, record.HeightCm)
		case MetricWeight:
			values = append(values, record.WeightKg)
		case MetricGrip:
			values = append(values, record.GripKg)
		case MetricReaction:
			values = append(values, float64(record.ReactionMs))
		}
	}
	return values
}

func RangeForMetric(records []MeasurementRecord, metric MetricName) MetricRange {
	values := ValuesForMetric(records, metric)
	rangeValue := MetricRange{Metric: metric, Count: len(values)}
	if len(values) == 0 {
		return rangeValue
	}
	rangeValue.Min = values[0]
	rangeValue.Max = values[0]
	for _, value := range values[1:] {
		if value < rangeValue.Min {
			rangeValue.Min = value
		}
		if value > rangeValue.Max {
			rangeValue.Max = value
		}
	}
	rangeValue.Mean = Mean(values)
	return rangeValue
}

func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func Median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	ordered := append([]float64{}, values...)
	sort.Float64s(ordered)
	middle := len(ordered) / 2
	if len(ordered)%2 == 0 {
		return (ordered[middle-1] + ordered[middle]) / 2
	}
	return ordered[middle]
}

func StandardDeviation(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	mean := Mean(values)
	variance := 0.0
	for _, value := range values {
		delta := value - mean
		variance += delta * delta
	}
	return math.Sqrt(variance / float64(len(values)))
}

func BuildMetricProfile(schoolID, classID string, records []MeasurementRecord) ClassMetricProfile {
	metrics := []MetricName{MetricHeight, MetricWeight, MetricGrip, MetricReaction}
	profile := ClassMetricProfile{SchoolID: schoolID, ClassID: classID, Ranges: make([]MetricRange, 0, len(metrics)), Records: len(records)}
	for _, metric := range metrics {
		profile.Ranges = append(profile.Ranges, RangeForMetric(records, metric))
	}
	return profile
}

func SeriesForClass(records []MeasurementRecord, metric MetricName) MetricSeries {
	values := ValuesForMetric(records, metric)
	return MetricSeries{Metric: metric, Points: values}
}

func Trend(values []float64) string {
	if len(values) < 2 {
		return "flat"
	}
	first := values[0]
	last := values[len(values)-1]
	if math.Abs(last-first) < 0.0001 {
		return "flat"
	}
	if last > first {
		return "up"
	}
	return "down"
}
