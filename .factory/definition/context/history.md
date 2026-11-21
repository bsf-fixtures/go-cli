<!--
  What has gone wrong before, and what was learned. Seeded by setup from the
  mined corrections, then added to by the scheduled review.

  This file is for failures whose lesson could not be turned into a check. A
  correction that became a rule does not belong here: the rule is the record,
  and repeating it in prose means two things to keep in step.
-->

# History — go-cli

**Nothing has gone wrong here yet, and this file says so rather than
inventing lessons.** The repository has a single commit (`1df2dd4`, "Upgrade
the factory to build-software-factory ecee54e (#7)"), no merged pull requests
to mine, no revert, no issue, and no factory run. There is no corrections
corpus to seed from. Every table below is empty on purpose; an invented row
would be worse than a blank one, because an agent would read it as evidence.

## Lessons

| What happened | What to do instead | Where it came from |
| --- | --- | --- |
| | | |

*Empty: no past failure is recorded in this repository's history. The one
belief this codebase does hold about failure — that a malformed line must be
an error naming its line number rather than a silent zero — is not a lesson
learned here. It is stated up front in the package comment of
`internal/tally/tally.go` and in the README as the rule the fixture exists to
enforce, and it is already enforced by
`TestParseNamesTheMalformedLineInsteadOfCountingZero`. Because a test
enforces it, it belongs there and in* `fundamentals.md` *under product
intent, not here.*

## Problem classes already eliminated

Kept so nobody proposes a rule that exists, and so a regression is
recognisable as one.

| Problem class | How it was eliminated | When |
| --- | --- | --- |
| A malformed timesheet line silently contributing zero minutes to a total | `tally.Parse` returns on the first bad line with an error naming the line number; `TestParseNamesTheMalformedLineInsteadOfCountingZero` fails if it does not | Present in the initial commit `1df2dd4`; not the result of an incident recorded here |
| Non-deterministic output order from ranging over a map | `Totals` sorts explicitly with `sort.Slice` — minutes descending, then project name ascending — and `TestTotalsAreLargestFirstWithTiesByName` pins it | Present in the initial commit `1df2dd4` |
| Misaligned or ambiguous durations in the output column | `Format` renders `%dh%02dm`, pinned by `TestFormatPadsMinutes` (65 → `1h05m`), and `.github/workflows/ci.yml` greps the live binary's output for `^site *1h30m$` | Present in the initial commit `1df2dd4` |

*These three are listed because a change that reintroduces one should be
recognisable as a regression, not because any of them was ever observed
failing here. Each is backed by a test, so the test is the record; this table
is the index to it.*

## Escalations

Work the factory could not finish, and why. A pattern across several of these
is a signal about the factory rather than about the work.

| Work item | Where it stopped | Cause |
| --- | --- | --- |
| | | |

*Empty: the factory has not run on this repository. The first entries will
arrive from real runs. One limit is known in advance and worth expecting
here — the sandbox tier is `build_and_test`, so a role that needs to see a
rendered page will not get one; this repo has no web surface at all, so that
is a property of the application, not a factory failure, and it should not be
escalated as one.*
