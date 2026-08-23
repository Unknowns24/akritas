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
    /* [MOCK DOCS]
    if (typeof window !== "undefined" && localStorage.getItem("mock_auth") === "true") {
      console.warn("[MOCK] API failed, but mock_auth is set. Returning mock session.");
      return {
        administrator: { id: "mock-id", email: "admin@example.com", display_name: "Mock Admin", created_at: "", updated_at: "" },
        authenticated_at: new Date().toISOString(),
        idle_expires_at: "",
        absolute_expires_at: "",
      };
    }
    console.warn("[MOCK] API failed, throwing to simulate no active session");
    */
    throw e;
  }
}

export async function startAdministratorSetupService(body: SetupRequest): Promise<TotpEnrollment> {
  try {
    const { data, error } = await api.POST("/auth/setup", { body });
    if (error) throw error;
    return data.data;
  } catch (e) {
    /* [MOCK DOCS]
    console.warn("[MOCK] API failed, returning mock TOTP enrollment");
    return {
      enrollment_id: "mock-enrollment-id-123",
      otpauth_uri: "otpauth://totp/Akritas:admin@example.com?secret=JBSWY3DPEHPK3PXP&issuer=Akritas",
      manual_entry_key: "JBSWY3DPEHPK3PXP",
      expires_at: new Date(Date.now() + 1000 * 60 * 10).toISOString(),
    };
    */
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
    /* [MOCK DOCS]
    if (totp_code !== "123456") {
      throw new Error("Invalid authenticator code. Try 123456 for testing.");
    }
    if (typeof window !== "undefined") localStorage.setItem("mock_auth", "true");
    console.warn("[MOCK] API failed, returning mock session");
    return {
      administrator: { id: "mock-id", email: "admin@example.com", display_name: "Mock Admin", created_at: "", updated_at: "" },
      authenticated_at: new Date().toISOString(),
      idle_expires_at: "",
      absolute_expires_at: "",
    };
    */
    throw e;
  }
}

export async function loginAdministratorService(body: LoginRequest): Promise<SessionResponse["data"]> {
  try {
    const { data, error } = await api.POST("/auth/login", { body });
    if (error) throw error;
    return data.data;
  } catch (e) {
    /* [MOCK DOCS]
    if (body.totp_code !== "123456") {
      throw new Error("Invalid authenticator code. Try 123456 for testing.");
    }
    if (typeof window !== "undefined") localStorage.setItem("mock_auth", "true");
    console.warn("[MOCK] API failed, returning mock session");
    return {
      administrator: { id: "mock-id", email: body.email, display_name: "Mock Admin", created_at: "", updated_at: "" },
      authenticated_at: new Date().toISOString(),
      idle_expires_at: "",
      absolute_expires_at: "",
    };
    */
    throw e;
  }
}

export async function startAdministratorRecoveryService(body: RecoveryRequest): Promise<TotpEnrollment> {
  try {
    const { data, error } = await api.POST("/auth/recovery", { body });
    if (error) throw error;
    return data.data;
  } catch (e) {
    /* [MOCK DOCS]
    console.warn("[MOCK] API failed, returning mock TOTP enrollment");
    return {
      enrollment_id: "mock-recovery-enrollment-id",
      otpauth_uri: "otpauth://totp/Akritas:admin@example.com?secret=RECOVERY3DPEHPK3PXP&issuer=Akritas",
      manual_entry_key: "RECOVERY3DPEHPK3PXP",
      expires_at: new Date(Date.now() + 1000 * 60 * 10).toISOString(),
    };
    */
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
    /* [MOCK DOCS]
    if (totp_code !== "123456") {
      throw new Error("Invalid authenticator code. Try 123456 for testing.");
    }
    console.warn("[MOCK] API failed, returning mock session");
    return {
      administrator: { id: "mock-id", email: "admin@example.com", display_name: "Mock Admin", created_at: "", updated_at: "" },
      authenticated_at: new Date().toISOString(),
      idle_expires_at: "",
      absolute_expires_at: "",
    };
    */
    throw e;
  }
}

export async function logoutAdministratorService(): Promise<void> {
  try {
    const { error } = await api.DELETE("/auth/session");
    if (error) throw error;
  } catch (e) {
    /* [MOCK DOCS]
    console.warn("[MOCK] API failed, simulating logout by removing mock_auth");
    if (typeof window !== "undefined") {
      localStorage.removeItem("mock_auth");
    }
    */
    throw e;
  }
}
