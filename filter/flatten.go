package filter

import (
	"fmt"
	"strings"
)

// Flattener flattens nested JSON fields into dot-notation keys.
type Flattener struct {
	prefix    string
	maxDepth  int
}

// NewFlattener creates a Flattener. maxDepth <= 0 means unlimited.
func NewFlattener(maxDepth int) (*Flattener, error) {
	if maxDepth < 0 {
		return nil, fmt.Errorf("maxDepth must be >= 0 (0 = unlimited)")
	}
	return &Flattener{maxDepth: maxDepth}, nil
}

// Flatten returns a new record with nested maps expanded to dot-notation keys.
func (f *Flattener) Flatten(record map[string]any) map[string]any {
	out := make(map[string]any)
	f.flattenMap("", record, out, 0)
	return out
}

func (f *Flattener) flattenMap(prefix string, src map[string]any, dst map[string]any, depth int) {
	for k, v := range src {
		key := k
		if prefix != "" {
			key = strings.Join([]string{prefix, k}, ".")
		}
		if nested, ok := v.(map[string]any); ok && (f.maxDepth == 0 || depth < f.maxDepth) {
			f.flattenMap(key, nested, dst, depth+1)
		} else {
			dst[key] = v
		}
	}
}

// RunFlatten reads lines from src, flattens each record, and writes to dst.
func RunFlatten(src []string, maxDepth int) ([]map[string]any, error) {
	fl, err := NewFlattener(maxDepth)
	if err != nil {
		return nil, err
	}
	var results []map[string]any
	for _, line := range src {
		rec, err := ParseRecord(line)
		if err != nil {
			continue
		}
		results = append(results, fl.Flatten(rec))
	}
	return results, nil
}
