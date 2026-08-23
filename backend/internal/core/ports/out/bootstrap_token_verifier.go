package out

// BootstrapTokenVerifier compares a candidate token against the configured
// AKRITAS_BOOTSTRAP_TOKEN in constant time.
type BootstrapTokenVerifier interface {
	Verify(candidate string) bool
}
