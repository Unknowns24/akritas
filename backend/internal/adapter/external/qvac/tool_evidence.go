package qvac

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/service/evidencesafety"
	"github.com/google/uuid"
)

func (r *Runner) toolEvidence(investigationID uuid.UUID, name, arguments, content string) (*domain.Evidence, error) {
	var args struct {
		Path  string `json:"path"`
		Query string `json:"query"`
		SHA   string `json:"sha"`
	}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return nil, err
	}
	evidenceType := domain.EvidenceCodeLocation
	summary := "Hallazgo de repositorio obtenido mediante " + name + "."
	switch name {
	case "search_code":
		summary = fmt.Sprintf("Resultados reales de búsqueda para %q.", evidencesafety.RedactAndLimit(args.Query, 200))
	case "read_file":
		summary = fmt.Sprintf("Archivo %s leído desde el repositorio configurado.", evidencesafety.RedactAndLimit(args.Path, 500))
	case "list_recent_commits", "read_commit":
		evidenceType = domain.EvidenceCommit
	case "read_diff":
		evidenceType = domain.EvidenceDiff
	default:
		return nil, ErrUnknownTool
	}
	evidence, err := domain.NewEvidence(r.newID(), investigationID, evidenceType, summary, content, r.now().UTC())
	if err != nil {
		return nil, err
	}
	evidence.FilePath = evidencesafety.RedactAndLimit(strings.TrimSpace(args.Path), 4096)
	evidence.CommitSHA = evidencesafety.RedactAndLimit(strings.TrimSpace(args.SHA), 64)
	if evidenceType == domain.EvidenceDiff {
		evidence.Patch = evidencesafety.RedactAndLimit(content, maximumToolPayloadBytes)
		evidence.Content = ""
	}
	return evidence, evidence.Validate()
}
