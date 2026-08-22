# Frontend UI Quality Policy

The UI must be consistent, understandable and resilient.

## Required remote-data states

Every screen or section depending on remote data should consider:

- loading
- empty state
- error state
- success state when applicable
- disabled state
- permission denied / restricted access when applicable

## Tables and lists

For operational lists:

- Use visible and clear filters.
- Support loading state.
- Support empty state.
- Support error state.
- Put row actions in a menu/dropdown when there is more than one action.
- Avoid too many inline buttons.
- Use pagination according to backend/OpenAPI contract.
- Keep columns consistent across modules.

## Forms

- Use visible labels.
- Placeholders do not replace labels.
- Show clear human error messages.
- Add loading state to submit buttons.
- Prevent double submit.
- Confirm destructive actions.
- Split large forms into sections.

## Icons

- Use Lucide Icons or the existing consistent icon set.
- Do not mix many icon families.
- Keep consistent stroke and base size.
- Pair icons with text when actions may not be obvious.
- Destructive actions should be visually clear.

## Charts and reports

When using charts:

- Use a modern maintainable chart library already present or approved by the project.
- Respect design tokens.
- Avoid saturated palettes without a reason.
- Prioritize readability.
- Show loading, empty and error states.
- Use clear tooltips.
- Do not pass raw backend DTOs directly to charts; map data first.
