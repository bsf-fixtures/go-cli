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

> **This file is nearly empty on purpose.** A codebase's dominant patterns are
> read off the codebase, and as of 2026-11-21 there is none: no `go.mod`, no
> `.go` file, no test, no configuration. There is nothing to be dominant.
>
> Every row below is labelled with where it came from. Rows sourced from "the
> repo" are observed. Rows sourced from "Go's own guidance" are the language's
> defaults, offered so the first change has something to follow rather than
> inventing a house style in private — they are weak, and the first real
> pattern the codebase establishes beats them. Delete or rewrite them once the
> code says otherwise; that is what the scheduled review is for.

## Dominant patterns

For each recurring concern, the pattern most of this codebase already uses.
Follow it rather than introducing a second way of doing the same thing.

| Concern | The pattern here | Source |
| --- | --- | --- |
| Formatting | `gofmt`, unconditionally. No house style, no debate, no configuration. `gofmt -l .` printing nothing is the whole rule. | Go's own guidance; `gofmt` is installed in the sandbox |
| Static analysis | `go vet ./...` is the only linter that runs. `bsf gates static` invokes it directly and counts its findings per file against `.factory/guardrails/baselines.json`, which is currently empty — so the first vet finding introduced is a new one. | Observed: `runGoVet` in `.factory/bin/runner.tgz`; `baselines.json` has `"rules": {}` |
| Errors | Return `error` as the last value and let it travel up. Wrap with `fmt.Errorf("doing x: %w", err)` so the chain survives `errors.Is` / `errors.As`. A library package does not log and does not call `os.Exit`; only `main` decides what a failure means to the user. | Go's own guidance |
| Exit codes | A CLI's contract is its exit code. `main` returns non-zero when the command failed, and the error text goes to stderr while real output goes to stdout. This is the part no test catches by accident and no reviewer sees in a diff. | Go's own guidance |
| Package visibility | `internal/` for the CLI's own packages, `cmd/<name>/` for `main`. `internal/` is compiler-enforced, so it keeps the CLI's guts from becoming someone's accidental library API. | Go's own guidance; see `feature-map.md` |
| Tests | `_test.go` beside the code, table-driven where there is more than one case, `t.Run` per case so a failure names itself. No test framework dependency unless something genuinely needs one — the standard library's `testing` is the default. | Go's own guidance |
| Dependencies | None yet, and the bar for the first one is high: a CLI that is `go build`-able with nothing but the standard library is a CLI anyone can build. See "Taste". | Observed: no `go.mod`, so zero dependencies |

## Product intent

What this product is trying to be, where that would change a decision an
agent has to make.

**Almost nothing is known, and that should be said rather than filled in.**
The only statement of intent anywhere in the repo is `README.md`, in full:

> # go-cli
> Build Software Factory: Fixture of a Go CLI to install factory on

So: it is a Go command-line program, and this repository exists as a fixture
for installing the Build Software Factory on. There is no product
specification, no issue history (the repo has one commit, `c16628e`, which
added the factory scaffolding), no design document and no user.

What follows from that, for an agent making a decision today:

- **Nothing here constrains what the CLI should do.** If a work item arrives
  that implies a product direction, the direction is in the work item, not in
  this file. Do not infer intent from the repo, because the repo has none to
  give.
- **It is a CLI, and that is load-bearing.** It has no server, no UI and no
  page. A change that adds an HTTP server is a change of kind, not of scale,
  and the sandbox cannot verify it (see `runbook.md`) — that is an amendment,
  not something to work around.
- **The first few changes set the conventions.** With no existing code to
  follow, an early change is deciding the layout, the error style and the
  flag-parsing approach for everything after it. Whichever change does that
  should say so in its description and update this file, so the next agent
  inherits a decision rather than re-making it.

## Taste

Things a reviewer would ask for that no rule states.

- **Earn each dependency.** The standard library's `flag` handles a simple
  command. Reach for `cobra`, `urfave/cli` or `kong` when there are
  subcommands, nested flags and generated help to maintain — not before. The
  first dependency added to an empty `go.mod` is a decision about the whole
  project's build, not a convenience. *(Go's own guidance; the repo has no
  dependencies today.)*
- **Thin `main`, testable everything else.** `cmd/<name>/main.go` parses flags,
  wires things together and maps an error to an exit code. It should be short
  enough that nobody needs to test it, because everything it calls is already
  tested. *(Go's own guidance.)*
- **Don't reach for the real world in a test.** No network, no clock, no
  `$HOME`, no global state. The sandbox's environment is closed and its egress
  is limited to the Go module hosts, so a test that reaches out does not fail
  subtly — it fails loudly, somewhere far from the change that caused it. Pass
  an `io.Writer`, a `fs.FS`, a `time.Time` instead. *(Observed: the closed
  environment in `bsf`'s `shellBase`.)*
- **`[no test files]` is a finding.** `go test ./...` exits 0 for a package
  with no tests at all, so a green run is not evidence. A verifier looking at a
  change should read the `go test` output for the package the change touched,
  not just its exit code. *(Observed: Go's behaviour, confirmed against a
  scratch module on 2026-11-21.)*
- **Write the help text as if it is the documentation, because it is.** For a
  CLI with no README beyond one line, `--help` is the entire interface a user
  reads. *(Judgement; the README is two lines long.)*
- **Prefer a boring name.** Go's convention is short, lowercase, no
  underscores, no stutter: `greet.New`, not `greet.NewGreeter` in package
  `greet`. *(Go's own guidance.)*

## Deliberate exceptions

Places that look like they break a rule and should stay as they are, with the
reason. Without this section, cleanup keeps rediscovering them.

| Place | Looks like | Why it stays |
| --- | --- | --- |
| `commands.install`, `commands.build` and `commands.test` in `.factory/sandbox/manifest.json` are each wrapped in `if [ -f go.mod ]; then ...; else echo ...; fi` | A build that cannot fail, which is the classic way to fake a green gate | The repo has no `go.mod` yet, and an unguarded `go build ./...` here exits 1 with "go.mod file not found", which would measure the sandbox at tier `static` and leave the factory unable to run tests at all. The guard is a statement about the repo, not about the build: each command is the ordinary Go command the instant a `go.mod` lands, with nothing else changed. **Remove the guards once `go.mod` is committed** — keeping them after that would hide a deleted module file. |
| `commands.run` and `app_url` are `null` | An app nobody bothered to start | A command-line program binds no port and serves no page, so there is nothing for an HTTP check to reach. Filling these in would have meant guessing an entrypoint path that does not exist. See `runbook.md` → Run. |
| `env` is `[]` in the manifest | A missed step — surely something needs configuring | Nothing in this repo reads an environment variable, because there is no source to read one. The environment is closed, so a variable that later turns out to be needed fails loudly rather than quietly reaching something real. |
| `.factory/guardrails/baselines.json` is `{"rules": {}}` | No static rules are enforced | There is no code to have findings. The static gate still runs `go vet`; an empty baseline means the first finding counts as new, which is the strict direction. |
