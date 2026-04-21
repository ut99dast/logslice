package filter

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// WhereOp represents a comparison operator.
type WhereOp string

const (
	WhereEq  WhereOp = "eq"
	WhereNe  WhereOp = "ne"
	WhereLt  WhereOp = "lt"
	WhereLte WhereOp = "lte"
	WhereGt  WhereOp = "gt"
	WhereGte WhereOp = "gte"
)

// WhereFilter filters records by comparing a numeric or string field against a value.
type WhereFilter struct {
	field string
	op    WhereOp
	value string
}

// NewWhereFilter parses an expression like "latency gt 200" or "status eq error".
func NewWhereFilter(expr string) (*WhereFilter, error) {
	parts := strings.Fields(expr)
	if len(parts) != 3 {
		return nil, fmt.Errorf("where: invalid expression %q, expected '<field> <op> <value>'", expr)
	}
	op := WhereOp(strings.ToLower(parts[1]))
	switch op {
	case WhereEq, WhereNe, WhereLt, WhereLte, WhereGt, WhereGte:
	default:
		return nil, fmt.Errorf("where: unknown operator %q, valid: eq ne lt lte gt gte", parts[1])
	}
	return &WhereFilter{field: parts[0], op: op, value: parts[2]}, nil
}

// Match returns true if the record satisfies the where condition.
func (w *WhereFilter) Match(rec map[string]interface{}) bool {
	v, ok := rec[w.field]
	if !ok {
		return false
	}
	recVal := fmt.Sprintf("%v", v)

	// Try numeric comparison first.
	recNum, recErr := strconv.ParseFloat(recVal, 64)
	cmpNum, cmpErr := strconv.ParseFloat(w.value, 64)
	if recErr == nil && cmpErr == nil {
		return compareFloat(recNum, cmpNum, w.op)
	}
	// Fall back to string comparison.
	return compareString(recVal, w.value, w.op)
}

func compareFloat(a, b float64, op WhereOp) bool {
	switch op {
	case WhereEq:
		return a == b
	case WhereNe:
		return a != b
	case WhereLt:
		return a < b
	case WhereLte:
		return a <= b
	case WhereGt:
		return a > b
	case WhereGte:
		return a >= b
	}
	return false
}

func compareString(a, b string, op WhereOp) bool {
	switch op {
	case WhereEq:
		return a == b
	case WhereNe:
		return a != b
	case WhereLt:
		return a < b
	case WhereLte:
		return a <= b
	case WhereGt:
		return a > b
	case WhereGte:
		return a >= b
	}
	return false
}

// RunWhere reads JSON lines from r, writes matching records to w.
func RunWhere(r io.Reader, w io.Writer, expr string, format OutputFormat) error {
	filter, err := NewWhereFilter(expr)
	if err != nil {
		return err
	}
	writer, err := NewWriter(w, format)
	if err != nil {
		return err
	}
	defer writer.Flush()

	scanner := NewScanner(r)
	for scanner.Scan() {
		rec, err := ParseRecord(scanner.Bytes())
		if err != nil {
			continue
		}
		if filter.Match(rec) {
			if err := writer.Write(rec); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}
