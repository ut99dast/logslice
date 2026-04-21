package filter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// Counter counts occurrences of unique values for a given field.
type Counter struct {
	field  string
	counts map[string]int
}

// NewCounter creates a Counter for the specified field.
func NewCounter(field string) (*Counter, error) {
	if field == "" {
		return nil, fmt.Errorf("count: field name must not be empty")
	}
	return &Counter{
		field:  field,
		counts: make(map[string]int),
	}, nil
}

// Add records the value of the tracked field from a parsed record.
func (c *Counter) Add(record map[string]interface{}) {
	v, ok := record[c.field]
	if !ok {
		return
	}
	key := fmt.Sprintf("%v", v)
	c.counts[key]++
}

// Results returns a sorted slice of [value, count] pairs descending by count.
func (c *Counter) Results() []CountEntry {
	entries := make([]CountEntry, 0, len(c.counts))
	for k, n := range c.counts {
		entries = append(entries, CountEntry{Value: k, Count: n})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Value < entries[j].Value
	})
	return entries
}

// CountEntry holds a single value/count pair.
type CountEntry struct {
	Value string
	Count int
}

// RunCount reads JSON lines from r, counts occurrences of field values,
// and writes a JSON-lines summary to w.
func RunCount(r io.Reader, w io.Writer, field string) error {
	counter, err := NewCounter(field)
	if err != nil {
		return err
	}

	scanner := NewScanner(r)
	for scanner.Scan() {
		record, err := ParseRecord(scanner.Bytes())
		if err != nil {
			continue
		}
		counter.Add(record)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("count: scanner error: %w", err)
	}

	for _, entry := range counter.Results() {
		out := map[string]interface{}{
			field:   entry.Value,
			"count": entry.Count,
		}
		line, err := json.Marshal(out)
		if err != nil {
			return fmt.Errorf("count: marshal error: %w", err)
		}
		fmt.Fprintln(w, string(line))
	}
	return nil
}
