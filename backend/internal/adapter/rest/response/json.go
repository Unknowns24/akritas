package response

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
)

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, code, message, userMessage string, details ...dto.ErrorDetail) {
	WriteJSON(w, status, dto.ErrorResponse{Error: dto.Error{
		Code:        code,
		Message:     message,
		UserMessage: userMessage,
		RequestID:   middleware.GetReqID(r.Context()),
		Details:     details,
	}})
}
