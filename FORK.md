# htdocs-dev fork of Notifuse

This fork tracks `upstream/main` (https://github.com/Notifuse/notifuse) and
carries a short list of patches. Upstream does not accept pull requests at the
moment (see the v40.0 changelog entry), so every patch lives here.

## How to update from upstream

```bash
git fetch upstream
git rebase upstream/main   # or: git merge upstream/main
go build ./... && go test ./...
```

Patches touch upstream files as little as possible. New behaviour sits in new
files; upstream files get one-line hooks. When a rebase conflicts, the list
below says what each patch needs to keep working.

## Patches

### 1. Spam complaints suppress the contact

Upstream marks the message as complained but leaves the contact `active` on
every list, so the next broadcast mails the complainer again.

- `internal/domain/contact.go`: `ContactRepository.MarkEmailsAsComplained`.
- `internal/repository/contact_postgres.go`: `MarkEmailsAsBounced` and
  `MarkEmailsAsComplained` share `markEmailsAsTerminalListStatus`.
- `internal/service/inbound_webhook_event_service.go`: `ProcessWebhook`
  collects complaint recipients and calls `MarkEmailsAsComplained`.
- `internal/domain/mocks/mock_contact_repository.go`: mock method (hand-added,
  same shape as mockgen output).

### 2. Disposable (throw-away) email block

- `pkg/disposable_emails`: `IsDisposableEmail` accepts a full address or a bare
  domain, lower-cases it and uses a binary search. Upstream compared the full
  address against a domain list, so the public subscribe check never matched.
- `internal/domain/workspace.go`: `WorkspaceSettings.BlockDisposableEmails`
  (`block_disposable_emails`), listed in `preservableWorkspaceSettingKeys`.
- `internal/service/contact_service.go`: `UpsertContact` and
  `BatchImportContacts` reject disposable addresses when the setting is on
  (`disposableEmailFilter`, one workspace read per call).
- Console: switch in `components/settings/GeneralSettings.tsx`, type in
  `services/api/workspace.ts`.

### 3. Resend to non-openers

A processed broadcast gets a "resend" action. It creates (or rebuilds) a
segment `nonopen_<broadcast id>` = sent the broadcast AND opened it 0 times,
then creates a draft broadcast on the same list, limited to that segment, with
the winning (or only) template. Nothing is sent until the draft is scheduled.

- `internal/domain/broadcast_resend.go` (new): request type, segment ID and
  tree recipe. Hook: `BroadcastService.ResendToNonOpeners` in the interface in
  `internal/domain/broadcast.go`.
- `internal/service/broadcast_resend.go` (new): service method and
  `SetSegmentService`. Hook: `segmentService` field in `broadcast_service.go`,
  one `SetSegmentService` call in `internal/app/app.go`.
- `internal/http/broadcast_resend_handler.go` (new): `POST
  /api/broadcasts.resendToNonOpeners`. Hook: one route line in
  `broadcast_handler.go`.
- `internal/domain/mocks/mock_broadcast_service.go`: mock method.
- Console: `services/api/broadcast.ts` (`resendToNonOpeners`) and a button in
  `pages/BroadcastsPage.tsx`.
- `openapi/paths/broadcasts.yaml` + `openapi/components/schemas/broadcast.yaml`, bundled into `openapi.json` with `make openapi-bundle`.

### 4. Dark mode in the console

- `console/src/themeMode.ts` (new): mode light/dark/system, localStorage key
  `notifuse_theme`, stamps `data-theme` on `<html>`.
- `console/src/contexts/ThemeContext.tsx` (new) and
  `console/src/components/ThemeSwitcher.tsx` (new): provider and the header button.
- `console/src/theme.ts` (new): the Ant Design theme, moved out of `App.tsx`, picks
  `darkAlgorithm` from the resolved theme. Hooks in `App.tsx`: `ThemeProvider` around
  `LocaleProvider`, `useTheme` in `AppContent`.
- `console/src/index.css`: `@custom-variant dark` on `data-theme`, `--nf-surface` /
  `--nf-border` variables, and a `:root[data-theme="dark"]` block that reverses the
  Tailwind grey scale and remaps `--color-white`. This is what turns every hard-coded
  `text-gray-500` / `bg-gray-50` / `bg-white` dark without editing components.
- `WorkspaceLayout.tsx`: the `#F9F9F9` / `#f0f0f0` literals became the two variables;
  the logo gets class `workspace-logo` (inverted in dark). `FilterBuilder.tsx`: two
  `text-white` pinned to `#ffffff`. `SetupWizard.tsx`: logo class.
- Known gaps: inline hex colours in charts, automation nodes and the email builder
  chrome still render their light values. Fix them file by file as they annoy you.

### 5. Dynamic segments and faster segment builds

Broadcast audiences and the contacts search used to read `contact_segments`,
which the build task fills 100 contacts at a time, so a new segment was empty
for hours on a large workspace. Now they run the segment's stored SQL inline.

- `internal/repository/segment_dynamic.go` (new): `segmentMembershipExpr` loads
  `generated_sql` / `generated_args` for the segment IDs, rewrites `$n` to
  squirrel `?` and returns `c.email IN (<segment sql>)` OR-ed across segments.
  A segment with no generated SQL falls back to its `contact_segments` rows.
  Hooks: three `if len(...Segments) > 0` blocks in `contact_postgres.go`
  (`GetContacts`, `GetContactsForBroadcast`, `CountContactsForBroadcast`).
  `contact_segments` still feeds automations, the per-contact chips and the
  segment page counts.
- Batch size 100 → 1000 in `segment_build_processor.go`,
  `contact_segment_queue_processor.go` and the three task states in
  `segment_service.go`. The queue task still runs 50 s per minute, so a
  workspace drains about 1000 queued contacts a minute instead of 100.

### 6. Console theme revamp

One token set in `console/src/index.css` drives everything: `--nf-page`,
`--nf-surface`, `--nf-elevated`, `--nf-border(-strong)`, `--nf-text(-2,-3)`,
`--primary` (mint `#12A66F` light / `#3ECF8E` dark), `--nf-accent-soft`,
`--nf-ring`. The `:root[data-theme="dark"]` block swaps the values and
reverses the Tailwind grey scale. `console/src/theme.ts` maps the same
variables onto antd tokens (backgrounds, borders, text, Menu, Button, Table,
Input, Tabs, Tag, Modal, Drawer) and sets Inter (from
`@fontsource-variable/inter`), 8/12 px radii, 34 px controls.

- `WorkspaceLayout.tsx`: outer `Layout` has class `workspace-page` (page
  gradient), `Header` has class `workspace-header` (blur); inner layouts are
  transparent. Version text uses `--nf-text-3`.
- `index.css` also carries the global polish: title tracking, focus ring, thin
  scrollbars, active nav bar, card shadow, uppercase table headers, primary
  button glow. `BaseNode.tsx`, `web_analytics/lib/types.ts` and the
  contacts table fixed columns read the variables instead of literals.
- `pages/ContactsPage.tsx`: the unsubscribed list chip uses antd's `default`
  colour (`gray` is not a preset and rendered as a white block in dark mode).

### 7. Contact tags

Tags are a flat JSON array of strings in `custom_json_1`
(`domain.ContactTagsField`). No schema change: the API reads and writes them as
a custom field, and segments match them with `in_array` on that field.

- `internal/domain/contact_tags.go` (new): `ContactTagsField`,
  `TagContactsRequest`, `ContactTagger` interface. Hook: `ContactTagger`
  embedded in `ContactRepository` and `ContactService` in `contact.go`.
- `internal/repository/contact_tags_postgres.go` (new): `AddContactTags`
  (sorted union, skips contacts that already have every tag) and
  `RemoveContactTags`, both one `UPDATE ... WHERE email = ANY($1)`.
- `internal/service/contact_tags.go` (new): contacts:write check, then the
  repository.
- `internal/http/contact_tags_handler.go` (new): `POST /api/contacts.tag` and
  `/api/contacts.untag`, body `{workspace_id, emails[], tags[]}`, reply
  `{success, updated}`. Hook: two route lines in `contact_handler.go`.
- Mocks: `mock_contact_repository.go`, `mock_contact_service.go`.
- Console: `components/contacts/TagChips.tsx` + `stringList.ts` (new); a JSON
  field that is a list of strings renders as chips in the contacts table
  (`JsonViewer`) and the contact drawer (`InlineEditableField`). Label the
  field "Tags" in workspace settings → custom field labels.
- `openapi/paths/contacts.yaml`, `openapi/components/schemas/contact.yaml`,
  `openapi/openapi.yaml`, bundled into `openapi.json`.
