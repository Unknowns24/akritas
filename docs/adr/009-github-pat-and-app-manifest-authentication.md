# ADR-009 — GitHub mediante PAT y GitHub App Manifest

## Estado

Accepted

## Contexto

Akritas necesita leer repositorios y crear Issues, branches, commits y Pull
Requests. El MVP debe permitir una conexión rápida con Personal Access Token y
también una integración GitHub App con permisos explícitos y reutilizables.

`account_type` describe la cuenta objetivo, mientras que PAT/GitHub App describe
cómo Akritas se autentica. Modelarlos como el mismo concepto impediría representar
una GitHub App instalada tanto en una cuenta personal como en una organización.

## Decisión

`GitHubAccount` expondrá dos clasificaciones independientes:

```text
account_type
  personal | organization

authentication_method
  personal_access_token | github_app
```

### Personal Access Token

El administrador entrega el PAT una sola vez. El backend valida la identidad y
permisos, cifra el token mediante el Credential Store y sólo devuelve una
proyección segura con `credential_configured` y estado de autenticación.

La rotación reemplaza el secreto asociado sin cambiar el ID de `GitHubAccount`.

### GitHub App Manifest

Akritas implementará el flujo oficial GitHub App Manifest:

```text
Administrator inicia registro
        ↓
Akritas genera manifest + state de un solo uso
        ↓
navegador POST a github.com/settings/apps/new
        ↓
GitHub callback con code + state
        ↓
backend intercambia code y cifra private key/webhook secret
        ↓
redirect a instalación
        ↓
setup callback con installation_id + state
        ↓
backend valida la instalación usando autenticación de la App
        ↓
GitHubAccount conectado
```

`state` será impredecible, tendrá expiración máxima de una hora, estará asociado
al intento de conexión y será de un solo uso. `installation_id` nunca se acepta
sin verificar que pertenece a la App creada.

Los secretos generados por GitHub no pasan por el frontend ni se devuelven por la
API. El callback los cifra antes de persistir cualquier referencia estable.

La App solicitará únicamente:

- metadata: read;
- contents: read/write;
- issues: write;
- pull requests: write.

Los webhooks permanecen inactivos en el MVP porque el pipeline se apoya en
polling de Dokploy y operaciones explícitas, no en eventos GitHub.

## Reglas

- GitHub.com es el único host soportado en el MVP.
- Una conexión puede seleccionar únicamente repositorios accesibles mediante esa
  cuenta/instalación.
- Una integración referenciada por un Project no puede eliminarse.
- Nunca se registran PAT, PEM, webhook secret, installation token o payloads que
  puedan contenerlos.
- Las respuestas de error normalizan fallos del proveedor.

## Consecuencias

### Positivas

- PAT mantiene un camino de setup corto para la demo.
- GitHub App permite permisos explícitos y tokens de instalación efímeros.
- El dominio diferencia correctamente identidad de método de autenticación.

### Negativas

- El flujo Manifest requiere redirects y dos callbacks coordinados.
- Akritas debe custodiar private key y webhook secret aunque los webhooks estén
  desactivados.
- La experiencia local depende de que las URLs de callback sean alcanzables por
  el navegador que completa la instalación.

## Fuera de alcance del MVP

- GitHub Enterprise Server;
- OAuth App tradicional;
- múltiples instalaciones de una misma App administradas como un marketplace;
- webhooks GitHub;
- merge automático.

## Referencias

- [GitHub Docs — Registering a GitHub App from a manifest](https://docs.github.com/en/apps/sharing-github-apps/registering-a-github-app-from-a-manifest)
