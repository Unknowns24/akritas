package domain

import (
	"time"

	"github.com/google/uuid"
)

type AdministratorSession struct {
	ID                uuid.UUID  `gorm:"column:id;type:uuid;primaryKey"`
	AdministratorID   uuid.UUID  `gorm:"column:administrator_id;type:uuid"`
	AuthenticatedAt   time.Time  `gorm:"column:authenticated_at"`
	IdleExpiresAt     time.Time  `gorm:"column:idle_expires_at"`
	AbsoluteExpiresAt time.Time  `gorm:"column:absolute_expires_at"`
	RevokedAt         *time.Time `gorm:"column:revoked_at"`
}

func NewAdministratorSession(id, administratorID uuid.UUID, authenticatedAt, idleExpiresAt, absoluteExpiresAt time.Time) (*AdministratorSession, error) {
	session := &AdministratorSession{
		ID: id, AdministratorID: administratorID, AuthenticatedAt: authenticatedAt,
		IdleExpiresAt: idleExpiresAt, AbsoluteExpiresAt: absoluteExpiresAt,
	}
	if err := session.Validate(); err != nil {
		return nil, err
	}
	return session, nil
}

func (s AdministratorSession) Validate() error {
	invalid := s.ID == uuid.Nil || s.AdministratorID == uuid.Nil || !validTime(s.AuthenticatedAt) ||
		!s.IdleExpiresAt.After(s.AuthenticatedAt) || !s.AbsoluteExpiresAt.After(s.AuthenticatedAt) ||
		s.IdleExpiresAt.After(s.AbsoluteExpiresAt)
	if s.RevokedAt != nil && s.RevokedAt.Before(s.AuthenticatedAt) {
		invalid = true
	}
	if invalid {
		return ErrInvalidAdministratorSession.Wrap(validationCause("administrator session"))
	}
	return nil
}

func (s AdministratorSession) IsActive(now time.Time) bool {
	return s.Validate() == nil && s.RevokedAt == nil && now.Before(s.IdleExpiresAt) && now.Before(s.AbsoluteExpiresAt)
}

func (s *AdministratorSession) Revoke(at time.Time) error {
	if s == nil || at.Before(s.AuthenticatedAt) {
		return ErrAdministratorSessionTransition.Wrap(validationCause("revocation time"))
	}
	if s.RevokedAt != nil {
		return nil
	}
	s.RevokedAt = &at
	return nil
}

// ExtendIdle slides the idle expiration forward on an active session,
// capped so it never exceeds AbsoluteExpiresAt.
func (s *AdministratorSession) ExtendIdle(now time.Time, idleTTL time.Duration) error {
	if s == nil || !s.IsActive(now) {
		return ErrAdministratorSessionTransition.Wrap(validationCause("idle extension"))
	}
	newIdle := now.Add(idleTTL)
	if newIdle.After(s.AbsoluteExpiresAt) {
		newIdle = s.AbsoluteExpiresAt
	}
	s.IdleExpiresAt = newIdle
	return nil
}
