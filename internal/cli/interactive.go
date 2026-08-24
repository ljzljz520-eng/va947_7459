package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"example.com/childfitness/internal/domain"
)

type Prompt struct {
	Reader io.Reader
	Writer io.Writer
}

type EnrollmentInput struct {
	SchoolID  string
	ClassID   string
	ChildID   string
	Name      string
	BirthYear int
	Consent   bool
}

func NewPrompt(reader io.Reader, writer io.Writer) Prompt {
	return Prompt{Reader: reader, Writer: writer}
}

func (p Prompt) ReadLine(label string) (string, error) {
	if p.Reader == nil || p.Writer == nil {
		return "", errors.New("prompt streams are required")
	}
	if _, err := fmt.Fprint(p.Writer, label); err != nil {
		return "", err
	}
	line, err := bufio.NewReader(p.Reader).ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func ParseConsent(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "y", "yes", "true", "1", "agree":
		return true
	default:
		return false
	}
}

func (p Prompt) ReadEnrollment() (EnrollmentInput, error) {
	school, err := p.ReadLine("school: ")
	if err != nil {
		return EnrollmentInput{}, err
	}
	classID, err := p.ReadLine("class: ")
	if err != nil {
		return EnrollmentInput{}, err
	}
	childID, err := p.ReadLine("child id: ")
	if err != nil {
		return EnrollmentInput{}, err
	}
	name, err := p.ReadLine("name: ")
	if err != nil {
		return EnrollmentInput{}, err
	}
	return EnrollmentInput{SchoolID: school, ClassID: classID, ChildID: childID, Name: name, BirthYear: 2014, Consent: true}, nil
}

func (p Prompt) SaveEnrollment(input EnrollmentInput) error {
	if input.BirthYear == 0 {
		input.BirthYear = 2014
	}
	profile := domain.ChildProfile{ID: input.ChildID, SchoolID: input.SchoolID, ClassID: input.ClassID, Name: input.Name, BirthYear: input.BirthYear, Consent: input.Consent}
	return p.validateProfile(profile)
}

func (p Prompt) validateProfile(profile domain.ChildProfile) error {
	return domain.ValidateProfile(profile)
}
