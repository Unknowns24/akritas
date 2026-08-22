package router

import (
	"errors"
	"testing"
)

func TestRouterFailsClosedWithoutAdministratorMiddleware(t *testing.T) {
	handler, err := New(Config{})
	if !errors.Is(err, ErrAdminMiddlewareUnavailable) || handler != nil {
		t.Fatalf("router must not be mountable without PB-061..063 middleware: handler=%v err=%v", handler, err)
	}
}
