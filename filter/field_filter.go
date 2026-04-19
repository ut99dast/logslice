package filter

import (
	"fmt"
	"strings"
)

// FieldFilter matches log records by a specific field value.
type FieldFilter struct {
	Field    string
	Operator string
	Value    string
}

// NewFieldFilter parses an expression like "level=error" or "status!=200".
func NewFieldFilter(expr string) (*FieldFilter, error) {
	for _, op := range []string{"!=", "="} {
		parts := strings.SplitN(expr, op, 2)
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if field == "" {
				return nil, fmt.Errorf("field name cannot be empty in expression %q", expr)
			}
			return &FieldFilter{Field: field, Operator: op, Value: value}, nil
		}
	}
	return nil, fmt.Errorf("invalid field filter expression %q: expected field=value or field!=value", expr)
}

// Match returns true if the record satisfies the field filter.
func (f *FieldFilter) Match(record map[string]interface{}) bool {
	raw, ok := record[f.Field]
	if !ok {
		return false
	}
	actual := fmt.Sprintf("%v", raw)
	switch f.Operator {
	case "=":
		return actual == f.Value
	case "!=":
		return actual != f.Value
	}
	return false
}
