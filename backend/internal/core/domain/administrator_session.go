package domain

import (
	"time"

	"github.com/google/uuid"
)

type AdministratorSession struct {
	ID                uuid.UUID
	AdministratorID   uuid.UUID
	AuthenticatedAt   time.Time
	IdleExpiresAt     time.Time
	AbsoluteExpiresAt time.Time
	RevokedAt         *time.Time
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
