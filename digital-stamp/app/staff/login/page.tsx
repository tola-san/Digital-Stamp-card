import type { Metadata } from "next";
import Link from "next/link";
import { Gift, ShieldCheck } from "lucide-react";

import { StaffLoginForm } from "@/components/staff/staff-login-form";

export const metadata: Metadata = {
  title: "Staff login | Digital Stamp Card",
  description: "Secure access for Digital Stamp Card staff.",
};

export default function StaffLoginPage() {
  return (
    <main className="relative min-h-screen overflow-hidden bg-[#f3f8fc] text-zinc-950">
      <div aria-hidden="true" className="absolute -top-40 -right-32 size-[30rem] rounded-full bg-sky-200/50 blur-3xl" />
      <div aria-hidden="true" className="absolute -bottom-52 -left-40 size-[32rem] rounded-full bg-amber-100/75 blur-3xl" />

      <div className="relative mx-auto flex min-h-screen w-full max-w-6xl flex-col px-5 py-6 sm:px-8 lg:px-10">
        <header className="flex items-center justify-between">
          <Link className="group flex items-center gap-2.5 rounded-xl outline-none focus-visible:ring-4 focus-visible:ring-sky-600/15" href="/">
            <span className="flex size-10 items-center justify-center rounded-xl bg-sky-700 text-white shadow-sm transition-transform duration-150 group-hover:-rotate-3">
              <Gift aria-hidden="true" className="size-5" strokeWidth={2} />
            </span>
            <span className="text-sm font-semibold tracking-tight sm:text-base">Digital Stamp</span>
          </Link>
          <div className="flex items-center gap-2 text-sm font-medium text-zinc-600">
            <ShieldCheck aria-hidden="true" className="size-4 text-sky-700" />
            Staff portal
          </div>
        </header>

        <div className="grid flex-1 items-center gap-12 py-12 lg:grid-cols-[1fr_29rem] lg:gap-20 lg:py-16">
          <section className="mx-auto max-w-xl lg:mx-0">
            <p className="text-sm font-semibold tracking-wide text-sky-700 uppercase">Shop operations</p>
            <h1 className="mt-4 text-4xl font-semibold tracking-[-0.035em] text-balance sm:text-5xl lg:text-6xl lg:leading-[1.05]">
              Welcome back to your shop.
            </h1>
            <p className="mt-5 max-w-lg text-lg leading-8 text-zinc-600">
              Sign in to manage loyalty activity and serve returning customers.
            </p>
          </section>

          <section aria-label="Staff login"><StaffLoginForm /></section>
        </div>
      </div>
    </main>
  );
}
