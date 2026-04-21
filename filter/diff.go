package filter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Differ compares two JSON log records and reports added, removed, or changed fields.
type Differ struct {
	ignoreFields map[string]bool
}

// DiffResult holds the comparison between two records.
type DiffResult struct {
	Added   map[string]interface{}
	Removed map[string]interface{}
	Changed map[string][2]interface{} // [before, after]
}

// NewDiffer creates a Differ that optionally ignores certain fields.
// ignoreFields is a comma-separated list of field names to skip during comparison.
func NewDiffer(ignoreFields string) (*Differ, error) {
	ignore := map[string]bool{}
	if ignoreFields != "" {
		for _, f := range strings.Split(ignoreFields, ",") {
			f = strings.TrimSpace(f)
			if f == "" {
				return nil, fmt.Errorf("diff: empty field name in ignore list")
			}
			ignore[f] = true
		}
	}
	return &Differ{ignoreFields: ignore}, nil
}

// Compare returns a DiffResult describing differences between before and after.
func (d *Differ) Compare(before, after map[string]interface{}) DiffResult {
	result := DiffResult{
		Added:   map[string]interface{}{},
		Removed: map[string]interface{}{},
		Changed: map[string][2]interface{}{},
	}

	for k, v := range after {
		if d.ignoreFields[k] {
			continue
		}
		if bv, ok := before[k]; !ok {
			result.Added[k] = v
		} else if fmt.Sprintf("%v", bv) != fmt.Sprintf("%v", v) {
			result.Changed[k] = [2]interface{}{bv, v}
		}
	}

	for k, v := range before {
		if d.ignoreFields[k] {
			continue
		}
		if _, ok := after[k]; !ok {
			result.Removed[k] = v
		}
	}

	return result
}

// HasChanges returns true if the DiffResult contains any differences.
func (r DiffResult) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0
}

// RunDiff reads pairs of consecutive lines from in, computes their diff,
// and writes a JSON summary of changes to out. Lines that are not valid
// JSON are skipped and counted. Only pairs where at least one change
// exists are written.
func RunDiff(in io.Reader, out io.Writer, ignoreFields string) error {
	differ, err := NewDiffer(ignoreFields)
	if err != nil {
		return err
	}

	scanner := NewScanner(in)
	var prev map[string]interface{}
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		var rec map[string]interface{}
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			prev = nil
			continue
		}
		lineNum++

		if prev != nil {
			diff := differ.Compare(prev, rec)
			if diff.HasChanges() {
				summary := buildSummary(diff, lineNum)
				b, err := json.Marshal(summary)
				if err != nil {
					return fmt.Errorf("diff: marshal error: %w", err)
				}
				fmt.Fprintln(out, string(b))
			}
		}
		prev = rec
	}
	return scanner.Err()
}

// buildSummary converts a DiffResult into a plain map suitable for JSON output.
func buildSummary(diff DiffResult, lineNum int) map[string]interface{} {
	summary := map[string]interface{}{
		"_diff_at_line": lineNum,
	}

	if len(diff.Added) > 0 {
		summary["added"] = sortedMap(diff.Added)
	}
	if len(diff.Removed) > 0 {
		summary["removed"] = sortedMap(diff.Removed)
	}
	if len(diff.Changed) > 0 {
		changed := make(map[string]interface{}, len(diff.Changed))
		keys := make([]string, 0, len(diff.Changed))
		for k := range diff.Changed {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := diff.Changed[k]
			changed[k] = map[string]interface{}{"before": v[0], "after": v[1]}
		}
		summary["changed"] = changed
	}
	return summary
}

// sortedMap returns a copy of m with keys in sorted order represented as a
// regular map (JSON marshalling handles ordering independently).
func sortedMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
