# Golang Development Manual

[中文说明](README.md) · [Read online](https://foreversrc.github.io/golang-dev-manual/)

![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)
![Go 1.27.1](https://img.shields.io/badge/Go-1.27.1-00ADD8.svg)

A personal Golang development manual: 113 language-level clauses graded `MUST` / `SHOULD` / `MAY`, each carrying a summary, details, a rationale, an upstream quote, good and bad examples, and references. Every clause is objectively checkable, so it can be turned into a lint rule or a review item. The manual ships as bilingual markdown and is published as a static site.

## Design intent

- Turn Go coding experience scattered across daily work into a checkable set of conventions.
- Follow the form of the *Alibaba Java Development Manual*: three grades, numbered clauses, good and bad examples.
- Its users are the author and his agents, with no organization-wide rollout.
- Friendly to people: readable chapter by chapter as learning material.
- Friendly to agents: code review gets a low-cost three-level clause list to check item by item, with little context cost. This one is the emphasis.
- Still useful in the AI coding era, for two reasons:
  - Correcting AI's default deviations: several clauses were added while fixing code written by AI, such as error assertions comparing error text, assertions projecting the result into a subset, exported and unexported declarations out of order, and comments restating what the code already says. These deviations are easy to judge, yet AI makes them by default and people do not always catch them, so locking them into numbered, checkable clauses beats an ad-hoc reminder each time.
  - Getting AI to absorb the latest idioms and follow best practices: this matches the purpose behind the go fix and modernize clauses.

## Relationship to other style guides

- This manual is a checklist curated on top of Effective Go, Go Code Review Comments, Google Go Style, and the Uber Go Style Guide, with personal selection and wording.
- It keeps clauses that are short, checkable, and supported by a directly corresponding upstream reference.
- Those guides stay the reference for anything this manual does not cover.

## Quick start

Install the CLI:

```bash
go install github.com/ForeverSRC/golang-dev-manual/gdm/cmd/gdm-cli@latest
```

```bash
gdm-cli list --level MUST     # one line per clause
gdm-cli explain NAMING-001    # details, examples, and references
gdm-cli check ./some/package  # run the automatable checks
```

The clause data is embedded in the binary, so the installed `gdm-cli` needs no checkout.

The companion skill `skills/gdm-workflows/` wraps the CLI into agent workflows and installs with `npx skills`:

```bash
npx skills add ForeverSRC/golang-dev-manual
```

## Building from source

Contributors work from a checkout:

```bash
make build     # build bin/gdm-cli (published CLI) and bin/gdm-gen (manual generator)
make gen       # regenerate manual/zh/ and manual/en/ from gdm/data/
make check     # go vet + go test
make lint      # golangci-lint
make generate  # regenerate wire wiring code
```

Without building, run `go run ./gdm/cmd/gdm-cli <command>`.

## CLI usage

`--lang` picks the output language: `zh` (the default) or `en`. Both the clause text and the field labels follow it.

`gdm-cli list` prints one line per clause, cheap enough to read in full before deciding what to expand.

```console
$ gdm-cli list --level MUST --lang en
NAMING-001: Package names use lowercase words written together, without underscores, camel case, or plurals.
NAMING-002: Identifiers do not use underscores to separate words.
```

`gdm-cli explain` expands one or more clauses with details, examples, and references.

```console
$ gdm-cli explain NAMING-001 --lang en
【MUST】NAMING-001 Package names use lowercase words written together, without underscores, camel case, or plurals.

Category: Programming Conventions/Naming | Since: 1.0

Details:
The package name is the prefix that callers use to reference identifiers.……

Why:
Quote: In Go, package names must be concise and use only lowercase letters and numbers……
Source: https://google.github.io/styleguide/go/decisions#package-names

Good: the accepted form, as a Go snippet
Bad: the rejected form, as a Go snippet

References:
- https://go.dev/wiki/CodeReviewComments#package-names

Detection: golangci-lint stylecheck(ST1003) (Check the package name identifier)
```

`gdm-cli check` runs the grep-based checks against a code directory and reports the hits.

```console
$ gdm-cli check ./gdm --verbose --lang en
No hits found.

Clauses without an automated check:
  [NAMING-001] golangci-lint stylecheck(ST1003) (Check the package name identifier)
```

`gdm-cli list` also takes:

- `--level`: filter by level, e.g. `MUST,SHOULD`.
- `--category`: filter by category id, e.g. `programming-conventions` or `programming-conventions/naming`. An unknown value fails and lists every available id.

## Repository layout

- `manual/zh/`, `manual/en/` — the generated bilingual markdown manual. Committed, and read directly as the site source. Do not edit by hand.
- `gdm/` — the Go tooling.
  - `gdm/data/` — the single source of truth. `gdm/data/data.go` embeds it into the `gdm-cli` and `gdm-gen` binaries.
    - `gdm/data/manual.json` — version, Go baseline, table of contents.
    - `gdm/data/clauses/<chapter id>.json` — clauses, one file per chapter.
    - `gdm/data/i18n/<lang>/` — language overlays. Chinese is the source; other languages hang translations off it.
  - `gdm/cmd/gdm-cli` — the published CLI, layered as `internal/domain`, `internal/service`, `internal/repository`, `internal/adapter`, `internal/api`, `internal/server`, with wiring under `di/`.
  - `gdm/cmd/gdm-gen` — the manual generator, built only by maintainers and CI.
  - `gdm/CLAUSES.md` — the clause data format, writing rules, and self-checks.
- `REFERENCE.md` — the upstream material list behind the clauses.
- `skills/gdm-workflows/` — a distributable skill that wraps the CLI into agent workflows.
- `website/` — MkDocs configuration for the static site.

## Clause format

| Field | Meaning |
| --- | --- |
| `id` | Clause number, `PREFIX-NNN`, globally unique |
| `level` | `MUST` / `SHOULD` / `MAY` |
| `chapter` | Chapter id, must exist in `toc` |
| `section` | Section id, optional; omitted means the whole chapter |
| `since_go` | Lowest Go version the clause can be applied with |
| `summary` | One-line requirement, what `gdm-cli list` prints |
| `details` | Scope and exceptions |
| `rationale` | Why the clause is written this way |
| `quote` | Verbatim upstream excerpt and its source |
| `examples` | `good` and `bad` minimal snippets |
| `sources` | Upstream links, one per directly corresponding reference |
| `detect` | How the clause is checked, optional |

Levels:

- `MUST` — violating it is wrong; it should be detectable automatically or by mechanical review.
- `SHOULD` — follow by default; an exception needs a stated reason.
- `MAY` — direction and idiom, applied by judgement.

`detect.tool` takes three values:

- `golangci-lint` — the rule name in `detect.rule`.
- `grep` — a regex in `detect.pattern`, scanned over `.go` files by `gdm-cli check`.
- `manual` — human review, with the specifics in `detect.note`.

## Contributing

- Open an issue to report a defective clause, a missing clause, a broken or misattributed reference, or a wrong category.
- Send a pull request that changes `gdm/data/` only, and regenerate `manual/` with `make gen` in the same commit.
- Run `make check` and `make lint` before submitting.
- Project-specific conventions are out of scope, as is organization-wide rollout.

## License and acknowledgements

- Released under the MIT License, covering the CLI source and the clause text.
- Quoted upstream excerpts remain the property of their respective sources; each is attributed at `quote.source` and listed in the appendix.
- Upstream guides this manual draws on: Effective Go, Go Code Review Comments, Google Go Style, Uber Go Style Guide.
