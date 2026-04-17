# logslice

A CLI tool to filter and slice structured log files by time range or field values.

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
logslice --from "2024-01-15T10:00:00Z" --to "2024-01-15T11:00:00Z" app.log

# Filter by field value
logslice --field level=error app.log

# Combine filters
logslice --from "2024-01-15T10:00:00Z" --field service=api --field level=warn app.log

# Read from stdin
cat app.log | logslice --from "2024-01-15T10:00:00Z" --to "2024-01-15T11:00:00Z"
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start of time range (RFC3339) |
| `--to` | End of time range (RFC3339) |
| `--field` | Filter by field value (`key=value`), repeatable |
| `--format` | Input format: `json`, `logfmt` (default: `json`) |
| `--output` | Output format: `json`, `logfmt` (default: same as input) |

## Requirements

- Go 1.21 or later

## License

MIT — see [LICENSE](LICENSE) for details.