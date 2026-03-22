# Investors Section

## Описание
Секция "Our investors" с заголовком и grid из 5 инвесторов (аватар + имя + компания).

## Layout
- Wrapper: `pt-28 lg:pt-32`
- Grid: `mt-8 grid grid-cols-2 gap-12 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5`
- Gap: 48px (gap-12)

## Элементы

### H2 — "Our investors"
- Font: DM Sans 36px / 500 (lighter than other headings)

### Investor Cards (5)
| Name | Company | Image |
|------|---------|-------|
| Dennis Bouvard | Blackbird Ventures | /investors/1.webp |
| Renatus Gerard | Center Studies | /investors/2.webp |
| Leslie Alexander | TechNexus | /investors/3.webp |
| Matthew Stephens | Etymol Cap | /investors/4.webp |
| Josephine Newman | Vandenberg | /investors/5.webp |

Each card:
- Avatar: img, rounded, small (~80px)
- Name (h3): font-medium, foreground color
- Company (p): text-sm, muted-foreground

## Код компонента
```tsx
const investors = [
  { name: "Dennis Bouvard", company: "Blackbird Ventures", img: "/investors/1.webp" },
  { name: "Renatus Gerard", company: "Center Studies", img: "/investors/2.webp" },
  { name: "Leslie Alexander", company: "TechNexus", img: "/investors/3.webp" },
  { name: "Matthew Stephens", company: "Etymol Cap", img: "/investors/4.webp" },
  { name: "Josephine Newman", company: "Vandenberg", img: "/investors/5.webp" },
];

export function InvestorsSection() {
  return (
    <div className="pt-28 lg:pt-32">
      <h2 className="text-2xl font-medium tracking-tight md:text-4xl">Our investors</h2>
      <div className="mt-8 grid grid-cols-2 gap-12 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
        {investors.map((inv) => (
          <div key={inv.name}>
            <img src={inv.img} alt={inv.name} className="size-16 rounded-full object-cover" />
            <h3 className="mt-3 font-medium">{inv.name}</h3>
            <p className="text-sm text-muted-foreground">{inv.company}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
```
