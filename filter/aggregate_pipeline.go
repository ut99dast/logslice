package filter

import (
	"bufio"
	"fmt"
	"io"
)

// RunAggregate reads JSON log lines from r, aggregates counts for the given
// field among lines that match f, and writes a summary table to w.
// Lines that cannot be parsed as JSON are silently skipped.
// If f is nil, all parseable lines are included in the aggregation.
func RunAggregate(r io.Reader, w io.Writer, field string, f Filter) error {
	result := NewAggregateResult(field)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		record, err := ParseRecord(line)
		if err != nil {
			continue
		}
		if f != nil && !f.Match(record) {
			continue
		}
		result.Add(record)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	entries := result.Sorted()
	if err := writeAggregateTable(w, field, entries); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}

// writeAggregateTable formats and writes the aggregation results as a table.
func writeAggregateTable(w io.Writer, field string, entries []AggregateEntry) error {
	if _, err := fmt.Fprintf(w, "%-40s %s\n", field, "count"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "%-40s %s\n", "----------------------------------------", "-----"); err != nil {
		return err
	}
	for _, e := range entries {
		if _, err := fmt.Fprintf(w, "%-40s %d\n", e.Value, e.Count); err != nil {
			return err
		}
	}
	return nil
}

// Filter is the interface satisfied by all filter types in this package.
type Filter interface {
	Match(record map[string]interface{}) bool
}
