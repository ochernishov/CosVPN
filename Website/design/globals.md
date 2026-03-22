# Глобальные дизайн-токены -- Mainline Next.js Template

Источник: https://mainline-nextjs-template.vercel.app/
Тема: СВЕТЛАЯ (warm beige/yellow oklch)

---

## Layout-ширины (КРИТИЧНО)

| Использование | max-width | Padding | Tailwind-класс |
|---------------|-----------|---------|----------------|
| `.container` (основной) | 1220px (76.25rem) | 0 24px | `container` (кастомный) |
| `.container.max-w-5xl` (pricing, faq) | 1024px | 0 24px | `container max-w-5xl` |
| `.container.max-w-7xl` (bento bottom row) | 1280px | 0 24px | `container max-w-7xl` |
| Hero h1 max-width | 640px | -- | `max-w-160` |
| Hero feature desc max-width | 304px | -- | `max-w-76` |
| Description paragraphs | 448px | -- | `max-w-md` |
| Description centered | 576px | -- | `max-w-xl` |
| Features grid | 896px | -- | `max-w-4xl` |

### Container CSS (кастомный, в @layer components)
```css
.container {
  max-width: 76.25rem; /* 1220px */
  margin-left: auto;
  margin-right: auto;
  padding-left: 1.5rem; /* 24px */
  padding-right: 1.5rem; /* 24px */
}
```

### Section padding
| Класс | Mobile | Desktop (lg:) |
|-------|--------|---------------|
| `py-28 lg:py-32` | 112px top+bottom | 128px top+bottom |
| `pb-28 lg:pb-32` | 0 top, 112px bottom | 0 top, 128px bottom |
| `py-28 lg:py-32 lg:pt-44` | 112px top+bottom | 176px top, 128px bottom |

### Gradient section wrappers
```css
/* Top gradient wrapper (hero -> bento) */
/* Classes: relative mx-2.5 mt-2.5 lg:mx-4 from-primary/50 rounded-t-4xl rounded-b-2xl bg-linear-to-b via-20% via-muted to-muted/80 */
margin: 10px 10px 0; /* mobile: mx-2.5 mt-2.5 */
margin: 10px 16px 0; /* desktop: lg:mx-4 */
border-radius: 32px 32px 16px 16px;
background: linear-gradient(to bottom, oklch(92% .04 86.47 / 0.5), oklch(97% 0 0) 20%, oklch(97% 0 0 / 0.8));

/* Bottom gradient wrapper (pricing + faq) */
/* Classes: relative mx-2.5 mt-2.5 lg:mx-4 from-background via-background to-primary/50 rounded-t-2xl rounded-b-4xl bg-linear-to-b */
margin: 10px 10px 0;
margin: 10px 16px 0; /* desktop */
border-radius: 16px 16px 32px 32px;
background: linear-gradient(to bottom, oklch(100% 0 0), oklch(100% 0 0), oklch(92% .04 86.47 / 0.5));
```

---

## Цветовая палитра (oklch)

### Основные цвета
| Имя | oklch | HEX (примерный) | Описание |
|-----|-------|-----------------|----------|
| `--background` | `oklch(100% 0 0)` | #FFFFFF | Белый фон |
| `--foreground` | `oklch(14.5% 0 0)` | #0A0A0A | Почти чёрный текст |
| `--primary` | `oklch(92% .04 86.47)` | #F5E6C8 | Warm beige/yellow (АКЦЕНТ) |
| `--primary-foreground` | `oklch(31% .02 86.64)` | #3D3526 | Тёмный на primary |
| `--secondary` | `oklch(97% 0 0)` | #F5F5F5 | Светло-серый |
| `--muted` | `oklch(97% 0 0)` | #F5F5F5 | Muted background |
| `--muted-foreground` | `oklch(55.6% 0 0)` | #737373 | Серый текст |
| `--border` | `oklch(92.2% 0 0)` | #E5E5E5 | Светлая граница |
| `--input` | `oklch(92.2% 0 0)` | #E5E5E5 | Input border |
| `--ring` | `oklch(70.8% 0 0)` | #A3A3A3 | Focus ring |
| `--card` | `oklch(100% 0 0)` | #FFFFFF | Card background |
| `--card-foreground` | `oklch(14.5% 0 0)` | #0A0A0A | Card text |
| `--destructive` | `oklch(57.7% .245 27.325)` | #DC2626 | Красный ошибки |

