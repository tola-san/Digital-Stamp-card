import type { Metadata } from "next";
import Link from "next/link";
import { CheckCircle2, Gift, ShieldCheck, Sparkles } from "lucide-react";

import { CustomerRegistrationForm } from "@/components/customer/customer-registration-form";

export const metadata: Metadata = {
  title: "Create your stamp card | Digital Stamp Card",
  description: "Register for a digital loyalty card using your phone number.",
};

const benefits = [
  "Keep every stamp in one place",
  "See your progress toward rewards",
  "Return anytime with your phone number",
];

export default function RegisterPage() {
  return (
    <main className="relative min-h-screen overflow-hidden bg-[#f3f8fc] text-zinc-950">
      <div
        aria-hidden="true"
        className="absolute -top-44 -left-44 size-[28rem] rounded-full bg-sky-200/45 blur-3xl"
      />
      <div
        aria-hidden="true"
        className="absolute -right-32 -bottom-48 size-[32rem] rounded-full bg-amber-100/70 blur-3xl"
      />

      <div className="relative mx-auto flex min-h-screen w-full max-w-6xl flex-col px-5 py-6 sm:px-8 lg:px-10">
        <header className="flex items-center justify-between">
          <Link
            className="group flex items-center gap-2.5 rounded-xl outline-none focus-visible:ring-4 focus-visible:ring-sky-600/15"
            href="/"
          >
            <span className="flex size-10 items-center justify-center rounded-xl bg-sky-700 text-white shadow-sm transition-transform duration-150 group-hover:-rotate-3">
              <Gift aria-hidden="true" className="size-5" strokeWidth={2} />
            </span>
            <span className="text-sm font-semibold tracking-tight sm:text-base">
              Digital Stamp
            </span>
          </Link>
          <div className="hidden items-center gap-2 text-sm font-medium text-zinc-600 sm:flex">
            <ShieldCheck aria-hidden="true" className="size-4 text-sky-700" />
            Secure registration
          </div>
        </header>

        <div className="grid flex-1 items-center gap-12 py-12 lg:grid-cols-[1fr_29rem] lg:gap-20 lg:py-16">
          <section className="mx-auto max-w-xl lg:mx-0">
            <div className="inline-flex items-center gap-2 rounded-full bg-white/80 px-3 py-1.5 text-sm font-medium text-sky-800 shadow-sm ring-1 ring-sky-900/8 backdrop-blur-sm">
              <Sparkles aria-hidden="true" className="size-4" />
              Rewards made simple
            </div>
            <h1 className="mt-6 text-4xl font-semibold tracking-[-0.035em] text-balance sm:text-5xl lg:text-6xl lg:leading-[1.05]">
              Your next reward starts here.
            </h1>
            <p className="mt-5 max-w-lg text-lg leading-8 text-zinc-600">
              Create your free digital stamp card in seconds. No paper card to
              lose, and your progress is always ready when you return.
            </p>

            <ul className="mt-8 space-y-3.5">
              {benefits.map((benefit) => (
                <li className="flex items-center gap-3 text-sm font-medium text-zinc-700" key={benefit}>
                  <CheckCircle2
                    aria-hidden="true"
                    className="size-5 shrink-0 text-sky-700"
                    strokeWidth={1.8}
                  />
                  {benefit}
                </li>
              ))}
            </ul>
          </section>

          <section aria-label="Customer registration">
            <CustomerRegistrationForm />
          </section>
        </div>
      </div>
    </main>
  );
}
