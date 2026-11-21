<!--
  The conventions no check enforces. Everything here is something a script
  could not decide, which is why it is prose.

  The bar for a line in this file: it has to be a judgement. If a static rule
  could catch it, it belongs in .factory/guardrails/ instead, because a rule
  is reproducible and prose is something agents argue with.

  Each line carries where it came from, so a line nobody can justify can be
  removed later.
-->

# Fundamentals — go-cli

## Dominant patterns

For each recurring concern, the pattern most of this codebase already uses.
Follow it rather than introducing a second way of doing the same thing.

| Concern | The pattern here | Source |
| --- | --- | --- |
| Reporting a bad input | Return an `error` from the `internal/tally` function, built with `fmt.Errorf`, that names the line number first and then quotes the offending text: `line %d: %q is not a YYYY-MM-DD day`. Never log, never skip, never substitute a zero. | `internal/tally/tally.go` — all three error returns in `Parse` follow it. |
| Where the error is turned into an exit | Only in `cmd/tally/main.go`, as `fmt.Fprintln(os.Stderr, "tally:", err)` followed by `os.Exit(1)`. The `tally` prefix is added there, so error strings in `internal/tally` must not repeat the program name. | `cmd/tally/main.go` — both failure paths (open and parse) are written this way. |
| Failing fast vs. collecting | `Parse` returns on the *first* malformed line rather than accumulating a list. One clear error beats a report the caller has to read. | `internal/tally/tally.go`; the doc comment states the intent. |
| Dependencies | The standard library, and nothing else. `go.mod` has no `require` block and there is no `go.sum`. | `go.mod`; every import in both packages. |
| Where logic lives | In `internal/tally`, as functions over `io.Reader` and plain slices, so they are testable without a process. `cmd/tally/main.go` is 30 lines of wiring and does no parsing, summing or formatting of its own. | The split between the two packages. |
| Ordering that users see | Decided once in `Totals` with `sort.Slice`, largest minutes first, ties broken by project name so the output is deterministic. | `internal/tally/tally.go`. |
| Test naming | The test name is the sentence the test proves, not the function it calls: `TestParseNamesTheMalformedLineInsteadOfCountingZero`, `TestTotalsAreLargestFirstWithTiesByName`. A reader of the failure output learns the rule. | All four tests in `internal/tally/tally_test.go`. |
| How a test asserts | `t.Fatalf` with the actual value interpolated (`got %+v, want %+v`); no assertion library, no table unless the cases genuinely vary. | `internal/tally/tally_test.go`. |
| Doc comments | Package and exported-symbol comments are full sentences that say *why*, not just what. The `tally` package comment spends a clause justifying the fail-loudly rule. | `internal/tally/tally.go`, `cmd/tally/main.go`. |

## Product intent

`tally` is a bookkeeping tool for time that someone will bill or report. The
product's whole claim is that its number can be trusted. That makes the
failure mode asymmetric, and it should decide any close call:

**A total that is quietly wrong is worse than no total at all.** A timesheet
that loses an hour without saying so sends someone an incorrect invoice; a
timesheet that refuses to total sends someone back to their text editor for
thirty seconds. Where an agent has to choose between being permissive and
being noisy, be noisy.

Concretely, that means: do not add a "lenient" or "skip bad lines" mode, do
not default a missing field, do not coerce an unparseable duration to zero,
and do not round in a way that hides a parse that half-worked. If a new input
form genuinely should be accepted — say, `1.5h` — accept it *explicitly* in
`Parse` with a test, rather than by loosening the error path.

The second intent is smaller but real: this is a **pipe-friendly Unix tool**.
Results go to stdout, diagnostics go to stderr, a failure exits non-zero, and
the input can come from stdin. Anything that breaks `tally x.txt | sort` or
`cat x.txt | tally` is a regression even if nothing fails.

## Taste

Things a reviewer would ask for that no rule states.

- **Say the line number.** Any new validation error follows the existing
  shape: line number, then the quoted text, then what was expected in the
  form a user could copy (`is not a positive duration like 1h30m`). An error
  that says only "invalid input" would pass every check here and still be
  rejected.
- **Keep `main` thin.** If you find yourself writing a loop or a condition
  about timesheet content in `cmd/tally/main.go`, it belongs in
  `internal/tally` where a test can reach it. `main` exists to choose a
  reader, call `Parse`, and map an error to exit 1.
- **No dependency without an argument for it.** The empty `require` block is
  a feature of this repo, not an accident of its age. Adding a flag library,
  an assertion library or a table-driven test helper needs to be justified
  against the standard library doing the same job.
- **Keep the output columnar.** `%-20s` and the zero-padded `%dh%02dm` exist
  so a human can read a column of totals. Changing either breaks every script
  that greps this output, including `.github/workflows/ci.yml`.
- **A test name is a claim.** Write the sentence the test proves. If you
  cannot state the rule in the name, the test may not be testing a rule.
- **`cmd/tally` has no tests, and that is a known gap, not a licence.** If
  you change `cmd/tally/main.go`, add the test file — the coverage gate asks
  for 80% on changed lines and package `main` currently reports 0.0%.
- **Prefer extending `Parse`'s vocabulary over adding a flag.** There is no
  `flag` usage anywhere in this repo; a command-line option is a new surface
  with no precedent here and should be a deliberate decision, not a default
  way to make something configurable.

## Deliberate exceptions

Places that look like they break a rule and should stay as they are, with the
reason. Without this section, cleanup keeps rediscovering them.

| Place | Looks like | Why it stays |
| --- | --- | --- |
| `Parse` returns `nil, err` on the first bad line and discards entries already read | Throwing away work; a "collect all errors" refactor is the obvious cleanup | Deliberate. The package comment commits to it: one error, named precisely, beats a partial total the caller might use. Do not change this into an error slice. |
| `Entry.Day` is a `string`, not a `time.Time`, even though `Parse` already calls `time.Parse` on it | Parsing and then throwing the parsed value away | The parse is there to *validate*, not to convert. Nothing downstream groups or sorts by day — `Totals` keys on project only — so holding a `time.Time` would add a type with no reader. Revisit only when something actually needs the date. |
| `Entry.Minutes` is an `int`, not a `time.Duration` | Reinventing a standard type | Minutes are the unit the output is in (`%dh%02dm`), and `int` keeps `Totals`'s summation and the tests readable. Note the truncation this implies: `int(d.Minutes())` drops any sub-minute part, so `90s` becomes 1 minute. That is accepted for a timesheet. |
| `cmd/tally/main.go` has no test file at all | An untested package | Accepted for now because it is thin wiring and the behaviour is covered by running the binary (see the runbook's Verify section) and by the smoke step in `.github/workflows/ci.yml`. It is recorded under `capabilities.not_exercised` in the sandbox manifest rather than pretended away. |
| `defer f.Close()` sits before `os.Exit(1)` can be reached elsewhere in `main` | A deferred close that will not run if the process exits | Harmless here: the only `os.Exit` after the file is opened is the parse failure, and the process is ending anyway. Not worth restructuring. |
| `sc.Err()` is returned alongside a non-nil `out` slice | Returning both a value and an error | Correct as written — a scanner error after some entries is still a failure, and every caller checks `err` first. |
