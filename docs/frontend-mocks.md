# Frontend Mocks & Testing Workarounds

This document tracks all temporary mocks, overrides, and hardcoded values injected into the frontend codebase for the purpose of local testing without a complete backend API. 

**Goal:** Before releasing to production or integrating the real backend, search for these keys and revert them to their intended logic.

---

## 1. Authentication Service Mocks

**File:** `frontend/src/features/auth/services/auth.service.ts`

**Reasoning:** Since the backend API is not yet active, the `openapi-fetch` client fails with `ECONNREFUSED` or `404` when hitting `/auth/setup-status` and `/auth/session`. To allow testing the UI flows (Setup, Login, Recovery) and accessing the protected dashboard routes, all endpoints in this file have been wrapped in a `try/catch` block that returns fake data.

**Specific Workarounds:**
- **`/auth/setup-status` (`getAuthSetupStatusService`)**:
  - Returns `{ setup_required: false, registration_open: false }` upon failure. 
  - *Action needed:* Remove the `catch` block so the app genuinely relies on the backend response.
- **`/auth/setup` (`startAdministratorSetupService`)**:
  - Returns a hardcoded `TotpEnrollment` object with a fake `enrollment_id`, `otpauth_uri`, and `manual_entry_key`.
- **`/auth/recovery` (`startAdministratorRecoveryService`)**:
  - Returns a hardcoded `TotpEnrollment` similar to setup.
- **`/auth/setup/verify`, `/auth/login`, `/auth/recovery/verify`**:
  - **Magic Code:** They all check if `totp_code === "123456"`. If not, they throw a fake `Error("Invalid authenticator code")`.
  - **Session Persistence:** They execute `localStorage.setItem("mock_auth", "true")`.
  - *Action needed:* Remove the magic code check, the `localStorage` logic, and the mock `SessionResponse` payload.
- **`/auth/session` (`getCurrentSessionService`)**:
  - Checks `localStorage.getItem("mock_auth") === "true"`. If true, returns a fake active session. If false, throws an error to force redirect to login.
  - *Action needed:* Remove the `localStorage` read and the mock payload. Real sessions will be authenticated securely via HTTP-Only cookies.
- **`/auth/session` (DELETE) (`logoutAdministratorService`)**:
  - Catches the API failure and executes `localStorage.removeItem("mock_auth")` to simulate ending the session.
  - *Action needed:* Remove the `catch` block and the local storage clearing logic.

## 2. Incidents Service Mocks

**Files:** 
- `frontend/src/features/incidents/services/get-incidents.service.ts`
- `frontend/src/features/incidents/services/get-incident.service.ts`

**Reasoning:** To preview the incident UI components and detail views, the endpoints are mocked since the real backend is unavailable.

**Specific Workarounds:**
- **`/incidents` (`getIncidentsService`)**:
  - Catches API failure and returns hardcoded lists of `IncidentSummary` (one in `detected` phase, one in `failed` phase).
  - *Action needed:* Remove the `catch` block and the mock array.
- **`/incidents/{id}` (`getIncidentService`)**:
  - Catches API failure and returns a hardcoded full `Incident` object mimicking issue `AKR-184` (Nil pointer panic). Includes mock stack traces, root cause analysis, and a patch diff.
  - *Action needed:* Remove the `catch` block and mock object.

---

## 3. Settings - GitHub Integration Mocks

**Files:** 
- `frontend/src/features/settings/services/github/*.service.ts`

**Reasoning:** Since the backend API is not yet active, the openapi-fetch client fails with `ECONNREFUSED` or 404 when hitting `/integrations/github/accounts/*`. To allow testing the settings UI, all endpoints in this module have been wrapped in a check that returns fake data if the request fails.

**Specific Workarounds:**
- **`/integrations/github/accounts` (GET & POST)**:
  - Return a static list of connected accounts or a fake created account object to simulate PAT setup.
  - *Action needed:* Remove the `!data || error` condition returning fake objects.
- **`/integrations/github/accounts/{id}` (PATCH, DELETE)**:
  - Return fake updated objects or empty success responses instead of throwing.
  - *Action needed:* Remove fake early returns.
- **`/integrations/github/accounts/{id}/connection-test`**:
  - Automatically simulates a successful connection.
  - *Action needed:* Remove fake return.
- **`/integrations/github/accounts/{id}/repositories`**:
  - Generates 42 mock repositories.
  - *Action needed:* Remove the mock generation block.

---

*Note: Whenever a new mock or local workaround is added for testing in future prompts, it should be appended to this document.*
