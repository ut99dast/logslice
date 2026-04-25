package filter

import (
	"fmt"
	"io"
	"strings"
)

// RunCoalescePipeline parses CLI-style arguments and runs the coalesce
// operation. Expected args:
//
//	--fields field1,field2[,...] --out outField [--format json|pretty|csv]
//
// Example:
//
//	RunCoalescePipeline(r, w, []string{"--fields", "email,login", "--out", "user"})
func RunCoalescePipeline(r io.Reader, w io.Writer, args []string) error {
	var fields []string
	var outField string
	format := "json"

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--fields":
			if i+1 >= len(args) {
				return fmt.Errorf("--fields requires a value")
			}
			i++
			for _, f := range strings.Split(args[i], ",") {
				f = strings.TrimSpace(f)
				if f != "" {
					fields = append(fields, f)
				}
			}
		case "--out":
			if i+1 >= len(args) {
				return fmt.Errorf("--out requires a value")
			}
			i++
			outField = args[i]
		case "--format":
			if i+1 >= len(args) {
				return fmt.Errorf("--format requires a value")
			}
			i++
			format = args[i]
		default:
			return fmt.Errorf("unknown argument: %s", args[i])
		}
	}

	if len(fields) == 0 {
		return fmt.Errorf("--fields is required")
	}
	if outField == "" {
		return fmt.Errorf("--out is required")
	}

	return RunCoalesce(r, w, fields, outField, format)
}
