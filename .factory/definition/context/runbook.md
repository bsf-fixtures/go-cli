<!--
  How to build, run, test and verify this app. Written by setup's phase 2,
  with every command tested rather than copied from a README.

  This is the most load-bearing context file: a verifier agent that cannot
  start the app cannot verify anything, and then the factory produces
  confident unverified PRs.
-->

# Runbook — go-cli

> **Read this first.** As of 2026-11-21 this repo contains no Go code: no
> `go.mod`, no `.go` file, no test. Every command below is guarded with
> `[ -f go.mod ]` so that it exits 0 and says so, instead of failing with
> "go.mod file not found in the current directory or any parent directory".
> The guards are not a workaround for a broken build — they are the honest
> answer for an empty repo, and each command becomes the real Go command the
> moment the first `go.mod` is committed. Nothing else about them changes.
>
> The first change that adds code should add `go.mod` at the repo root. Until
> it does, a green build here proves only that the toolchain is reachable.

## Build

```bash
if [ -f go.mod ]; then go build ./...; else echo 'no go.mod in this repo yet: nothing to build'; fi
```

Tested on 2026-11-21, twice: against this repo as it stands (printed its
message, exit 0) and against a scratch module with a `cmd/` entrypoint, an
`internal/` package and a third-party dependency (compiled, exit 0). Both runs
used an environment reduced to `PATH`, `HOME`, `LANG` and `CI`, which is what
the sandbox gives.

`go build ./...` compiles every package and throws the binaries away. To get a
binary you can run, build the command package explicitly — e.g.
`go build -o bin/ ./cmd/...` once `cmd/` exists. There is no `Makefile` and no
build script in this repo; `go build` is the whole build.

Go 1.26.11 is installed in the agent's sandbox; the trial and the gates run on
`ubuntu-latest`, which supplies its own Go. Pinning a version is `go.mod`'s
job, not the sandbox's — with the default `GOTOOLCHAIN=auto`, Go fetches the
toolchain a `go.mod` asks for.

## Run

```bash
# There is nothing to run. See below.
```

The app is then at **nowhere — there is no URL.**

Two separate reasons, and both matter:

1. This is a command-line program. It binds no port and serves no page, so
   there is no address for a browser or an HTTP check to visit. The manifest's
   `app_url` is `null` and its `commands.run` is `null` for this reason, which
   is why the sandbox tier is `build_and_test` and not `renders`. That will not
   change when the code arrives — a CLI never renders.
2. Right now there is also no `main` package, so there is not even a binary to
   invoke.

Once a command package exists, running it is:

```bash
go run ./cmd/<binary-name> --help
# or, from a built binary:
go build -o bin/ ./cmd/... && ./bin/<binary-name> --help
```

`commands.run` was deliberately left `null` rather than filled with a guess at
an entrypoint path that does not exist. If a later change needs the sandbox to
actually start something, that is an amendment: `bsf sandbox report`.

## Test

```bash
if [ -f go.mod ]; then go test ./...; else echo 'no go.mod in this repo yet: no tests to run'; fi
```

Tested on 2026-11-21 the same two ways as the build. Against the scratch
module it reported `ok example.com/probe/internal/greet` and
`? example.com/probe/cmd/probe [no test files]`, exit 0.

Notes:

- `go test ./...` exits 0 for a package with no test files. That is Go's
  behaviour, not a swallowed failure — but it does mean a green test run says
  nothing about whether tests exist.
- This is the command `bsf gates tests` runs on every pull request, read
  straight out of the frozen manifest. Keep it working.
- Tests go in `_test.go` files next to the code they test, in the same
  package (or `package foo_test` for black-box tests). `policy.yaml` requires
  edits to *existing* test files to be declared in the spec's
  `tests_changed`; new test files are always allowed.
- **Coverage is not gated.** `bsf gates coverage` reads only
  `coverage/coverage-summary.json`, `coverage/lcov.info` or `coverage.xml`.
  Go's `-coverprofile` format is none of those and no converter is installed,
  so with `floor_pct` at 0 the coverage gate reports *neutral* and the
  `diff_min_pct` of 80 in `.factory/guardrails/coverage.json` is not actually
  enforced. For a human number, run
  `go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out`.

