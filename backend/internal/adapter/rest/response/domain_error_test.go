package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestWriteDomainErrorMapsCodeSuffixToStatus(t *testing.T) {
	t.Parallel()

	cases := map[byte]int{
		'V': http.StatusBadRequest,
		'U': http.StatusUnauthorized,
		'F': http.StatusForbidden,
		'N': http.StatusNotFound,
		'C': http.StatusConflict,
		'I': http.StatusInternalServerError,
	}

	for letter, wantStatus := range cases {
		letter, wantStatus := letter, wantStatus
		t.Run(string(letter), func(t *testing.T) {
			t.Parallel()
			err := &domain.Error{Code: "0x999999" + string(letter), Message: "m", UserMessage: "u"}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			response.WriteDomainError(rec, req, err)

			if rec.Code != wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, wantStatus)
			}
		})
	}
}