### Shadows
| Имя | Значение |
|-----|----------|
| `--shadow-2xs` | `0 1px 3px 0px #0000000d` |
| `--shadow-xs` | `0 1px 3px 0px #0000000d` |
| `--shadow-sm` | `0 1px 3px 0px #0000001a, 0 1px 2px -1px #0000001a` |
| `--shadow-md` | `0 1px 3px 0px #0000001a, 0 2px 4px -1px #0000001a` |
| `--shadow-lg` | `0 1px 3px 0px #0000001a, 0 4px 6px -1px #0000001a` |
| `--shadow-xl` | `0 1px 3px 0px #0000001a, 0 8px 10px -1px #0000001a` |
| `--shadow-2xl` | `0 1px 3px 0px #00000040` |

---

## Типографика

### Шрифты
| Назначение | Шрифт | CSS переменная | Шрифт загружается из |
|-----------|-------|----------------|---------------------|
| Display (h1-h6) | DM Sans | `--font-dm-sans` / `--display-family` | /fonts/dm-sans/ (local) |
| Body (p, a, span) | Inter | `--font-inter` / `--text-family` | next/font/google |
| Display weight | 600 (semibold) | `--display-weight` | |
| Body weight | 400 (regular) | `--text-weight` | |

### Размеры заголовков (desktop)
| Элемент | font-size | font-weight | line-height | letter-spacing | font-family |
|---------|-----------|-------------|-------------|----------------|-------------|
| H1 (hero, `text-5xl`) | 48px | 600 | 48px | -1.2px | DM Sans |
| H1 (about, `text-6xl`) | 60px | 600 | 60px | -1.5px | DM Sans |
| H2 (section title, `text-5xl`) | 48px | 600 | 48px | -1.2px | DM Sans |
| H2 (bento title, `text-6xl`) | 60px | 600 | 60px | -1.5px | DM Sans |
| H2 (logo cloud, `text-3xl`) | 30px | 600 | 36px | normal | DM Sans |
| H2 (sub-heading, `font-text font-semibold`) | 16px | 600 | 24px | normal | Inter |
| H3 (feature card, `text-2xl`) | 24px | 600 | 30px | -0.6px | DM Sans |
| H3 (bento card, `font-semibold`) | 16px | 600 | 24px | normal | Inter |
| H3 (pricing title) | 16px | 600 | 24px | normal | Inter |
| H3 (faq category) | 16px | 400 | 24px | normal | Inter (muted color) |

### Tailwind размеры текста
| Класс | font-size | line-height | Tracking |
|-------|-----------|-------------|----------|
| `text-sm` | 14px | 20px | normal |
| `text-base` | 16px | 24px | normal |
| `text-xl` | 20px | 28px | normal |
| `text-2xl` | 24px | 32px | normal |
| `text-3xl` | 30px | 36px | `tracking-tight` = -0.025em |
| `text-4xl` | 36px | 40px | `tracking-tight` |
| `text-5xl` | 48px | 48px | `tracking-tight` |
| `text-6xl` | 60px | 60px | `tracking-tight` |

### Body text
| Элемент | font-size | font-weight | line-height | color |
|---------|-----------|-------------|-------------|-------|
| Body | 16px | 400 | 24px | oklch(14.5% 0 0) |
| Hero subtitle (`text-3xl`) | 30px | 400 | 36px | oklch(55.6% 0 0) |
| About subtitle (`text-4xl`) | 36px | 400 | 40px | oklch(55.6% 0 0) |
| Description (`leading-snug`) | 16px | 400 | 22px | oklch(55.6% 0 0) |
| Small text (`text-sm`) | 14px | 400 | 20px | oklch(55.6% 0 0) |
| Nav link | 14px | 500 | 20px | oklch(14.5% 0 0) |
| Button text | 14px | 500 | 20px | varies |
| FAQ trigger | 14px | 500 | -- | oklch(14.5% 0 0) |
| Footer link | 16px | 500 | -- | oklch(14.5% 0 0) |
| Footer privacy | 14px | 400 | -- | oklch(55.6% 0 0) |

---

## Border Radius

