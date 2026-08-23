package evidencesafety

import (
	"strings"
	"testing"
)

func TestRedactRemovesCredentialShapesAndPreservesBounds(t *testing.T) {
	t.Parallel()
	input := "Authorization: Bearer abc.def.ghi token=supersecret github_pat_ABCDEFGHIJKLMNOPQRSTUVWXYZ123456 " +
		"DOKPLOY_API_KEY=dokploy-value DATABASE_PASSWORD=db-value postgres://admin:connection-password@db/private " +
		"eyJheader.payload.signature\n-----BEGIN PRIVATE KEY-----\nsecret\n-----END PRIVATE KEY-----"
	redacted := RedactAndLimit(input, 512)
	for _, secret := range []string{"abc.def.ghi", "supersecret", "github_pat_", "dokploy-value", "db-value", "connection-password", "eyJheader", "BEGIN PRIVATE KEY", "\nsecret\n"} {
		if strings.Contains(redacted, secret) {
			t.Fatalf("secret %q leaked in %q", secret, redacted)
		}
	}
	if len(redacted) > 524 {
		t.Fatalf("redacted bounded value unexpectedly large: %d", len(redacted))
	}
}
