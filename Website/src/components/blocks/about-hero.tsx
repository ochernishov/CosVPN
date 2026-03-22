import { DashedLine } from "@/components/dashed-line";

const stats = [
  {
    value: "50+",
    label: "Countries",
  },
  {
    value: "10K+",
    label: "Businesses",
  },
  {
    value: "99.99%",
    label: "Uptime SLA",
  },
  {
    value: "500K+",
    label: "Protected devices",
  },
];

export function AboutHero() {
  return (
    <section className="">
      <div className="container flex max-w-5xl flex-col justify-between gap-8 md:gap-20 lg:flex-row lg:items-center lg:gap-24 xl:gap-24">
        <div className="flex-[1.5]">
          <h1 className="text-3xl tracking-tight sm:text-4xl md:text-5xl lg:text-6xl">
            Securing businesses worldwide
          </h1>

          <p className="text-muted-foreground mt-5 text-2xl md:text-3xl lg:text-4xl">
            {process.env.NEXT_PUBLIC_SITE_NAME} is building the most secure and fastest enterprise VPN platform.
          </p>

          <p className="text-muted-foreground mt-8 hidden max-w-lg space-y-6 text-lg text-balance md:block lg:mt-12">
            At {process.env.NEXT_PUBLIC_SITE_NAME}, we are dedicated to transforming how businesses protect
            their networks and data. Our mission is to provide enterprise-grade
            security without compromising on speed or usability. We combine
            cutting-edge WireGuard protocol with zero-trust architecture to
            deliver the fastest, most secure VPN experience available.
            <br />
            <br />
            We&apos;re customer-obsessed — investing the time to understand every
            aspect of your security needs so that we can help you protect your
            organization better than ever before. Your security is our priority.
            In our history as a company, we&apos;ve maintained a 99.99% uptime
            record, because when your network is secure, your business thrives.
          </p>
        </div>

        <div
          className={`relative flex flex-1 flex-col justify-center gap-3 pt-10 lg:pt-0 lg:pl-10`}
        >
          <DashedLine
            orientation="vertical"
            className="absolute top-0 left-0 max-lg:hidden"
          />
          <DashedLine
            orientation="horizontal"
            className="absolute top-0 lg:hidden"
          />
          {stats.map((stat) => (
            <div key={stat.label} className="flex flex-col gap-1">
              <div className="font-display text-4xl tracking-wide md:text-5xl">
                {stat.value}
              </div>
              <div className="text-muted-foreground">{stat.label}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
