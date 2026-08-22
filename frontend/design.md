---
name: Akritas
colors:
    surface: "#131316"
    surface-dim: "#131316"
    surface-bright: "#39393c"
    surface-container-lowest: "#0e0e11"
    surface-container-low: "#1b1b1e"
    surface-container: "#1f1f22"
    surface-container-high: "#2a2a2d"
    surface-container-highest: "#353438"
    on-surface: "#e4e1e6"
    on-surface-variant: "#c4c7c8"
    inverse-surface: "#e4e1e6"
    inverse-on-surface: "#303033"
    outline: "#8e9192"
    outline-variant: "#444748"
    surface-tint: "#c6c6c7"
    primary: "#ffffff"
    on-primary: "#2f3131"
    primary-container: "#e2e2e2"
    on-primary-container: "#636565"
    inverse-primary: "#5d5f5f"
    secondary: "#c0c1ff"
    on-secondary: "#1000a9"
    secondary-container: "#3131c0"
    on-secondary-container: "#b0b2ff"
    tertiary: "#ffffff"
    on-tertiary: "#2f3131"
    tertiary-container: "#e2e2e2"
    on-tertiary-container: "#636565"
    error: "#ffb4ab"
    on-error: "#690005"
    error-container: "#93000a"
    on-error-container: "#ffdad6"
    primary-fixed: "#e2e2e2"
    primary-fixed-dim: "#c6c6c7"
    on-primary-fixed: "#1a1c1c"
    on-primary-fixed-variant: "#454747"
    secondary-fixed: "#e1e0ff"
    secondary-fixed-dim: "#c0c1ff"
    on-secondary-fixed: "#07006c"
    on-secondary-fixed-variant: "#2f2ebe"
    tertiary-fixed: "#e2e2e2"
    tertiary-fixed-dim: "#c6c6c7"
    on-tertiary-fixed: "#1a1c1c"
    on-tertiary-fixed-variant: "#454747"
    background: "#131316"
    on-background: "#e4e1e6"
    surface-variant: "#353438"
typography:
    display:
        fontFamily: Inter
        fontSize: 32px
        fontWeight: "600"
        lineHeight: 40px
        letterSpacing: -0.02em
    headline-sm:
        fontFamily: Inter
        fontSize: 18px
        fontWeight: "600"
        lineHeight: 28px
    body-md:
        fontFamily: Inter
        fontSize: 14px
        fontWeight: "400"
        lineHeight: 20px
    body-sm:
        fontFamily: Inter
        fontSize: 13px
        fontWeight: "400"
        lineHeight: 18px
    label-caps:
        fontFamily: Inter
        fontSize: 11px
        fontWeight: "600"
        lineHeight: 16px
        letterSpacing: 0.05em
    code-md:
        fontFamily: JetBrains Mono
        fontSize: 13px
        fontWeight: "400"
        lineHeight: 20px
rounded:
    sm: 0.125rem
    DEFAULT: 0.25rem
    md: 0.375rem
    lg: 0.5rem
    xl: 0.75rem
    full: 9999px
spacing:
    unit: 4px
    xs: 4px
    sm: 8px
    md: 16px
    lg: 24px
    xl: 32px
    gutter: 16px
    margin: 24px
---

## Brand & Style

The design system is engineered for **Akritas**, an autonomous production incident response platform. The visual narrative is built upon the "Safe" identity—a reliable, high-fidelity execution that mirrors the precision of modern infrastructure tools like Vercel and Linear.

The aesthetic is strictly **Dark Mode**, utilizing deep blacks and structural grays to create a focused, low-distraction environment for high-stakes incident remediation. It prioritizes information density without sacrificing clarity, using a developer-first approach that favors utility over ornamentation. The emotional response is one of calm control amidst systemic failure, achieved through rigorous alignment, subtle borders, and a systematic color language.

## Colors

This design system utilizes a monochromatic foundation with purposeful semantic highlights.

- **Primary Canvas:** The background uses `#09090b` for deep immersion, with surfaces elevated using `#18181b`.
- **Borders:** All structural divisions use a low-contrast `#27272a` to define hierarchy without visual noise.
- **Accents:** White is used for primary actions and high-contrast text. Indigo (`#6366f1`) is reserved exclusively for AI-driven insights and autonomous remediation features.
- **Semantic Logic:** Success (Emerald), Warning (Amber), and Error (Rose) are used for infrastructure status and incident severity, following industry standard patterns for immediate recognition.

## Typography

The typography system uses **Inter** for all UI elements to ensure maximum legibility and a neutral, systematic feel. It follows a tight scale to accommodate high-density data views.

**JetBrains Mono** is introduced as a secondary functional font. It must be used for all technical data, including:

- Server logs and stack traces.
- File paths and Git SHA hashes.
- Code diffs and configuration snippets.

Use `label-caps` for table headers and section titles in sidebars to differentiate them from interactive content.

## Layout & Spacing

The design system employs a **12-column fixed grid** for main dashboard views, maximizing screen real estate while maintaining rigorous alignment.

- **Sidebar:** Fixed at 240px width.
- **Density:** Elements use a 4px base unit. For technical views (incident logs), use `sm` (8px) padding to increase information density. For high-level overviews, use `md` (16px).
- **Responsive:** On mobile, the grid collapses to a single column. The sidebar transitions to a hidden drawer. Content padding reduces from 24px to 16px.

## Elevation & Depth

Depth is conveyed through **Tonal Layering** and **Low-Contrast Outlines** rather than traditional shadows. This ensures the UI feels like a single, integrated piece of hardware.

1. **Base (Level 0):** `#09090b` - The main application background.
2. **Surface (Level 1):** `#18181b` - Used for cards, sidebar, and navigation panels.
3. **Overlay (Level 2):** `#27272a` - Used for modals and dropdown menus, often with a 1px border of the same color or slightly lighter to define the edge against Level 1.

Do not use blurs or shadows. Separation is achieved purely through color value shifts and hairline borders.

## Shapes

In alignment with the "Safe" visual identity, this design system uses a **Soft** roundedness level (0.25rem).

- **Standard Elements:** Buttons, input fields, and tags use `0.25rem`.
- **Containers:** Large cards and modals use `rounded-lg` (0.5rem).
- **Status Pips:** Status indicators (e.g., "System Healthy") are perfectly circular.
- **Code Blocks:** Maintain a sharp or minimally rounded edge to preserve the technical aesthetic.

## Components

### Buttons & Inputs

Buttons are high-contrast. The primary button is white text on a subtle gray background or a solid white background for peak emphasis. Inputs use a 1px border (`#27272a`) and transition to a white border on focus.

### Status Indicators & Badges

Badges are compact and use low-saturation backgrounds with high-saturation text for semantic clarity.

- **Error:** Background: `rgba(244, 63, 94, 0.1)`, Text: `#f43f5e`.
- **AI/Intel:** Background: `rgba(99, 102, 241, 0.1)`, Text: `#a5b4fc`.

### Metrics & Incident Cards

Metric cards feature a large `headline-sm` value with a small monochrome sparkline. Incident cards are grouped by severity, using a vertical status line on the left edge of the card to indicate state (e.g., a 2px red line for "Critical").

### AI Diagnosis Cards

Specialized Indigo-themed components. They should feature a subtle "Intel" icon and use a slightly distinct background tint (e.g., `#1e1b4b`) to differentiate autonomous suggestions from manual logs.

### Code Diffs

Use JetBrains Mono. Deletions are marked with a muted red background; additions with a muted green. Text remains legible with high contrast.
