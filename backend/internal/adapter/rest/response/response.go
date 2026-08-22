package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	resterrors "github.com/Unknowns24/akritas/backend/internal/adapter/rest/errors"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

const (
	SessionCookieName    = "akritas_session"
	CodeMalformedRequest = "0x102001V"
	CodeInternalError    = "0x102002I"
	CodeRateLimited      = "0x102003R"
)

func JSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	JSON(w, status, value)
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, code, message, userMessage string, details ...commondto.ErrorDetailDTO) {
	JSON(w, status, commondto.ErrorResponseDTO{Error: commondto.ErrorDTO{
		Code: code, Message: message, UserMessage: userMessage, RequestID: requestID(r), Details: details,
	}})
}

func WriteDomainError(w http.ResponseWriter, r *http.Request, err *domain.Error) {
	if len(err.Code) > 0 && err.Code[len(err.Code)-1] == 'R' {
		w.Header().Set("Retry-After", "60")
	}
	Error(w, r, err)
}

func WriteInternalError(w http.ResponseWriter, r *http.Request) {
	Error(w, r, resterrors.ErrRequestFailed)
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	var stable *domain.Error
	if !errors.As(err, &stable) {
		stable = resterrors.ErrRequestFailed
	}
	status := http.StatusInternalServerError
	if len(stable.Code) > 0 {
		switch stable.Code[len(stable.Code)-1] {
		case 'V':
			status = http.StatusBadRequest
		case 'U':
			status = http.StatusUnauthorized
		case 'F':
			status = http.StatusForbidden
		case 'N':
			status = http.StatusNotFound
		case 'C':
			status = http.StatusConflict
		case 'R':
			status = http.StatusTooManyRequests
		}
	}
	JSON(w, status, commondto.ErrorResponseDTO{Error: commondto.ErrorDTO{Code: stable.Code, Message: stable.Message, UserMessage: stable.UserMessage, RequestID: requestID(r)}})
}

func Invalid(w http.ResponseWriter, r *http.Request) {
	stable := resterrors.ErrInvalidRequest
	JSON(w, http.StatusBadRequest, commondto.ErrorResponseDTO{Error: commondto.ErrorDTO{
		Code: stable.Code, Message: stable.Message, UserMessage: stable.UserMessage, RequestID: requestID(r),
	}})
}

func requestID(r *http.Request) string {
	if r != nil {
		value := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if len(value) >= 8 && len(value) <= 100 {
			return value
		}
	}
	return "req-" + uuid.NewString()
}
