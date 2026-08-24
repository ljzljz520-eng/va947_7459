package childfitness

import (
	"path/filepath"
	"testing"

	"example.com/childfitness/internal/domain"
)

func TestChildFitnessImportIsolatesBadRow(t *testing.T) {
	project, err := Open(filepath.Join(t.TempDir(), "fitness.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer project.Close()
	if err := project.Register(domain.ChildProfile{ID: "C001", SchoolID: "S1", ClassID: "5A", Name: "A Child", BirthYear: 2014, Consent: true}); err != nil {
		t.Fatal(err)
	}
	input := "child_id,height,weight,grip,reaction,received_at\nC001,132.0,30.0,18.0,400,2026-01-01T00:00:00Z\nC001,133.0,31.0,not-a-number,410,2026-01-01T00:00:01Z\n"
	batch, issues, err := project.Import("S1", "5A", "device", input)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected one issue, got %d", len(issues))
	}
	if batch.Status != domain.BatchPartial {
		t.Fatalf("expected partial batch, got %s", batch.Status)
	}
	if len(batch.Records) != 1 || batch.Records[0].GripKg != 18 {
		t.Fatalf("valid result was not isolated: %#v", batch.Records)
	}
}
