package git

import "testing"

func TestValidateRefName(t *testing.T) {
	valid := []string{"main", "master", "akritas/remediation/3f2504e0-4f89-11d3-9a0c-0305e82c3301", "feature-123"}
	for _, name := range valid {
		if err := validateRefName(name); err != nil {
			t.Fatalf("expected %q to be valid, got %v", name, err)
		}
	}

	invalid := []string{
		"",
		"-Xupload-pack=/bin/sh",
		"--upload-pack=touch /tmp/pwned",
		"..",
		"a..b",
		"a b",
		"a\tb",
		"a\nb",
		"-",
	}
	for _, name := range invalid {
		if err := validateRefName(name); err == nil {
			t.Fatalf("expected %q to be rejected", name)
		}
	}
}
