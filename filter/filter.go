package filter

import (
	"encoding/json"
	"strings"
	"time"
)

// Record represents a single parsed log line.
type Record struct {
	Raw    string
	Fields map[string]interface{}
}

// ParseRecord attempts to parse a raw log line as JSON.
func ParseRecord(line string) (*Record, error) {
	r := &Record{Raw: line, Fields: make(map[string]interface{})}
	if err := json.Unmarshal([]byte(line), &r.Fields); err != nil {
		return nil, err
	}
	return r, nil
}

// TimeField extracts a time.Time from the record using the given field name.
func (r *Record) TimeField(key string) (time.Time, bool) {
	v, ok := r.Fields[key]
	if !ok {
		return time.Time{}, false
	}
	s, ok := v.(string)
	if !ok {
		return time.Time{}, false
	}
	formats := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05"}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// FieldEquals checks whether a record field equals the given value (string comparison).
func (r *Record) FieldEquals(key, value string) bool {
	v, ok := r.Fields[key]
	if !ok {
		return false
	}
	s := strings.TrimSpace(strings.ToLower(fmt.Sprintf("%v", v)))
	return s == strings.TrimSpace(strings.ToLower(value))
}
