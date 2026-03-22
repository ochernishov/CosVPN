# Карта блоков shadcnblocks -- Mainline Next.js Template

Источник: https://mainline-nextjs-template.vercel.app/
Дата парсинга: 2026-03-20
Репозиторий: https://github.com/shadcnblocks/mainline-nextjs-template

## ВАЖНО: Статус блоков

Mainline -- это **open-source шаблон от shadcnblocks** (MIT лицензия). Блоки в нём **кастомные**, написанные специально для шаблона. Они **НЕ** являются Pro-блоками из shadcnblocks registry и **НЕ** устанавливаются через CLI `npx shadcn@latest add`.

Однако, каждый блок визуально и структурно соответствует определённой **категории** блоков на shadcnblocks.com. Ниже -- маппинг секций шаблона на категории и ближайшие аналоги Pro-блоков.

## Маппинг секций на категории shadcnblocks

### Главная страница (/)

| # | Секция шаблона | Категория shadcnblocks | Ближайший Pro-аналог | URL категории |
|---|---------------|----------------------|---------------------|---------------|
| 1 | Floating Pill Navbar | Navbar | navbar-1 (floating pill) | https://www.shadcnblocks.com/blocks/navbar |
| 2 | Hero (h1 + features + image) | Hero | hero-1 (split layout with features) | https://www.shadcnblocks.com/blocks/hero |
| 3 | Logo Cloud (marquee) | Logos | logos-1 (marquee carousel) | https://www.shadcnblocks.com/blocks/logos |
| 4 | Features "Modern Teams" (badge + card grid) | Feature | feature-1 (card grid with images) | https://www.shadcnblocks.com/blocks/feature |
| 5 | Bento Grid "Resource Allocation" | Bento | bento-1 (2-col + 3-col layout) | https://www.shadcnblocks.com/blocks/bento |
| 6 | Testimonials (carousel) | Testimonial | testimonial-1 (carousel + photo) | https://www.shadcnblocks.com/blocks/testimonial |
| 7 | Pricing (3 cards) | Pricing | pricing-1 (3 tiers + toggle) | https://www.shadcnblocks.com/blocks/pricing |
| 8 | FAQ (categorized accordion) | Faq | faq-1 (multi-column accordion) | https://www.shadcnblocks.com/blocks/faq |
| 9 | Footer (CTA + nav + big logo) | Footer + CTA | footer-1 + cta-1 | https://www.shadcnblocks.com/blocks/footer |

### About (/about)

| # | Секция шаблона | Категория shadcnblocks | Ближайший Pro-аналог |
|---|---------------|----------------------|---------------------|
| 1 | Hero (h1 + stats) | About + Stats | about-1 + stats-1 |
| 2 | Team section (images + text) | About | about-2 (content + images) |
| 3 | Investors (avatar grid) | Team | team-1 (avatar + name + role) |

### Pricing (/pricing)

| # | Секция шаблона | Категория shadcnblocks | Ближайший Pro-аналог |
|---|---------------|----------------------|---------------------|
| 1 | Pricing cards (3 tiers) | Pricing | pricing-1 |
| 2 | Comparison table | Compare | compare-1 (feature comparison) |

### FAQ (/faq)

| # | Секция шаблона | Категория shadcnblocks | Ближайший Pro-аналог |
|---|---------------|----------------------|---------------------|
| 1 | FAQ (categorized accordion) | Faq | faq-1 |
| 2 | Testimonials (carousel) | Testimonial | testimonial-1 |

### Contact (/contact)

| # | Секция шаблона | Категория shadcnblocks | Ближайший Pro-аналог |
|---|---------------|----------------------|---------------------|
| 1 | Contact info (3-col) | Contact | contact-1 (info + form) |
| 2 | Contact form | Contact | contact-1 |

### Login (/login)

| # | Секция шаблона | Категория shadcnblocks | Ближайший Pro-аналог |
|---|---------------|----------------------|---------------------|
| 1 | Login form (card) | Login | login-1 (card + social) |

## shadcn/ui компоненты (используются в шаблоне)

Эти компоненты устанавливаются стандартно через shadcn CLI:

```bash
npx shadcn@latest add accordion
npx shadcn@latest add button
npx shadcn@latest add carousel
npx shadcn@latest add checkbox
npx shadcn@latest add input
npx shadcn@latest add label
npx shadcn@latest add navigation-menu
npx shadcn@latest add select
npx shadcn@latest add separator
npx shadcn@latest add switch
npx shadcn@latest add textarea
```

## Дополнительные зависимости

```bash
npm install embla-carousel-react embla-carousel-autoplay
npm install lucide-react
npm install tw-animate-css
npm install next-themes
```

## Как использовать Pro-блоки shadcnblocks

Если вы хотите использовать **платные Pro-блоки** shadcnblocks вместо кастомных (более полные, с вариациями):

```bash
# Требует подписки All Access на shadcnblocks.com
npx shadcn@latest add https://www.shadcnblocks.com/r/navbar-1
npx shadcn@latest add https://www.shadcnblocks.com/r/hero-1
npx shadcn@latest add https://www.shadcnblocks.com/r/logos-1
npx shadcn@latest add https://www.shadcnblocks.com/r/feature-1
npx shadcn@latest add https://www.shadcnblocks.com/r/bento-1
npx shadcn@latest add https://www.shadcnblocks.com/r/testimonial-1
npx shadcn@latest add https://www.shadcnblocks.com/r/pricing-1
npx shadcn@latest add https://www.shadcnblocks.com/r/faq-1
npx shadcn@latest add https://www.shadcnblocks.com/r/footer-1
npx shadcn@latest add https://www.shadcnblocks.com/r/cta-1
npx shadcn@latest add https://www.shadcnblocks.com/r/contact-1
npx shadcn@latest add https://www.shadcnblocks.com/r/login-1
npx shadcn@latest add https://www.shadcnblocks.com/r/about-1
npx shadcn@latest add https://www.shadcnblocks.com/r/stats-1
npx shadcn@latest add https://www.shadcnblocks.com/r/team-1
npx shadcn@latest add https://www.shadcnblocks.com/r/compare-1
```

**Примечание:** Точные имена блоков (hero-1, pricing-3 и т.д.) могут отличаться. Проверьте актуальный каталог на https://www.shadcnblocks.com/blocks
