package domain

import (
	"fmt"
	"strings"
)

const (
	BatchPending   = "pending"
	BatchComplete  = "complete"
	BatchPartial   = "partial"
	RecordValid    = "valid"
	RecordRejected = "rejected"
)

type ChildProfile struct {
	ID        string `json:"id"`
	SchoolID  string `json:"school_id"`
	ClassID   string `json:"class_id"`
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	Consent   bool   `json:"consent"`
}

type MeasurementRecord struct {
	ID         string  `json:"id"`
	BatchID    string  `json:"batch_id"`
	ChildID    string  `json:"child_id"`
	HeightCm   float64 `json:"height_cm"`
	WeightKg   float64 `json:"weight_kg"`
	GripKg     float64 `json:"grip_kg"`
	ReactionMs int     `json:"reaction_ms"`
	ReceivedAt string  `json:"received_at"`
	Status     string  `json:"status"`
}

type ImportIssue struct {
	Row     int    `json:"row"`
	ChildID string `json:"child_id"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type FitnessBatch struct {
	ID            string              `json:"id"`
	SchoolID      string              `json:"school_id"`
	ClassID       string              `json:"class_id"`
	Source        string              `json:"source"`
	Records       []MeasurementRecord `json:"records"`
	Issues        []ImportIssue       `json:"issues"`
	Status        string              `json:"status"`
	ImportedCount int                 `json:"imported_count"`
	ErrorCount    int                 `json:"error_count"`
}

type ClassSummary struct {
	SchoolID          string  `json:"school_id"`
	ClassID           string  `json:"class_id"`
	TotalChildren     int     `json:"total_children"`
	ValidRecords      int     `json:"valid_records"`
	RejectedRecords   int     `json:"rejected_records"`
	AverageHeightCm   float64 `json:"average_height_cm"`
	AverageWeightKg   float64 `json:"average_weight_kg"`
	AverageGripKg     float64 `json:"average_grip_kg"`
	AverageReactionMs float64 `json:"average_reaction_ms"`
	ReviewedMetrics   int     `json:"reviewed_metrics"`
	GeneratedFrom     string  `json:"generated_from"`
}

type ClassRoster struct {
	SchoolID string         `json:"school_id"`
	ClassID  string         `json:"class_id"`
	Children []ChildProfile `json:"children"`
}

func (p ChildProfile) Key() string {
	return strings.Join([]string{p.SchoolID, p.ClassID, p.ID}, ":")
}

func (r MeasurementRecord) Key() string {
	return strings.Join([]string{r.BatchID, r.ChildID}, ":")
}

func (b FitnessBatch) Key() string {
	return b.ID
}

func (b FitnessBatch) String() string {
	return fmt.Sprintf("%s/%s %s (%d records, %d issues)", b.SchoolID, b.ClassID, b.ID, len(b.Records), len(b.Issues))
}

func (s ClassSummary) HasData() bool {
	return s.ValidRecords > 0 || s.RejectedRecords > 0
}
