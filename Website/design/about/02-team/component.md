# Team Section

## Описание
Двухколоночная секция о команде: слева — текст + кнопка "View open roles", справа — 4 фотографии в grid. Также есть вторая подсекция с текстом справа и фото слева.

## Layout
- Section classes: `container mt-10 flex max-w-5xl flex-col-reverse gap-8 md:mt-14 md:gap-14 lg:mt-20 lg:flex-row lg:items-end`
- Container: max-w-5xl (1024px)
- Flex: column-reverse mobile, row desktop
- Gap: 32px / 56px

## Элементы

### H2 — "The team"
- Font: DM Sans 36px / 600

### Body Text (3 paragraphs)
- Font: Inter 16px / 400
- Color: muted-foreground

### "View open roles" Button
- Background: primary (oklch(0.92 0.04 86.47))
- Color: primary-foreground
- Border-radius: 6px
- Padding: 0 24px
- Links to /careers

### Images
4 team photos:
- `/about/1.webp` — "Team collaboration"
- `/about/2.webp` — "Team workspace"
- `/about/3.webp` — "Modern workspace"
- `/about/4.webp` — "Team collaboration"
- Rounded corners
- Arranged in asymmetric grid/layout

## Код компонента
```tsx
import Link from "next/link";
import { Button } from "@/components/ui/button";

export function TeamSection() {
  return (
    <section className="container mt-10 flex max-w-5xl flex-col-reverse gap-8 md:mt-14 md:gap-14 lg:mt-20 lg:flex-row lg:items-end">
      <div className="flex-1">
        <div className="grid grid-cols-2 gap-4">
          <img src="/about/1.webp" alt="Team collaboration" className="rounded-xl object-cover" />
          <img src="/about/2.webp" alt="Team workspace" className="rounded-xl object-cover" />
        </div>
        <h2 className="mt-8 text-2xl font-semibold tracking-tight md:text-4xl">The team</h2>
        <div className="text-muted-foreground mt-4 space-y-4">
          <p>We started building Mainline in 2019 and launched in 2022...</p>
          <p>We are 100% founder and team-owned, profitable...</p>
          <p>If you're interested in building the future of PM, check out our open roles below.</p>
        </div>
        <Button asChild className="mt-6">
          <Link href="/careers">View open roles</Link>
        </Button>
      </div>
      <div className="flex-1 space-y-4">
        <div className="text-muted-foreground space-y-4">
          <p>At Mainline, we are dedicated to transforming...</p>
          <p>We're customer-obsessed...</p>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <img src="/about/3.webp" alt="Modern workspace" className="rounded-xl object-cover" />
          <img src="/about/4.webp" alt="Team collaboration" className="rounded-xl object-cover" />
        </div>
      </div>
    </section>
  );
}
```
