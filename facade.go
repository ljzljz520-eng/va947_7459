package childfitness

import (
	"example.com/childfitness/internal/cli"
	"example.com/childfitness/internal/domain"
	"example.com/childfitness/internal/privacy"
	"example.com/childfitness/internal/report"
	"example.com/childfitness/internal/service"
	"example.com/childfitness/internal/store"
)

type Project struct {
	app *cli.Application
}

func Open(path string) (*Project, error) {
	app, err := cli.NewApplication(path)
	if err != nil {
		return nil, err
	}
	return &Project{app: app}, nil
}

func (p *Project) Close() error {
	if p == nil || p.app == nil {
		return nil
	}
	return p.app.Close()
}

func (p *Project) Register(profile domain.ChildProfile) error {
	return p.app.Enrollment.RegisterChild(profile)
}

func (p *Project) Import(schoolID, classID, source, input string) (domain.FitnessBatch, []domain.ImportIssue, error) {
	return p.app.Intake.ImportDeviceBatch(schoolID, classID, source, input)
}

func (p *Project) Review(schoolID, classID string) (domain.ClassSummary, error) {
	return p.app.Review.ReviewClass(schoolID, classID)
}

func (p *Project) ClassReport(schoolID, classID string) (report.ClassReport, error) {
	context, err := p.app.Review.AssembleClassContext(schoolID, classID)
	if err != nil {
		return report.ClassReport{}, err
	}
	return service.BuildClassReport(context), nil
}

func (p *Project) RepairBatch(batchID string) (domain.FitnessBatch, error) {
	batch, _, err := p.app.Intake.RepairBatch(batchID)
	return batch, err
}

func (p *Project) AuditEvents() []privacy.AuditEvent {
	return p.app.Review.AuditEvents()
}

func (p *Project) Stats() (store.DatabaseStats, error) {
	return p.app.Store.Stats()
}

func (p *Project) Batch(batchID string) (domain.FitnessBatch, error) {
	return p.app.Intake.LoadBatch(batchID)
}

func (p *Project) Healthy() bool {
	return p != nil && p.app != nil && p.app.Store.Healthy()
}
