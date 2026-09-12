"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { CircleAlert, Gift, LoaderCircle, LogOut, Users } from "lucide-react";

import { Button } from "@/components/ui/button";
import { ApiError } from "@/lib/api";
import { getCurrentStaff, logoutStaff } from "@/services/staff-service";
import type { Staff } from "@/types/domain";

export function StaffDashboard() {
  const router = useRouter();
  const [staff, setStaff] = useState<Staff | null>(null);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  useEffect(() => {
    let isActive = true;

    getCurrentStaff()
      .then(({ staff: currentStaff }) => {
        if (isActive) setStaff(currentStaff);
      })
      .catch((requestError: unknown) => {
        if (!isActive) return;
        if (requestError instanceof ApiError && requestError.status === 401) {
          router.replace("/staff/login");
          return;
        }
        setError(requestError instanceof ApiError ? "មិនអាចបង្ហាញគណនីបុគ្គលិកបានទេ។" : "មិនអាចភ្ជាប់ទៅសេវាកម្មបានទេ។");
      })
      .finally(() => {
        if (isActive) setIsLoading(false);
      });

    return () => { isActive = false; };
  }, [router]);

  async function handleLogout() {
    setIsLoggingOut(true);
    setError("");
    try {
      await logoutStaff();
      router.replace("/staff/login");
    } catch (requestError) {
      setError(requestError instanceof ApiError ? "មិនអាចចាកចេញពីគណនីបានទេ។" : "មិនអាចភ្ជាប់ទៅសេវាកម្មបានទេ។");
      setIsLoggingOut(false);
    }
  }

  if (isLoading) {
    return <div className="flex min-h-screen items-center justify-center bg-[#f3f8fc]" role="status"><div className="flex flex-col items-center gap-3 text-sky-800"><LoaderCircle aria-hidden="true" className="size-8 animate-spin" /><p className="text-sm font-medium">កំពុងពិនិត្យការចូលប្រើរបស់បុគ្គលិក…</p></div></div>;
  }

  if (!staff) {
    return <main className="flex min-h-screen items-center justify-center bg-[#f3f8fc] px-5"><div className="max-w-md rounded-3xl bg-white p-8 text-center shadow-lg ring-1 ring-zinc-950/8"><CircleAlert aria-hidden="true" className="mx-auto size-10 text-red-600" /><h1 className="mt-4 text-2xl font-semibold">មិនអាចបង្ហាញផ្ទាំងគ្រប់គ្រងបានទេ</h1><p className="mt-2 text-zinc-600">{error || "សូមចូលគណនីម្តងទៀត។"}</p><Link className="mt-6 inline-flex h-11 items-center rounded-xl bg-sky-700 px-5 font-semibold text-white transition-[background-color,transform] duration-150 hover:bg-sky-800 active:scale-[0.96]" href="/staff/login">ត្រឡប់ទៅការចូលគណនីបុគ្គលិក</Link></div></main>;
  }

  return (
    <main className="min-h-screen bg-[#f3f8fc] text-zinc-950">
      <header className="border-b border-sky-950/8 bg-white/85 backdrop-blur-xl">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-5 py-4 sm:px-8 lg:px-10">
          <Link className="flex items-center gap-2.5 rounded-xl outline-none focus-visible:ring-4 focus-visible:ring-sky-600/15" href="/">
            <span className="flex size-10 items-center justify-center rounded-xl bg-sky-700 text-white shadow-sm"><Gift aria-hidden="true" className="size-5" /></span>
            <span className="font-semibold tracking-tight">កាតត្រាឌីជីថល</span>
          </Link>
          <Button className="rounded-xl text-zinc-600 transition-[background-color,color,transform] duration-150 hover:bg-zinc-100 hover:text-zinc-950 active:scale-[0.96] active:translate-y-0" disabled={isLoggingOut} onClick={() => void handleLogout()} size="lg" variant="ghost">
            {isLoggingOut ? <LoaderCircle aria-hidden="true" className="animate-spin" /> : <LogOut aria-hidden="true" />}
            <span className="hidden sm:inline">ចាកចេញ</span>
          </Button>
        </div>
      </header>

      <div className="mx-auto max-w-6xl px-5 py-10 sm:px-8 lg:px-10">
        <p className="text-sm font-semibold tracking-wide text-sky-700 uppercase">ផ្ទាំងគ្រប់គ្រងបុគ្គលិក</p>
        <h1 className="mt-2 text-3xl font-semibold tracking-tight sm:text-4xl">សួស្តី {staff.name}</h1>
        <p className="mt-2 text-zinc-600">បានចូលគណនីជា {staff.email}</p>

        {error && <div aria-live="polite" className="mt-6 flex gap-2.5 rounded-xl bg-red-50 px-4 py-3 text-sm text-red-700 ring-1 ring-inset ring-red-600/10" role="alert"><CircleAlert aria-hidden="true" className="size-4 shrink-0" />{error}</div>}

        <section className="mt-9 rounded-3xl bg-white p-7 shadow-[0_20px_60px_-30px_oklch(0.32_0.08_235/0.28)] ring-1 ring-sky-950/8 sm:p-9">
          <span className="flex size-11 items-center justify-center rounded-2xl bg-sky-100 text-sky-800 ring-1 ring-inset ring-sky-900/8"><Users aria-hidden="true" className="size-5" /></span>
          <h2 className="mt-5 text-xl font-semibold">ការផ្ទៀងផ្ទាត់បុគ្គលិកត្រូវបានភ្ជាប់រួចរាល់</h2>
          <p className="mt-2 max-w-2xl leading-7 text-zinc-600">ផ្ទាំងគ្រប់គ្រងបានផ្ទៀងផ្ទាត់ការចូលប្រើ Gin ប្រកបដោយសុវត្ថិភាពរបស់អ្នកតាមរយៈ <code className="rounded bg-zinc-100 px-1.5 py-0.5 text-sm text-zinc-800">GET /api/staff/me</code>។ ការគ្រប់គ្រងអតិថិជន និងស្ថិតិនឹងបង្ហាញនៅទីនេះ បន្ទាប់ពី API ផ្នែកខាងក្រោយត្រូវបានបង្កើតរួច។</p>
        </section>
      </div>
    </main>
  );
}