| Tailwind | Значение |
|----------|----------|
| `rounded-sm` | 4px |
| `rounded-md` | 6px |
| `rounded-lg` | 8px |
| `rounded-xl` | 12px |
| `rounded-2xl` | 16px |
| `rounded-3xl` | 24px |
| `rounded-4xl` | 32px |
| `rounded-full` | 50% |

---

## Кнопки

### Primary Button (default variant)
```css
background: oklch(92% .04 86.47); /* --primary warm beige */
color: oklch(31% .02 86.64); /* --primary-foreground */
font: Inter 14px / font-weight: 500;
padding: 8px 16px; /* px-4 py-2 */
height: 36px; /* h-9 */
border-radius: 6px; /* rounded-md */
border: none;
transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
/* Hover: opacity 0.9 */
```

### Primary Large (footer CTA)
```css
/* Same as primary but h-10 px-6 */
height: 40px;
padding: 0 24px;
```

### Secondary Button (outline variant)
```css
background: oklch(100% 0 0); /* white */
color: oklch(14.5% 0 0); /* foreground */
font: Inter 14px / font-weight: 500;
padding: 8px 12px; /* px-3 py-2 */
height: 36px; /* h-9 */
border: 1px solid oklch(92.2% 0 0); /* border */
border-radius: 6px; /* rounded-md */
box-shadow: 0 1px 3px 0px rgba(0,0,0,0.1), 0 2px 4px -1px rgba(0,0,0,0.1); /* shadow-md */
transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
/* Hover: bg-accent */
```

### Ghost Button (nav items)
```css
background: transparent;
color: oklch(14.5% 0 0);
font: Inter 14px / font-weight: 500;
padding: 8px 6px;
height: 36px;
border-radius: 6px;
transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
/* Hover: bg-accent */
```

---

## Header / Floating Pill Navbar

### Desktop
```css
position: absolute; /* becomes fixed on scroll */
width: min(90%, 700px);
top: 48px; /* lg:top-12 */
left: 50%;
transform: translateX(-50%);
z-index: 50;
height: 62px;
border-radius: 32px; /* rounded-4xl */
background: oklch(100% 0 0 / 0.7); /* bg-background/70 */
backdrop-filter: blur(12px); /* backdrop-blur-md */
border: 1px solid oklch(92.2% 0 0); /* border */
transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); /* transition-all duration-300 */

/* Inner container */
display: flex;
align-items: center;
justify-content: space-between;
padding: 12px 24px; /* px-6 py-3 */
```

### Mobile
```css
top: 20px; /* top-5 */
```

### Dropdown (Features)
```css
width: 400px;
background: oklch(100% 0 0); /* bg-popover */
border: 1px solid oklch(92.2% 0 0);
border-radius: 6px;
box-shadow: shadow-sm;
padding: 0;
/* animate-in, zoom-in-95 */

/* Dropdown link item */
display: flex;
gap: 16px;
padding: 12px;
border-radius: 6px;
font: Inter 16px / weight 400;
transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
/* Hover: bg-accent */
```

---

## Dashed Line Separators

### Горизонтальная
```css
.dashed-line-horizontal {
  height: 1px;
  width: 100%;
  color: oklch(55.6% 0 0); /* text-muted-foreground */
  background: repeating-linear-gradient(
    90deg,
    transparent,
    transparent 4px,
    currentColor 4px,
    currentColor 10px
  );
  mask-image: linear-gradient(
    90deg,
    transparent,
    black 25%,
    black 75%,
    transparent
  );
}
```

### Вертикальная
```css
.dashed-line-vertical {
  width: 1px;
  height: 100%;
  color: oklch(55.6% 0 0);
  background: repeating-linear-gradient(
    180deg,
    transparent,
    transparent 4px,
    currentColor 4px,
    currentColor 10px
  );
  mask-image: linear-gradient(
    180deg,
    transparent,
    black 10%,
    black 90%,
    transparent
  );
}
```

---

## Анимации и Transitions

### Глобальные transitions
```
Кнопки/ссылки: all 0.15s cubic-bezier(0.4, 0, 0.2, 1)
Header: all 0.3s cubic-bezier(0.4, 0, 0.2, 1)
Accordion: 0.2s ease-out
```

