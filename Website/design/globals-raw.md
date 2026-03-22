# Сырые CSS-переменные -- Mainline Next.js Template

## :root переменные (извлечены через getComputedStyle)
```json
{
  "--background": "oklch(100% 0 0)",
  "--foreground": "oklch(14.5% 0 0)",
  "--card": "oklch(100% 0 0)",
  "--card-foreground": "oklch(14.5% 0 0)",
  "--popover": "oklch(100% 0 0)",
  "--popover-foreground": "oklch(14.5% 0 0)",
  "--primary": "oklch(92% .04 86.47)",
  "--primary-foreground": "oklch(31% .02 86.64)",
  "--secondary": "oklch(97% 0 0)",
  "--secondary-foreground": "oklch(20.5% 0 0)",
  "--muted": "oklch(97% 0 0)",
  "--muted-foreground": "oklch(55.6% 0 0)",
  "--accent": "oklch(97% 0 0)",
  "--accent-foreground": "oklch(20.5% 0 0)",
  "--destructive": "oklch(57.7% .245 27.325)",
  "--border": "oklch(92.2% 0 0)",
  "--input": "oklch(92.2% 0 0)",
  "--ring": "oklch(70.8% 0 0)",
  "--radius": "8px",
  "--shadow-2xs": "0 1px 3px 0px #0000000d",
  "--shadow-xs": "0 1px 3px 0px #0000000d",
  "--shadow-sm": "0 1px 3px 0px #0000001a,0 1px 2px -1px #0000001a",
  "--shadow-md": "0 1px 3px 0px #0000001a,0 2px 4px -1px #0000001a",
  "--shadow-lg": "0 1px 3px 0px #0000001a,0 4px 6px -1px #0000001a",
  "--shadow-xl": "0 1px 3px 0px #0000001a,0 8px 10px -1px #0000001a",
  "--shadow-2xl": "0 1px 3px 0px #00000040",
  "--font-dm-sans": "\"dmSans\",sans-serif",
  "--font-inter": "\"inter\",sans-serif",
  "--font-sans": "var(--font-dm-sans),var(--font-inter),ui-sans-serif,system-ui,sans-serif",
  "--font-mono": "ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,\"Liberation Mono\",\"Courier New\",monospace",
  "--display-family": "var(--font-sans)",
  "--text-family": "var(--font-inter)",
  "--display-weight": "600",
  "--text-weight": "400",
  "--chart-1": "oklch(64.6% .222 41.116)",
  "--chart-2": "oklch(60% .118 184.704)",
  "--chart-3": "oklch(39.8% .07 227.392)",
  "--chart-4": "oklch(82.8% .189 84.429)",
  "--chart-5": "oklch(76.9% .188 70.08)",
  "--sidebar": "oklch(98.5% 0 0)",
  "--sidebar-foreground": "oklch(14.5% 0 0)",
  "--sidebar-primary": "oklch(20.5% 0 0)",
  "--sidebar-primary-foreground": "oklch(98.5% 0 0)",
  "--sidebar-accent": "oklch(97% 0 0)",
  "--sidebar-accent-foreground": "oklch(20.5% 0 0)",
  "--sidebar-border": "oklch(92.2% 0 0)",
  "--sidebar-ring": "oklch(70.8% 0 0)"
}
```

## @theme inline (Tailwind v4 token mapping)
```css
@theme inline {
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-card: var(--card);
  --color-card-foreground: var(--card-foreground);
  --color-popover: var(--popover);
  --color-popover-foreground: var(--popover-foreground);
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-foreground);
  --color-secondary: var(--secondary);
  --color-secondary-foreground: var(--secondary-foreground);
  --color-muted: var(--muted);
  --color-muted-foreground: var(--muted-foreground);
  --color-accent: var(--accent);
  --color-accent-foreground: var(--accent-foreground);
  --color-destructive: var(--destructive);
  --color-border: var(--border);
  --color-input: var(--input);
  --color-ring: var(--ring);

  --radius-sm: calc(var(--radius) - 4px);
  --radius-md: calc(var(--radius) - 2px);
  --radius-lg: var(--radius);
  --radius-xl: calc(var(--radius) + 4px);
  --radius-2xl: calc(var(--radius) + 8px);
  --radius-4xl: 2rem;

  --font-sans: var(--font-sans);
  --font-mono: var(--font-mono);

  --animate-accordion-down: accordion-down 0.2s ease-out;
  --animate-accordion-up: accordion-up 0.2s ease-out;
  --animate-marquee: scroll 40s linear infinite;
  --animate-marquee-reverse: scroll-reverse 40s linear infinite;
  --animate-marquee-slow: scroll 60s linear infinite;
}
```

## Keyframes
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

## @layer base
```css
@layer base {
  * {
    @apply border-border outline-ring/50;
    box-sizing: border-box;
  }
  body {
    @apply bg-background text-foreground;
    font-family: var(--text-family);
  }
  h1, h2, h3, h4, h5, h6 {
    font-family: var(--display-family);
    font-weight: var(--display-weight);
  }
}
```

