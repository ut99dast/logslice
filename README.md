# logslice

A CLI tool to filter and slice structured log files by time range or field value.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

## Usage

```bash
# Filter logs by time range
logslice --file app.log --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z"

# Filter by field value
logslice --file app.log --field level=error

# Combine filters
logslice --file app.log --from "2024-01-15T08:00:00Z" --field service=api

# Read from stdin
cat app.log | logslice --field level=warn
```

### Flags

| Flag | Description |
|------|-------------|
| `--file` | Path to the log file (defaults to stdin) |
| `--from` | Start of time range (RFC3339 format) |
| `--to` | End of time range (RFC3339 format) |
| `--field` | Filter by field value in `key=value` format |
| `--format` | Log format: `json` or `logfmt` (default: `json`) |

## Supported Formats

- **JSON** — Newline-delimited JSON logs (e.g., produced by `zerolog`, `zap`, `logrus`)
- **logfmt** — Key-value pair log format

## License

MIT © [yourusername](https://github.com/yourusername)