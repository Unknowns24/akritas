# Flutter Mobile Project Structure Policy

Expected mobile stack by default:

- Flutter.
- Dart.
- Material 3.
- Riverpod.
- GoRouter.
- Dio.
- CookieJar / dio_cookie_manager when the backend uses cookies.
- flutter_secure_storage when bearer tokens or other sensitive local secrets are used.
- Feature-based architecture inside `lib/features/`.
- Hexagonal layering adapted to Flutter.

For existing projects, preserve the current stack and conventions. Do not introduce a new state manager, router, HTTP client, UI kit or generated architecture unless the project already uses it or the task explicitly requests it.

## Sources of truth

Respect these sources in this order:

1. `docs/openapi.yaml` — backend contract: endpoints, DTOs, responses, errors and auth.
2. `DESIGN.md`, `design.md`, Figma or existing screens — visual and UX contract.
3. `AGENTS.md` and `.harness/kernel/**` — technical architecture, workflow and quality rules.
4. Existing code patterns.

## Base structure

```text
lib/
  main.dart
  app/
    app.dart
    router.dart
    shell/
    theme/
  core/
    api/
    auth/
    config/
    errors/
    providers/
    ui/
    utils/
  features/
    <feature_name>/
      domain/
      application/
      data/
        dto/
        mapper/
        remote/
      presentation/
        pages/
        screens/
        widgets/
      shared/
```

Use only the folders needed by the feature, but do not collapse layers in a way that breaks dependency direction.

## Architecture direction

Allowed direction:

```text
presentation → application → domain
data → domain
data → core/api
application → core providers when needed
presentation → core/ui
```

Forbidden:

- `domain` importing `data`.
- `domain` importing `presentation`.
- `domain` importing Flutter widgets.
- `presentation` importing DTOs.
- `presentation` calling `ApiClient` directly.
- `presentation` creating `Dio()`.
- A feature importing internals of another feature.
- Business logic living inside widgets.

## `lib/app/`

`app` contains app bootstrap, router, shell and theme.

Rules:

- Centralize routes in `app/router.dart` or the existing router location.
- Keep auth redirects close to the router/auth state.
- Keep visual tokens in `app/theme` or the existing theme location.
- Do not put feature business logic in `app`.

## `lib/core/`

`core` contains cross-feature infrastructure and UI primitives.

Suggested responsibilities:

- `core/api`: shared `ApiClient`, envelopes, pagination types, API errors.
- `core/auth`: token/session storage abstractions.
- `core/providers`: global providers such as API client, router or config.
- `core/ui`: reusable app widgets, not feature-specific business UI.
- `core/utils`: generic helpers.

Rules:

- `core` must not import feature internals.
- Do not place one-feature-only logic in `core`.
- Shared UI must not contain business rules.

## `lib/features/`

Each feature should be autonomous and expose a clear boundary.

Recommended structure:

```text
lib/features/<feature>/
  domain/
    <entity>.dart
    <feature>_repository.dart
  application/
    <feature>_providers.dart
    <feature>_state.dart
    <feature>_notifier.dart
  data/
    dto/
      <entity>_dto.dart
    mapper/
      <feature>_mapper.dart
    remote/
      <feature>_remote_data_source.dart
    <feature>_repository_impl.dart
  presentation/
    screens/
    pages/
    widgets/
```

### domain

Contains:

- Entities.
- Value objects.
- Enums.
- Repository contracts.
- Pure domain rules.

Rules:

- No DTOs.
- No Dio.
- No ApiClient.
- No Flutter widgets.
- No Riverpod unless the existing project intentionally keeps simple providers there; prefer application.

### application

Contains:

- Use cases.
- Services.
- Controllers/notifiers.
- Providers.
- State objects.

Rules:

- Depends on `domain`.
- Coordinates repositories.
- No JSON parsing.
- No DTOs.
- No direct `ApiClient` calls.

### data

Contains:

- DTOs.
- Mappers.
- Remote/local datasources.
- Repository implementations.

Rules:

- This is the feature layer that can know JSON.
- This is the feature layer that can call `ApiClient`.
- Map DTOs into domain entities before returning to application/presentation.
- Do not expose DTOs to UI.

### presentation

Contains:

- Screens.
- Pages.
- Widgets.
- Small screen controllers when purely visual.

Rules:

- Consume application providers/notifiers.
- Do not call `ApiClient`.
- Do not import DTOs.
- Do not parse JSON.
- Do not invent backend data.
- Handle loading, empty, error and retry for remote data.

## Feature public surface

For larger projects, expose a public surface per feature when useful:

```text
lib/features/<feature>/<feature>.dart
```

Avoid importing deep internals of another feature. If cross-feature reuse becomes necessary, move truly shared code to `core` or model the dependency explicitly in application/domain.

## Naming conventions

- Widgets/screens: `PascalCase` classes, snake_case filenames.
- Providers: suffix `Provider`.
- Notifiers/controllers: suffix `Notifier` or `Controller` according to existing patterns.
- States: suffix `State`.
- DTOs: suffix `Dto`.
- Mappers: suffix `Mapper`.
- Repository contracts: `<Feature>Repository` in `domain`.
- Repository implementations: `<Feature>RepositoryImpl` in `data`.

Prefer existing project language conventions, but keep naming consistent inside each repository.
