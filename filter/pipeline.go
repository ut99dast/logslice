package filter

import (
	"fmt"
	"io"
)

// Pipeline ties together a Scanner and a Writer to process
// a log stream end-to-end, returning the number of records written.
type Pipeline struct {
	scanner *Scanner
	writer  *Writer
}

// NewPipeline constructs a Pipeline from reader/writer and options.
func NewPipeline(r io.Reader, w io.Writer, f Filter, format OutputFormat) *Pipeline {
	return &Pipeline{
		scanner: NewScanner(r, f),
		writer:  NewWriter(w, format),
	}
}

// Run processes all records and returns the count written or an error.
func (p *Pipeline) Run() (int, error) {
	count := 0
	for {
		record, ok := p.scanner.Next()
		if !ok {
			break
		}
		if err := p.writer.Write(record); err != nil {
			return count, fmt.Errorf("write error: %w", err)
		}
		count++
	}
	if err := p.scanner.Err(); err != nil {
		return count, fmt.Errorf("scan error: %w", err)
	}
	return count, nil
}
