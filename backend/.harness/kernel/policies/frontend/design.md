# Frontend Design Policy

Before implementing UI, check for design guidance.

Priority:

1. `src/features/<feature>/DESIGN.md`
2. `src/features/<feature>/design.md`
3. `DESIGN.md`
4. `design.md`
5. Existing UI patterns

If a design file exists, treat it as a visual and UX contract.

## Mandatory rules

- Follow layout, spacing, colors, typography, interaction patterns and component guidelines from the design file.
- Do not introduce contradictory visual patterns.
- If the requirement conflicts with design guidance, stop and report the conflict.
- If something is not covered, infer the smallest consistent solution and document it.
- Do not invent a new visual language.
- Avoid scattered hardcoded colors.
- Use design tokens from a central location when present.

## Default visual direction

Unless a project-specific design file says otherwise:

- Modern, clean, minimal dashboard style.
- White and soft gray surfaces.
- Controlled primary color usage.
- Avoid saturating screens with the primary color.
- Use semantic colors for success, warning, error and info.
- Use icons consistently.

Do not hardcode colors everywhere. Prefer tokens in one of these locations when available:

```text
src/core/ui/tokens.ts
src/core/ui/theme.css
src/core/ui/design-tokens.ts
```

## Minimum token expectations

When defining or extending design tokens, cover:

- primary
- primary hover/active
- background
- surface
- border
- text primary
- text secondary
- muted
- success
- warning
- danger/error
- info
- spacing
- border radius
- shadows
- font sizes
- font weights
- layout sizes
- z-index when needed

## UI states

Any screen or section that depends on remote data should consider:

- loading
- empty state
- error state
- success state when applicable
- disabled state
- permission denied / restricted access when applicable
