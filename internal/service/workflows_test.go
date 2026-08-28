package service

import (
	"path/filepath"
	"testing"

	"example.com/childfitness/internal/domain"
	"example.com/childfitness/internal/privacy"
	"example.com/childfitness/internal/store"
)

func openWorkflowStore(t *testing.T) (*store.Store, *EnrollmentService, *IntakeService, *ReviewService) {
	t.Helper()
	storage, err := store.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	return storage, NewEnrollmentService(storage), NewIntakeService(storage), NewReviewService(storage)
}

func TestWorkflowEnrollmentAndCapture(t *testing.T) {
	storage, enrollment, intake, review := openWorkflowStore(t)
	defer storage.Close()
	for _, profile := range []domain.ChildProfile{{ID: "C1", SchoolID: "S1", ClassID: "5A", Name: "One", BirthYear: 2014, Consent: true}, {ID: "C2", SchoolID: "S1", ClassID: "5A", Name: "Two", BirthYear: 2014, Consent: true}} {
		if err := enrollment.RegisterChild(profile); err != nil {
			t.Fatal(err)
		}
	}
	input := "child_id,height,weight,grip,reaction,received_at\nC1,130,30,15,400,2026-01-01\nC2,131,31,16,410,2026-01-01\n"
	if _, issues, err := intake.ImportDeviceBatch("S1", "5A", "device", input); err != nil || len(issues) != 0 {
		t.Fatalf("import failed issues=%v err=%v", issues, err)
	}
	summary, err := review.ReviewClass("S1", "5A")
	if err != nil || summary.ValidRecords != 2 {
		t.Fatalf("review failed %#v %v", summary, err)
	}
}

func TestWorkflowClassReviewAndPrivacy(t *testing.T) {
	storage, enrollment, intake, review := openWorkflowStore(t)
	defer storage.Close()
	if err := enrollment.RegisterChild(domain.ChildProfile{ID: "C1", SchoolID: "S1", ClassID: "5B", Name: "Private Child", BirthYear: 2013, Consent: true}); err != nil {
		t.Fatal(err)
	}
	input := "child_id,height,weight,grip,reaction,received_at\nC1,140,35,20,390,2026-01-01\n"
	if _, _, err := intake.ImportDeviceBatch("S1", "5B", "device", input); err != nil {
		t.Fatal(err)
	}
	children, err := review.MaskedRoster("S1", "5B", privacy.AudienceTeacher)
	if err != nil || len(children) != 1 || children[0].Name == "Private Child" {
		t.Fatalf("privacy review failed %#v %v", children, err)
	}
}

func TestWorkflowImportAndRecovery(t *testing.T) {
	path := filepath.Join(t.TempDir(), "recovery.db")
	storage, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	intake := NewIntakeService(storage)
	if _, _, err := intake.ImportDeviceBatch("S2", "6A", "device", "child_id,height,weight,grip,reaction,received_at\nC9,145,40,22,360,2026-02-02\n"); err != nil {
		t.Fatal(err)
	}
	if err := storage.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if _, err := NewIntakeService(reopened).LoadBatch("s2-6a-device"); err != nil {
		t.Fatal(err)
	}
}

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "persist.db")
	storage, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.SaveProfile(domain.ChildProfile{ID: "C7", SchoolID: "S7", ClassID: "7A", Name: "Persisted", BirthYear: 2012, Consent: true}); err != nil {
		t.Fatal(err)
	}
	if err := storage.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	profile, err := reopened.GetProfile("S7", "7A", "C7")
	if err != nil || profile.Name != "Persisted" {
		t.Fatalf("reopen failed %#v %v", profile, err)
	}
}
