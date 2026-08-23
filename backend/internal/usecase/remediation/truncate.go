package remediation

const truncationMarker = "...[truncated]"

// truncateWithMarker caps s at max bytes, leaving a visible marker so a
// truncated excerpt is never mistaken for a complete one.
func truncateWithMarker(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= len(truncationMarker) {
		return truncationMarker[:max]
	}
	return s[:max-len(truncationMarker)] + truncationMarker
}
