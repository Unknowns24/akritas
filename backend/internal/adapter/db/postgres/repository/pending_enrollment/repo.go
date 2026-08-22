package pendingenrollment

import (
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) out.PendingEnrollmentRepository {
	return &repository{db: db}
}
