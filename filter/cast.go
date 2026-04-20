package filter

import (
	"fmt"
	"strconv"
)

// CastType represents the target type for casting a field value.
type CastType string

const (
	CastString  CastType = "string"
	CastInt     CastType = "int"
	CastFloat   CastType = "float"
	CastBool    CastType = "bool"
)

// Caster converts the value of a specified field to the given type.
type Caster struct {
	field    string
	target   CastType
}

// NewCaster creates a Caster that will convert the named field to the target type.
// Returns an error if field is empty or target type is unrecognised.
func NewCaster(field string, target CastType) (*Caster, error) {
	if field == "" {
		return nil, fmt.Errorf("cast: field name must not be empty")
	}
	switch target {
	case CastString, CastInt, CastFloat, CastBool:
	default:
		return nil, fmt.Errorf("cast: unsupported target type %q", target)
	}
	return &Caster{field: field, target: target}, nil
}

// Apply returns a copy of the record with the named field cast to the target type.
// If the field is absent the record is returned unchanged.
// If conversion fails an error is returned.
func (c *Caster) Apply(record map[string]interface{}) (map[string]interface{}, error) {
	val, ok := record[c.field]
	if !ok {
		return record, nil
	}

	str := fmt.Sprintf("%v", val)

	var converted interface{}
	var err error

	switch c.target {
	case CastString:
		converted = str
	case CastInt:
		converted, err = strconv.ParseInt(str, 10, 64)
	case CastFloat:
		converted, err = strconv.ParseFloat(str, 64)
	case CastBool:
		converted, err = strconv.ParseBool(str)
	}

	if err != nil {
		return nil, fmt.Errorf("cast: field %q value %q cannot be cast to %s: %w", c.field, str, c.target, err)
	}

	out := make(map[string]interface{}, len(record))
	for k, v := range record {
		out[k] = v
	}
	out[c.field] = converted
	return out, nil
}

// RunCast reads all records from src, applies the caster, and writes results to dst.
// Lines that fail JSON parsing are counted as invalid and skipped.
// Lines where casting fails are also skipped and the error is printed to stderr.
func RunCast(src []string, caster *Caster) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	for _, line := range src {
		record, err := ParseRecord(line)
		if err != nil {
			continue
		}
		out, err := caster.Apply(record)
		if err != nil {
			return nil, err
		}
		results = append(results, out)
	}
	return results, nil
}
