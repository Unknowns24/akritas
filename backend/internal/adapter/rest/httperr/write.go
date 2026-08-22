package httperr

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type ErrorBody struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	UserMessage string `json:"user_message"`
	RequestID   string `json:"request_id"`
}

type Envelope struct {
	Error ErrorBody `json:"error"`
}

func Write(w http.ResponseWriter, err error, requestID string) {
	status := http.StatusInternalServerError
	body := ErrorBody{
		Code: "0x101001I", Message: "No se pudo procesar la operación.",
		UserMessage: "No se pudo procesar la operación.", RequestID: requestID,
	}
	var enriched *domain.Error
	if errors.As(err, &enriched) {
		body.Code = enriched.Code
		body.Message = enriched.Message
		body.UserMessage = enriched.UserMessage
		status = statusFor(enriched.Type())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Error: body})
}

func WriteStatus(w http.ResponseWriter, status int, code, message, userMessage, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Error: ErrorBody{
		Code: code, Message: message, UserMessage: userMessage, RequestID: requestID,
	}})
}

func statusFor(errorType byte) int {
	switch errorType {
	case 'V':
		return http.StatusBadRequest
	case 'U':
		return http.StatusUnauthorized
	case 'F':
		return http.StatusForbidden
	case 'N':
		return http.StatusNotFound
	case 'C':
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
