# GhanaCalendar

`GhanaCalendar` is an independent Digital Ghana public-infrastructure product for evidence-backed public holidays and deterministic working-day calculations.

- Web and sandbox: <https://calendar.digitalghana.dev>
- REST, GraphQL and exports: <https://api-calendar.digitalghana.dev>
- Timezone: `Africa/Accra`
- Dataset: 2024–2026, version `2026.09.01`

## Current beta

The beta includes REST/GraphQL working-day operations, JSON/CSV/ICS holiday exports, a TypeScript client, a public register, and a browser calculator. Movable dates remain pending until the Ghana Ministry of the Interior publishes confirmation.

Stable status remains blocked on the reviewed admin publication workflow, draft-only source watcher, package releases, and mature operational evidence. See `agent_plan.md` and `docs/runbooks/release-evidence.md`.

## Verification

Run `ruby scripts/validate.rb`, `go test ./...`, `go vet ./...`, `pnpm typecheck`, `pnpm test`, and `pnpm build`.
