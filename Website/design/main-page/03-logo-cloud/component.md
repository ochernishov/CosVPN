# Logo Cloud (Marquee)

## Описание
Секция с логотипами партнёров/клиентов. Заголовок по центру, два ряда логотипов (верхний и нижний) с marquee анимацией в разных направлениях.

## shadcnblocks
Категория: **Logos** | Ближайший аналог: logos-1 | URL: https://www.shadcnblocks.com/blocks/logos

## Section
```
Classes: pb-28 lg:pb-32 overflow-hidden
Padding: 0 top, 128px bottom
Overflow: hidden (для marquee)
```

## Container
```
Classes: container space-y-10 lg:space-y-16
max-width: 1220px, padding: 0 24px
Spacing: 40px (mobile) / 64px (desktop) between title and logos
```

## H2
- **Classes:** `mb-4 text-xl text-balance md:text-2xl lg:text-3xl`
- **Font:** DM Sans 30px / 600
- **Line-height:** 36px
- **Color:** oklch(14.5% 0 0)
- **Text-align:** center
- **Content:** "Powering the world's best product teams." + muted line "From next-gen startups to established enterprises."

## Logo rows
- **Layout:** flex w-full flex-col items-center gap-8
- **Gap between logos:** 32px
- **Logos:** Mercury, Watershed, Retool, Descript (row 1); Perplexity, Monzo, Ramp, Raycast, Arc (row 2)
- **Logo opacity:** 0.5 (50%)
- **Logo sizes:** ~110-155px width, ~20-32px height
- **Marquee animation:** scroll 40s linear infinite / scroll-reverse 40s linear infinite
