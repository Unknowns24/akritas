# Security Review Policy

Security review must happen after implementation and before marking ready for human review.

Check at minimum:

- Authentication and authorization.
- RBAC/permissions.
- Tenant/user scoping.
- Input validation.
- Error leakage.
- Secrets and environment variables.
- File uploads/downloads.
- Path traversal.
- SSRF risks.
- SQL injection risks.
- XSS risks in frontend.
- Unsafe HTML rendering.
- CORS changes.
- Token/cookie/session handling.

If a security issue is found, generate a fix plan and request human approval before applying non-trivial changes.
