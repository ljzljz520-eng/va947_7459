package domain

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateProfile(p ChildProfile) error {
	if strings.TrimSpace(p.ID) == "" {
		return errors.New("child id is required")
	}
	if strings.TrimSpace(p.SchoolID) == "" || strings.TrimSpace(p.ClassID) == "" {
		return errors.New("school and class are required")
	}
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("child name is required")
	}
	if p.BirthYear < 2000 || p.BirthYear > 2030 {
		return fmt.Errorf("birth year %d is outside supported range", p.BirthYear)
	}
	if !p.Consent {
		return errors.New("consent is required")
	}
	return nil
}

func ValidateMeasurement(r MeasurementRecord) error {
	if strings.TrimSpace(r.ChildID) == "" || strings.TrimSpace(r.BatchID) == "" {
		return errors.New("measurement identifiers are required")
	}
	if r.HeightCm < 50 || r.HeightCm > 230 {
		return fmt.Errorf("height %.2f is outside supported range", r.HeightCm)
	}
	if r.WeightKg < 5 || r.WeightKg > 250 {
		return fmt.Errorf("weight %.2f is outside supported range", r.WeightKg)
	}
	if r.GripKg < 0 || r.GripKg > 150 {
		return fmt.Errorf("grip %.2f is outside supported range", r.GripKg)
	}
	if r.ReactionMs < 50 || r.ReactionMs > 10000 {
		return fmt.Errorf("reaction %d is outside supported range", r.ReactionMs)
	}
	return nil
}

func ValidateBatch(b FitnessBatch) error {
	if strings.TrimSpace(b.ID) == "" {
		return errors.New("batch id is required")
	}
	if strings.TrimSpace(b.SchoolID) == "" || strings.TrimSpace(b.ClassID) == "" {
		return errors.New("batch school and class are required")
	}
	if len(b.Records) == 0 && len(b.Issues) == 0 {
		return errors.New("batch has no observations")
	}
	if b.Status != BatchPending && b.Status != BatchComplete && b.Status != BatchPartial {
		return fmt.Errorf("unsupported batch status %q", b.Status)
	}
	return nil
}

func NormalizeName(name string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(name)), " ")
}

func NormalizeClassID(classID string) string {
	return strings.ToUpper(strings.TrimSpace(classID))
}