## Verify

How to see for yourself that the app works, as a developer would:

1. **Check there is something to verify.** `ls go.mod` — if it is missing, the
   repo still has no code and every command below is a no-op. Say so rather
   than reporting a pass.
2. **Compile it.** `go build ./...`. Exit 0 and no output means every package
   compiles.
3. **Vet it.** `go vet ./...`. This is what `bsf gates static` runs for Go, so
   anything it reports will show up in the gate. Exit 0 and no output is clean.
4. **Format it.** `gofmt -l .` must print nothing. Any file it lists is
   misformatted; `gofmt -w .` fixes it.
5. **Test it.** `go test ./...`. Watch for `[no test files]` on a package the
   change touched — that is the signal a change arrived untested.
6. **Run it for real.** Build the binary and invoke it the way a user would,
   including the failure paths: `go build -o bin/ ./cmd/... && ./bin/<name> --help`,
   then the actual subcommand with real arguments, then with bad arguments.
   Check the exit code (`echo $?`): a CLI that prints an error and exits 0 is
   broken, and no test catches that if nobody looks.
7. **Read what it wrote.** If the command creates or edits files, look at them.

There is no browser step and no screenshot to take. A verifier that reports
"rendered correctly" for this repo is reporting on something that does not
exist.

## Seed data

```bash
# None.
```

There is no datastore, no fixture set and no seed script in this repo, so
`commands.seed` is `null`. Nothing needs loading before the tests run.

## Commands that change persistent state

These are kept out of the build path on purpose. Each one writes to something
that outlives the process, so running it against the wrong target is the kind
of mistake that cannot be undone by reverting a commit.

**None found.** There is no migration, no deploy, no publish and no release
step in this repo — nothing it can run writes to anything that outlives the
process. The table is empty because the repo is empty, not because the search
was skipped.

| Command | Writes to | Only when |
| --- | --- | --- |
| | | |

Two things to watch for as code arrives, because they would belong here:
`go mod tidy` rewrites `go.mod` and `go.sum` in the working tree, and any
command the CLI itself grows that writes outside the repo — a user's home
directory, a config file, a remote API.

## Services this app needs

**None.** No database, cache, queue, object store or mail catcher. The
manifest's `services` list is empty, and the trial brings up no container.

| Service | How it starts | Reached at |
| --- | --- | --- |
| | | |

## The sandbox's environment

Built from this list, and from nothing else. It does not inherit a
developer's environment and it does not read any `.env` file, because a
documented build command that runs a migration would otherwise migrate
whatever a developer's credentials happen to point at.

**The list is empty.** No code in this repo reads an environment variable,
because there is no code. The manifest's `env` is `[]` — a fact about the
repo, established by looking for config readers and finding no source to
search, not a gap.

| Variable | Value in the sandbox |
| --- | --- |
| | | 

What the sandbox *does* carry, from the runner rather than from this list, is
fixed by `bsf` itself and not something a manifest chooses: `PATH`, `HOME`,
`TMPDIR`, `LANG`, `CI` and `RUNNER_TEMP`. `HOME` is what lets Go find its
module cache and build cache in their default locations, so no `GOPATH`,
`GOCACHE` or `GOMODCACHE` needs declaring.

Note that `HTTP_PROXY`/`HTTPS_PROXY` are *not* carried. That is fine where
this runs — the trial and the gates both run directly on `ubuntu-latest` with
open network — but it is why `go mod download` fails inside an agent's own
proxied sandbox when the environment is stripped. If you hit "network is
unreachable" from `go` while a `curl` to the same host works, that is the
cause, and the fix is to keep the proxy variables, never to add them to the
manifest.

The sandbox may reach `proxy.golang.org` and `sum.golang.org` for modules, and
`go.dev` / `pkg.go.dev` for documentation. A test that reaches anything else
will fail, which is intended: the environment is closed so that a missing
dependency is loud.
