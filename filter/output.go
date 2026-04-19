package filter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// OutputFormat defines how matched records are written.
type OutputFormat int

const (
	FormatJSON OutputFormat = iota
	FormatPretty
	FormatCSV
)

// ParseOutputFormat converts a string to an OutputFormat.
func ParseOutputFormat(s string) (OutputFormat, error) {
	switch strings.ToLower(s) {
	case "json", "":
		return FormatJSON, nil
	case "pretty":
		return FormatPretty, nil
	case "csv":
		return FormatCSV, nil
	}
	return 0, fmt.Errorf("unknown output format: %q", s)
}

// Writer writes log records to an io.Writer in the specified format.
type Writer struct {
	out    io.Writer
	format OutputFormat
	csvKeys []string
	headerWritten bool
}

// NewWriter creates a new Writer.
func NewWriter(out io.Writer, format OutputFormat) *Writer {
	return &Writer{out: out, format: format}
}

// Write outputs a single record.
func (w *Writer) Write(record map[string]interface{}) error {
	switch w.format {
	case FormatPretty:
		return w.writePretty(record)
	case FormatCSV:
		return w.writeCSV(record)
	default:
		return w.writeJSON(record)
	}
}

func (w *Writer) writeJSON(record map[string]interface{}) error {
	b, err := json.Marshal(record)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w.out, "%s\n", b)
	return err
}

func (w *Writer) writePretty(record map[string]interface{}) error {
	b, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w.out, "%s\n", b)
	return err
}

func (w *Writer) writeCSV(record map[string]interface{}) error {
	if !w.headerWritten {
		for k := range record {
			w.csvKeys = append(w.csvKeys, k)
		}
		_, err := fmt.Fprintln(w.out, strings.Join(w.csvKeys, ","))
		if err != nil {
			return err
		}
		w.headerWritten = true
	}
	vals := make([]string, len(w.csvKeys))
	for i, k := range w.csvKeys {
		if v, ok := record[k]; ok {
			vals[i] = fmt.Sprintf("%v", v)
		}
	}
	_, err := fmt.Fprintln(w.out, strings.Join(vals, ","))
	return err
}
