package cli

import (
	"errors"
	"fmt"
	"io"

	"example.com/childfitness/internal/service"
	"example.com/childfitness/internal/store"
)

type Application struct {
	Store      *store.Store
	Enrollment *service.EnrollmentService
	Intake     *service.IntakeService
	Review     *service.ReviewService
}

func NewApplication(databasePath string) (*Application, error) {
	storage, err := store.Open(databasePath)
	if err != nil {
		return nil, err
	}
	return &Application{Store: storage, Enrollment: service.NewEnrollmentService(storage), Intake: service.NewIntakeService(storage), Review: service.NewReviewService(storage)}, nil
}

func (a *Application) Close() error {
	if a == nil || a.Store == nil {
		return nil
	}
	return a.Store.Close()
}

func (a *Application) Execute(args []string, output io.Writer) error {
	if a == nil {
		return errors.New("application is nil")
	}
	if len(args) == 0 {
		return a.writeUsage(output)
	}
	switch args[0] {
	case "demo":
		return a.runDemo(output)
	case "import":
		return a.runImport(args[1:], output)
	case "review":
		return a.runReview(args[1:], output)
	case "help":
		return a.writeUsage(output)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func (a *Application) writeUsage(output io.Writer) error {
	_, err := fmt.Fprintln(output, "childfitness commands: demo | import | review")
	return err
}
