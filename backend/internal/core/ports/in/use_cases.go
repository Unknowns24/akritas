package in

// UseCases groups the input ports required by the REST adapter. It is a wiring
// container only; behavior remains owned by each cohesive input port.
type UseCases struct {
	GetSetupStatus           GetSetupStatusUseCase
	StartAdministratorSetup  StartAdministratorSetupUseCase
	VerifyAdministratorSetup VerifyAdministratorSetupUseCase
	LoginAdministrator       LoginAdministratorUseCase
	AuthenticateSession      AuthenticateSessionUseCase
	GetCurrentSession        GetCurrentSessionUseCase
	LogoutAdministrator      LogoutAdministratorUseCase
	GitHubAccount            GitHubAccountUseCase
	GitHubApp                GitHubAppUseCase
	DokployServer            DokployServerUseCase
	Project                  ProjectUseCase
	Incident                 IncidentUseCase
}
