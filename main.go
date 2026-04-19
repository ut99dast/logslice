package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

const timeLayout = time.RFC3339

type Config struct {
	InputFile  string
	TimeField  string
	FilterField string
	FilterValue string
	From       string
	To         string
}

func parseFlags() Config {
	var cfg Config
	flag.StringVar(&cfg.InputFile, "f", "", "Path to the log file (default: stdin)")
	flag.StringVar(&cfg.TimeField, "time-field", "timestamp", "JSON field to use as the time key")
	flag.StringVar(&cfg.From, "from", "", "Start of time range (RFC3339, e.g. 2024-01-01T00:00:00Z)")
	flag.StringVar(&cfg.To, "to", "", "End of time range (RFC3339, e.g. 2024-01-02T00:00:00Z)")
	flag.StringVar(&cfg.FilterField, "field", "", "Field name to filter on")
	flag.StringVar(&cfg.FilterValue, "value", "", "Value the filter field must match")
	flag.Parse()
	return cfg
}

func matchesTimeRange(entry map[string]interface{}, timeField, from, to string) (bool, error) {
	if from == "" && to == "" {
		return true, nil
	}

	raw, ok := entry[timeField]
	if !ok {
		return false, nil
	}

	timeStr, ok := raw.(string)
	if !ok {
		return false, fmt.Errorf("time field %q is not a string", timeField)
	}

	t, err := time.Parse(timeLayout, timeStr)
	if err != nil {
		return false, fmt.Errorf("cannot parse time %q: %w", timeStr, err)
	}

	if from != "" {
		start, err := time.Parse(timeLayout, from)
		if err != nil {
			return false, fmt.Errorf("invalid --from value: %w", err)
		}
		if t.Before(start) {
			return false, nil
		}
	}

	if to != "" {
		end, err := time.Parse(timeLayout, to)
		if err != nil {
			return false, fmt.Errorf("invalid --to value: %w", err)
		}
		if t.After(end) {
			return false, nil
		}
	}

	return true, nil
}

func matchesFieldFilter(entry map[string]interface{}, field, value string) bool {
	if field == "" {
		return true
	}
	v, ok := entry[field]
	if !ok {
		return false
	}
	return fmt.Sprintf("%v", v) == value
}

func run(cfg Config) error {
	var input *os.File
	if cfg.InputFile != "" {
		f, err := os.Open(cfg.InputFile)
		if err != nil {
			return fmt.Errorf("cannot open file: %w", err)
		}
		defer f.Close()
		input = f
	} else {
		input = os.Stdin
	}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil t// Skip non-JSON lines silently
			continue
		}

, err := matchesTimeRange(entry, cfg.TimeField, cfg.From, cfg.To)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: %v\n", err)
			continue
		}
		if !ok {
			continue
		}

		if !matchesFieldFilter(entry, cfg.FilterField, cfg.FilterValue) {
			continue
		}

		fmt.Println(line)
	}

	return scanner.Err()
}

func main() {
	cfg := parseFlags()
	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
