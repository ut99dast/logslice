package filter

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// PivotResult holds aggregated counts keyed by (row, col) pairs.
type PivotResult struct {
	RowField string
	ColField string
	Cells    map[string]map[string]int
	RowKeys  []string
	ColKeys  []string
}

// NewPivotResult creates an empty PivotResult for the given row and column fields.
func NewPivotResult(rowField, colField string) *PivotResult {
	return &PivotResult{
		RowField: rowField,
		ColField: colField,
		Cells:    make(map[string]map[string]int),
	}
}

// Add records a single observation for the given row and column values.
func (p *PivotResult) Add(rowVal, colVal string) {
	if _, ok := p.Cells[rowVal]; !ok {
		p.Cells[rowVal] = make(map[string]int)
	}
	p.Cells[rowVal][colVal]++
}

// Finalize sorts row and column keys for deterministic output.
func (p *PivotResult) Finalize() {
	rowSet := make(map[string]struct{})
	colSet := make(map[string]struct{})
	for r, cols := range p.Cells {
		rowSet[r] = struct{}{}
		for c := range cols {
			colSet[c] = struct{}{}
		}
	}
	for r := range rowSet {
		p.RowKeys = append(p.RowKeys, r)
	}
	for c := range colSet {
		p.ColKeys = append(p.ColKeys, c)
	}
	sort.Strings(p.RowKeys)
	sort.Strings(p.ColKeys)
}

// WritePivotTable writes the pivot table to w in a tab-aligned format.
func WritePivotTable(w io.Writer, p *PivotResult) {
	tw := tabwriter.NewWriter(w, 4, 0, 2, ' ', 0)
	// header row
	fmt.Fprintf(tw, "%s\\%s", p.RowField, p.ColField)
	for _, c := range p.ColKeys {
		fmt.Fprintf(tw, "\t%s", c)
	}
	fmt.Fprintln(tw)
	// data rows
	for _, r := range p.RowKeys {
		fmt.Fprintf(tw, "%s", r)
		for _, c := range p.ColKeys {
			fmt.Fprintf(tw, "\t%d", p.Cells[r][c])
		}
		fmt.Fprintln(tw)
	}
	tw.Flush()
}

// RunPivot reads JSON log lines from r, counts occurrences of each
// (rowField, colField) pair, and writes a pivot table to w.
func RunPivot(r io.Reader, w io.Writer, rowField, colField string) error {
	if rowField == "" || colField == "" {
		return fmt.Errorf("pivot: row and col fields must not be empty")
	}
	result := NewPivotResult(rowField, colField)
	scanner := newLineScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		rec, err := ParseRecord(line)
		if err != nil {
			continue
		}
		rowVal, rok := rec[rowField]
		colVal, cok := rec[colField]
		if !rok || !cok {
			continue
		}
		result.Add(fmt.Sprintf("%v", rowVal), fmt.Sprintf("%v", colVal))
	}
	result.Finalize()
	WritePivotTable(w, result)
	return nil
}
