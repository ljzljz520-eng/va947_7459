package privacy

import (
	"sort"
	"strings"

	"example.com/childfitness/internal/domain"
)

type ViewOptions struct {
	Audience       Audience
	IncludeMetrics bool
	SortBy         string
}

type ClassView struct {
	SchoolID string          `json:"school_id"`
	ClassID  string          `json:"class_id"`
	Audience string          `json:"audience"`
	Children []DisplayChild  `json:"children"`
	Records  []DisplayRecord `json:"records"`
	Denied   int             `json:"denied"`
}

func BuildClassView(schoolID, classID string, profiles []domain.ChildProfile, records []domain.MeasurementRecord, options ViewOptions) ClassView {
	policy := PolicyFor(options.Audience)
	view := ClassView{SchoolID: schoolID, ClassID: classID, Audience: string(options.Audience), Children: make([]DisplayChild, 0), Records: make([]DisplayRecord, 0)}
	allowed := make(map[string]bool, len(profiles))
	for _, profile := range profiles {
		if CanViewClass(policy, schoolID, classID, profile) {
			view.Children = append(view.Children, MaskChild(profile, policy))
			allowed[profile.ID] = true
		} else {
			view.Denied++
		}
	}
	if options.IncludeMetrics {
		for _, record := range records {
			if allowed[record.ChildID] {
				view.Records = append(view.Records, MaskRecord(record, policy))
			}
		}
	}
	SortClassView(&view, options.SortBy)
	return view
}

func SortClassView(view *ClassView, sortBy string) {
	if view == nil {
		return
	}
	if strings.EqualFold(sortBy, "name") {
		sort.SliceStable(view.Children, func(i, j int) bool { return view.Children[i].Name < view.Children[j].Name })
	} else {
		sort.SliceStable(view.Children, func(i, j int) bool { return view.Children[i].ID < view.Children[j].ID })
	}
	sort.SliceStable(view.Records, func(i, j int) bool { return view.Records[i].ChildID < view.Records[j].ChildID })
}

func VisibleChildCount(view ClassView) int {
	return len(view.Children)
}

func ViewHasDenied(view ClassView) bool {
	return view.Denied > 0
}

func AudienceLabel(audience Audience) string {
	switch audience {
	case AudienceAdmin:
		return "Administrator"
	case AudienceParent:
		return "Parent"
	default:
		return "Teacher"
	}
}
