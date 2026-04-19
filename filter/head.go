package filter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// Header returns the first N valid records from the input.
type Header struct {
	n int
}

// NewHeader creates a new Header that returns the first n records.
// n must be greater than zero.
func NewHeader(n int) (*Header, error) {
	if n <= 0 {
		return nil, fmt.Errorf("head: n must be greater than zero, got %d", n)
	}
	return &Header{n: n}, nil
}

// Take reads up to n valid JSON records from r and writes them to w.
func (h *Header) Take(r io.Reader, w io.Writer) (int, error) {
	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() && count < h.n {
		line := scanner.Text()
		var rec map[string]interface{}
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			continue
		}
		fmt.Fprintln(w, line)
		count++
	}
	return count, scanner.Err()
}

// RunHead reads from r, writes the first n valid JSON lines to w.
func RunHead(r io.Reader, w io.Writer, n int) (int, error) {
	h, err := NewHeader(n)
	if err != nil {
		return 0, err
	}
	return h.Take(r, w)
}
