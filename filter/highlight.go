package filter

import (
	"fmt"
	"regexp"
	"strings"
)

// HighlightMode controls how matched values are annotated.
type HighlightMode string

const (
	HighlightBracket HighlightMode = "bracket" // wraps value in [[ ]]
	HighlightUpper   HighlightMode = "upper"   // uppercases matched substring
	HighlightMark    HighlightMode = "mark"    // prefixes field value with >>>
)

// Highlighter annotates matching field values in log records.
type Highlighter struct {
	field  string
	re     *regexp.Regexp
	mode   HighlightMode
}

// NewHighlighter creates a Highlighter for the given field, regex pattern, and mode.
// Valid modes: bracket, upper, mark.
func NewHighlighter(field, pattern string, mode HighlightMode) (*Highlighter, error) {
	if field == "" {
		return nil, fmt.Errorf("highlight: field must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("highlight: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("highlight: invalid pattern %q: %w", pattern, err)
	}
	switch mode {
	case HighlightBracket, HighlightUpper, HighlightMark:
	default:
		return nil, fmt.Errorf("highlight: unknown mode %q", mode)
	}
	return &Highlighter{field: field, re: re, mode: mode}, nil
}

// Apply annotates the field value in the record if the pattern matches.
// Returns a shallow copy with the annotated field; original is unchanged.
func (h *Highlighter) Apply(record map[string]interface{}) map[string]interface{} {
	val, ok := record[h.field]
	if !ok {
		return record
	}
	str, ok := val.(string)
	if !ok {
		return record
	}
	if !h.re.MatchString(str) {
		return record
	}
	out := shallowCopy(record)
	switch h.mode {
	case HighlightBracket:
		out[h.field] = h.re.ReplaceAllStringFunc(str, func(m string) string {
			return "[[" + m + "]]"
		})
	case HighlightUpper:
		out[h.field] = h.re.ReplaceAllStringFunc(str, func(m string) string {
			return strings.ToUpper(m)
		})
	case HighlightMark:
		out[h.field] = ">>>" + str
	}
	return out
}
