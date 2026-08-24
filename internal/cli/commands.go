package cli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"

	"example.com/childfitness/internal/domain"
	"example.com/childfitness/internal/privacy"
	"example.com/childfitness/internal/report"
)

func (a *Application) runDemo(output interface{ Write([]byte) (int, error) }) error {
	profiles := []domain.ChildProfile{{ID: "C001", SchoolID: "S001", ClassID: "5A", Name: "Lin Mei", BirthYear: 2014, Consent: true}, {ID: "C002", SchoolID: "S001", ClassID: "5A", Name: "Wang Jun", BirthYear: 2014, Consent: true}}
	for _, profile := range profiles {
		if err := a.Enrollment.RegisterChild(profile); err != nil {
			return err
		}
	}
	input := "child_id,height,weight,grip,reaction,received_at\nC001,132.5,31.2,18.0,420,2026-01-05T09:00:00Z\nC002,130.0,29.4,17.5,450,2026-01-05T09:00:01Z\n"
	batch, issues, err := a.Intake.ImportDeviceBatch("S001", "5A", "demo", input)
	if err != nil {
		return err
	}
	summary, err := a.Review.ReviewClass("S001", "5A")
	if err != nil {
		return err
	}
	masked, err := a.Review.MaskedRoster("S001", "5A", privacy.AudienceTeacher)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(output, "Imported %s; issues=%d; children=%d; valid=%d\n", batch.ID, len(issues), len(masked), summary.ValidRecords)
	return err
}

func (a *Application) runImport(args []string, output interface{ Write([]byte) (int, error) }) error {
	flags := flag.NewFlagSet("import", flag.ContinueOnError)
	flags.SetOutput(bytes.NewBuffer(nil))
	school := flags.String("school", "", "school id")
	classID := flags.String("class", "", "class id")
	source := flags.String("source", "device", "source name")
	file := flags.String("file", "", "csv file")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *file == "" {
		return errors.New("-file is required")
	}
	input, err := os.ReadFile(*file)
	if err != nil {
		return err
	}
	batch, issues, err := a.Intake.ImportDeviceBatch(*school, *classID, *source, string(input))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(output, "%s imported=%d issues=%d status=%s\n", batch.ID, batch.ImportedCount, len(issues), batch.Status)
	return err
}

func (a *Application) runReview(args []string, output interface{ Write([]byte) (int, error) }) error {
	flags := flag.NewFlagSet("review", flag.ContinueOnError)
	flags.SetOutput(bytes.NewBuffer(nil))
	school := flags.String("school", "", "school id")
	classID := flags.String("class", "", "class id")
	audienceValue := flags.String("audience", "teacher", "audience")
	format := flags.String("format", "text", "text or json")
	if err := flags.Parse(args); err != nil {
		return err
	}
	audience, err := privacy.ValidateAudience(*audienceValue)
	if err != nil {
		return err
	}
	profiles, err := a.Enrollment.ListClass(*school, *classID)
	if err != nil {
		return err
	}
	records, err := a.Store.RecordsForClass(*school, *classID)
	if err != nil {
		return err
	}
	summary, err := a.Review.ReviewClass(*school, *classID)
	if err != nil {
		return err
	}
	children, err := a.Review.MaskedRoster(*school, *classID, audience)
	if err != nil {
		return err
	}
	_ = children
	classReport := report.BuildReport(summary, profiles, records, nil)
	if *format == "json" {
		data, err := report.JSON(classReport)
		if err != nil {
			return err
		}
		_, err = output.Write(append(data, '\n'))
		return err
	}
	_, err = fmt.Fprint(output, report.Text(classReport))
	return err
}
