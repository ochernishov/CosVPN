# Реестр компонентов -- Mainline Next.js Template

Источник: https://mainline-nextjs-template.vercel.app/
Дата парсинга: 2026-03-20
Стек оригинала: Next.js 15 + React 19 + shadcn/ui + Tailwind CSS v4

## Сводка
- Страниц: 7 (/, /about, /pricing, /faq, /contact, /login, /privacy)
- Секций всего: 22
- Общие компоненты: Header, Footer (CTA + nav + big logo)
- Тема: СВЕТЛАЯ (warm beige primary oklch(92% .04 86.47))

## Страницы и секции

### Главная (/)

| # | Секция | Папка | shadcnblocks категория | Статус |
|---|--------|-------|----------------------|--------|
| 1 | Floating Pill Navbar | main-page/01-header/ | navbar | done |
| 2 | Hero (h1 + features + image) | main-page/02-hero/ | hero | done |
| 3 | Logo Cloud (marquee) | main-page/03-logo-cloud/ | logos | done |
| 4 | Features "Modern Teams" | main-page/04-features/ | feature | done |
| 5 | Bento Grid "Resource Allocation" | main-page/05-bento-grid/ | bento | done |
| 6 | Testimonials (carousel) | main-page/06-testimonials/ | testimonial | done |
| 7 | Pricing (3 cards) | main-page/07-pricing/ | pricing | done |
| 8 | FAQ (categorized accordion) | main-page/08-faq/ | faq | done |
| 9 | Footer (CTA + nav + big logo) | main-page/09-footer/ | footer + cta | done |

### About (/about)

| # | Секция | Папка | shadcnblocks категория | Статус |
|---|--------|-------|----------------------|--------|
| 1 | Hero + Stats (h1 + stats grid) | about/01-hero/ | about + stats | done |
| 2 | Team (images + text + button) | about/02-team/ | about | done |
| 3 | Investors (avatar grid) | about/03-investors/ | team | done |

### Pricing (/pricing)

| # | Секция | Папка | shadcnblocks категория | Статус |
|---|--------|-------|----------------------|--------|
| 1 | Pricing cards (3 tiers) | pricing/01-pricing-cards/ | pricing | done |
| 2 | Comparison table | pricing/02-comparison-table/ | compare | done |

### FAQ (/faq)

| # | Секция | Папка | shadcnblocks категория | Статус |
|---|--------|-------|----------------------|--------|
| 1 | FAQ (h1 + categorized accordion) | faq/01-faq/ | faq | done |
| 2 | Testimonials (carousel, reuse) | faq/02-testimonials/ | testimonial | done |

### Contact (/contact)

| # | Секция | Папка | shadcnblocks категория | Статус |
|---|--------|-------|----------------------|--------|
| 1 | Contact info (3-col) | contact/01-contact-info/ | contact | done |
| 2 | Contact form | contact/02-contact-form/ | contact | done |

### Login (/login)

| # | Секция | Папка | shadcnblocks категория | Статус |
|---|--------|-------|----------------------|--------|
| 1 | Login form (card) | login/01-login-form/ | login | done |

### Privacy (/privacy)

| # | Секция | Папка | shadcnblocks категория | Статус |
|---|--------|-------|----------------------|--------|
| 1 | MDX content page | privacy/ | content | done |

## Общие компоненты (используются на нескольких страницах)

| Компонент | Где используется | Папка с эталоном |
|-----------|-----------------|------------------|
| Header (Floating Pill Navbar) | ВСЕ страницы | main-page/01-header/ |
| Footer (CTA + nav + big logo) | ВСЕ страницы | main-page/09-footer/ |
| Testimonials carousel | /, /faq | main-page/06-testimonials/ |
| Pricing cards | /, /pricing | main-page/07-pricing/ |
| FAQ accordion | /, /faq | main-page/08-faq/ |
| Gradient wrapper top | /, /about, /contact, /login | globals.md |
| Gradient wrapper bottom | /, /pricing | globals.md |
| Dashed line separator | /, /contact | globals.md |

## Ключевые CSS классы шаблона (кастомные)

| Класс | Описание | Где используется |
|-------|----------|-----------------|
| `.container` | max-w 1220px, px 24px, mx auto | Все секции |
| `.dashed-line` | Horizontal dashed separator | Hero, Bento, Contact |
| `.gradient-section-top` | Top gradient wrapper | /, /about |
| `.gradient-section-bottom` | Bottom gradient wrapper | /, /pricing |
| `.section-padding` | py-28 lg:py-32 | Utility |
| `font-text` | Inter (body font) | Heading overrides |

## shadcn/ui компоненты

| Компонент | Использование |
|-----------|---------------|
| Accordion | FAQ sections |
| Button | CTA, forms, nav |
| Carousel (Embla) | Testimonials |
| Checkbox | Contact form |
| Input | Contact/Login forms |
| Label | Forms |
| NavigationMenu (Radix) | Header dropdown |
| Select | Contact form |
| Separator | Various |
| Switch | Pricing toggle |
| Textarea | Contact form |
