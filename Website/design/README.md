# Дизайн-система: Mainline Next.js Template

## Источник
- URL: https://mainline-nextjs-template.vercel.app/
- Repo: https://github.com/shadcnblocks/mainline-nextjs-template
- Дата парсинга: 2026-03-20
- Автор шаблона: shadcnblocks.com (Rob Austin, Callum Flack, Yassine Zaanouni)

## Стек оригинала
- Next.js 15 (App Router)
- React 19
- Tailwind CSS 4
- shadcn/ui
- TypeScript
- Embla Carousel
- Lucide React icons
- next-themes (dark/light mode)
- DM Sans + Inter fonts

## Тема
**СВЕТЛАЯ** -- warm beige/yellow oklch primary. Это НЕ тёмная тема!
- Primary accent: oklch(92% .04 86.47) -- warm beige/sand
- Background: oklch(100% 0 0) -- pure white
- Text: oklch(14.5% 0 0) -- near black
- Muted: oklch(55.6% 0 0) -- gray

## Страницы
| Страница | URL | Секций | Уникальные блоки |
|----------|-----|--------|-----------------|
| Главная | / | 9 | Hero, Logo Cloud, Features, Bento, Testimonials, Pricing, FAQ |
| О нас | /about | 3 | Hero+Stats, Team, Investors |
| Цены | /pricing | 2 | Pricing Cards, Comparison Table |
| FAQ | /faq | 2 | FAQ, Testimonials (reuse) |
| Контакты | /contact | 2 | Contact Info, Contact Form |
| Вход | /login | 1 | Login Form |
| Приватность | /privacy | 1 | MDX Content |

## Ключевые файлы
| Файл | Описание |
|------|----------|
| [globals.md](./globals.md) | Полная дизайн-система (цвета, типографика, layout, кнопки, анимации) |
| [globals-raw.md](./globals-raw.md) | Сырые CSS переменные и computed styles |
| [components.md](./components.md) | Центральный реестр всех компонентов |
| [SHADCNBLOCKS-MAP.md](./SHADCNBLOCKS-MAP.md) | Маппинг секций на блоки shadcnblocks.com |
| [_sitemap.md](./_sitemap.md) | Карта сайта и навигация |

## Дизайн-токены (краткая сводка)
- **Шрифты:** DM Sans (display, 600), Inter (body, 400/500)
- **Container:** max-width 1220px, padding 0 24px
- **Section padding:** py-28 lg:py-32 (112px / 128px)
- **Border radius:** 6px (buttons) / 12px (cards) / 24px (large cards) / 32px (navbar, wrappers)
- **Shadows:** Very subtle rgba(0,0,0,0.1) shadows
- **Transitions:** 0.15s cubic-bezier(0.4, 0, 0.2, 1) for interactions
- **Header:** Floating pill, w=700px, backdrop-blur(12px), rounded-4xl

## Gradient Wrappers (уникальная особенность)
Шаблон оборачивает группы секций в rounded divs с градиентами:
- **Top wrapper** (hero->bento): rounded-t-4xl rounded-b-2xl, from-primary/50 via-muted to-muted/80
- **Bottom wrapper** (pricing+faq): rounded-t-2xl rounded-b-4xl, from-background via-background to-primary/50
- Margin: mx-2.5 mt-2.5 lg:mx-4 (10px mobile, 16px desktop)

## Dashed Line Separators (уникальная особенность)
Кастомные dashed линии через repeating-linear-gradient + mask-image. Используются как вертикальные и горизонтальные разделители в hero и bento секциях.

## shadcnblocks
Шаблон Mainline -- **open-source** от shadcnblocks (MIT). Блоки кастомные, НЕ из Pro registry. Подробный маппинг на категории Pro-блоков: [SHADCNBLOCKS-MAP.md](./SHADCNBLOCKS-MAP.md)
