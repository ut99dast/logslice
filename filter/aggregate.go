package filter

import (
	"fmt"
	"sort"
)

// AggregateResult holds counted occurrences of a field value.
type AggregateResult struct {
	Field  string
	Counts map[string]int
}

// NewAggregateResult creates an AggregateResult for the given field.
func NewAggregateResult(field string) *AggregateResult {
	return &AggregateResult{
		Field:  field,
		Counts: make(map[string]int),
	}
}

// Add records the value of the tracked field from a log record.
func (a *AggregateResult) Add(record map[string]interface{}) {
	val, ok := record[a.Field]
	if !ok {
		a.Counts["<missing>"]++
		return
	}
	a.Counts[fmt.Sprintf("%v", val)]++
}

// Sorted returns field values sorted by count descending.
func (a *AggregateResult) Sorted() []AggregateEntry {
	entries := make([]AggregateEntry, 0, len(a.Counts))
	for k, v := range a.Counts {
		entries = append(entries, AggregateEntry{Value: k, Count: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Value < entries[j].Value
	})
	return entries
}

// AggregateEntry is a single value/count pair.
type AggregateEntry struct {
	Value string
	Count int
}
