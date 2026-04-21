package filter

import (
	"fmt"
	"io"
	"strings"
)

// Formatter rewrites a string field using a Go-style template with {field} placeholders.
type Formatter struct {
	field    string
	template string
	parts    []fmtPart
}

type fmtPart struct {
	literal string // literal text before the placeholder
	key     string // field name to substitute, empty if no placeholder follows
}

// NewFormatter creates a Formatter that sets field to the result of expanding
// template. Placeholders are written as {fieldName}.
func NewFormatter(field, template string) (*Formatter, error) {
	if field == "" {
		return nil, fmt.Errorf("formatter: field name must not be empty")
	}
	if template == "" {
		return nil, fmt.Errorf("formatter: template must not be empty")
	}
	parts, err := parseTemplate(template)
	if err != nil {
		return nil, fmt.Errorf("formatter: %w", err)
	}
	return &Formatter{field: field, template: template, parts: parts}, nil
}

func parseTemplate(tmpl string) ([]fmtPart, error) {
	var parts []fmtPart
	for tmpl != "" {
		open := strings.Index(tmpl, "{")
		if open == -1 {
			parts = append(parts, fmtPart{literal: tmpl})
			break
		}
		close := strings.Index(tmpl[open:], "}")
		if close == -1 {
			return nil, fmt.Errorf("unclosed '{' in template")
		}
		close += open
		key := strings.TrimSpace(tmpl[open+1 : close])
		if key == "" {
			return nil, fmt.Errorf("empty placeholder in template")
		}
		parts = append(parts, fmtPart{literal: tmpl[:open], key: key})
		tmpl = tmpl[close+1:]
	}
	return parts, nil
}

// Apply returns a new record with the target field set to the expanded template.
func (f *Formatter) Apply(rec map[string]interface{}) (map[string]interface{}, error) {
	var sb strings.Builder
	for _, p := range f.parts {
		sb.WriteString(p.literal)
		if p.key != "" {
			val, ok := rec[p.key]
			if !ok {
				return nil, fmt.Errorf("formatter: field %q not found in record", p.key)
			}
			fmt.Fprintf(&sb, "%v", val)
		}
	}
	out := shallowCopy(rec)
	out[f.field] = sb.String()
	return out, nil
}

// RunFormat reads JSON lines from r, applies the formatter, and writes results to w.
func RunFormat(r io.Reader, w io.Writer, field, template string) error {
	fmt_, err := NewFormatter(field, template)
	if err != nil {
		return err
	}
	scanner := NewScanner(r)
	writer := NewWriter(w, FormatJSON)
	for scanner.Scan() {
		rec, err := scanner.Record()
		if err != nil {
			continue
		}
		out, err := fmt_.Apply(rec)
		if err != nil {
			continue
		}
		if err := writer.Write(out); err != nil {
			return err
		}
	}
	return scanner.Err()
}
