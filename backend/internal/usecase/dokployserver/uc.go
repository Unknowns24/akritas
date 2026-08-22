package dokployserver

import (
	"time"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

type UseCase struct {
	store   portsout.DokployServerStore
	gateway portsout.DokployGateway
	usage   portsout.IntegrationUsageReader
	newID   func() uuid.UUID
	now     func() time.Time
}

func New(store portsout.DokployServerStore, gateway portsout.DokployGateway, usage portsout.IntegrationUsageReader, newID func() uuid.UUID, now func() time.Time) portsin.DokployServerUseCase {
	return &UseCase{store: store, gateway: gateway, usage: usage, newID: newID, now: now}
}

func wipe(value []byte) {
	for index := range value {
		value[index] = 0
	}
}
