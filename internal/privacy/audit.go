package privacy

import (
	"fmt"
	"sort"
	"strings"

	"example.com/childfitness/internal/domain"
)

type AuditEvent struct {
	Audience string `json:"audience"`
	SchoolID string `json:"school_id"`
	ClassID  string `json:"class_id"`
	ChildID  string `json:"child_id"`
	Action   string `json:"action"`
	Outcome  string `json:"outcome"`
}

type AuditTrail struct {
	Events []AuditEvent `json:"events"`
}

func (a *AuditTrail) Record(policy Policy, schoolID, classID string, profile domain.ChildProfile, outcome string) {
	if a == nil {
		return
	}
	a.Events = append(a.Events, AuditEvent{Audience: string(policy.Audience), SchoolID: schoolID, ClassID: classID, ChildID: MaskIdentifier(profile.ID), Action: "view-child", Outcome: outcome})
}

func (a AuditTrail) ForAudience(audience Audience) []AuditEvent {
	filtered := make([]AuditEvent, 0)
	for _, event := range a.Events {
		if Audience(event.Audience) == audience {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func (a AuditTrail) Summary() string {
	counts := make(map[string]int)
	for _, event := range a.Events {
		counts[event.Outcome]++
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", key, counts[key]))
	}
	return strings.Join(parts, ",")
}

func MaskingReason(policy Policy, profile domain.ChildProfile) string {
	if policy.Audience == AudienceAdmin {
		return "administrator-authorized"
	}
	if policy.Audience == AudienceParent && profile.Consent {
		return "consented-parent-view"
	}
	return "classroom-minimum-identity"
}

func IsSensitiveField(field string) bool {
	switch strings.ToLower(strings.TrimSpace(field)) {
	case "name", "child_id", "birth_year", "address", "phone":
		return true
	default:
		return false
	}
}
