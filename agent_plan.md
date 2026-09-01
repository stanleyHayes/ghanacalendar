# GhanaCalendar execution ledger

Last updated: 2026-09-01  
Status: Public beta live; stable gates remain
Canonical hosts: `calendar.digitalghana.dev`, `api-calendar.digitalghana.dev`

## Product boundary

GhanaCalendar provides evidence-backed Ghana holiday occurrences and deterministic working-day calculations. It is independent open-source infrastructure and does not claim government endorsement. The timezone is always `Africa/Accra`.

### Non-goals for beta

- Predicting movable Islamic holidays before an official notice.
- Treating commemorative days as non-working days.
- Rewriting historical answers after a law or naming change.
- Accepting unreviewed web scraping as a canonical update.
- Building shared portfolio identity, billing, or a cross-product database.

## Definition evidence

- [x] Ghana Ministry of the Interior annual notices for 2024, amended 2025, and 2026 reviewed.
- [x] Act 601 legal basis recorded.
- [x] Source documents linked rather than redistributed; facts transcribed with attribution and an explicit rights note.
- [x] 2025 amendment modeled prospectively; 2024 names and observed dates remain historical.
- [x] Unannounced 2026 Eid-ul-Fitr, Shaqq Day, and Eid-ul-Adha remain pending with null dates.

## Beta acceptance

- [x] Versioned 2024–2026 occurrence dataset and source IDs.
- [x] `isWorkingDay`, next/previous/add working days, and end-exclusive days-between engine.
- [x] Official 2026 fixture, overlap, historical, pending-movable, and arithmetic tests.
- [x] REST JSON plus CSV and ICS exports.
- [x] Constrained GraphQL `holidays` and `isWorkingDay` operations.
- [x] TypeScript client surface.
- [x] Responsive public calendar and in-browser working-day sandbox.
- [ ] Draft → review → publish admin workflow with audit trail.
- [ ] React Query hooks package and package-registry release.
- [ ] Automated source watcher that creates drafts only.

## Live task board

| ID | Task | Status | Owner | Dependency | Evidence |
|---|---|---|---|---|---|
| CAL-0.1 | Source and legal review | Done | Codex | — | `docs/governance/source-register.json` |
| CAL-1.1 | Historical occurrence dataset | Done | Codex | CAL-0.1 | `data/holidays.json`; null pending dates |
| CAL-1.2 | Working-day engine | Done | Codex | CAL-1.1 | `go test ./...`, `go vet ./...` |
| CAL-2.1 | REST/GraphQL/export interfaces | Done | Codex | CAL-1.2 | Contracts and API handlers |
| CAL-2.2 | TypeScript SDK | Done locally | Codex | CAL-2.1 | `sdk/typescript/index.ts`; registry publication pending |
| CAL-3.1 | Public website and sandbox | Done locally | Codex | CAL-1.2 | typecheck, dataset tests, production build pass |
| CAL-3.2 | Admin review/publish workflow | Pending | Unassigned | CAL-0.1 | Required before stable; beta remains source-file reviewed |
| CAL-4.1 | Web production release | Done | Codex | CAL-3.1 | Vercel `dpl_Ca1QeBj5KXhJBiqbchTKfx1L6hjB`; canonical host HTTP 200 and TLS verified |
| CAL-4.2 | API production release | Done | Codex | CAL-2.1 | Render `srv-dabdeaf40ujc73aji6n0`, deploy `dep-dabdeb740ujc73aji93g`; custom domain verified |
| CAL-4.3 | Production smoke and rollback | Done | Codex | CAL-4.1, CAL-4.2 | Health, REST, GraphQL, CSV, ICS, web, sitemap and robots passed; provider URLs retained for rollback |

## Release rule

Beta may be declared only after web and API TLS/HTTP checks, REST/GraphQL parity fixtures, export smoke, immutable deployment IDs, and rollback instructions are recorded. Stable remains blocked on the admin review workflow, automated draft-only source monitoring, SDK/hooks publication, security/load evidence, and named ongoing operations.
