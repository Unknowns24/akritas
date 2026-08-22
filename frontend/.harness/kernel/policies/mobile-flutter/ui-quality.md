# Flutter Mobile UI Quality Policy

The mobile UI must remain consistent with the existing app design.

## Visual source of truth

Use, in order:

1. Existing implemented screens.
2. `DESIGN.md` / `design.md` if present.
3. Figma/screenshots explicitly attached to the task.
4. Theme tokens in `lib/app/theme` or equivalent.

## Rules

- Do not redesign existing screens unless explicitly requested.
- Do not change the layout just because backend integration is being added.
- Do not replace existing visual components if they already satisfy the use case.
- Use centralized theme/colors/typography.
- Do not scatter hardcoded colors.
- Do not use demo text/data as productive fallback.
- Do not show mock data when an API call fails.

## Remote states

Every screen or widget backed by remote data must handle:

- Loading.
- Empty.
- Error.
- Retry.
- Data.

Rules:

- Empty states should explain the situation in Spanish unless the project uses another language.
- Error states should use normalized user-facing API messages.
- Retry should re-run the relevant use case/provider, not reload unrelated app state.
- Skeletons/placeholders are acceptable when they match existing visual patterns.

## Forms

Forms should:

- Validate required fields before calling application logic.
- Keep validation messages user-friendly.
- Avoid duplicating validation rules already centralized in application/domain.
- Disable submit while a request is in progress when duplicate submissions would be harmful.
- Show API errors near the form or in the established app feedback pattern.

## Navigation

- Use GoRouter or the existing router.
- Centralize route paths.
- Prefer route helper methods for parameterized paths.
- Do not duplicate string paths across widgets.
- Protected routes should rely on centralized auth redirects/guards.

## Accessibility and UX basics

- Tap targets should be reasonably sized.
- Text should not overflow in common mobile widths.
- Important async operations should provide feedback.
- Destructive actions should require confirmation when data loss is possible.
- Do not block the UI with silent long-running operations.

## Images and assets

- Use backend image URLs when provided.
- If no image exists, use the current project placeholder.
- Do not use random network images as productive fallback.
- Keep asset paths centralized or close to the widget when asset is feature-specific.

## Quality bar

Before considering a Flutter UI task complete:

- Existing visual design is preserved.
- No productive mocks remain.
- No DTOs are imported by presentation.
- No direct HTTP calls exist in widgets.
- Loading/empty/error/retry are covered.
- Route usage follows the centralized router.
- Theme tokens are used where available.
- User-facing errors are friendly.
- The implementation works on small mobile widths.
