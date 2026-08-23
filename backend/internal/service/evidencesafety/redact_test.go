package evidencesafety

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestRedactRemovesCredentialFormats(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		leaks   []string
		markers []string
	}{
		{
			name:    "json password string",
			input:   `{"password":"supersecret"}`,
			leaks:   []string{"supersecret"},
			markers: []string{`"password":"[REDACTED]"`},
		},
		{
			name:    "json token string with space",
			input:   `{"token": "secret-value"}`,
			leaks:   []string{"secret-value"},
			markers: []string{`"token": "[REDACTED]"`},
		},
		{
			name:    "json api key case insensitive",
			input:   `{"API_KEY":"secret value"}`,
			leaks:   []string{"secret value"},
			markers: []string{`"API_KEY":"[REDACTED]"`},
		},
		{
			name:    "double quoted assignment with spaces",
			input:   `PASSWORD="two words"`,
			leaks:   []string{"two words"},
			markers: []string{`PASSWORD="[REDACTED]"`},
		},
		{
			name:    "single quoted assignment with spaces",
			input:   `TOKEN='secret value'`,
			leaks:   []string{"secret value"},
			markers: []string{`TOKEN='[REDACTED]'`},
		},
		{
			name:    "mixed case secret assignment",
			input:   `AppSecret=secret-value`,
			leaks:   []string{"secret-value"},
			markers: []string{"AppSecret=[REDACTED]"},
		},
		{
			name:    "cookie assignment",
			input:   `session_cookie='session secret'`,
			leaks:   []string{"session secret"},
			markers: []string{"session_cookie='[REDACTED]'"},
		},
		{
			name:    "authorization bearer",
			input:   `Authorization: Bearer bearer-secret-value`,
			leaks:   []string{"bearer-secret-value"},
			markers: []string{"Authorization: Bearer [REDACTED]"},
		},
		{
			name:    "authorization basic",
			input:   `Authorization: Basic basic-secret-value`,
			leaks:   []string{"basic-secret-value"},
			markers: []string{"Authorization: Basic [REDACTED]"},
		},
		{
			name:    "github pat",
			input:   `github_pat_11ABCDEFGHijklmnopQRSTUVwxYZ1234567890`,
			leaks:   []string{"github_pat_11ABCDEFGHijklmnopQRSTUVwxYZ1234567890"},
			markers: []string{"[REDACTED_GITHUB_TOKEN]"},
		},
		{
			name:    "github app token",
			input:   `ghs_ABCDEFGHIJKLMNOPQRSTUVWXYZ123456`,
			leaks:   []string{"ghs_ABCDEFGHIJKLMNOPQRSTUVWXYZ123456"},
			markers: []string{"[REDACTED_GITHUB_TOKEN]"},
		},
		{
			name:    "jwt",
			input:   `eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJhZG1pbiJ9.sflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c`,
			leaks:   []string{"eyJhbGciOiJIUzI1NiJ9", "eyJzdWIiOiJhZG1pbiJ9", "sflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
			markers: []string{"[REDACTED_SESSION_TOKEN]"},
		},
		{
			name:    "session token",
			input:   `SESSION_TOKEN=session-secret-value`,
			leaks:   []string{"session-secret-value"},
			markers: []string{"SESSION_TOKEN=[REDACTED]"},
		},
		{
			name:    "cookie header",
			input:   `Cookie: sessionid=session-secret; theme=dark`,
			leaks:   []string{"session-secret"},
			markers: []string{"Cookie: [REDACTED]"},
		},
		{
			name:    "postgres dsn",
			input:   `postgres://admin:db-password@db/private`,
			leaks:   []string{"admin:db-password", "db-password"},
			markers: []string{"postgres://[REDACTED]@db/private"},
		},
		{
			name:    "mysql dsn",
			input:   `mysql://user:mysql-password@db:3306/service`,
			leaks:   []string{"user:mysql-password", "mysql-password"},
			markers: []string{"mysql://[REDACTED]@db:3306/service"},
		},
		{
			name:    "redis dsn password only",
			input:   `redis://:redis-password@cache/0`,
			leaks:   []string{":redis-password", "redis-password"},
			markers: []string{"redis://[REDACTED]@cache/0"},
		},
		{
			name: "pem private key",
			input: "-----BEGIN PRIVATE KEY-----\n" +
				"private-key-material\n" +
				"-----END PRIVATE KEY-----",
			leaks:   []string{"BEGIN PRIVATE KEY", "private-key-material", "END PRIVATE KEY"},
			markers: []string{"[REDACTED_PRIVATE_KEY]"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			redacted := Redact(tt.input)
			if !utf8.ValidString(redacted) {
				t.Fatal("redacted output must remain valid UTF-8")
			}
			for _, leak := range tt.leaks {
				if strings.Contains(redacted, leak) {
					t.Fatalf("case %q leaked sensitive content", tt.name)
				}
			}
			for _, marker := range tt.markers {
				if !strings.Contains(redacted, marker) {
					t.Fatalf("case %q missed expected redaction marker", tt.name)
				}
			}
		})
	}
}

func TestRedactKeepsNormalMentionsWithoutValues(t *testing.T) {
	t.Parallel()
	input := "password policy changed; token rotation needed; cookie settings documented"
	if got := Redact(input); got != input {
		t.Fatalf("normal security vocabulary without values must not be redacted: %q", got)
	}
}

func TestRedactAndLimitPreservesUTF8AndMaximumBytes(t *testing.T) {
	t.Parallel()
	input := "TOKEN='secret value' " + strings.Repeat("á", 80) + string([]byte{0xff, 0xfe})
	redacted := RedactAndLimit(input, 64)
	if len(redacted) > 64 {
		t.Fatalf("redacted output exceeded deterministic byte limit: %d", len(redacted))
	}
	if !utf8.ValidString(redacted) {
		t.Fatal("redacted output must remain valid UTF-8")
	}
	if strings.Contains(redacted, "secret value") {
		t.Fatal("redacted output leaked sensitive content")
	}
	if !strings.Contains(redacted, "[TRUNCATED]") {
		t.Fatal("bounded redaction should mark truncation")
	}
}
