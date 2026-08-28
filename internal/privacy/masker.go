package privacy

import (
	"fmt"
	"strings"

	"example.com/childfitness/internal/domain"
)

type DisplayChild struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
}

type DisplayRecord struct {
	ChildID    string  `json:"child_id"`
	HeightCm   float64 `json:"height_cm"`
	WeightKg   float64 `json:"weight_kg"`
	GripKg     float64 `json:"grip_kg"`
	ReactionMs int     `json:"reaction_ms"`
	Status     string  `json:"status"`
}

func MaskChild(profile domain.ChildProfile, policy Policy) DisplayChild {
	display := DisplayChild{ID: profile.ID, Name: profile.Name, BirthYear: profile.BirthYear}
	if !policy.ShowChildIDs || policy.MaskIdentifiers {
		display.ID = MaskIdentifier(profile.ID)
	}
	if !policy.ShowNames {
		display.Name = MaskName(profile.Name)
	}
	if !policy.ShowBirthYears {
		display.BirthYear = 0
	}
	return display
}

func MaskRecord(record domain.MeasurementRecord, policy Policy) DisplayRecord {
	id := record.ChildID
	if policy.MaskIdentifiers || !policy.ShowChildIDs {
		id = MaskIdentifier(id)
	}
	return DisplayRecord{ChildID: id, HeightCm: record.HeightCm, WeightKg: record.WeightKg, GripKg: record.GripKg, ReactionMs: record.ReactionMs, Status: record.Status}
}

func MaskIdentifier(identifier string) string {
	clean := strings.TrimSpace(identifier)
	if clean == "" {
		return "***"
	}
	if len(clean) <= 2 {
		return "**"
	}
	return clean[:1] + strings.Repeat("*", len(clean)-2) + clean[len(clean)-1:]
}

func MaskName(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "anonymous"
	}
	masked := make([]string, len(parts))
	for i, part := range parts {
		masked[i] = fmt.Sprintf("%c*", []rune(part)[0])
	}
	return strings.Join(masked, " ")
}
