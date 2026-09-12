# Contributing to GhanaCalendar

GhanaCalendar is independent open-source public-interest infrastructure: evidence-backed Ghana holiday occurrences and deterministic working-day arithmetic. Contributions to the dataset, the engine, the interfaces, the documentation and the provenance record are all welcome.

Read [`AGENTS.md`](AGENTS.md), [`agent_plan.md`](agent_plan.md) and [`docs/adr/0001-product-boundary.md`](docs/adr/0001-product-boundary.md) before substantive work. The task board in `agent_plan.md` is the live record of what is claimed, blocked and done.

## Prerequisites

| Tool | Version | Where it is pinned |
|---|---|---|
| Node.js | 24 | [`.github/workflows/quality.yml`](.github/workflows/quality.yml) |
| pnpm | 11.11.0 | `packageManager` in [`package.json`](package.json) |
| Go | 1.23 | [`go.mod`](go.mod), [`Dockerfile`](Dockerfile) |
| Ruby | 3.4 | [`.github/workflows/quality.yml`](.github/workflows/quality.yml) — only for `scripts/validate.rb`, which needs no gems |

## Local setup

```sh
git clone https://github.com/stanleyHayes/ghanacalendar.git
cd ghanacalendar
pnpm install --frozen-lockfile
```

Run the two surfaces separately, both from the repository root:

```sh
pnpm dev            # public calendar and sandbox on http://localhost:3000
go run ./cmd/api    # REST/GraphQL API on http://localhost:8080 (set PORT to change)
```

The Go service reads `data/holidays.json` relative to the working directory, so start it from the repository root. A quick smoke check:

```sh
curl -s "http://localhost:8080/v1/working-day?date=2026-07-03"
curl -s "http://localhost:8080/v1/working-days-between?start=2026-04-01&end=2026-04-08"
```

## Verification

Run the same sequence CI runs before opening a pull request. All six must pass.

```sh
ruby scripts/validate.rb
go test ./...
go vet ./...
pnpm typecheck
pnpm test
pnpm build
```

There is no lint script in this repository; do not add one to a pull request that is about something else. `ruby scripts/validate.rb` checks that the required governance files exist, that the source register is non-empty, and that no unresolved template token or private key has been committed.

> Note: the validator reads Markdown with Ruby's default external encoding. On a machine whose default locale is not UTF-8 — an older system Ruby on macOS, for example — it aborts with `invalid byte sequence in US-ASCII` as soon as a documentation file contains an en dash. Run it under Ruby 3.4 as CI does, or prefix the command with `RUBYOPT="-E UTF-8"`.

Tests live in [`internal/calendar/calendar_test.go`](internal/calendar/calendar_test.go) (engine behaviour: official 2026 fixtures, weekend/holiday determinism, historical rules, pending-movable handling, arithmetic) and [`tests/dataset.test.mjs`](tests/dataset.test.mjs) (dataset invariants: timezone, unique identifiers, null pending dates). A change to the engine or the dataset needs a test that would have failed before it.

## Commits and pull requests

This repository uses Conventional Commits with a lowercase imperative subject, matching the existing history:

```text
feat: align Calendar UI and social metadata
docs: record GhanaCalendar beta release
chore: ignore Vercel project metadata
```

Use `feat`, `fix`, `docs`, `chore`, `test` or `refactor`. Keep the subject short and specific to this product. Append `[skip ci]` only for documentation-only commits that cannot affect the build.

A pull request should state:

- the scope and which surface it touches (dataset, engine, API, web, contracts, SDK, docs);
- the source or ADR references behind any factual change;
- the verification you ran, including output for anything the six commands above do not cover;
- migration and rollback impact for a contract change;
- any external gate it depends on (for example an unpublished Ministry notice).

Breaking changes to [`contracts/openapi.yaml`](contracts/openapi.yaml) or [`contracts/graphql/schema.graphql`](contracts/graphql/schema.graphql) need a version decision recorded in the pull request. Stable identifiers in `data/holidays.json` must not be reused for a different occurrence.

## Proposing a data correction

