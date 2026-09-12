# GhanaCalendar

Evidence-backed Ghana public holidays and deterministic working-day arithmetic, served over REST, GraphQL, CSV and ICS — with every date traceable to the notice it came from.

[![Licence: Apache-2.0](https://img.shields.io/badge/licence-Apache--2.0-blue)](LICENSE)
[![Quality](https://img.shields.io/github/actions/workflow/status/stanleyHayes/ghanacalendar/quality.yml?branch=main&label=quality)](.github/workflows/quality.yml)
[![Web](https://img.shields.io/badge/web-calendar.digitalghana.dev-brightgreen)](https://calendar.digitalghana.dev)
[![API](https://img.shields.io/badge/api-api--calendar.digitalghana.dev-brightgreen)](https://api-calendar.digitalghana.dev/health)
[![Runtime](https://img.shields.io/badge/runtime-Go%201.23%20%C2%B7%20Next.js%2016-informational)](go.mod)

Independent open-source public-interest infrastructure from the [Digital Ghana](https://digitalghana.dev) portfolio. Timezone is always `Africa/Accra`. Dataset covers 2024–2026 at version `2026.09.01`.

## Live now

| Surface | URL | State |
|---|---|---|
| Public calendar and in-browser working-day sandbox | <https://calendar.digitalghana.dev> | Live |
| REST + GraphQL + exports | <https://api-calendar.digitalghana.dev> | Live |

The API runs on a free Render instance. A cold start can take **30–60 seconds** on the first request; subsequent calls are fast.

```sh
# Is 3 July 2026 a working day in Ghana?
curl -s "https://api-calendar.digitalghana.dev/v1/working-day?date=2026-07-03"
```

```json
{"date":"2026-07-03","holiday":{"id":"republic-2026","name":"Republic Day","date":"2026-07-01","observedDate":"2026-07-03","kind":"statutory","status":"confirmed","sourceId":"mint-2026-calendar"},"timezone":"Africa/Accra","workingDay":false}
```

```sh
# Working days in [2026-04-01, 2026-04-08) — the Good Friday / Easter Monday interval
curl -s "https://api-calendar.digitalghana.dev/v1/working-days-between?start=2026-04-01&end=2026-04-08"
# {"end":"2026-04-08","endExclusive":true,"start":"2026-04-01","workingDays":3}

# Health and the dataset version the service is actually running
curl -s "https://api-calendar.digitalghana.dev/health"
# {"status":"ok","version":"2026.09.01"}
```

Browser `fetch` from other origins is not supported: the API sends `Access-Control-Allow-Origin: https://calendar.digitalghana.dev` only ([`cmd/api/main.go`](cmd/api/main.go)). Call it server-side, or run your own instance.

## The problem this solves

### For developers

Ghana's public holidays are published by the Ministry of the Interior as an annual notice — a PDF or a web page, written for people, not for programs. There is no official machine-readable feed. So every team that needs a Ghana business calendar ends up doing the same work:

- **Transcribing a notice by hand, once a year, into a constants file.** It is a five-minute job that silently rots. Nobody diffs it when the Ministry amends a date mid-year.
- **Guessing the movable Islamic holidays.** Eid-ul-Fitr, Shaqq Day and Eid-ul-Adha are confirmed close to the date. A generic holidays library will happily compute an astronomical estimate and hand you a date the government did not declare. Payroll runs on that estimate.
- **Re-implementing substitution rules.** When the Ministry declares a substitute day it is recorded as `observedDate`. Usually the holiday fell on a weekend — Boxing Day 2026 is Saturday 26 December, observed Monday 28 December; Constitution Day 2024 was Sunday 7 January, observed Monday 8 January — but not always: the 2026 notice observes Republic Day (Wednesday 1 July) on Friday 3 July. Two dates, one holiday, and both must be non-working.
- **Getting history wrong after an amendment.** The 2025 amendment changed names and the calendar shape; the 2024 record still has Founders' Day on 4 August and Kwame Nkrumah Memorial Day on 21 September. A library that applies today's rules retroactively will give you the wrong answer for a 2024 payroll audit or an SLA dispute.
- **Writing the same off-by-one loop again.** "Add two working days", "days between", "next working day" — each is a small loop, and each team gets the inclusive/exclusive boundary wrong at least once.

The cost is not dramatic; it is quiet. A wrong non-working day shifts a value date, misprices an SLA credit, or books a delivery into a closed warehouse.

GhanaCalendar gives you, instead:

| You need | You get |
|---|---|
| A holiday list you can diff | [`data/holidays.json`](data/holidays.json) — 39 versioned occurrences across 2024–2026, each with a stable `id` and a `sourceId` |
| A boolean, not a guess | `GET /v1/working-day` — weekends, statutory dates and confirmed observed dates all resolve to `workingDay: false` |
| Honest unknowns | Unconfirmed movable holidays stay `status: "pending"` with `date: null`. The service never invents a date |
| Correct history | 2024 answers keep 2024 rules. An amendment applies from its own effective period forward |
| The boring arithmetic, done once | `next-working-day`, `previous-working-day`, `add-working-days`, and an explicitly end-exclusive `working-days-between` |
| Exports for the rest of the org | JSON for services, CSV for analysts, ICS for the calendars people actually look at |
| A typed client | [`sdk/typescript/index.ts`](sdk/typescript/index.ts) |

### For the community

Public holiday dates are public legal facts. They decide when courts sit, when banks clear, when wages are due, and when a statutory deadline lands. That those facts exist only as a yearly PDF is a small, real gap in public infrastructure — and it is usually filled by closed commercial APIs or by unattributed scraped lists nobody can audit.

This project keeps that layer open and checkable:

- **Provenance, not assertion.** Every occurrence names the notice or Act behind it in [`docs/governance/source-register.json`](docs/governance/source-register.json). You can verify a date against the authority without trusting us.
- **Reproducible.** The dataset is a flat, versioned JSON file in git. Any answer the API gives can be recomputed offline from the same file and the engine in [`internal/calendar/calendar.go`](internal/calendar/calendar.go).
- **No lock-in.** Apache-2.0, no key, no account, no rate-limited free tier. Vendor the JSON if you prefer.
- **Visible limits.** Pending dates are shown as pending on the public site rather than filled in to look complete.
- **Independent.** Nothing here proxies a government transaction, form, credential or payment. Official business stays on official sites.

### What this is not

Drawn from the stated non-goals in [`agent_plan.md`](agent_plan.md):

- Not a predictor of movable Islamic holidays before an official notice.
- Not a commemorative-day calendar — commemorative days are not treated as non-working days.
- Not a service that rewrites historical answers after a law or naming change.
- Not fed by unreviewed web scraping; no scraper may publish a canonical update.
- Not a government system, and not endorsed by or affiliated with the Government of Ghana.
- Not comprehensive beyond 2024–2026. Requests outside that range are rejected, not extrapolated.

## Quickstart

Prerequisites: Node 24, pnpm 11.11.0, Go 1.23, Ruby 3.4 (for the foundation validator only).

```sh
git clone https://github.com/stanleyHayes/ghanacalendar.git
cd ghanacalendar
pnpm install --frozen-lockfile

pnpm dev            # public calendar + sandbox on http://localhost:3000
go run ./cmd/api    # REST/GraphQL API on http://localhost:8080
```

Run `go run ./cmd/api` from the repository root — it loads `data/holidays.json` relative to the working directory. Override the port with `PORT`. Then `curl -s "http://localhost:8080/v1/next-working-day?date=2026-12-25"` returns `{"date":"2026-12-25","result":"2026-12-29","timezone":"Africa/Accra"}` — Christmas Day, the weekend, and the observed Boxing Day on Monday 28 December are all skipped.

## Usage

Contracts: [`contracts/openapi.yaml`](contracts/openapi.yaml) and [`contracts/graphql/schema.graphql`](contracts/graphql/schema.graphql).

| Method | Path | Query parameters | Returns |
|---|---|---|---|
| `GET` | `/health` | — | `status`, dataset `version` |
| `GET` | `/v1/holidays` | `year` (2024–2026, required), `includePending`, `format` (`csv`, `ics`) | Occurrences for the year |
| `GET` | `/v1/working-day` | `date` (`YYYY-MM-DD`) | Working-day decision plus the matched occurrence |
| `GET` | `/v1/next-working-day` | `date` | Next working date |
| `GET` | `/v1/previous-working-day` | `date` | Previous working date |
| `GET` | `/v1/add-working-days` | `date`, `count` (−3660…3660) | Shifted working date |
| `GET` | `/v1/working-days-between` | `start`, `end` | End-exclusive working-day count |
| `POST` | `/graphql` | body: `query`, `variables` | `holidays(year:)` and `isWorkingDay(date:)` only |

Occurrence fields: `id`, `name`, `date` (nullable), `observedDate` (optional), `kind` (`statutory` or `additional`), `status` (`confirmed` or `pending`), `sourceId`.

```sh
curl -s "https://api-calendar.digitalghana.dev/v1/holidays?year=2026&includePending=true"
```

```json
{"holidays":[{"id":"new-year-2026","name":"New Year's Day","date":"2026-01-01","kind":"statutory","status":"confirmed","sourceId":"mint-2026-calendar"}],"timezone":"Africa/Accra","version":"2026.09.01"}
```

Abbreviated: the live response lists every 2026 occurrence, including the three pending ones with `"date":null`.

```sh
curl -s -X POST https://api-calendar.digitalghana.dev/graphql \
  -H 'content-type: application/json' \
  -d '{"query":"query($date:String!){isWorkingDay(date:$date){date workingDay timezone}}","variables":{"date":"2026-07-03"}}'
```

GraphQL returns the same decision and the same occurrence as REST — that parity is part of the release smoke suite.

Exports: `?format=csv` emits the header `id,name,date,observedDate,kind,status,sourceId`; `?format=ics` emits `VEVENT` records with `UID:<id>@calendar.digitalghana.dev` and skips undated pending holidays.

TypeScript:

```ts
import { GhanaCalendarClient } from "./sdk/typescript/index";
const calendar = new GhanaCalendarClient();
await calendar.isWorkingDay("2026-07-03");
const { result } = await calendar.addWorkingDays("2026-04-02", 2); // result === "2026-04-08"
```

The client is source-only for now; it is not yet published to a package registry.

## Data and provenance

| Property | Value |
|---|---|
| Coverage | 2024, 2025, 2026 |
| Records | 39 occurrences — 12 dated in 2024, 13 in 2025, 11 dated plus 3 pending in 2026 |
| Dataset version | `2026.09.01` |
| Sources reviewed | 2026-09-01 |
| Timezone | `Africa/Accra`, enforced at load time |

Sources are recorded in [`docs/governance/source-register.json`](docs/governance/source-register.json) against [`source-register.schema.json`](docs/governance/source-register.schema.json): the Ministry of the Interior notices for 2024, the amended 2025 calendar, and 2026, plus the Public Holidays and Commemorative Days Act, 2001 (Act 601) as the legal basis.

**Licence boundary.** Dates are treated as public legal facts, transcribed with attribution. The source documents themselves are linked, not redistributed, and their rights are unaffected by inclusion here — see [`NOTICE`](NOTICE). Original code and the compiled dataset structure are Apache-2.0.

**Corrections.** Open an issue or pull request with the affected `id`, the current value, the proposed value, the authoritative source and its publication date, and whether the change touches historical records. Automation may draft a correction; a human reviewer approves publication. The reviewed draft → publish workflow with an audit trail is a stable gate that is not yet built — see [ROADMAP.md](ROADMAP.md). Full contributor guidance is in [CONTRIBUTING.md](CONTRIBUTING.md).

## Project layout

```
app/                  Next.js 16 public calendar, exports and working-day sandbox
cmd/api/              Go HTTP service: REST, GraphQL, CSV and ICS
internal/calendar/    Dataset loader and the working-day engine
data/holidays.json    Versioned 2024-2026 occurrence dataset
contracts/            OpenAPI 3.1 and GraphQL schema
sdk/typescript/       Typed client (unpublished)
docs/adr/             Architecture decision records
docs/governance/      Source register and its schema
docs/runbooks/        Operations baseline and release evidence
scripts/validate.rb   Product-foundation validator (no gems required)
tests/                Node dataset invariant tests
infra/                Vercel headers configuration
```

## Verification

These are the commands CI runs, in [`.github/workflows/quality.yml`](.github/workflows/quality.yml):

```sh
ruby scripts/validate.rb   # required files, source register, template/secret scan
go test ./...              # fixtures, weekend/holiday determinism, history, pending-movable, arithmetic
go vet ./...
pnpm typecheck
pnpm test                  # dataset invariants
pnpm build
```

There is no lint script in this repository. `scripts/validate.rb` needs no gems, but it reads Markdown with Ruby's default external encoding: on a machine whose default locale is not UTF-8 it aborts with `invalid byte sequence in US-ASCII`. Run it under Ruby 3.4 as CI does, or prefix the command with `RUBYOPT="-E UTF-8"`.

Immutable deployment IDs, the production smoke results and the rollback path are in [`docs/runbooks/release-evidence.md`](docs/runbooks/release-evidence.md).

## Status and roadmap

**Public beta.** Web and API are live and verified; the engine, interfaces, exports, SDK source and public site are feature-complete for beta. The dataset is a deliberately bounded 2024–2026 subset and is not complete — three 2026 movable holidays remain undated pending official confirmation. Stable is blocked on the reviewed admin publication workflow, a draft-only source watcher, package-registry releases, and mature operational evidence.

Directional plan: [ROADMAP.md](ROADMAP.md). Machine-readable state and the live task board: [`agent_plan.md`](agent_plan.md).

## Contributing, conduct, security, licence

- [CONTRIBUTING.md](CONTRIBUTING.md) — setup, verification, commit conventions, data corrections
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) — Contributor Covenant v2.1; reports to the private channel documented there
- [SECURITY.md](SECURITY.md) — report privately to the private channel documented there, never in a public issue
- [AGENTS.md](AGENTS.md) — coordination rules for automated contributors
- [LICENSE](LICENSE) and [NOTICE](NOTICE) — Apache License 2.0; inbound=outbound

## Independence

GhanaCalendar is an independent open-source project. It is **not** operated by, endorsed by, or affiliated with the Government of Ghana, the Ministry of the Interior, or any public authority. It does not proxy government transactions, credentials, forms or payments. Where an official answer is required, consult the authority's own published notice — the links are in the source register.
