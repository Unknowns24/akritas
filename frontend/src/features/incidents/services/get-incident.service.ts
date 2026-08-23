import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type Incident = components["schemas"]["Incident"];

export async function getIncidentService(id: string): Promise<Incident> {
  const { data, error } = await api.GET("/incidents/{incident_id}", {
    params: {
      path: { incident_id: id },
    },
  });

  if (error) {
    if (typeof window === "undefined") {
      return {} as unknown as Incident;
    }
    throw error;
  }

  if (!data) {
    if (typeof window === "undefined") {
      return {} as unknown as Incident;
    }
    throw new Error("No data returned");
  }

  /* [MOCK DOCS]
  if (error || !data) {
    // Return mock data for local testing
    return {
      id: "inc-1",
      key: "AKR-184",
      project: { id: "1", name: "SENTINEL-API" },
      fingerprint: "nil_pointer_panic_get_users",
      severity: "critical",
      title: "Nil pointer panic on GET /users/:id",
      summary: "FindByID can return a nil user when the requested user does not exist, but service.go dereferences user.Name without validating the result.",
      phase: "remediating",
      root_cause_status: "identified",
      resolution_status: "fixable",
      confidence: 0.94,
      occurrence_count: 37,
      first_seen_at: new Date(Date.now() - 8 * 60000).toISOString(),
      last_seen_at: new Date().toISOString(),
      deployment_correlation: {
        id: "dep-123",
        deployment_id: "abc1234",
        first_incident_at: new Date(Date.now() - 5 * 60000).toISOString(),
        occurred_before_incident_seconds: 180,
      },
      latest_investigation: {
        id: "inv-1",
        status: "completed",
        started_at: new Date(Date.now() - 7 * 60000).toISOString(),
        completed_at: new Date(Date.now() - 6 * 60000).toISOString(),
        root_cause_summary: "FindByID can return a nil user when the requested user does not exist, but service.go dereferences user.Name without validating the result.",
        stack_traces: [
          {
            id: "trace-1",
            raw_content: "// internal/users/service.go:81\n  panic: runtime error: invalid memory address or nil pointer dereference",
            lang: "go"
          }
        ],
        files_analyzed: [
          "internal/users/service.go",
          "internal/users/repository.go",
          "internal/users/handler.go"
        ],
      },
      remediation: {
        id: "rem-1",
        status: "draft",
        patch_diff: `@@ -78,6 +78,9 @@
    78 user, err := s.repo.FindByID(ctx, id)
    79 if err != nil {
    80 return "", err
    81 }
  + 82 + if user == nil {
  + 83 + return "", errors.New("user not found")
  + 84 + }
    85 return user.Name, nil`,
        tests_passed: true,
        validation_passed: true,
      }
    } as any;
  }
  */

  return data.data;
}
