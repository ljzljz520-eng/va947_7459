package domain

import "testing"

func TestValidationRejectsUnsafeMeasurements(t *testing.T) {
	record := MeasurementRecord{BatchID: "B", ChildID: "C", HeightCm: 20, WeightKg: 30, GripKg: 10, ReactionMs: 300}
	if err := ValidateMeasurement(record); err == nil {
		t.Fatal("expected height validation error")
	}
	profile := ChildProfile{ID: "C", SchoolID: "S", ClassID: "5A", Name: "Child", BirthYear: 2014, Consent: false}
	if err := ValidateProfile(profile); err == nil {
		t.Fatal("expected consent validation error")
	}
	if NormalizeClassID(" 5a ") != "5A" || NormalizeName(" A   B ") != "A B" {
		t.Fatal("normalization failed")
	}
}