Holiday dates are public legal facts. Corrections need evidence, not recollection.

Open an issue or pull request that gives:

1. the affected stable `id` from [`data/holidays.json`](data/holidays.json) — for example `republic-2026`;
2. the current stored value and the proposed value, including `observedDate` where a substitution applies;
3. the authoritative source: the Ministry of the Interior notice, gazette or Act, with a URL and its publication date;
4. whether the correction changes a historical record — if so, say which year's answers move and why the change is retrospective rather than prospective;
5. a new or updated entry in [`docs/governance/source-register.json`](docs/governance/source-register.json) if the evidence is a source not already registered, including `licence`, `publicationDecision` and `reviewedAt`.

Rules that are not negotiable:

- A movable holiday that has not been officially announced stays `status: "pending"` with `date: null`. Do not compute, estimate or infer it.
- An amendment applies prospectively. Historical occurrences keep the names and dates that were in force at the time.
- Commemorative days are not non-working days and must not be added as holidays.
- Scraped or aggregated third-party lists are not acceptable evidence on their own. Cite the authority.
- Source documents are linked, not redistributed. Do not commit a copy of a Ministry PDF.

Automation may draft a correction; a human reviewer approves canonical publication. The reviewed draft → review → publish workflow with an audit trail is still a stable gate — see [ROADMAP.md](ROADMAP.md).

## Review expectations

- Claim one path-bounded task in `agent_plan.md` before substantive work, so two contributors do not rewrite the same file.
- Expect review on evidence quality first, then correctness, then style. A change that is right but unsourced will be held.
- Status words in documentation must stay honest: proposed, building, beta, stable, externally blocked, retired, deferred. Do not describe something as live until deployment evidence exists in [`docs/runbooks/release-evidence.md`](docs/runbooks/release-evidence.md).
- Never claim or imply government endorsement, affiliation or official status.
- Never commit credentials, provider environment files, private exports or personal data. Report security issues privately per [`SECURITY.md`](SECURITY.md) rather than in a public issue.

## Good first contributions

Real gaps, drawn from the task board in [`agent_plan.md`](agent_plan.md) and the beta limitations in [`docs/runbooks/release-evidence.md`](docs/runbooks/release-evidence.md):

1. **Complete the OpenAPI response schemas.** [`contracts/openapi.yaml`](contracts/openapi.yaml) currently documents responses with a description only, and `/v1/next-working-day`, `/v1/previous-working-day`, `/v1/add-working-days` and `/v1/working-days-between` declare no query parameters at all. Add component schemas and parameters that match what [`cmd/api/main.go`](cmd/api/main.go) actually returns.
2. **Add automated REST ↔ GraphQL parity tests.** Parity is currently verified by hand during release smoke. A Go test that asserts `/v1/working-day` and the GraphQL `isWorkingDay` return the same decision and the same occurrence would make that a gate instead of a ritual.
3. **Harden the ICS export.** The generator in `cmd/api/main.go` emits `VEVENT` records without `DTSTAMP`, which RFC 5545 requires. Add it, add `DTEND` for the all-day form, and add a test that the 2026 export parses.
4. **Publish the TypeScript client.** [`sdk/typescript/index.ts`](sdk/typescript/index.ts) exists in source but has no package manifest, build or registry release (task `CAL-2.2`). Proposing the packaging shape is a useful first design contribution.
5. **React Query hooks package.** A pending beta-acceptance item: thin hooks over the client surface, with the caching policy written down.
6. **Draft-only source watcher.** A checker that notices when a registered Ministry URL changes and opens a draft — never a canonical update. The draft-only constraint is the whole point; see the non-goals in `agent_plan.md`.

Items 4, 5 and 6, plus the admin review workflow, are stable gates. Please comment on the task board before starting one, as they need a design decision agreed first.

## Licensing

Unless explicitly stated otherwise, contributions intentionally submitted for inclusion are provided under the Apache License 2.0 on the licence's inbound=outbound terms — see [`LICENSE`](LICENSE) and [`NOTICE`](NOTICE). Contributions must not include third-party data or documents without recorded permission.
