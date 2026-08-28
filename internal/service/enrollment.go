package service

import (
	"errors"
	"strings"

	"example.com/childfitness/internal/domain"
)

type EnrollmentService struct {
	repo Repository
}

func NewEnrollmentService(repo Repository) *EnrollmentService {
	return &EnrollmentService{repo: repo}
}

func (s *EnrollmentService) RegisterChild(profile domain.ChildProfile) error {
	profile.Name = domain.NormalizeName(profile.Name)
	profile.ClassID = domain.NormalizeClassID(profile.ClassID)
	if err := domain.ValidateProfile(profile); err != nil {
		return err
	}
	return s.repo.SaveProfile(profile)
}

func (s *EnrollmentService) FindChild(schoolID, classID, childID string) (domain.ChildProfile, error) {
	if strings.TrimSpace(schoolID) == "" || strings.TrimSpace(classID) == "" || strings.TrimSpace(childID) == "" {
		return domain.ChildProfile{}, errors.New("school, class, and child are required")
	}
	return s.repo.GetProfile(schoolID, domain.NormalizeClassID(classID), childID)
}

func (s *EnrollmentService) ListClass(schoolID, classID string) ([]domain.ChildProfile, error) {
	if strings.TrimSpace(schoolID) == "" || strings.TrimSpace(classID) == "" {
		return nil, errors.New("school and class are required")
	}
	return s.repo.ListProfiles(schoolID, domain.NormalizeClassID(classID))
}
