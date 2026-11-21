# Print a grand total line

## plan

Today, without --project, tally prints one %-20s line per project (largest first) and stops; with --project it prints only that project's line. After this change, the no-flag path still prints those per-project lines, then one extra line labelled total whose minutes are the sum of every entry, formatted with the existing Format helper and the same %-20s column; the --project path is unchanged and must not print that grand-total line. Add GrandTotal in internal/tally (summing Entry.Minutes, returning Total{Project: "total", ...} so main can reuse the same Fprintf as ProjectTotal) and call it only on the all-projects branch in cmd/tally/main.go — no new flags, dependencies, or markdown.

## files

- internal/tally/tally.go
- internal/tally/tally_test.go
- cmd/tally/main.go
- cmd/tally/main_test.go

## feature folders

- internal/tally
- cmd/tally

## acceptance checks

- Build with `go build -o tally ./cmd/tally`. Pipe `printf '2026-11-20 site 1h30m\n2026-11-20 docs 45m\n# note\n2026-11-21 site 2h\n' | ./tally`. Stdout must be exactly three lines: `site` then `3h30m`, `docs` then `0h45m`, then `total` then `4h15m`, each using the existing `%-20s` column and `Format` duration; exit status 0.
- With the same timesheet on stdin, run `./tally --project site`. Stdout must be exactly one line (`site` then `3h30m`) and must not contain the word `total`; exit status 0.

## tests added

- TestGrandTotalSumsEveryEntry — GrandTotal over mixed-project entries equals the sum of their Minutes and uses Project "total".
- TestRunWithoutProjectFlagEndsWithGrandTotal — run with no --project prints per-project lines then a total line whose duration is their sum (update or replace TestRunWithoutProjectFlagPrintsEveryProjectLargestFirst's expected stdout).
- TestRunWithProjectFlagOmitsGrandTotal — run with --project still prints only that project and never a total line (covered by tightening the existing --project exact-stdout tests if they already pin a single line).

## tests changed

- internal/tally/tally_test.go
- cmd/tally/main_test.go

## research

- `{"claim":"Without --project, main loops tally.Totals and prints each row with fmt.Fprintf(stdout, \"%-20s %s\\n\", t.Project, tally.Format(t.Minutes)); with --project it prints ProjectTotal once and returns — that is the only place a grand-total line can be appended, and only on the all-projects branch.","source":"cmd/tally/main.go:49"}`
- `{"claim":"ProjectTotal already returns a Total for one project (including zero when absent); a GrandTotal that returns Total{Project: \"total\", Minutes: sum} matches that pattern and keeps summing out of main.","source":"internal/tally/tally.go:77"}`
- `{"claim":"Format already renders minutes as NhMMm with zero-padded minutes (65 → 1h05m), so the grand-total duration must go through Format rather than inventing a second renderer.","source":"internal/tally/tally.go:91"}`
- `{"claim":"TestRunWithoutProjectFlagPrintsEveryProjectLargestFirst pins exact stdout without a total line (site 3h30m then docs 0h45m), so that expectation must be updated when the grand total is added; TestRunWithProjectFlagPrintsOnlyThatProject already pins a single-line stdout that must stay free of total.","source":"cmd/tally/main_test.go:34"}`
- `{"claim":"Product intent and column contract: keep %-20s and Format's %dh%02dm; prefer extending domain helpers over new flags; keep main as wiring that maps readers/writers and exit codes.","source":".factory/definition/context/fundamentals.md:52"}`
- `{"claim":"Sandbox can build the binary, run it on stdin/file args, and assert stdout/stderr/exit — so binary-level acceptance checks for the total line and --project omission are exercisable here; there is no web surface to check.","source":".factory/sandbox/manifest.json:56"}`
- `{"claim":"On the current tree, the runbook timesheet prints site 3h30m then docs 0h45m (210+45=255 minutes → 4h15m), and --project site prints only the site line — the numbers the acceptance checks use.","source":"ran: printf '2026-11-20 site 1h30m\\n2026-11-20 docs 45m\\n# note\\n2026-11-21 site 2h\\n' | ./tally  (and the same with --project site)"}`
- `{"claim":"Go's time.ParseDuration (already used by Parse) accepts the timesheet duration forms; Minutes() is the unit Format divides into hours/minutes — no new parsing API is needed for summing.","source":"https://pkg.go.dev/time#ParseDuration"}`
