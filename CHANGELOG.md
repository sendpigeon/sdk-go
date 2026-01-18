# Changelog

## 0.5.0

- Add Contacts API (`Contacts.List`, `Create`, `Batch`, `Get`, `Update`, `Delete`, `Unsubscribe`, `Resubscribe`, `Stats`, `Tags`)
- Add Broadcasts API (`Broadcasts.List`, `Create`, `Get`, `Update`, `Delete`, `Send`, `Schedule`, `Cancel`, `Test`, `Recipients`, `Analytics`)
- Broadcast targeting: `IncludeTags` and `ExcludeTags` options

## 0.4.0

- Per-email tracking: `Tracking: &TrackingOptions{Opens: &opens, Clicks: &clicks}` in send requests
- Response `Warnings` field for non-fatal issues (e.g., tracking disabled at org level)
- Updated `TrackingDefaults` to use `TrackingEnabled` master toggle

## 0.3.0

- Add Suppressions API (`Suppressions.List`, `Suppressions.Delete`)

## 0.2.0

- Add `email.opened` and `email.clicked` webhook events
- Add typed `WebhookPayload` and `WebhookPayloadData` structs
- Add `ParseWebhookPayload()` helper function
- Add webhook event constants

## 0.1.0

- Initial release
- Send emails (single + batch)
- Templates API
- Domains API
- API Keys API
- Webhook signature verification
