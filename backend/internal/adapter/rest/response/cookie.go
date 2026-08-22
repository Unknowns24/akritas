package response

// SessionCookieName is the opaque session cookie name, per the OpenAPI
// cookieAuth security scheme. Shared between handlers (that set/expire it)
// and the auth middleware (that reads it), so both stay in sync without
// depending on each other.
const SessionCookieName = "akritas_session"
