package response

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func TestRequestIDUsesValidatedHeaderThenChiContext(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		want      string
		forbidden string
	}{
		{name: "valid header", header: "req-client-12345", want: "req-client-12345"},
		{name: "missing header"},
		{name: "invalid header", header: "bad", forbidden: "bad"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got string
			handler := chimiddleware.RequestID(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				got = requestID(r)
			}))
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.header != "" {
				request.Header.Set("X-Request-ID", test.header)
			}
			handler.ServeHTTP(httptest.NewRecorder(), request)

			if test.want != "" && got != test.want {
				t.Fatalf("requestID() = %q, want %q", got, test.want)
			}
			if test.forbidden != "" && got == test.forbidden {
				t.Fatalf("requestID() trusted invalid value %q", got)
			}
			if trimmed := strings.TrimSpace(got); len(trimmed) < 8 || len(trimmed) > 100 {
				t.Fatalf("requestID() = %q, want 8..100 chars", got)
			}
		})
	}
}
