import { describe, it } from "node:test";
import assert from "node:assert/strict";
import type { Incident } from "../services/get-incident.service";
import { buildIncidentTraceabilityChain } from "./traceability.utils.ts";

describe("Traceability Chain Builder - FE-H5-04", () => {
  const baseIncident: Incident = {
    id: "inc-101",
    key: "INC-101",
    project: {
      id: "proj-1",
      name: "Core API",
    },
    fingerprint: "fp-12345",
    title: "Database connection pool exhausted",
    severity: "error",
    phase: "completed",
    occurrence_count: 1,
    first_seen_at: "2026-08-23T07:00:00Z",
    last_seen_at: "2026-08-23T07:00:00Z",
    resolution_status: "fixable",
    latest_investigation: {
      id: "inv-101",
      incident_id: "inc-101",
      status: "completed",
      created_at: "2026-08-23T07:01:00Z",
      summary: "Identified leaked connection handle in db.go",
      root_cause_status: "identified",
      resolution_status: "fixable",
      relevant_commits: ["a1b2c3d4e5f6"],
      relevant_files: ["pkg/db/pool.go"],
      hypotheses: [],
      evidence_ids: [],
      recommended_actions: ["Apply connection pooling"],
    },
    github_issue_reference: {
      number: 42,
      url: "https://github.com/org/repo/issues/42",
      repository: "org/repo",
      created_at: "2026-08-23T07:02:00Z",
    },
    remediation: {
      id: "rem-101",
      incident_id: "inc-101",
      status: "pull_request_created",
      branch_name: "akritas/fix-inc-101",
      changes_summary: "Added defer close() in pool.go",
      changes: [
        {
          file_path: "pkg/db/pool.go",
          change_type: "modified",
          patch: "@@ -10,3 +10,4 @@\n+ defer conn.Close()",
          redacted: true,
        },
      ],
      validation_summary: {
        total: 2,
        passed: 2,
        failed: 0,
      },
      pull_request_reference: {
        number: 105,
        url: "https://github.com/org/repo/pull/105",
        repository: "org/repo",
        branch: "akritas/fix-inc-101",
        created_at: "2026-08-23T07:05:00Z",
      },
      created_at: "2026-08-23T07:03:00Z",
    },
    pull_request_reference: {
      number: 105,
      url: "https://github.com/org/repo/pull/105",
      repository: "org/repo",
      branch: "akritas/fix-inc-101",
      created_at: "2026-08-23T07:05:00Z",
    },
  };

  it("should generate a complete 7-step traceability chain for a fixable incident with PR", () => {
    const chain = buildIncidentTraceabilityChain(baseIncident);

    assert.equal(chain.length, 7);

    // Step 1: Incident
    assert.equal(chain[0].id, "incident");
    assert.equal(chain[0].label, "INC-101");
    assert.equal(chain[0].status, "completed");

    // Step 2: Investigation
    assert.equal(chain[1].id, "investigation");
    assert.equal(chain[1].status, "completed");

    // Step 3: GitHub Issue
    assert.equal(chain[2].id, "issue");
    assert.equal(chain[2].label, "Issue #42");
    assert.equal(chain[2].url, "https://github.com/org/repo/issues/42");
    assert.equal(chain[2].status, "completed");

    // Step 4: Remediation
    assert.equal(chain[3].id, "remediation");
    assert.equal(chain[3].status, "completed");

    // Step 5: Branch
    assert.equal(chain[4].id, "branch");
    assert.equal(chain[4].label, "akritas/fix-inc-101");
    assert.equal(chain[4].status, "completed");

    // Step 6: Commit
    assert.equal(chain[5].id, "commit");
    assert.equal(chain[5].status, "completed");

    // Step 7: Pull Request
    assert.equal(chain[6].id, "pull_request");
    assert.equal(chain[6].label, "PR #105");
    assert.equal(chain[6].url, "https://github.com/org/repo/pull/105");
    assert.equal(chain[6].status, "completed");
  });

  it("should correctly halt chain when resolution_status is requires_human", () => {
    const humanIncident: Incident = {
      ...baseIncident,
      resolution_status: "requires_human",
      remediation: undefined,
      pull_request_reference: undefined,
    };

    const chain = buildIncidentTraceabilityChain(humanIncident);
    assert.equal(chain.length, 7);

    const remStep = chain.find((s) => s.id === "remediation");
    assert.equal(remStep?.status, "halted");

    const prStep = chain.find((s) => s.id === "pull_request");
    assert.equal(prStep?.status, "not_applicable");
  });

  it("should mark commit and PR as blocked when remediation fails", () => {
    const failedIncident: Incident = {
      ...baseIncident,
      remediation: {
        ...baseIncident.remediation!,
        status: "failed",
        failure_user_message: "Tests failed on branch",
        pull_request_reference: undefined,
      },
      pull_request_reference: undefined,
    };

    const chain = buildIncidentTraceabilityChain(failedIncident);
    assert.equal(chain.length, 7);

    const remStep = chain.find((s) => s.id === "remediation");
    assert.equal(remStep?.status, "failed");

    const commitStep = chain.find((s) => s.id === "commit");
    assert.equal(commitStep?.status, "failed");

    const prStep = chain.find((s) => s.id === "pull_request");
    assert.equal(prStep?.status, "failed");
  });
});

