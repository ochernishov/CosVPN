import Image from "next/image";

import {
  ArrowRight,
  Globe,
  Lock,
  Shield,
  Zap,
} from "lucide-react";

import { DashedLine } from "@/components/dashed-line";
import { Button } from "@/components/ui/button";

const features = [
  {
    title: "AES-256 encryption",
    description: "Military-grade encryption protects all your traffic.",
    icon: Lock,
  },
  {
    title: "Zero-log policy",
    description: "We never store, track, or share your browsing data.",
    icon: Shield,
  },
  {
    title: "50+ server locations",
    description: "Connect from anywhere with our global server network.",
    icon: Globe,
  },
  {
    title: "99.99% uptime SLA",
    description: "Enterprise-grade reliability your team can count on.",
    icon: Zap,
  },
];

export const Hero = () => {
  return (
    <section className="py-28 lg:py-32 lg:pt-44">
      <div className="container flex flex-col justify-between gap-8 md:gap-14 lg:flex-row lg:gap-20">
        {/* Left side - Main content */}
        <div className="flex-1">
          <h1 className="text-foreground text-3xl tracking-tight md:text-4xl lg:text-5xl">
            Enterprise VPN for modern businesses
          </h1>

          <p className="text-muted-foreground text-1xl mt-5 md:text-3xl">
            {process.env.NEXT_PUBLIC_SITE_NAME} delivers zero-trust network access with military-grade
            encryption. Built for teams that need speed, security, and
            centralized control.
          </p>

          <div className="mt-8 flex flex-wrap items-center gap-4 lg:flex-nowrap">
            <Button asChild>
              <a href="/signup">
                Start free trial
              </a>
            </Button>
            <Button
              variant="outline"
              className="from-background dark:from-card h-auto gap-2 bg-linear-to-r to-transparent shadow-md dark:shadow-white/5"
              asChild
            >
              <a
                href="/pricing"
                className="max-w-56 truncate text-start md:max-w-none"
              >
                View pricing
                <ArrowRight className="stroke-3" />
              </a>
            </Button>
          </div>
        </div>

        {/* Right side - Features */}
        <div className="relative flex flex-1 flex-col justify-center space-y-5 max-lg:pt-10 lg:pl-10">
          <DashedLine
            orientation="vertical"
            className="absolute top-0 left-0 max-lg:hidden"
          />
          <DashedLine
            orientation="horizontal"
            className="absolute top-0 lg:hidden"
          />
          {features.map((feature) => {
            const Icon = feature.icon;
            return (
              <div key={feature.title} className="flex gap-2.5 lg:gap-5">
                <Icon className="text-foreground mt-1 size-4 shrink-0 lg:size-5" />
                <div>
                  <h2 className="font-text text-foreground font-semibold">
                    {feature.title}
                  </h2>
                  <p className="text-muted-foreground max-w-76 text-sm">
                    {feature.description}
                  </p>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      <div className="mt-10 max-lg:ml-6 max-lg:h-[550px] max-lg:overflow-hidden md:mt-14 lg:container lg:mt-16">
        <div className="relative h-[793px] w-full">
          <Image
            src="/hero.webp"
            alt="CosVPN dashboard interface"
            fill
            className="rounded-2xl object-cover object-left-top shadow-lg max-lg:rounded-tr-none"
          />
        </div>
      </div>
    </section>
  );
};
