package security

import "testing"

func TestBootstrapTokenVerifierVerify(t *testing.T) {
	t.Parallel()

	verifier := NewBootstrapTokenVerifier("deployment-bootstrap-secret-not-a-totp-seed")

	if !verifier.Verify("deployment-bootstrap-secret-not-a-totp-seed") {
		t.Fatal("matching token must verify")
	}
	if verifier.Verify("wrong-token") {
		t.Fatal("mismatched token must not verify")
	}
	if verifier.Verify("") {
		t.Fatal("empty candidate must not verify")
	}
	if verifier.Verify("deployment-bootstrap-secret-not-a-totp-see") {
		t.Fatal("a candidate shorter by one character must not verify")
	}
}

func TestBootstrapTokenVerifierRejectsEmptyConfiguredToken(t *testing.T) {
	t.Parallel()

	verifier := NewBootstrapTokenVerifier("")

	if verifier.Verify("") {
		t.Fatal("an empty configured token must never verify, even against an empty candidate")
	}
	if verifier.Verify("anything") {
		t.Fatal("an empty configured token must never verify")
	}
}
