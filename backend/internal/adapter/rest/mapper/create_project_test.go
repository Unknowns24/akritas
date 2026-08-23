package mapper

import (
	"testing"

	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestCreateProjectInitialLogIngestionDefaultAndValidation(t *testing.T) {
	enabled, patterns, grouping, context := false, []string{}, "PT30M", 20
	configuration := &projectdto.MonitoringConfigurationRequestDTO{Enabled: &enabled, ErrorPatterns: &patterns, IgnoredPatterns: &patterns, GroupingWindow: &grouping, ContextBefore: &context, ContextAfter: &context}
	command, err := CreateProjectToCommand(projectdto.CreateProjectRequestDTO{MonitoringConfiguration: configuration})
	if err != nil || command.InitialLogIngestion != domain.InitialLogIngestionFromNow {
		t.Fatalf("command = %+v, %v", command, err)
	}
	_, err = CreateProjectToCommand(projectdto.CreateProjectRequestDTO{MonitoringConfiguration: configuration, InitialLogIngestion: "everything"})
	if err == nil {
		t.Fatal("invalid initial ingestion accepted")
	}
}
