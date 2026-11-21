# go-cli

Build Software Factory fixture: a command-line tool with no web surface, so
verification has to run the binary and read its output rather than open a
page.

`tally` totals the time per project in a plain-text timesheet. The rule worth
enforcing is honesty about bad input: a malformed line is an error naming its
line number, never a silent zero.

## Build, test, run

```bash
go build ./...
go test ./...
go vet ./...
printf '2026-11-20 site 1h30m\n2026-11-20 docs 45m\n' | go run ./cmd/tally
```
