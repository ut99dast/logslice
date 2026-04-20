package filter

import (
	"fmt"
	"regexp"
)

// Matcher filters log records by matching a field value against a regular expression.
type Matcher struct {
	field  string
	regexp *regexp.Regexp
	invert bool
}

// NewMatcher creates a Matcher that checks whether the given field matches the
// provided regular expression pattern. If invert is true the record matches
// only when the pattern does NOT match (equivalent to grep -v).
func NewMatcher(field, pattern string, invert bool) (*Matcher, error) {
	if field == "" {
		return nil, fmt.Errorf("regex: field name must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("regex: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("regex: invalid pattern %q: %w", pattern, err)
	}
	return &Matcher{field: field, regexp: re, invert: invert}, nil
}

// Match returns true when the record satisfies the regex condition.
func (m *Matcher) Match(record map[string]interface{}) bool {
	val, ok := record[m.field]
	if !ok {
		// missing field never matches
		return m.invert
	}
	s := fmt.Sprintf("%v", val)
	matched := m.regexp.MatchString(s)
	if m.invert {
		return !matched
	}
	return matched
}

// RunRegex reads records from scanner, writes those matching the regex filter
// to writer, and returns the total and matched counts.
func RunRegex(scanner *Scanner, writer *Writer, matcher *Matcher) (total, matched int, err error) {
	for scanner.Scan() {
		total++
		record, parseErr := ParseRecord(scanner.Text())
		if parseErr != nil {
			continue
		}
		if matcher.Match(record) {
			matched++
			if writeErr := writer.Write(record); writeErr != nil {
				return total, matched, writeErr
			}
		}
	}
	return total, matched, scanner.Err()
}
