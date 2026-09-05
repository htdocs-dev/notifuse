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
