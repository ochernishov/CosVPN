import Image from "next/image";

import { ArrowRight } from "lucide-react";

import { DashedLine } from "../dashed-line";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
import { cn } from "@/lib/utils";

const items = [
  {
    quote: "CosVPN cut our security incident response time by 60%. The centralized dashboard is a game-changer.",
    author: "Amy Chase",
    role: "CTO",
    company: "Nextera Systems",
    image: "/testimonials/amy-chase.webp",
  },
  {
    quote: "We deployed CosVPN across 500 endpoints in under an hour. The admin controls are incredibly intuitive.",
    author: "Jonas Kotara",
    role: "CISO",
    company: "Vertex Capital",
    image: "/testimonials/jonas-kotara.webp",
  },
  {
    quote: "Finally, a VPN that doesn't make my engineering team complain about speed. WireGuard performance is outstanding.",
    author: "Kevin Yam",
    role: "VP Engineering",
    company: "Meridian Labs",
    image: "/testimonials/kevin-yam.webp",
  },
  {
    quote: "The zero-log policy and SOC 2 compliance made our security audit a breeze. Highly recommended.",
    author: "Kundo Marta",
    role: "Head of Security",
    company: "Arclight Partners",
    image: "/testimonials/kundo-marta.webp",
  },
  {
    quote: "CosVPN cut our security incident response time by 60%. The centralized dashboard is a game-changer.",
    author: "Amy Chase",
    role: "CTO",
    company: "Nextera Systems",
    image: "/testimonials/amy-chase.webp",
  },
  {
    quote: "We deployed CosVPN across 500 endpoints in under an hour. The admin controls are incredibly intuitive.",
    author: "Jonas Kotara",
    role: "CISO",
    company: "Vertex Capital",
    image: "/testimonials/jonas-kotara.webp",
  },
  {
    quote: "Finally, a VPN that doesn't make my engineering team complain about speed. WireGuard performance is outstanding.",
    author: "Kevin Yam",
    role: "VP Engineering",
    company: "Meridian Labs",
    image: "/testimonials/kevin-yam.webp",
  },
  {
    quote: "The zero-log policy and SOC 2 compliance made our security audit a breeze. Highly recommended.",
    author: "Kundo Marta",
    role: "Head of Security",
    company: "Arclight Partners",
    image: "/testimonials/kundo-marta.webp",
  },
];

export const Testimonials = ({
  className,
  dashedLineClassName,
}: {
  className?: string;
  dashedLineClassName?: string;
}) => {
  return (
    <>
      <section className={cn("overflow-hidden py-28 lg:py-32", className)}>
        <div className="container">
          <div className="space-y-4">
            <h2 className="text-2xl tracking-tight md:text-4xl lg:text-5xl">
              Trusted by IT leaders
            </h2>
            <p className="text-muted-foreground max-w-md leading-snug">
              {process.env.NEXT_PUBLIC_SITE_NAME} protects thousands of businesses worldwide — from lean startups
              to enterprise security teams managing global distributed
              workforces.
            </p>
            <Button variant="outline" className="shadow-md">
              Read case studies <ArrowRight className="size-4" />
            </Button>
          </div>

          <div className="relative mt-8 -mr-[max(3rem,calc((100vw-80rem)/2+3rem))] md:mt-12 lg:mt-20">
            <Carousel
              opts={{
                align: "start",
                loop: true,
              }}
              className="w-full"
            >
              <CarouselContent className="">
                {items.map((testimonial, index) => (
                  <CarouselItem
                    key={index}
                    className="xl:basis-1/3.5 grow basis-4/5 sm:basis-3/5 md:basis-2/5 lg:basis-[28%] 2xl:basis-[24%]"
                  >
                    <Card className="bg-muted h-full overflow-hidden border-none dark:border dark:border-border">
                      <CardContent className="flex h-full flex-col p-0">
                        <div className="relative h-[288px] lg:h-[328px]">
                          <Image
                            src={testimonial.image}
                            alt={testimonial.author}
                            fill
                            className="object-cover object-top"
                          />
                        </div>
                        <div className="flex flex-1 flex-col justify-between gap-10 p-6">
                          <blockquote className="font-display text-lg leading-none! font-medium md:text-xl lg:text-2xl">
                            {testimonial.quote}
                          </blockquote>
                          <div className="space-y-0.5">
                            <div className="text-primary font-semibold">
                              {testimonial.author}, {testimonial.role}
                            </div>
                            <div className="text-muted-foreground text-sm">
                              {testimonial.company}
                            </div>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </CarouselItem>
                ))}
              </CarouselContent>
              <div className="mt-8 flex gap-3">
                <CarouselPrevious className="bg-muted hover:bg-muted/80 dark:border dark:border-border static size-14.5 translate-x-0 translate-y-0 transition-colors [&>svg]:size-6 lg:[&>svg]:size-8" />
                <CarouselNext className="bg-muted hover:bg-muted/80 dark:border dark:border-border static size-14.5 translate-x-0 translate-y-0 transition-colors [&>svg]:size-6 lg:[&>svg]:size-8" />
              </div>
            </Carousel>
          </div>
        </div>
      </section>
      <DashedLine
        orientation="horizontal"
        className={cn("mx-auto max-w-[80%]", dashedLineClassName)}
      />
    </>
  );
};
