package capture

import "testing"

func TestParserReportsMalformedRowAndParsesValidRow(t *testing.T) {
	input := "child_id,height,weight,grip,reaction,received_at\nC1,130,30,15,400,2026-01-01\nC2,130,30\n"
	rows, issues, err := NewParser().ReadRows(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || len(issues) != 1 {
		t.Fatalf("rows=%d issues=%d", len(rows), len(issues))
	}
	record, issue := ParseRow(rows[0], "B1")
	if issue != nil || record.Status != "valid" || record.HeightCm != 130 {
		t.Fatalf("unexpected parse result %#v %#v", record, issue)
	}
}
