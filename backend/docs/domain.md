# Akritas Domain Notes

## Incident, Investigation and Evidence

- An `Incident` groups observed occurrences for one Project and fingerprint.
- An `Investigation` belongs to exactly one Incident and records QVAC conclusions once completed.
- `Evidence` belongs to an Investigation and stores observed, redacted material used during investigation.
- GitHub Issue publication uses persisted Evidence as observed facts and Investigation fields as QVAC-generated conclusions. The Issue body must keep these sections visibly separated.

## GitHub Issue Reference

- A `GitHubIssueReference` records the GitHub Issue created for a completed Investigation.
- Idempotency is scoped to Investigation: one Investigation can publish at most one IssueReference.
- An Incident may have multiple Investigation attempts over time; Incident detail projects the latest IssueReference for that Incident.
- PostgreSQL enforces that the IssueReference `incident_id` and `investigation_id` are consistent: the Investigation referenced by `investigation_id` must belong to the same Incident recorded in `incident_id`.

## Evidence Safety

- Redaction removes secret values, not only known prefixes.
- Secret markers are deterministic and do not preserve partial secret material.
- Normal prose mentioning security terms such as token or password is preserved when no associated value is present.
- Redaction preserves valid UTF-8 and bounded content sizes before persistence/publication.

