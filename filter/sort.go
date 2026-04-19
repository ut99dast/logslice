package filter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// SortOrder defines ascending or descending sort direction.
type SortOrder int

const (
	SortAsc  SortOrder = iota
	SortDesc SortOrder = iota
)

// Sorter buffers records and emits them sorted by a field.
type Sorter struct {
	field   string
	order   SortOrder
	records []map[string]interface{}
}

// NewSorter creates a Sorter for the given field and order ("asc" or "desc").
func NewSorter(field, order string) (*Sorter, error) {
	if field == "" {
		return nil, fmt.Errorf("sort field must not be empty")
	}
	var o SortOrder
	switch order {
	case "", "asc":
		o = SortAsc
	case "desc":
		o = SortDesc
	default:
		return nil, fmt.Errorf("unknown sort order %q: want asc or desc", order)
	}
	return &Sorter{field: field, order: o}, nil
}

// Add buffers a record for later sorting.
func (s *Sorter) Add(record map[string]interface{}) {
	s.records = append(s.records, record)
}

// Sort sorts the buffered records and writes them to w as JSON lines.
func (s *Sorter) Sort(w io.Writer) error {
	sort.SliceStable(s.records, func(i, j int) bool {
		vi := fmt.Sprintf("%v", s.records[i][s.field])
		vj := fmt.Sprintf("%v", s.records[j][s.field])
		if s.order == SortDesc {
			return vi > vj
		}
		return vi < vj
	})
	enc := json.NewEncoder(w)
	for _, r := range s.records {
		if err := enc.Encode(r); err != nil {
			return err
		}
	}
	return nil
}

// RunSort reads all JSON lines from r, sorts by field/order, writes to w.
func RunSort(r io.Reader, w io.Writer, field, order string) error {
	sorter, err := NewSorter(field, order)
	if err != nil {
		return err
	}
	scanner := NewScanner(r, nil)
	for scanner.Scan() {
		sorter.Add(scanner.Record())
	}
	return sorter.Sort(w)
}
