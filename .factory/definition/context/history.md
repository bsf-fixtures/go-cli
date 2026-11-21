<!--
  What has gone wrong before, and what was learned. Seeded by setup from the
  mined corrections, then added to by the scheduled review.

  This file is for failures whose lesson could not be turned into a check. A
  correction that became a rule does not belong here: the rule is the record,
  and repeating it in prose means two things to keep in step.
-->

# History — go-cli

> **There is no history to mine.** This repository has a single commit,
> `c16628e` "Upgrade the factory to e29d8e0 (#3)", which added the factory
> scaffolding. There is no application code, no previous pull request against
> application code, no reverted change, no incident and no review comment to
> learn from. Nothing in the sections below was invented to fill the file.
>
> The only entries are the two things setup itself ran into, recorded because
> they are the kind of thing the next agent will otherwise hit again.

## Lessons

| What happened | What to do instead | Where it came from |
| --- | --- | --- |
| In a repo with no `go.mod`, every Go command fails identically — `go build ./...`, `go test ./...` and `go mod download` all exit 1 with "go.mod file not found in the current directory or any parent directory". Read as a build failure it looks like a broken toolchain; it is an empty repo. | Check `ls go.mod` before concluding anything about a failing Go command here. The sandbox's commands are guarded for exactly this, and they will say "no go.mod in this repo yet" rather than failing. If you see that message, there is nothing to verify, and saying so is the correct report. | Setup, 2026-11-21, attempt 3 |
| `go mod download` failed with "dial tcp ...: network is unreachable" inside an agent's own sandbox, while `curl` to the same host returned 200. The cause was not the network: the command had been run with the environment stripped, and this sandbox routes HTTPS through a proxy that `go` could no longer see. | When `go` cannot reach `proxy.golang.org` but `curl` can, suspect missing `HTTPS_PROXY`/`HTTP_PROXY` before suspecting the firewall. Never put a proxy variable in the sandbox manifest to fix it — the trial and the gates run directly on `ubuntu-latest` with open network, where that value would be wrong. | Setup, 2026-11-21, attempt 5 |

## Problem classes already eliminated

Kept so nobody proposes a rule that exists, and so a regression is
recognisable as one.

**None.** No problem class has been eliminated here, because none has been
encountered — there has been no application code to have problems. The table
stays empty until something is actually fixed for good.

| Problem class | How it was eliminated | When |
| --- | --- | --- |
| | | |

## Escalations

Work the factory could not finish, and why. A pattern across several of these
is a signal about the factory rather than about the work.

**None.** The factory has not run a work item on this repo yet; setup is the
first thing to execute here.

| Work item | Where it stopped | Cause |
| --- | --- | --- |
| | | |

## Known at setup, not yet a problem

Not failures — things setup established that will otherwise be rediscovered
the hard way. Move a row into "Lessons" if it ever actually bites.

| What | Why it matters |
| --- | --- |
| The sandbox tier is `build_and_test` and will stay there. | This is a CLI: it binds no port and renders no page, so `renders` is unreachable by design. A verifier that reports a rendered screen for this repo is reporting on something that does not exist. |
| Diff coverage is not enforced, despite `diff_min_pct: 80` in `.factory/guardrails/coverage.json`. | `bsf gates coverage` reads only istanbul, lcov or cobertura files; Go's `-coverprofile` is none of them and no converter is installed. With `floor_pct: 0` the gate reports *neutral*. Do not read a green coverage gate as evidence of coverage. |
| `go test ./...` exits 0 for a package with no test files. | A green test gate is not evidence that a change was tested. Read the output for `[no test files]` against the package the change touched. |
| The first change to add code also chooses the layout, the error style and the flag-parsing approach. | With no existing code, there is no pattern to follow and the decision gets made implicitly. Whichever change makes it should say so and update `fundamentals.md`, so the next agent inherits a decision instead of re-making it. |
