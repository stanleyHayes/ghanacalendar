# GhanaCalendar roadmap

This roadmap is directional, not a commitment. It contains no dates, because none are promised. The authoritative, machine-readable state — the live task board with owners, dependencies and evidence — is [`agent_plan.md`](agent_plan.md); where the two disagree, the ledger wins.

Current lifecycle state: **public beta**. Web and API are live and verified. Stable is blocked on the gates listed below.

## Now — shipped in the current beta

Items the ledger records as done, with the evidence that supports them.

**Evidence and definition**

- [x] Ghana Ministry of the Interior annual notices for 2024, amended 2025, and 2026 reviewed.
- [x] Public Holidays and Commemorative Days Act, 2001 (Act 601) recorded as the legal basis.
- [x] Source documents linked rather than redistributed; facts transcribed with attribution and an explicit rights note — [`docs/governance/source-register.json`](docs/governance/source-register.json).
- [x] The 2025 amendment modelled prospectively; 2024 names and observed dates remain historical.
- [x] Unannounced 2026 Eid-ul-Fitr, Shaqq Day and Eid-ul-Adha held as pending with null dates.

**Product**

- [x] Versioned 2024–2026 occurrence dataset with source IDs — [`data/holidays.json`](data/holidays.json), version `2026.09.01`.
- [x] Working-day engine: `isWorkingDay`, next/previous working day, add working days, end-exclusive days between — [`internal/calendar/calendar.go`](internal/calendar/calendar.go).
- [x] Official 2026 fixture, weekend/holiday determinism, historical, pending-movable and arithmetic tests — [`internal/calendar/calendar_test.go`](internal/calendar/calendar_test.go).
- [x] REST JSON plus CSV and ICS exports — [`cmd/api/main.go`](cmd/api/main.go).
- [x] Constrained GraphQL `holidays` and `isWorkingDay` operations — [`contracts/graphql/schema.graphql`](contracts/graphql/schema.graphql).
- [x] TypeScript client surface, in source — [`sdk/typescript/index.ts`](sdk/typescript/index.ts).
- [x] Responsive public calendar and in-browser working-day sandbox, with product favicon, manifest and canonical/social metadata — [`app/`](app).

**Release**

- [x] Web production release on the canonical host, UI/SEO/TLS verified.
- [x] API production release with a verified custom domain.
- [x] Production smoke and a retained rollback path — [`docs/runbooks/release-evidence.md`](docs/runbooks/release-evidence.md).

## Next — the gates that block stable

Declaring stable requires all of the following. None is started as a claimed task at the time of writing.

| Gate | Ledger reference | Blocking dependency |
|---|---|---|
| Draft → review → publish admin workflow with an audit trail | `CAL-3.2`, status Pending, owner unassigned | `CAL-0.1` source and legal review (done). Beta remains source-file reviewed until this exists |
| React Query hooks package and package-registry release | Follow-on to `CAL-2.2` ("Done locally", registry publication pending) | `CAL-2.1` interfaces (done); a packaging and versioning decision |
| Automated source watcher that creates drafts only | Beta acceptance item, unchecked | The admin draft workflow (`CAL-3.2`) — a watcher must never publish a canonical update |
| Security and load evidence | Stable release rule | A live API (`CAL-4.2`, done); a named test and threshold set |
| Named ongoing operations | Stable release rule | An acknowledged owner for alert routes, per [`docs/runbooks/operations.md`](docs/runbooks/operations.md) |

## Later — deferred or externally gated

- **The three pending 2026 movable holidays.** Eid-ul-Fitr, Shaqq Day and Eid-ul-Adha stay `date: null` until the Ministry of the Interior publishes confirmation. This is externally gated, not scheduled, and the dates will not be estimated in the meantime.
- **Years beyond 2026.** The API accepts 2024–2026 only. A year is added when its annual notice exists and has been reviewed into the source register — never before.
- **Contract completeness.** [`contracts/openapi.yaml`](contracts/openapi.yaml) documents responses by description rather than schema, and the arithmetic endpoints declare no parameters. Fuller schemas and drift detection follow the contract expectations in [`contracts/README.md`](contracts/README.md).
- **Automated REST ↔ GraphQL parity fixtures.** Parity is currently checked by hand in the release smoke suite; making it an automated gate is pending.
- **Operations runbook reconciliation.** [`docs/runbooks/operations.md`](docs/runbooks/operations.md) still describes a pre-production baseline and needs rewriting against the live service and the signals now in place.
- **Shared portfolio primitives.** Per [ADR-0001](docs/adr/0001-product-boundary.md), small primitives may be repeated until two proven consumers justify a versioned shared package. Extraction is deliberately deferred until then.

## Explicitly out of scope

These are stated non-goals in [`agent_plan.md`](agent_plan.md), not backlog items. A pull request that implements one will be declined.

- Predicting movable Islamic holidays before an official notice.
- Treating commemorative days as non-working days.
- Rewriting historical answers after a law or naming change.
- Accepting unreviewed web scraping as a canonical update.
- Building shared portfolio identity, billing, or a cross-product database.
- Proxying government transactions, credentials, forms or payments, or claiming government endorsement.

## How this changes

Status words in this repository are limited to: proposed, building, beta, stable, externally blocked, retired, deferred. A status only moves when evidence is recorded — deployment IDs, smoke results and rollback steps in [`docs/runbooks/release-evidence.md`](docs/runbooks/release-evidence.md), and the corresponding row updated in [`agent_plan.md`](agent_plan.md). To propose a change of direction, open a pull request against this file with the reasoning, or an ADR under [`docs/adr/`](docs/adr) if it changes the product boundary. See [CONTRIBUTING.md](CONTRIBUTING.md).
