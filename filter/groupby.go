package filter

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// GroupByResult holds aggregated records grouped by a field value.
type GroupByResult struct {
	Key    string
	Field  string
	Counts map[string]int
}

// NewGroupByResult creates a new GroupByResult for the given field.
func NewGroupByResult(field string) *GroupByResult {
	return &GroupByResult{
		Field:  field,
		Counts: make(map[string]int),
	}
}

// Add records a value for the group-by field from the given record.
func (g *GroupByResult) Add(record map[string]interface{}) {
	val, ok := record[g.Field]
	if !ok {
		g.Counts["(missing)"]++
		return
	}
	g.Counts[fmt.Sprintf("%v", val)]++
}

// Sorted returns the group keys sorted by count descending.
func (g *GroupByResult) Sorted() []string {
	keys := make([]string, 0, len(g.Counts))
	for k := range g.Counts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if g.Counts[keys[i]] != g.Counts[keys[j]] {
			return g.Counts[keys[i]] > g.Counts[keys[j]]
		}
		return keys[i] < keys[j]
	})
	return keys
}

// RunGroupBy reads JSON lines from r, groups by the given field, and writes a
// summary table to w.
func RunGroupBy(r io.Reader, w io.Writer, field string) error {
	if strings.TrimSpace(field) == "" {
		return fmt.Errorf("groupby: field name must not be empty")
	}
	result := NewGroupByResult(field)
	scanner := NewScanner(r)
	for scanner.Scan() {
		record, err := ParseRecord(scanner.Bytes())
		if err != nil {
			continue
		}
		result.Add(record)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("groupby: scanner error: %w", err)
	}
	writeGroupByTable(w, result)
	return nil
}

func writeGroupByTable(w io.Writer, result *GroupByResult) {
	fmt.Fprintf(w, "%-40s %s\n", result.Field, "count")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 50))
	for _, k := range result.Sorted() {
		fmt.Fprintf(w, "%-40s %d\n", k, result.Counts[k])
	}
}
