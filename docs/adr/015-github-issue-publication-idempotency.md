# ADR-015 - GitHub Issue publication and idempotency

## Status

Accepted

## Context

H4 turns a completed QVAC Investigation into auditable content and a GitHub Issue. The Issue is created in an external system, while Akritas stores workflow state in PostgreSQL. A single transaction cannot cover both systems, and retrying blindly after an interruption can create duplicate remote Issues.

## Decision

Each completed Investigation may publish at most one GitHub Issue. `GitHubIssueReference` is a durable PostgreSQL entity linked to both Incident and Investigation, with a unique Investigation reference and a unique `(repository, issue_number)` pair.

Incident detail exposes a singular `github_issue_reference`, but that is only a projection of the most recently created Issue for the Incident. The Incident itself is not unique-constrained to one Issue because repeated completed Investigations can each produce their own audit record.

The workflow uses two short PostgreSQL transactions around the external call:

1. Persist QVAC completion and move the Incident to `publishing_issue`.
2. Close the transaction and call GitHub.
3. Persist `GitHubIssueReference`, update Incident projection/outcome, and finish the Operation.

GitHub Issue bodies include a safe HTML marker containing the Investigation UUID. This is not used for success today, but enables future reconciliation when GitHub accepted the Issue and local persistence failed.

For `resolution_status = requires_human`, the Incident reaches `completed/requires_human` after IssueReference persistence. For `resolution_status = fixable`, H4 leaves the Incident in `publishing_issue`; H5 owns the transition to `remediating`.

## Consequences

- A preexisting IssueReference for the same Investigation is the idempotency boundary and prevents a second publish.
- A GitHub failure leaves the Investigation completed, fails the Incident with `issue_publication_failed`, fails the Operation with a safe public message, and creates no IssueReference.
- A PostgreSQL failure after GitHub accepted the Issue does not fabricate local success; the workflow fails explicitly and relies on the marker for H6 reconciliation.
- Startup recovery must never republish a completed Investigation without a durable IssueReference.
- External GitHub calls remain outside database transactions, preserving the short-transaction rule from ADR-014.