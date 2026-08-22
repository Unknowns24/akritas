# Frontend Security Policy

This policy complements the global security policy for frontend-specific concerns.

## Credentials and sensitive data

- Never expose backend or integration secrets to browser code.
- Never place GitHub, Dokploy, QVAC or internal-service credentials in `NEXT_PUBLIC_*` variables.
- Do not store sensitive integration tokens in localStorage, sessionStorage or client-visible application state.
- Prefer safe backend projections such as connection status instead of returning credentials to the frontend.

## Authentication and permissions

- Private routes should live under a private route group such as `src/app/(private)` when the project uses that convention.
- Public routes should live under `src/app/(public)` when applicable.
- Do not show actions the user cannot execute when permissions are known.
- Hiding actions in frontend is not security; backend remains the source of truth.
- Handle session expiration consistently through the shared API client.
- Avoid storing sensitive tokens in unsafe locations if the project has a safer existing pattern.
- Authentication payload handling must follow the OpenAPI contract; do not invent client-side cryptographic protocols.

## Errors

- Do not show raw technical errors, stack traces or infrastructure details to users.
- Use normalized human-friendly messages.
- Preserve request IDs or technical details only where appropriate.

## Uploads

When files are uploaded, validate type and size before sending when the contract defines constraints.

Do not trust client-side validation as security; backend must still validate.
