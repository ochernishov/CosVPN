import Link from "next/link";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { cn } from "@/lib/utils";

const categories = [
  {
    title: "Security",
    questions: [
      {
        question: "What encryption does CosVPN use?",
        answer:
          "CosVPN uses AES-256 encryption, the same standard used by governments and financial institutions worldwide. All traffic is encrypted end-to-end using the WireGuard protocol.",
      },
      {
        question: "Does CosVPN keep logs of my activity?",
        answer:
          "No. CosVPN operates under a strict zero-log policy. We do not store, track, or share any browsing data, connection timestamps, or IP addresses. Our infrastructure is independently audited annually.",
      },
      {
        question: "Will CosVPN slow down my internet speed?",
        answer:
          "CosVPN is built on the WireGuard protocol, which is significantly faster than legacy VPN protocols like OpenVPN and IPSec. Most users see less than 5% speed reduction on typical connections.",
      },
    ],
  },
  {
    title: "Your account",
    questions: [
      {
        question: "How many devices can I connect simultaneously?",
        answer:
          "The Starter plan supports up to 5 simultaneous connections per user. Business and Enterprise plans offer unlimited simultaneous connections per user.",
      },
      {
        question: "Which platforms are supported?",
        answer:
          "CosVPN supports Windows, macOS, Linux, iOS, Android, and Chrome OS. We also offer browser extensions for Chrome and Firefox, plus router-level configuration for network-wide protection.",
      },
    ],
  },
  {
    title: "Enterprise",
    questions: [
      {
        question: "Can I integrate CosVPN with my existing SSO provider?",
        answer:
          "Yes. Enterprise plans include SSO/SAML integration with providers like Okta, Azure AD, Google Workspace, and OneLogin. Setup typically takes less than 30 minutes.",
      },
      {
        question: "What is the uptime SLA for Enterprise plans?",
        answer:
          "Enterprise plans include a 99.99% uptime SLA with dedicated server infrastructure and 24/7 priority support. We also provide a dedicated account manager for onboarding and ongoing support.",
      },
    ],
  },
];

export const FAQ = ({
  headerTag = "h2",
  className,
  className2,
}: {
  headerTag?: "h1" | "h2";
  className?: string;
  className2?: string;
}) => {
  return (
    <section className={cn("py-28 lg:py-32", className)}>
      <div className="container max-w-5xl">
        <div className={cn("mx-auto grid gap-16 lg:grid-cols-2", className2)}>
          <div className="space-y-4">
            {headerTag === "h1" ? (
              <h1 className="text-2xl tracking-tight md:text-4xl lg:text-5xl">
                Got Questions?
              </h1>
            ) : (
              <h2 className="text-2xl tracking-tight md:text-4xl lg:text-5xl">
                Got Questions?
              </h2>
            )}
            <p className="text-muted-foreground max-w-md leading-snug lg:mx-auto">
              If you can't find what you're looking for,{" "}
              <Link href="/contact" className="underline underline-offset-4">
                get in touch
              </Link>
              .
            </p>
          </div>

          <div className="grid gap-6 text-start">
            {categories.map((category, categoryIndex) => (
              <div key={category.title} className="">
                <h3 className="text-muted-foreground border-b py-4">
                  {category.title}
                </h3>
                <Accordion type="single" collapsible className="w-full">
                  {category.questions.map((item, i) => (
                    <AccordionItem key={i} value={`${categoryIndex}-${i}`}>
                      <AccordionTrigger>{item.question}</AccordionTrigger>
                      <AccordionContent className="text-muted-foreground">
                        {item.answer}
                      </AccordionContent>
                    </AccordionItem>
                  ))}
                </Accordion>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
};
