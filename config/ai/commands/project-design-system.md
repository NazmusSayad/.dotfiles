# project-design-system

## Explore & Extract Design System

Scan the entire project and document every single design element that already exists:

**Typography**:

- Font families in use
- Font sizes (px, rem values)
- Font weights (100, 200, 300... 900)
- Line heights
- Letter spacing
- All text styles from display to caption

**Spacing**:

- Base unit size
- Scale progression (extra small, small, medium, large, extra large)
- Padding amounts
- Margin amounts
- Gaps between elements

**Colors**:

- Primary color and variations
- Secondary colors
- Neutrals/grays scale
- Status colors (success, error, warning, info)
- Text colors
- Background colors
- Border colors

**Visual Patterns**:

- Button looks and variations
- Input field appearance
- Card design
- Modal/dialog design
- Navigation design
- Menu/dropdown appearance
- Loading state designs
- Empty state designs
- Error state designs

**Other Design Elements**:

- Border radius amounts
- Shadow depths and blur
- Responsive breakpoints
- Animation speeds
- Spacing hierarchy

## CRITICAL RULE

**NEVER CREATE YOUR OWN DESIGN SYSTEM.**

**NEVER INVENT NEW DESIGN TOKENS.**

**NEVER DEVIATE FROM EXISTING DESIGN.**

Only use what exists in the project. Follow the extracted design system exactly. Every design decision must reference the documented system.

## Apply Design System

When designing new UI:

1. Use typography from extracted scale - match existing font sizes and weights
2. Use spacing from extracted scale - match existing padding and margins
3. Use colors from extracted palette - only colors from the system
4. Use border radius from extracted values
5. Use shadow patterns from extracted designs
6. Use responsive breakpoints from extracted system
7. Design multiple states matching existing patterns
8. Maintain visual consistency with all existing designs

## How to Use

- "Extract and document the design system"
- "Review the design system"
- "Show typography scale"
- "Show color palette"
- "Show spacing scale"
- "Apply design system to new design"
