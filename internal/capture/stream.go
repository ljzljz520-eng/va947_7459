package capture

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"example.com/childfitness/internal/domain"
)

type StreamOptions struct {
	MaxRows       int
	RequireHeader bool
	KeepEmpty     bool
}

type StreamResult struct {
	RowsRead      int
	RowsDelivered int
	RowsSkipped   int
	Issues        []domain.ImportIssue
}

type LineStream struct {
	options StreamOptions
}

func NewLineStream(options StreamOptions) LineStream {
	if options.MaxRows < 0 {
		options.MaxRows = 0
	}
	return LineStream{options: options}
}

func (s LineStream) Read(input string, handler func(int, string) error) (StreamResult, error) {
	result := StreamResult{Issues: make([]domain.ImportIssue, 0)}
	reader := bufio.NewScanner(strings.NewReader(input))
	lineNumber := 0
	for reader.Scan() {
		lineNumber++
		result.RowsRead++
		line := strings.TrimSpace(reader.Text())
		if lineNumber == 1 && s.options.RequireHeader {
			if !strings.HasPrefix(strings.ToLower(line), "child_id") {
				return result, fmt.Errorf("line one must contain child_id header")
			}
			continue
		}
		if line == "" && !s.options.KeepEmpty {
			result.RowsSkipped++
			continue
		}
		if s.options.MaxRows > 0 && result.RowsDelivered >= s.options.MaxRows {
			result.RowsSkipped++
			continue
		}
		if err := handler(lineNumber, line); err != nil {
			result.Issues = append(result.Issues, domain.ImportIssue{Row: lineNumber, Field: "stream", Message: err.Error()})
			continue
		}
		result.RowsDelivered++
	}
	if err := reader.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func ReadNonEmptyLines(input string) []string {
	lines := make([]string, 0)
	stream := NewLineStream(StreamOptions{KeepEmpty: false})
	_, _ = stream.Read(input, func(_ int, line string) error {
		lines = append(lines, line)
		return nil
	})
	return lines
}

func WriteRows(writer io.Writer, rows []DeviceRow) error {
	if _, err := io.WriteString(writer, "child_id,height,weight,grip,reaction,received_at\n"); err != nil {
		return err
	}
	for _, row := range rows {
		line := strings.Join([]string{row.ChildID, row.Height, row.Weight, row.Grip, row.Reaction, row.ReceivedAt}, ",") + "\n"
		if _, err := io.WriteString(writer, line); err != nil {
			return err
		}
	}
	return nil
}

func GroupRowsByChild(rows []DeviceRow) map[string][]DeviceRow {
	groups := make(map[string][]DeviceRow)
	for _, row := range rows {
		groups[row.ChildID] = append(groups[row.ChildID], row)
	}
	return groups
}
