package filter

import (
	"fmt"
	"strings"
)

// MaskMode controls how a field value is masked.
type MaskMode string

const (
	MaskFull    MaskMode = "full"    // replace entire value with stars
	MaskPartial MaskMode = "partial" // keep first/last N chars, mask middle
)

// Masker redacts sensitive field values in log records.
type Masker struct {
	field   string
	mode    MaskMode
	keep    int    // chars to keep on each side (partial mode)
	maskStr string // replacement string
}

// NewMasker creates a Masker for the given field.
// expr format: "field:mode" or "field:mode:keep"
// Examples: "password:full", "email:partial:3"
func NewMasker(expr string) (*Masker, error) {
	parts := strings.SplitN(expr, ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("mask: invalid expression %q, expected field:mode", expr)
	}
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return nil, fmt.Errorf("mask: field name must not be empty")
	}
	mode := MaskMode(strings.TrimSpace(parts[1]))
	if mode != MaskFull && mode != MaskPartial {
		return nil, fmt.Errorf("mask: unknown mode %q, use 'full' or 'partial'", mode)
	}
	keep := 2
	if len(parts) == 3 {
		_, err := fmt.Sscanf(parts[2], "%d", &keep)
		if err != nil || keep < 0 {
			return nil, fmt.Errorf("mask: invalid keep value %q", parts[2])
		}
	}
	return &Masker{field: field, mode: mode, keep: keep, maskStr: "***"}, nil
}

// Apply returns a copy of the record with the target field masked.
// If the field is missing or not a string, the record is returned unchanged.
func (m *Masker) Apply(record map[string]interface{}) map[string]interface{} {
	out := shallowCopy(record)
	v, ok := out[m.field]
	if !ok {
		return out
	}
	str, ok := v.(string)
	if !ok {
		return out
	}
	switch m.mode {
	case MaskFull:
		out[m.field] = m.maskStr
	case MaskPartial:
		out[m.field] = maskPartial(str, m.keep, m.maskStr)
	}
	return out
}

func maskPartial(s string, keep int, mask string) string {
	runes := []rune(s)
	n := len(runes)
	if n <= keep*2 {
		return mask
	}
	return string(runes[:keep]) + mask + string(runes[n-keep:])
}

func shallowCopy(r map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(r))
	for k, v := range r {
		out[k] = v
	}
	return out
}