## @layer components (кастомные)
```css
@layer components {
  .container {
    max-width: 76.25rem;
    margin-left: auto;
    margin-right: auto;
    padding-left: 1.5rem;
    padding-right: 1.5rem;
  }

  .section-padding {
    padding-top: 7rem;
    padding-bottom: 7rem;
  }
  @media (min-width: 1024px) {
    .section-padding {
      padding-top: 8rem;
      padding-bottom: 8rem;
    }
  }

  .dashed-line {
    height: 1px;
    width: 100%;
    color: var(--muted-foreground);
    background: repeating-linear-gradient(90deg, transparent, transparent 4px, currentColor 4px, currentColor 10px);
    mask-image: linear-gradient(90deg, transparent, black 25%, black 75%, transparent);
  }
}
```

## Computed font styles (по элементам)
```json
{
  "h1": {
    "fontFamily": "dmSans, sans-serif, inter, sans-serif, ui-sans-serif, system-ui, sans-serif",
    "fontSize": "48px",
    "fontWeight": "600",
    "lineHeight": "48px",
    "letterSpacing": "-1.2px",
    "color": "oklch(0.145 0 0)"
  },
  "h2-section": {
    "fontFamily": "dmSans, sans-serif",
    "fontSize": "48px",
    "fontWeight": "600",
    "lineHeight": "48px",
    "letterSpacing": "-1.2px",
    "color": "oklch(0.145 0 0)"
  },
  "h2-bento": {
    "fontSize": "60px",
    "fontWeight": "600",
    "lineHeight": "60px",
    "letterSpacing": "-1.5px"
  },
  "h2-logo": {
    "fontSize": "30px",
    "fontWeight": "600",
    "lineHeight": "36px",
    "letterSpacing": "normal"
  },
  "h3-card": {
    "fontSize": "24px",
    "fontWeight": "600",
    "lineHeight": "30px",
    "letterSpacing": "-0.6px"
  },
  "h3-bento-card": {
    "fontFamily": "inter, sans-serif",
    "fontSize": "16px",
    "fontWeight": "600",
    "lineHeight": "24px"
  },
  "body": {
    "fontFamily": "inter, sans-serif",
    "fontSize": "16px",
    "fontWeight": "400",
    "lineHeight": "24px",
    "color": "oklch(0.145 0 0)"
  },
  "p-muted": {
    "fontSize": "16px",
    "fontWeight": "400",
    "lineHeight": "22px",
    "color": "oklch(0.556 0 0)"
  },
  "p-small": {
    "fontSize": "14px",
    "fontWeight": "400",
    "lineHeight": "20px",
    "color": "oklch(0.556 0 0)"
  },
  "nav-link": {
    "fontSize": "14px",
    "fontWeight": "500",
    "color": "oklch(0.145 0 0)"
  },
  "button": {
    "fontSize": "14px",
    "fontWeight": "500",
    "lineHeight": "20px"
  },
  "blockquote": {
    "fontFamily": "dmSans, sans-serif",
    "fontSize": "24px",
    "fontWeight": "600",
    "lineHeight": "24px"
  },
  "stat-number": {
    "fontSize": "48px",
    "fontWeight": "600",
    "lineHeight": "48px",
    "letterSpacing": "1.2px"
  }
}
```

## Header element computed styles
```json
{
  "navbar-pill": {
    "position": "absolute",
    "width": "700px",
    "height": "62px",
    "top": "48px",
    "background": "oklab(1 0 0 / 0.7)",
    "backdropFilter": "blur(12px)",
    "border": "1px solid oklch(0.922 0 0)",
    "borderRadius": "32px",
    "transition": "0.3s cubic-bezier(0.4, 0, 0.2, 1)",
    "zIndex": "50"
  },
  "inner-container": {
    "display": "flex",
    "alignItems": "center",
    "justifyContent": "space-between",
    "padding": "12px 24px"
  },
  "logo-image": {
    "width": "94px",
    "height": "17.8px"
  },
  "nav-link-transition": "color 0.15s cubic-bezier(0.4, 0, 0.2, 1), background-color 0.15s ...",
  "dropdown": {
    "width": "400px",
    "background": "oklch(1 0 0)",
    "border": "1px solid oklch(0.922 0 0)",
    "borderRadius": "6px",
    "boxShadow": "shadow-sm"
  }
}
```

## Logo cloud images (opacity)
```json
{
  "Mercury": { "width": "143px", "height": "26px", "opacity": "0.5" },
  "Watershed": { "width": "154px", "height": "31.6px", "opacity": "0.5" },
  "Retool": { "width": "113px", "height": "21.8px", "opacity": "0.5" },
  "Descript": { "width": "112px", "height": "27px", "opacity": "0.5" }
}
```
