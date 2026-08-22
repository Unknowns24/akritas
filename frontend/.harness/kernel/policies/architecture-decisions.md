# Architecture Decision Policy

## ADR

Accepted ADRs are constraints.

Before making an architectural decision:

1. Search existing ADRs.
2. If an accepted ADR covers the decision, follow it.
3. Do not silently contradict an accepted ADR.
4. If the task requires changing an accepted decision, stop and propose
   the change for human approval.

## RFC

RFCs describe proposed or future capabilities.

An RFC must not be treated as approved implementation scope.

Do not implement an RFC unless:

- the task explicitly requests it; or
- its status has been promoted to an accepted decision/specification.

## New architectural decisions

If implementation requires a significant architectural decision not
covered by existing documentation:

1. Do not silently decide it.
2. Document the decision/proposal.
3. Request human approval when appropriate.
4. Continue implementation only after the decision is resolved.