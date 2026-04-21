package filter

import (
	"fmt"
	"io"
	"strings"
)

// RunWherePipeline applies one or more where expressions (AND logic) to the
// input stream, writing matching JSON records to w.
// Each expression in exprs must follow the form "<field> <op> <value>".
func RunWherePipeline(r io.Reader, w io.Writer, exprs []string, format OutputFormat) error {
	if len(exprs) == 0 {
		return fmt.Errorf("where: at least one expression is required")
	}

	filters := make([]*WhereFilter, 0, len(exprs))
	for _, expr := range exprs {
		expr = strings.TrimSpace(expr)
		if expr == "" {
			continue
		}
		f, err := NewWhereFilter(expr)
		if err != nil {
			return err
		}
		filters = append(filters, f)
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
		if allMatch(filters, rec) {
			if werr := writer.Write(rec); werr != nil {
				return werr
			}
		}
	}
	return scanner.Err()
}

// allMatch returns true only when every filter matches the record.
func allMatch(filters []*WhereFilter, rec map[string]interface{}) bool {
	for _, f := range filters {
		if !f.Match(rec) {
			return false
		}
	}
	return true
}
