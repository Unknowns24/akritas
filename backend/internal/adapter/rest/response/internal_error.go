package response

import "net/http"

// These two codes live outside the domain error catalog (internal/core/domain):
// they are REST-transport-layer failures, not domain business rules, so they
// use component layer "1" (REST) instead of "4" (domain) per the layer scheme
// in docs/errors/aaa-map.md. Registered in docs/errors/aaa-map.md alongside
// the domain catalog for traceability.
const (
	CodeMalformedRequest = "0x100002V"
	CodeInternalError    = "0x100001I"
	// CodeRateLimited has no fitting type letter (V/U/F/N/C/I): the OpenAPI
	// ErrorCode pattern has no letter for 429. 'C' is the closest fit (the
	// request conflicts with a configured rate constraint) and is never used
	// to derive the HTTP status here — the caller sets 429 explicitly.
	CodeRateLimited = "0x100003C"
)

func WriteInternalError(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusInternalServerError, CodeInternalError,
		"unexpected internal error", "Ocurrió un error inesperado. Intentá nuevamente más tarde.")
}
