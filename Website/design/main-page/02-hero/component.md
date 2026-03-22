# Hero Section

![Screenshot](../01-header/screenshot.png)

## Описание
Split-layout hero: слева H1 + subtitle + 2 CTA кнопки, справа 4 feature items с иконками. Ниже -- полноширинный screenshot приложения. Разделитель между колонками -- вертикальная dashed линия.

## shadcnblocks
Категория: **Hero** | Ближайший аналог: hero-1 (split) | URL: https://www.shadcnblocks.com/blocks/hero

## Section wrapper
```
Classes: py-28 lg:py-32 lg:pt-44
Padding: 176px top (desktop), 128px bottom
Inside gradient-section-top wrapper
```

## Container
```
Classes: container flex flex-col justify-between gap-8 md:gap-14 lg:flex-row lg:gap-20
max-width: 1220px
padding: 0 24px
display: flex (row on lg+)
gap: 80px (lg:gap-20)
```

## DOM структура
```
section.py-28.lg:py-32.lg:pt-44
  div.container.flex.flex-col.lg:flex-row.gap-20
    div.flex-1                          // Left column
      h1.text-5xl                       // Title
      p.text-3xl.text-muted-foreground  // Subtitle
      div.flex.gap-4                    // CTA buttons
        a (primary button)              // "Get template"
        a (outline button)             // "Built by shadcnblocks.com"
    div.relative.flex-1.flex-col.lg:pl-10  // Right column
      div.w-px.absolute (dashed vertical line)
      div.flex.gap-5 x4                 // Feature items
        svg (lucide icon)
        div
          h2.font-text.font-semibold    // Feature title
          p.text-sm.text-muted-foreground // Feature desc
  div.mt-24 (image wrapper)
    img.rounded-2xl.shadow-lg           // Hero screenshot
```

## Элементы

### H1
- **Classes:** `text-foreground max-w-160 text-3xl tracking-tight md:text-4xl lg:text-5xl xl:whitespace-nowrap`
- **Font:** DM Sans 48px (lg) / 36px (md) / 30px (base)
- **Weight:** 600 (semibold via display-weight)
- **Line-height:** 48px (1:1 ratio)
- **Letter-spacing:** -1.2px (tracking-tight at text-5xl)
- **Color:** oklch(14.5% 0 0) -- foreground
- **Max-width:** 640px (max-w-160)

### Subtitle (p)
- **Classes:** `text-muted-foreground text-1xl mt-5 md:text-3xl`
- **Font:** Inter 30px (md+) / 20px (base)
- **Weight:** 400
- **Line-height:** 36px
- **Color:** oklch(55.6% 0 0) -- muted-foreground
- **Margin-top:** 20px (mt-5)

### CTA "Get template" (Primary)
- **Classes:** `inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium`
- **Background:** oklch(92% .04 86.47) -- primary warm beige
- **Color:** oklch(31% .02 86.64) -- primary-foreground
- **Font:** Inter 14px / 500
- **Padding:** 8px 16px
- **Height:** 36px
- **Border-radius:** 6px
- **Transition:** all 0.15s cubic-bezier(0.4, 0, 0.2, 1)

### CTA "Built by shadcnblocks.com" (Outline)
- **Background:** oklch(100% 0 0) -- white
- **Color:** oklch(14.5% 0 0) -- foreground
- **Border:** 1px solid oklch(92.2% 0 0)
- **Box-shadow:** shadow-md
- **Padding:** 8px 12px
- **Height:** 38px
- **Has arrow icon:** ArrowRight (lucide)

### Dashed Vertical Line (separator left-right)
- **Wrapper:** `text-muted-foreground h-full w-px absolute top-0 left-0 max-lg:hidden`
- **Inner:** repeating-linear-gradient(180deg, transparent, transparent 4px, currentColor 4px, currentColor 10px)
- **Mask:** linear-gradient(180deg, transparent, black 10%, black 90%, transparent)
- **On mobile (<lg):** horizontal dashed line instead (top-0)

### Feature Items (x4: Tailored workflows, Cross-team projects, Milestones, Progress insights)
- **Container:** `div.flex.gap-2.5.lg:gap-5`
- **Gap:** 10px (mobile) / 20px (desktop)
- **Icon:** Lucide SVG 20x20px, stroke-width 2, color oklch(14.5% 0 0)
- **Title (h2):**
  - Classes: `font-text text-foreground font-semibold`
  - Font: Inter 16px / 600 (font-semibold)
  - Line-height: 24px
  - Color: oklch(14.5% 0 0)
- **Description (p):**
  - Classes: `text-muted-foreground max-w-76 text-sm`
  - Font: Inter 14px / 400
  - Line-height: 20px
  - Color: oklch(55.6% 0 0)
  - Max-width: 304px (max-w-76)
- **Spacing between items:** space-y-5 (20px)

### Hero Image
- **Classes:** `rounded-2xl object-cover object-left-top shadow-lg max-lg:rounded-tr-none`
- **Border-radius:** 16px (rounded-2xl), top-right 0 on mobile
- **Box-shadow:** shadow-lg
- **Object-fit:** cover
- **Object-position:** left top
- **Container:**
  - Classes: `mt-12 max-lg:ml-6 max-lg:h-[550px] max-lg:overflow-hidden md:mt-20 lg:container lg:mt-24`
  - Margin-top: 96px (lg:mt-24) / 80px (md:mt-20) / 48px (mt-12)
  - On mobile: left margin 24px, height capped at 550px, overflow hidden

## Responsive
- **Mobile (<lg):** Stack vertically, horizontal dashed line instead of vertical, hero image crops at 550px height
- **md:** text-4xl for h1, text-3xl for subtitle
- **lg:** Side-by-side layout, vertical dashed line, full image
- **xl:** h1 whitespace-nowrap
