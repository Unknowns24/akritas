import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type SetupStatusResponse = components["schemas"]["SetupStatusResponse"];
export type SetupRequest = components["schemas"]["SetupRequest"];
export type LoginRequest = components["schemas"]["LoginRequest"];
export type RecoveryRequest = components["schemas"]["RecoveryRequest"];
export type TotpEnrollment = components["schemas"]["TotpEnrollment"];
export type SessionResponse = components["schemas"]["SessionResponse"];

export async function getAuthSetupStatusService(): Promise<SetupStatusResponse["data"]> {
  try {
    const { data, error } = await api.GET("/auth/setup-status");
    if (error) throw error;
    return data.data;
  } catch (e) {
    // [MOCK DOCS] If the user wants to test setup again, they can change this to true,
    // but false is required to be able to access the Login screen and Dashboard.
    // console.warn("[MOCK] API failed, returning setup_required = false to allow login");
    // return { setup_required: false, registration_open: false };
    throw e;
  }
}

export async function getCurrentSessionService(): Promise<SessionResponse["data"]> {
  try {
    const { data, error } = await api.GET("/auth/session");
    if (error) throw error;
    return data.data;
  } catch (e) {
    
    throw e;
  }
}

export async function startAdministratorSetupService(body: SetupRequest): Promise<TotpEnrollment> {
  try {
    const { data, error } = await api.POST("/auth/setup", { body });
    if (error) throw error;
    return data.data;
  } catch (e) {
    
    throw e;
  }
}

export async function verifyAdministratorSetupService(totp_code: string, enrollment_id: string): Promise<SessionResponse["data"]> {
  try {
    const { data, error } = await api.POST("/auth/setup/verify", {
      body: { totp_code, enrollment_id } as any,
    });
    if (error) throw error;
    return data.data;
  } catch (e) {
    
    throw e;
  }
}

export async function loginAdministratorService(body: LoginRequest): Promise<SessionResponse["data"]> {
  try {
    const { data, error } = await api.POST("/auth/login", { body });
    if (error) throw error;
    return data.data;
  } catch (e) {
    
    throw e;
  }
}

export async function startAdministratorRecoveryService(body: RecoveryRequest): Promise<TotpEnrollment> {
  try {
    const { data, error } = await api.POST("/auth/recovery", { body });
    if (error) throw error;
    return data.data;
  } catch (e) {
    
    throw e;
  }
}

export async function verifyAdministratorRecoveryService(totp_code: string, enrollment_id: string): Promise<SessionResponse["data"]> {
  try {
    const { data, error } = await api.POST("/auth/recovery/verify", {
      body: { totp_code, enrollment_id } as any,
    });
    if (error) throw error;
    return data.data;
  } catch (e) {
    
    throw e;
  }
}

export async function logoutAdministratorService(): Promise<void> {
  try {
    const { error } = await api.DELETE("/auth/session");
    if (error) throw error;
  } catch (e) {
    
    throw e;
  }
}
