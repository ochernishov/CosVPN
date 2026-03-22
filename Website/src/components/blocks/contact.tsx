import React from "react";

import Link from "next/link";

import { Linkedin, Twitter } from "lucide-react";

import { ContactForm } from "@/components/blocks/contact-form";
import { DashedLine } from "@/components/dashed-line";

const careersEmail = process.env.NEXT_PUBLIC_CAREERS_EMAIL || "";
const pressEmail = process.env.NEXT_PUBLIC_PRESS_EMAIL || "";
const contactAddress = process.env.NEXT_PUBLIC_CONTACT_ADDRESS || "";
const twitterUrl = process.env.NEXT_PUBLIC_TWITTER_URL || "#";
const linkedinUrl = process.env.NEXT_PUBLIC_LINKEDIN_URL || "#";

const contactInfo = [
  {
    title: "Corporate office",
    content: (
      <p className="text-muted-foreground mt-3">
        {contactAddress}
      </p>
    ),
  },
  {
    title: "Email us",
    content: (
      <div className="mt-3">
        <div>
          <p className="">Careers</p>
          <Link
            href={`mailto:${careersEmail}`}
            className="text-muted-foreground hover:text-foreground"
          >
            {careersEmail}
          </Link>
        </div>
        <div className="mt-1">
          <p className="">Press</p>
          <Link
            href={`mailto:${pressEmail}`}
            className="text-muted-foreground hover:text-foreground"
          >
            {pressEmail}
          </Link>
        </div>
      </div>
    ),
  },
  {
    title: "Follow us",
    content: (
      <div className="mt-3 flex gap-6 lg:gap-10">
        <Link
          href={twitterUrl}
          className="text-muted-foreground hover:text-foreground"
        >
          <Twitter className="size-5" />
        </Link>
        <Link
          href={linkedinUrl}
          className="text-muted-foreground hover:text-foreground"
        >
          <Linkedin className="size-5" />
        </Link>
      </div>
    ),
  },
];

export default function Contact() {
  return (
    <section className="py-28 lg:py-32 lg:pt-44">
      <div className="container max-w-2xl">
        <h1 className="text-center text-2xl font-semibold tracking-tight md:text-4xl lg:text-5xl">
          Contact us
        </h1>
        <p className="text-muted-foreground mt-4 text-center leading-snug font-medium lg:mx-auto">
          Have questions about {process.env.NEXT_PUBLIC_SITE_NAME}? We&apos;d love to hear from you.
        </p>

        <div className="mt-10 flex justify-between gap-8 max-sm:flex-col md:mt-14 lg:mt-20 lg:gap-12">
          {contactInfo.map((info, index) => (
            <div key={index}>
              <h2 className="font-medium">{info.title}</h2>
              {info.content}
            </div>
          ))}
        </div>

        <DashedLine className="my-12" />

        {/* Inquiry Form */}
        <div className="mx-auto">
          <h2 className="mb-4 text-lg font-semibold">Inquiries</h2>
          <ContactForm />
        </div>
      </div>
    </section>
  );
}
