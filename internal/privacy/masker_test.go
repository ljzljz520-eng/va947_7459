package privacy

import (
	"testing"

	"example.com/childfitness/internal/domain"
)

func TestTeacherViewMasksIdentity(t *testing.T) {
	policy := PolicyFor(AudienceTeacher)
	child := MaskChild(domain.ChildProfile{ID: "C001", Name: "Lin Mei", BirthYear: 2014}, policy)
	if child.ID != "C**1" || child.Name != "L* M*" || child.BirthYear != 0 {
		t.Fatalf("unexpected masked child %#v", child)
	}
	if MaskIdentifier("") != "***" || MaskName("") != "anonymous" {
		t.Fatal("empty values were not masked")
	}
}
