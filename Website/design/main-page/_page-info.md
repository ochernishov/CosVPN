# Главная страница

- URL: https://mainline-nextjs-template.vercel.app/
- Title: Mainline - Modern Next.js Template
- Полный скриншот: [fullpage.png](./fullpage.png)

## Секции (9)
1. Header (Floating Navbar) — floating pill, backdrop-blur
2. Hero — заголовок + фичи + CTA + product screenshot
3. Logo Cloud — marquee логотипов партнёров
4. Features "Made for modern product teams" — badge + 3 карточки
5. Bento Grid "Resource Allocation" — 2+3 карточки
6. Testimonials "Trusted by product builders" — carousel
7. Pricing — 3 тарифных карточки
8. FAQ "Got Questions?" — accordion по категориям
9. Footer — CTA + nav + decorative logo

## Layout-обёртки
Страница использует 2 gradient wrapper div:
- Верхний: `from-primary/50 via-muted to-muted/80 rounded-t-4xl rounded-b-2xl` (секции 1-5)
- Разделитель: dashed line
- Нижний: `from-background via-background to-primary/50 rounded-t-2xl rounded-b-4xl` (секции 7-8)
- Секция 6 (testimonials) между обёртками
