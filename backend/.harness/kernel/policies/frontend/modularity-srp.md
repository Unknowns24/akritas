# Frontend Modularity and SRP Policy

The frontend must follow Single Responsibility Principle.

Each file, component, hook, service and function should have a clear responsibility.

## General rules

- A component should not handle layout, fetch, validation, mapping and complex rendering all at once.
- A hook should not mix unrelated UI state, requests, permissions, transformations and side effects.
- A service must not contain visual logic.
- A mapper must not perform requests.
- A schema must not navigate or render UI.
- A view may compose pieces, but should not become a large monolithic file.

## File size and modularity

Use judgment, but as a guideline:

- If a component exceeds ~150-200 lines, evaluate splitting it.
- If a view has many sections, create composites.
- If logic is reusable or complex, extract a hook.
- If data transformation is non-trivial, extract a mapper.
- If contract validation is needed, extract schemas.
- If constants repeat, move them to `*.constants.ts` inside the feature or to `core` if global.

## Recommended pattern for complex views

```text
views/
  SolicitudDetailView/
    index.tsx
    SolicitudDetailView.module.css
    useSolicitudDetailView.ts
    components/
      SolicitudHeader.tsx
      SolicitudFilesPanel.tsx
      SolicitudTimeline.tsx
      InformePanel.tsx
```

Responsibilities:

- View composes.
- Components render.
- Hooks orchestrate state and behavior.
- Services perform requests.
- Mappers transform data.
- Schemas validate contracts.

## View

Responsible for:

- Composing the screen.
- Connecting hooks with components.
- Defining the main layout for the view.

Not responsible for:

- Complex inline fetch calls.
- Heavy transformations.
- Inline contract validation.

## Component

Responsible for:

- Rendering UI.
- Receiving data through props.
- Emitting events through callbacks.

Not responsible for:

- Knowing endpoints.
- Knowing OpenAPI details.
- Handling global auth.
- Importing services, except in exceptional existing patterns.

## Composite

Responsible for:

- Building a functional section with multiple components.
- Receiving already-prepared data.
- Handling simple local UI state.

## Hook

Responsible for:

- Encapsulating state and behavior.
- Coordinating queries/mutations if applicable.
- Exposing a simple API to the view.

## Service

Responsible for:

- Calling the backend through the shared API client.
- Respecting `docs/openapi.yaml`.
- Returning typed and validated DTOs when relevant.

Not responsible for:

- Toasts.
- Navigation.
- Components.
- UI formatting.

## Mapper

Responsible for:

- Converting backend DTOs to UI-friendly models.
- Normalizing enums, dates and optional fields.

## Schema

Responsible for:

- Validating input/output data with Zod.
- Protecting the frontend from unexpected backend contract changes.
