# H6 Backend Demo Runbook

Task: `AKR-H6-BACKEND-HARDENING-DEMO`

Scope: backend only. The flow must stop after creating the Pull Request. Do not
merge, deploy, rollback or mutate production after PR creation.

## Fixture

Use `testdata/h6-demo-fixture`, an intentionally defective Go HTTP service. It
panics for repeated `/checkout?customer_id=abc` requests because the customer ID
is shorter than eight characters.

No secrets are required. Optional non-secret config:

```sh
AKRITAS_DEMO_DISCOUNT_PREFIX=H6
```

## Local Preparation

1. Create or reuse a non-production GitHub repository for the fixture.
2. Configure Akritas with a GitHub account and a Dokploy project pointing to
   that repository.
3. Run QVAC locally with the read-only repository tools enabled.
4. Use a workspace root controlled by Akritas; do not point cleanup at any path
   outside that root.

## Reproduce

```sh
cd testdata/h6-demo-fixture
go run .
curl "http://localhost:18080/checkout?customer_id=abc"
curl "http://localhost:18080/checkout?customer_id=abc"
curl "http://localhost:18080/checkout?customer_id=abc"
```

Expected Akritas stages:

1. Three safe `LogEvents` are stored with grouped `Incident` occurrence count.
2. Evidence is persisted without headers, cookies, tokens, DSNs or full request
   bodies.
3. Recent commits are attached as potentially related Evidence only, never as
   confirmed causality.
4. QVAC receives only sanitized, bounded context and read-only tool output.
5. The completed fixable Investigation maps to a single GitHub Issue.
6. The Remediation uses a deterministic branch for the same remediation ID.
7. The regression test starts red, the minimal fix bounds the short customer ID,
   validations pass, Akritas commits and pushes the branch.
8. Akritas creates or reconciles one PR for repository + base branch + head
   branch.
9. STOP: the PR remains open for human review.

## Retry Expectations

For the same logical run, reuse the same Incident/Investigation/Remediation
identifiers. Replays must return the existing Issue, Remediation branch and PR
instead of creating duplicates.

## Cleanup

Stop the fixture process and delete only disposable local workspaces under the
Akritas-managed workspace root. Close demo Issues/PRs manually if the remote
repository is temporary.
