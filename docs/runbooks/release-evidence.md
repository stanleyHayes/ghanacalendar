# GhanaCalendar beta release evidence

Verified: 2026-09-01<br>
Release commit: `54d7a7d1fbd6963437621a6e66b9d3a2a0b567c9`<br>
Dataset: `2026.09.01`

## Web

- Provider/project: Vercel `hayfordstanleys-projects/ghanacalendar`
- Deployment: `dpl_Ca1QeBj5KXhJBiqbchTKfx1L6hjB`
- Immutable URL: `ghanacalendar-cm8z9ig43-hayfordstanleys-projects.vercel.app`
- Canonical URL: <https://calendar.digitalghana.dev>
- `/`, `/sitemap.xml`, and `/robots.txt`: HTTP 200

## API

- Provider/service: Render `srv-dabdeaf40ujc73aji6n0`
- Deploy: `dep-dabdeb740ujc73aji93g`
- Provider URL: `ghanacalendar-api.onrender.com`
- Canonical URL: <https://api-calendar.digitalghana.dev>
- Custom domain: verified by Render
- `/health`: HTTP 200 with dataset version `2026.09.01`

## Production smoke

- `2026-07-03` resolves as non-working and identifies the observed Republic Day occurrence.
- Adding two working days to `2026-04-02` returns `2026-04-08` across the Good Friday/Easter Monday interval.
- Working days in `[2026-04-01, 2026-04-08)` returns `3`.
- GraphQL `isWorkingDay` returns the same decision and occurrence as REST.
- 2026 ICS export contains versioned VEVENT records.
- 2026 CSV export contains the declared schema and source IDs.
- GitHub Quality run `33515268948` passed.

## TLS and rollback

Both canonical hosts terminate valid HTTPS. Rollback keeps provider URLs immutable: Vercel can promote the prior deployment, while Render can redeploy the recorded commit/deploy. DNS remains provider-neutral and can be repointed only after the replacement passes the same smoke suite.

## Beta limitations

- The admin draft/review/publish workflow and audit trail are not yet implemented.
- The TypeScript client exists in source but is not yet published to a package registry; React Query hooks are pending.
- Automated Ministry source monitoring is pending.
- 2026 Eid-ul-Fitr, Shaqq Day, and Eid-ul-Adha remain deliberately undated until official confirmation.
