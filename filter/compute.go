package filter

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Compute operation types supported by NewComputer.
const (
	ComputeAdd = "add"
	ComputeSub = "sub"
	ComputeMul = "mul"
	ComputeDiv = "div"
)

// Computer derives a new numeric field from two existing numeric fields
// using a basic arithmetic operation and writes the result into the record.
type Computer struct {
	destField string
	leftField string
	rightField string
	op string
}

// NewComputer creates a Computer that computes:
//
//	<destField> = <leftField> <op> <rightField>
//
// op must be one of "add", "sub", "mul", or "div".
func NewComputer(destField, leftField, op, rightField string) (*Computer, error) {
	if destField == "" {
		return nil, fmt.Errorf("compute: dest field must not be empty")
	}
	if leftField == "" || rightField == "" {
		return nil, fmt.Errorf("compute: left and right fields must not be empty")
	}
	switch op {
	case ComputeAdd, ComputeSub, ComputeMul, ComputeDiv:
	default:
		return nil, fmt.Errorf("compute: unsupported operation %q (want add|sub|mul|div)", op)
	}
	return &Computer{
		destField:  destField,
		leftField:  leftField,
		rightField: rightField,
		op:         op,
	}, nil
}

// Apply evaluates the arithmetic expression and stores the result in the
// record under destField. The result is stored as a float64. If either
// source field is missing or non-numeric, an error is returned.
func (c *Computer) Apply(record map[string]interface{}) (map[string]interface{}, error) {
	left, err := toFloat(record, c.leftField)
	if err != nil {
		return nil, fmt.Errorf("compute: left field %q: %w", c.leftField, err)
	}
	right, err := toFloat(record, c.rightField)
	if err != nil {
		return nil, fmt.Errorf("compute: right field %q: %w", c.rightField, err)
	}

	var result float64
	switch c.op {
	case ComputeAdd:
		result = left + right
	case ComputeSub:
		result = left - right
	case ComputeMul:
		result = left * right
	case ComputeDiv:
		if right == 0 {
			return nil, fmt.Errorf("compute: division by zero")
		}
		result = left / right
	}

	out := shallowCopy(record)
	out[c.destField] = result
	return out, nil
}

// toFloat extracts a numeric value from the record field, accepting both
// native float64/int values and string representations.
func toFloat(record map[string]interface{}, field string) (float64, error) {
	v, ok := record[field]
	if !ok {
		return 0, fmt.Errorf("field not found")
	}
	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse %q as number", val)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("unsupported type %T", v)
	}
}

// RunCompute reads JSON log lines from r, applies the Computer to each valid
// record, and writes the transformed records to w. Lines that cannot be parsed
// or computed are skipped and counted as errors in the returned Stats.
func RunCompute(r io.Reader, w io.Writer, c *Computer, format OutputFormat) (*Stats, error) {
	stats := &Stats{}
	writer, err := NewWriter(w, format)
	if err != nil {
		return nil, err
	}

	scanner := NewScanner(r)
	for scanner.Scan() {
		record, err := scanner.Record()
		if err != nil {
			stats.Add(false, true)
			continue
		}
		result, err := c.Apply(record)
		if err != nil {
			stats.Add(false, false)
			continue
		}
		if err := writer.Write(result); err != nil {
			return stats, err
		}
		stats.Add(true, false)
	}
	writer.Flush()
	return stats, scanner.Err()
}
