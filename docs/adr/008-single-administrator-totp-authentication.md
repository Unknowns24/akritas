# ADR-008 — Autenticación single-admin con bootstrap, contraseña, TOTP y sesión opaca

## Estado

Accepted

## Contexto

Akritas administra credenciales de producción, código fuente, logs sensibles y
acciones de escritura sobre GitHub. Aunque el MVP no necesita usuarios
empresariales, RBAC ni multi-tenancy, exponer el dashboard sin autenticación
permitiría modificar integraciones y disparar operaciones sensibles.

La instalación debe seguir siendo simple y autocontenida. El mecanismo elegido
debe funcionar sin un proveedor de identidad externo y permitir usar aplicaciones
compatibles con TOTP como Google Authenticator o Apple Passwords.

## Decisión

Cada instalación de Akritas tendrá exactamente un `Administrator`.

El alta inicial se habilita mediante `AKRITAS_BOOTSTRAP_TOKEN`, un secreto de alta
entropía provisto al proceso desde el entorno. Ese valor:

- autoriza el setup inicial y la recuperación excepcional;
- nunca se persiste;
- nunca se devuelve por API;
- nunca se registra;
- no es el seed TOTP del administrador.

El flujo de setup será:

```text
bootstrap token + email + display name + password
        ↓
pending enrollment con expiración
        ↓
backend genera un secreto TOTP independiente
        ↓
otpauth URI / QR de una sola visualización
        ↓
confirmación con código TOTP
        ↓
Administrator activo + sesión
        ↓
registro cerrado
```

El login requiere dos factores:

```text
email + password + TOTP
```

TOTP sigue RFC 6238 con seis dígitos, período de 30 segundos, tolerancia máxima
de un período anterior/posterior y rechazo de reutilización de un período ya
aceptado.

La contraseña se almacena mediante Argon2id con salt único. La configuración
inicial sigue el mínimo vigente de OWASP: 19 MiB de memoria, dos iteraciones y
paralelismo 1 (`m=19456, t=2, p=1`). El mínimo del MVP es 12 caracteres y el
máximo 128, sin reglas de composición ni cambios periódicos forzados. Los
parámetros quedan versionados junto al hash para permitir upgrades futuros.

El secreto TOTP se cifra en reposo con `AKRITAS_MASTER_KEY`. Esta master key es
la misma raíz de cifrado de infraestructura definida para el Credential Store,
pero permanece separada de `AKRITAS_BOOTSTRAP_TOKEN`.

Las sesiones serán opacas y server-side. El navegador recibe únicamente una
cookie con:

```text
HttpOnly
Secure en producción
SameSite=Lax
Path=/
```

El lifetime por defecto será 12 horas de inactividad y 7 días absolutos. Login,
recovery y cambios de factor crean una sesión nueva; logout revoca la actual; un
recovery completado revoca todas las sesiones anteriores.

Las mutaciones autenticadas deben validar `Origin` contra el origen público
configurado. Producción sirve frontend y API bajo el mismo site. Desarrollo puede
usar una allowlist CORS explícita con credenciales.

Login y recovery aplican límites independientes por cuenta y por IP. Las
respuestas de autenticación son genéricas y no permiten inferir si falló el email,
password, TOTP o bootstrap token.

## Recovery

No existe reset por email en el MVP. El recovery requiere:

- email del único administrador;
- `AKRITAS_BOOTSTRAP_TOKEN`;
- nueva contraseña;
- nuevo enrollment TOTP confirmado.

El seed TOTP anterior queda invalidado al confirmar el nuevo factor.

## Consecuencias

### Positivas

- El dashboard deja de depender únicamente del aislamiento de red.
- El MVP no necesita un proveedor de identidad externo.
- Password y TOTP proporcionan dos factores para acciones administrativas.
- El frontend nunca accede al identificador opaco de sesión.
- Perder el autenticador puede resolverse desde la configuración segura del
  deployment.

### Negativas

- `AKRITAS_BOOTSTRAP_TOKEN` funciona como capacidad administrativa de recovery y
  debe protegerse como un secreto crítico.
- TOTP es susceptible a phishing y no reemplaza una passkey resistente a phishing.
- La instalación depende de sincronización horaria razonable.

## Fuera de alcance del MVP

- múltiples usuarios;
- invitaciones;
- roles o RBAC;
- SSO/OIDC/SAML;
- reset por email;
- recovery codes;
- passkeys WebAuthn;
- delegación de sesiones a clientes de terceros.

## Referencias

- [RFC 6238 — TOTP: Time-Based One-Time Password Algorithm](https://www.rfc-editor.org/info/rfc6238/)
- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [OWASP Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
