package filter

import "fmt"

// TransformFunc is a function that transforms a log record.
type TransformFunc func(record map[string]interface{}) (map[string]interface{}, error)

// Transformer applies a series of TransformFuncs to a record.
type Transformer struct {
	fns []TransformFunc
}

// NewTransformer creates a Transformer with the given functions.
func NewTransformer(fns ...TransformFunc) *Transformer {
	return &Transformer{fns: fns}
}

// Apply runs all transform functions on the record in order.
// Returns the transformed record or an error if any function fails.
func (t *Transformer) Apply(record map[string]interface{}) (map[string]interface{}, error) {
	current := record
	for _, fn := range t.fns {
		result, err := fn(current)
		if err != nil {
			return nil, err
		}
		current = result
	}
	return current, nil
}

// RenameField returns a TransformFunc that renames a field.
func RenameField(from, to string) TransformFunc {
	return func(record map[string]interface{}) (map[string]interface{}, error) {
		val, ok := record[from]
		if !ok {
			return record, nil
		}
		out := shallowCopy(record)
		delete(out, from)
		out[to] = val
		return out, nil
	}
}

// DropField returns a TransformFunc that removes a field from the record.
func DropField(field string) TransformFunc {
	return func(record map[string]interface{}) (map[string]interface{}, error) {
		out := shallowCopy(record)
		delete(out, field)
		return out, nil
	}
}

// AddField returns a TransformFunc that adds or overwrites a field with a static value.
func AddField(field string, value interface{}) TransformFunc {
	return func(record map[string]interface{}) (map[string]interface{}, error) {
		out := shallowCopy(record)
		out[field] = value
		return out, nil
	}
}

// RequireField returns a TransformFunc that errors if a field is missing.
func RequireField(field string) TransformFunc {
	return func(record map[string]interface{}) (map[string]interface{}, error) {
		if _, ok := record[field]; !ok {
			return nil, fmt.Errorf("required field %q missing", field)
		}
		return record, nil
	}
}

func shallowCopy(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
