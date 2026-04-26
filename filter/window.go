package filter

import (
	"fmt"
	"sort"
)

// Windower groups log records into fixed-size sliding or tumbling windows
// based on a numeric or string field value.
type Windower struct {
	field  string
	size   int
	overlap int
}

// WindowGroup holds a set of records belonging to one window.
type WindowGroup struct {
	Index   int
	Records []map[string]interface{}
}

// NewWindower creates a Windower that groups records by position into windows
// of the given size. overlap controls how many records are shared between
// consecutive windows (0 = tumbling, >0 = sliding).
func NewWindower(field string, size, overlap int) (*Windower, error) {
	if size <= 0 {
		return nil, fmt.Errorf("window size must be positive, got %d", size)
	}
	if overlap < 0 || overlap >= size {
		return nil, fmt.Errorf("overlap must be in [0, size), got %d", overlap)
	}
	return &Windower{field: field, size: size, overlap: overlap}, nil
}

// Apply partitions records into windows and returns them in order.
func (w *Windower) Apply(records []map[string]interface{}) []WindowGroup {
	if len(records) == 0 {
		return nil
	}

	step := w.size - w.overlap
	var groups []WindowGroup
	for start := 0; start < len(records); start += step {
		end := start + w.size
		if end > len(records) {
			end = len(records)
		}
		slice := make([]map[string]interface{}, end-start)
		copy(slice, records[start:end])
		groups = append(groups, WindowGroup{
			Index:   len(groups),
			Records: slice,
		})
		if end == len(records) {
			break
		}
	}
	return groups
}

// SortRecordsByField sorts a slice of records by the named field (string comparison).
func SortRecordsByField(records []map[string]interface{}, field string) {
	sort.SliceStable(records, func(i, j int) bool {
		vi, _ := records[i][field].(string)
		vj, _ := records[j][field].(string)
		return vi < vj
	})
}
