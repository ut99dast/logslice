package filter

import "fmt"

// MultiFilter combines multiple FieldFilters with AND logic.
type MultiFilter struct {
	Filters []*FieldFilter
}

// NewMultiFilter builds a MultiFilter from a slice of expression strings.
func NewMultiFilter(exprs []string) (*MultiFilter, error) {
	filters := make([]*FieldFilter, 0, len(exprs))
	for _, expr := range exprs {
		f, err := NewFieldFilter(expr)
		if err != nil {
			return nil, fmt.Errorf("invalid filter expression %q: %w", expr, err)
		}
		filters = append(filters, f)
	}
	return &MultiFilter{Filters: filters}, nil
}

// Match returns true only if all contained filters match the record.
func (m *MultiFilter) Match(record map[string]interface{}) bool {
	for _, f := range m.Filters {
		if !f.Match(record) {
			return false
		}
	}
	return true
}

// Empty returns true when no filters are configured.
func (m *MultiFilter) Empty() bool {
	return len(m.Filters) == 0
}
