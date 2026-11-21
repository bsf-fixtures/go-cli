# Technical

## `internal/tally/tally.go`

Adds `GrandTotal(entries []Entry) Total`. It walks every entry, sums `Minutes`, and returns `Total{Project: "total", Minutes: sum}`. Same shape as `ProjectTotal`, so `main` can reuse the existing `%-20s` / `Format` print path.

## `cmd/tally/main.go`

On the all-projects branch (no `--project`), after the per-project loop over `Totals`, prints one more row from `GrandTotal`. The `--project` branch is unchanged and still returns after a single `ProjectTotal` line.

## Tests

- `TestGrandTotalSumsEveryEntry` — mixed-project entries sum to 255 minutes under label `total`.
- `TestRunWithoutProjectFlagEndsWithGrandTotal` — stdout is site, docs, then `total 4h15m`.
- `TestRunWithProjectFlagOmitsGrandTotal` — `--project site` is still a single line with no `total`.
- File-argument exact stdout updated to include the grand-total line when no project flag is set.
