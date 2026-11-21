<!--
  How to build, run, test and verify this app. Written by setup's phase 2,
  with every command tested rather than copied from a README.

  This is the most load-bearing context file: a verifier agent that cannot
  start the app cannot verify anything, and then the factory produces
  confident unverified PRs.
-->

# Runbook — go-cli

## Build

```bash
go build ./... && go build -o tally ./cmd/tally
```

The first half compiles every package in the module, so a break in
`internal/tally` is caught even if `cmd/tally` does not reach it. The second
half links the binary at the repo root as `./tally` — the name `.gitignore`
already reserves, and the thing verification actually runs.

The module has no dependencies (`go.mod` has no `require` block, and there is
no `go.sum`), so the build needs no network. `go mod download` is the install
command and prints "no module dependencies to download" today; it is kept
because it is the right command the moment a dependency is added.

Go must be on `PATH`. `go.mod` asks for `go 1.22`; the toolchain measured
during setup was go1.26.11, which satisfies it.

## Run

```bash
printf '2026-11-20 site 1h30m\n2026-11-20 docs 45m\n' | ./tally
```

The app is then at **no URL**. This is a command-line program: it reads a
timesheet from stdin, or from a file path given as the single positional
argument, prints one line per project to stdout, and exits. It does not
listen on a port and there is nothing to open in a browser. `app_url` in the
manifest is null for this reason, and the sandbox trial does not execute the
run command at all — it has no URL to poll.

Both input routes work:

```bash
./tally timesheet.txt     # a file argument
cat timesheet.txt | ./tally   # or stdin
```

`go run ./cmd/tally` does the same thing without the build step and is what
`.github/workflows/ci.yml` uses.

## Test

```bash
go vet ./... && go test ./...
```

`go test ./...` runs the four tests in `internal/tally/tally_test.go` and
reports `[no test files]` for `cmd/tally`. With `-cover`, `internal/tally` is
**90.0% of statements** and `cmd/tally` is **0.0%** — package `main` has no
test file, so a change to `cmd/tally/main.go` has no coverage to count
against the 80% diff-coverage floor and will need a test written alongside
it.

`go vet` is run first because its failures are about the code being wrong in
a way the compiler allows, which is more useful to see before a test failure.

## Verify

How to see for yourself that the app works, as a developer would. Each of
these was run during setup and the output below is what it actually printed.

1. **It totals, largest first, with minutes zero-padded.**

   ```bash
   printf '2026-11-20 site 1h30m\n2026-11-20 docs 45m\n# note\n2026-11-21 site 2h\n' | ./tally
   ```

   ```
   site                 3h30m
   docs                 0h45m
   ```

   The `# note` line is skipped, `site` sums to 3h30m and sorts above `docs`,
   and 45 minutes renders as `0h45m` — hours always present, minutes always
   two digits, so the column stays aligned. Exit status 0.

2. **A malformed line is an error naming its line number, never a silent
   zero.** This is the behaviour the repo exists to protect.

   ```bash
   printf '2026-11-20 site 1h\n2026-11-21 site ninety\n' | ./tally; echo "exit=$?"
   ```

   ```
   tally: line 2: "ninety" is not a positive duration like 1h30m
   exit=1
   ```

   On stderr, not stdout, and nothing is totalled. A change that makes this
   print a total — or that drops the line number — is a regression however
   green the tests are.

3. **A file argument works the same as stdin.**

   ```bash
   printf '2026-11-20 site 1h30m\n' > /tmp/ts.txt && ./tally /tmp/ts.txt
   ```

   ```
   site                 1h30m
   ```

4. **A missing file fails loudly.**

   ```bash
   ./tally /nope.txt; echo "exit=$?"
   ```

   ```
   tally: open /nope.txt: no such file or directory
   exit=1
   ```

5. **The CI smoke check.** `.github/workflows/ci.yml` asserts the first line
   matches `^site *1h30m$` for a two-entry timesheet. Reproduce it with:

   ```bash
   printf '2026-11-20 site 1h30m\n2026-11-20 docs 45m\n' | go run ./cmd/tally | grep -q '^site *1h30m$'
   ```

## Seed data

```bash
# none — this repo ships no seed data and has nothing to seed.
```

There is no datastore and no fixture file. The app's input is a timesheet the
caller supplies, so "seeding" is writing one:

```bash
printf '2026-11-20 site 1h30m\n2026-11-20 docs 45m\n# a comment is skipped\n2026-11-21 site 2h\n' > /tmp/timesheet.txt
```

A timesheet line is `DAY PROJECT DURATION` — exactly three whitespace-separated
fields. `DAY` must parse as `2006-01-02`, `DURATION` must be a positive
`time.ParseDuration` string such as `1h30m`, `45m` or `2h`. Blank lines and
lines beginning with `#` are skipped. `seed` in the manifest is null because
inventing a seed step the repo does not have would be a fiction.

## Commands that change persistent state

These are kept out of the build path on purpose. Each one writes to something
that outlives the process, so running it against the wrong target is the kind
of mistake that cannot be undone by reverting a commit.

| Command | Writes to | Only when |
| --- | --- | --- |
| — | **There are none.** `tally` opens its input read-only and writes only to stdout and stderr. It has no migration, no cache warm, no upload, and no network call of any kind. | — |

The only command here that writes anything to disk is the build, and it
writes one file: `./tally`, which `.gitignore` already ignores.

## Services this app needs

| Service | How it starts | Reached at |
| --- | --- | --- |
| — | **None.** No database, cache, queue, object store or mail catcher. `services` in the manifest is empty. | — |

## The sandbox's environment

Built from this list, and from nothing else. It does not inherit a
developer's environment and it does not read any `.env` file, because a
documented build command that runs a migration would otherwise migrate
whatever a developer's credentials happen to point at.

| Variable | Value in the sandbox |
| --- | --- |
| — | **No variable is declared.** Every Go source file was searched for `os.Getenv`, `os.LookupEnv`, `os.Environ` and the `flag` package; there are no matches. The program reads its entire input from stdin or one positional argument, so there is nothing to configure and `env` in the manifest is empty. |

The sandbox still passes through `PATH`, `HOME`, `TMPDIR`, `LANG`, `CI` and
`RUNNER_TEMP` — that is the runner's fixed base, not a declaration of this
app's. `HOME` is what lets the Go build and module caches land in their usual
places; nothing else in the closed environment matters to this repo.

Egress is limited to `proxy.golang.org`, `sum.golang.org` and
`storage.googleapis.com`: the Go module proxy and checksum database. Nothing
is fetched from them today, because there are no dependencies. They are
allowed so that the first dependency added to `go.mod` does not fail the
sandbox before anyone has read why.
