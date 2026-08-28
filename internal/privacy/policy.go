package privacy

import (
	"errors"
	"strings"

	"example.com/childfitness/internal/domain"
)

type Audience string

const (
	AudienceTeacher Audience = "teacher"
	AudienceParent  Audience = "parent"
	AudienceAdmin   Audience = "admin"
)

type Policy struct {
	Audience        Audience
	ShowNames       bool
	ShowChildIDs    bool
	ShowBirthYears  bool
	MaskIdentifiers bool
}

func PolicyFor(audience Audience) Policy {
	switch audience {
	case AudienceAdmin:
		return Policy{Audience: audience, ShowNames: true, ShowChildIDs: true, ShowBirthYears: true}
	case AudienceParent:
		return Policy{Audience: audience, ShowNames: true, ShowChildIDs: false, ShowBirthYears: false, MaskIdentifiers: true}
	default:
		return Policy{Audience: AudienceTeacher, ShowNames: false, ShowChildIDs: false, ShowBirthYears: false, MaskIdentifiers: true}
	}
}

func ValidateAudience(value string) (Audience, error) {
	clean := Audience(strings.ToLower(strings.TrimSpace(value)))
	if clean != AudienceTeacher && clean != AudienceParent && clean != AudienceAdmin {
		return "", errors.New("audience must be teacher, parent, or admin")
	}
	return clean, nil
}

func CanViewClass(policy Policy, schoolID, classID string, profile domain.ChildProfile) bool {
	if schoolID == "" || classID == "" {
		return false
	}
	if profile.SchoolID != schoolID || profile.ClassID != classID {
		return false
	}
	if policy.Audience == AudienceParent && !profile.Consent {
		return false
	}
	return true
}
