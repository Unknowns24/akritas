import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type SetupStatusResponse = components["schemas"]["SetupStatusResponse"];
export type SetupRequest = components["schemas"]["SetupRequest"];
export type LoginRequest = components["schemas"]["LoginRequest"];
export type RecoveryRequest = components["schemas"]["RecoveryRequest"];
export type TotpEnrollment = components["schemas"]["TotpEnrollment"];
export type SessionResponse = components["schemas"]["SessionResponse"];
export type TotpEnrollmentVerificationRequest =
  components["schemas"]["TotpEnrollmentVerificationRequest"];

export async function getAuthSetupStatusService(): Promise<SetupStatusResponse["data"]> {
  const { data, error } = await api.GET("/auth/setup-status");
  return requireApiData(data, error).data;
}

export async function getCurrentSessionService(): Promise<SessionResponse["data"]> {
  const { data, error } = await api.GET("/auth/session");
  return requireApiData(data, error).data;
}

export async function startAdministratorSetupService(body: SetupRequest): Promise<TotpEnrollment> {
  const { data, error } = await api.POST("/auth/setup", { body });
  return requireApiData(data, error).data;
}

export async function verifyAdministratorSetupService(totp_code: string, enrollment_id: string): Promise<SessionResponse["data"]> {
  const body: TotpEnrollmentVerificationRequest = { totp_code, enrollment_id };
  const { data, error } = await api.POST("/auth/setup/verify", { body });
  return requireApiData(data, error).data;
}

export async function loginAdministratorService(body: LoginRequest): Promise<SessionResponse["data"]> {
  const { data, error } = await api.POST("/auth/login", { body });
  return requireApiData(data, error).data;
}

export async function startAdministratorRecoveryService(body: RecoveryRequest): Promise<TotpEnrollment> {
  const { data, error } = await api.POST("/auth/recovery", { body });
  return requireApiData(data, error).data;
}

export async function verifyAdministratorRecoveryService(totp_code: string, enrollment_id: string): Promise<SessionResponse["data"]> {
  const body: TotpEnrollmentVerificationRequest = { totp_code, enrollment_id };
  const { data, error } = await api.POST("/auth/recovery/verify", { body });
  return requireApiData(data, error).data;
}

export async function logoutAdministratorService(): Promise<void> {
  const { error } = await api.DELETE("/auth/session");
  if (error) throw error;
}
