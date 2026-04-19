package filter

import (
	"bufio"
	"fmt"
	"io"
)

// RunAggregate reads JSON log lines from r, aggregates counts for the given
// field among lines that match f, and writes a summary table to w.
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
	fmt.Fprintf(w, "%-40s %s\n", field, "count")
	fmt.Fprintf(w, "%-40s %s\n", "----------------------------------------", "-----")
	for _, e := range entries {
		fmt.Fprintf(w, "%-40s %d\n", e.Value, e.Count)
	}
	return nil
}

// Filter is the interface satisfied by all filter types in this package.
type Filter interface {
	Match(record map[string]interface{}) bool
}
