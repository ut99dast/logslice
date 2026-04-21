package filter

import (
	"bufio"
	"fmt"
	"io"
)

// newLineScanner returns a line-by-line bufio.Scanner for r.
// It is a thin wrapper used by pivot and other pipeline helpers.
func newLineScanner(r io.Reader) *bufio.Scanner {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	return sc
}

// RunPivotFromArgs parses "row=<field>" and "col=<field>" from args and
// delegates to RunPivot. It is the entry point called from main.
func RunPivotFromArgs(r io.Reader, w io.Writer, args []string) error {
	var rowField, colField string
	for _, a := range args {
		var key, val string
		if _, err := fmt.Sscanf(a, "row=%s", &val); err == nil {
			rowField = val
			continue
		}
		if _, err := fmt.Sscanf(a, "col=%s", &val); err == nil {
			colField = val
			continue
		}
		// Support positional: first arg = row, second = col
		if rowField == "" {
			key = a
			rowField = key
		} else if colField == "" {
			key = a
			colField = key
		}
	}
	if rowField == "" || colField == "" {
		return fmt.Errorf("pivot: usage: pivot <rowField> <colField> or row=<f> col=<f>")
	}
	return RunPivot(r, w, rowField, colField)
}
