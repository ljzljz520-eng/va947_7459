package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplicationDemoCommand(t *testing.T) {
	app, err := NewApplication(filepath.Join(t.TempDir(), "cli.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	var output bytes.Buffer
	if err := app.Execute([]string{"demo"}, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "valid=2") {
		t.Fatalf("unexpected demo output %q", output.String())
	}
}
