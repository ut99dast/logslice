package filter

import (
	"fmt"
	"io"
)

// Renamer renames one or more fields in a log record.
type Renamer struct {
	mappings map[string]string // old -> new
}

// NewRenamer creates a Renamer from a slice of "old=new" expressions.
func NewRenamer(exprs []string) (*Renamer, error) {
	if len(exprs) == 0 {
		return nil, fmt.Errorf("renamer: at least one mapping required")
	}
	mappings := make(map[string]string, len(exprs))
	for _, expr := range exprs {
		old, newName, err := splitKeyValue(expr)
		if err != nil {
			return nil, fmt.Errorf("renamer: invalid mapping %q: %w", expr, err)
		}
		if old == newName {
			return nil, fmt.Errorf("renamer: old and new names are the same: %q", old)
		}
		mappings[old] = newName
	}
	return &Renamer{mappings: mappings}, nil
}

// Apply renames fields in a copy of the record. Fields not in the mapping
// are left unchanged. If an old field is absent the mapping is silently skipped.
func (r *Renamer) Apply(record map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(record))
	for k, v := range record {
		if newName, ok := r.mappings[k]; ok {
			out[newName] = v
		} else {
			out[k] = v
		}
	}
	return out
}

// RunRename reads JSON lines from r, renames fields using renamer, and writes
// results to w using the provided Writer.
func RunRename(r io.Reader, w *Writer, renamer *Renamer) error {
	scanner := NewScanner(r)
	for scanner.Scan() {
		record, err := scanner.Record()
		if err != nil {
			continue
		}
		result := renamer.Apply(record)
		if err := w.Write(result); err != nil {
			return err
		}
	}
	return scanner.Err()
}
