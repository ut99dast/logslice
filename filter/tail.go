package filter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// Tailer keeps a rolling buffer of the last N records.
type Tailer struct {
	n      int
	buffer []map[string]interface{}
	pos    int
	full   bool
}

// NewTailer creates a Tailer that retains the last n records.
func NewTailer(n int) (*Tailer, error) {
	if n <= 0 {
		return nil, fmt.Errorf("tail: n must be greater than zero, got %d", n)
	}
	return &Tailer{
		n:      n,
		buffer: make([]map[string]interface{}, n),
	}, nil
}

// Add inserts a record into the rolling buffer.
func (t *Tailer) Add(record map[string]interface{}) {
	t.buffer[t.pos] = record
	t.pos = (t.pos + 1) % t.n
	if t.pos == 0 {
		t.full = true
	}
}

// Records returns the retained records in insertion order.
func (t *Tailer) Records() []map[string]interface{} {
	if !t.full {
		return t.buffer[:t.pos]
	}
	out := make([]map[string]interface{}, t.n)
	copy(out, t.buffer[t.pos:])
	copy(out[t.n-t.pos:], t.buffer[:t.pos])
	return out
}

// Reset clears the buffer.
func (t *Tailer) Reset() {
	t.buffer = make([]map[string]interface{}, t.n)
	t.pos = 0
	t.full = false
}

// RunTail reads all lines from r, keeps the last n valid JSON records, and writes them to w.
func RunTail(r io.Reader, w io.Writer, n int) error {
	tailer, err := NewTailer(n)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		var record map[string]interface{}
		if err := json.Unmarshal(line, &record); err != nil {
			continue
		}
		tailer.Add(record)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	for _, rec := range tailer.Records() {
		if err := enc.Encode(rec); err != nil {
			return err
		}
	}
	return nil
}
