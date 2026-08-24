package domain

import "testing"

func TestEntityKeysAreStable(t *testing.T) {
	profile := ChildProfile{ID: "C1", SchoolID: "S1", ClassID: "5A"}
	if profile.Key() != "S1:5A:C1" {
		t.Fatalf("unexpected profile key %q", profile.Key())
	}
	record := MeasurementRecord{BatchID: "B1", ChildID: "C1"}
	if record.Key() != "B1:C1" {
		t.Fatalf("unexpected record key %q", record.Key())
	}
	if !(ClassSummary{ValidRecords: 1}).HasData() {
		t.Fatal("summary should have data")
	}
}