### Keyframes
```css
@keyframes accordion-down {
  0% { height: 0px; }
  100% { height: var(--radix-accordion-content-height); }
}
@keyframes accordion-up {
  0% { height: var(--radix-accordion-content-height); }
  100% { height: 0px; }
}
@keyframes scroll {
  0% { transform: translateX(0%); }
  100% { transform: translateX(-50%); }
}
@keyframes scroll-reverse {
  0% { transform: translateX(-50%); }
  100% { transform: translateX(0%); }
}
```

### Marquee animations
```
--animate-marquee: scroll 40s linear infinite
--animate-marquee-reverse: scroll-reverse 40s linear infinite
--animate-marquee-slow: scroll 60s linear infinite
```

---

## Breakpoints

| Name | Width | Tailwind |
|------|-------|----------|
| sm | 640px | `sm:` |
| md | 768px | `md:` |
| lg | 1024px | `lg:` |
| xl | 1280px | `xl:` |
| 2xl | 1536px | `2xl:` |

---

## Карточки

### Standard Card (pricing, feature)
```css
background: oklch(100% 0 0); /* bg-card */
border: 1px solid oklch(92.2% 0 0); /* border */
border-radius: 12px; /* rounded-xl */
box-shadow: shadow-sm;
```

### Large Feature Card
```css
border-radius: 24px; /* rounded-3xl */
border: 1px solid oklch(92.2% 0 0);
background: oklch(100% 0 0);
box-shadow: shadow-sm;
margin-top: 80px; /* mt-20 */
```

### Testimonial Card
```css
border-radius: 12px; /* rounded-xl */
background: oklch(97% 0 0); /* bg-muted */
border: none;
height: 534px;
```

### Pricing Card (highlighted - Startup)
```css
/* Same as standard + */
outline: 4px solid oklch(92% .04 86.47); /* outline-primary outline-4 */
transform-origin: top;
```

---

## Carousel (Testimonials)

```css
/* Embla Carousel */
basis: 28% (lg:basis-[28%])
gap: 16px (pl-4)
overflow: hidden;

/* Navigation buttons */
width: 58px;
height: 58px;
border-radius: 50%; /* rounded-full */
border: 1px solid oklch(92.2% 0 0);
background: oklch(97% 0 0); /* bg-muted */
```

---

## Иконки

Библиотека: **Lucide React**
Размер по умолчанию: 20x20px
Stroke-width: 2
Color: oklch(14.5% 0 0) -- currentColor / foreground
Checkmark (pricing): oklch(14.5% 0 0)

---

## Footer

```css
display: flex;
flex-direction: column;
align-items: center;
gap: 56px; /* gap-14 */
padding-top: 128px; /* pt-28 lg:pt-32 */
```

### CTA section
```
h2: DM Sans 48px/600/-1.2px, text-center
p: Inter 16px/400, oklch(55.6% 0 0), max-w-xl, text-center
button: primary, h-10, px-6
margin-bottom: 12px between p and button
```

### Nav section
```
display: flex; flex-direction: column; gap: 16px;
Links: Inter 16px/500, oklch(14.5% 0 0)
Privacy: Inter 14px/400, oklch(55.6% 0 0)
```

### Big logo
```
Full-width wordmark "mainline" SVG at very bottom
opacity: ~0.1 (very faded)
```

---

## Stats (About page)

```css
/* Stat number */
font: DM Sans 48px / weight 600 / line-height 48px;
letter-spacing: 1.2px;
color: oklch(14.5% 0 0);

/* Stat label */
font: Inter 16px / weight 400;
color: oklch(55.6% 0 0);
```

---

## Form Elements (Contact page)

### Input
```css
border: 1px solid oklch(92.2% 0 0); /* border */
border-radius: 6px; /* rounded-md */
padding: 8px 12px;
font: Inter 14px;
background: oklch(100% 0 0);
/* Focus: ring oklch(70.8% 0 0) */
```

### Textarea
```css
/* Same as input */
min-height: 80px;
resize: vertical;
```

### Select
```css
/* Same as input + dropdown arrow */
```

### Checkbox
```css
width: 16px;
height: 16px;
border-radius: 4px;
border: 1px solid oklch(92.2% 0 0);
/* Checked: bg-primary */
```

### Submit button
```css
/* Primary button style */
background: oklch(92% .04 86.47);
```
