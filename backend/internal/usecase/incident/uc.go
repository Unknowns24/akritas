package incident

import portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"

type UseCase struct{ incidents portsout.IncidentStore }

func New(incidents portsout.IncidentStore) *UseCase { return &UseCase{incidents: incidents} }
