package store

import (
	"path/filepath"
	"testing"

	"example.com/childfitness/internal/domain"
)

func TestStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "store.db")
	storage, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	profile := domain.ChildProfile{ID: "C1", SchoolID: "S1", ClassID: "5A", Name: "Child", BirthYear: 2014, Consent: true}
	if err := storage.SaveProfile(profile); err != nil {
		t.Fatal(err)
	}
	batch := domain.FitnessBatch{ID: "B1", SchoolID: "S1", ClassID: "5A", Status: domain.BatchComplete, Records: []domain.MeasurementRecord{{ID: "R1", BatchID: "B1", ChildID: "C1", HeightCm: 130, WeightKg: 30, GripKg: 15, ReactionMs: 400, Status: domain.RecordValid}}}
	if err := storage.SaveBatch(batch); err != nil {
		t.Fatal(err)
	}
	if err := storage.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	got, err := reopened.GetBatch("B1")
	if err != nil || len(got.Records) != 1 {
		t.Fatalf("batch round trip failed %#v %v", got, err)
	}
}
