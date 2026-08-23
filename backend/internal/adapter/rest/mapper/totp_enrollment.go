package mapper

import (
	"time"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func TOTPEnrollment(output in.StartAdministratorSetupOutput) authdto.TOTPEnrollmentResponseDTO {
	return authdto.TOTPEnrollmentResponseDTO{Data: authdto.TOTPEnrollmentDTO{
		EnrollmentID: output.EnrollmentID.String(), OtpauthURI: output.OtpauthURI,
		ManualEntryKey: output.ManualEntryKey, ExpiresAt: output.ExpiresAt.Format(time.RFC3339),
	}}
}
