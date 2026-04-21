package filter

import (
	"fmt"
	"io"
	"strings"
)

// TemplateRenderer applies a Go-style text template to each log record,
// substituting {{.field}} placeholders with field values from the record.
type TemplateRenderer struct {
	tmpl string
	parts []templatePart
}

type templatePart struct {
	text  string
	field string // empty means literal text
}

// NewTemplateRenderer parses the template string and returns a TemplateRenderer.
// Returns an error if the template contains unclosed braces.
func NewTemplateRenderer(tmpl string) (*TemplateRenderer, error) {
	if tmpl == "" {
		return nil, fmt.Errorf("template: template string must not be empty")
	}
	parts, err := parseTemplateParts(tmpl)
	if err != nil {
		return nil, fmt.Errorf("template: %w", err)
	}
	return &TemplateRenderer{tmpl: tmpl, parts: parts}, nil
}

func parseTemplateParts(tmpl string) ([]templatePart, error) {
	var parts []templatePart
	for len(tmpl) > 0 {
		start := strings.Index(tmpl, "{{.")
		if start == -1 {
			parts = append(parts, templatePart{text: tmpl})
			break
		}
		if start > 0 {
			parts = append(parts, templatePart{text: tmpl[:start]})
		}
		rest := tmpl[start+3:]
		end := strings.Index(rest, "}}")
		if end == -1 {
			return nil, fmt.Errorf("unclosed '{{.' in template")
		}
		field := rest[:end]
		if field == "" {
			return nil, fmt.Errorf("empty field name in template")
		}
		parts = append(parts, templatePart{field: field})
		tmpl = rest[end+2:]
	}
	return parts, nil
}

// Apply renders the template against rec and returns the resulting string.
// Missing fields are rendered as "<nil>".
func (r *TemplateRenderer) Apply(rec map[string]interface{}) string {
	var sb strings.Builder
	for _, p := range r.parts {
		if p.field == "" {
			sb.WriteString(p.text)
		} else if val, ok := rec[p.field]; ok {
			sb.WriteString(fmt.Sprintf("%v", val))
		} else {
			sb.WriteString("<nil>")
		}
	}
	return sb.String()
}

// RunTemplate reads JSON records from r, renders each using the template,
// and writes one rendered line per record to w.
func RunTemplate(r io.Reader, w io.Writer, tmplStr string) error {
	renderer, err := NewTemplateRenderer(tmplStr)
	if err != nil {
		return err
	}
	scanner := NewScanner(r)
	for scanner.Scan() {
		rec, err := ParseRecord(scanner.Bytes())
		if err != nil {
			continue
		}
		line := renderer.Apply(rec)
		fmt.Fprintln(w, line)
	}
	return scanner.Err()
}
