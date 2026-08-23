package remediation

import "github.com/google/uuid"

// remediationBranchName is deterministic and traceable back to the
// Remediation identity alone — it never depends on IssueReference or a
// GitHub Issue existing.
func remediationBranchName(remediationID uuid.UUID) string {
	return "akritas/remediation/" + remediationID.String()
}
