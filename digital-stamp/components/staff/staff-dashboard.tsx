"use client";

import { FormEvent, useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  ArrowRight, BadgeCheck, Bell, Check, CircleAlert, Clock3, Gift,
  LayoutDashboard, LoaderCircle, LogOut, QrCode, ScanLine, Sparkles,
  UserRound, UsersRound,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { ApiError } from "@/lib/api";
import { confirmStamp, getCurrentStaff, logoutStaff, previewStampToken } from "@/services/staff-service";
import type { Staff } from "@/types/domain";
import type { StampConfirmation, StampScanPreview } from "@/types/staff";

function initials(name: string) {
  return name.trim().split(/\s+/).slice(0, 2).map((part) => part[0]).join("").toUpperCase();
}

function progress(card: StampScanPreview["card"]) {
  return Math.min(100, Math.round((card.stamp_count / card.required_stamps) * 100));
}

export function StaffDashboard() {
  const router = useRouter();
  const [staff, setStaff] = useState<Staff | null>(null);
  const [token, setToken] = useState("");
  const [preview, setPreview] = useState<StampScanPreview | null>(null);
  const [confirmation, setConfirmation] = useState<StampConfirmation | null>(null);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [isScanning, setIsScanning] = useState(false);
  const [isConfirming, setIsConfirming] = useState(false);
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  useEffect(() => {
    let active = true;
    getCurrentStaff()
      .then(({ staff: currentStaff }) => active && setStaff(currentStaff))
      .catch((requestError: unknown) => {
        if (!active) return;
        if (requestError instanceof ApiError && requestError.status === 401) {
          router.replace("/staff/login");
          return;
        }
        setError("មិនអាចភ្ជាប់ទៅសេវាកម្មបានទេ។");
      })
      .finally(() => active && setIsLoading(false));
    return () => { active = false; };
  }, [router]);

  async function handlePreview(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!token.trim()) {
      setError("សូមស្កេន ឬបញ្ចូលលេខកូដ QR របស់អតិថិជន។");
      return;
    }
    setError("");
    setConfirmation(null);
    setIsScanning(true);
    try {
      setPreview(await previewStampToken(token.trim()));
    } catch (requestError) {
      setError(requestError instanceof ApiError && requestError.code === "QR_TOKEN_EXPIRED"
        ? "កូដ QR នេះបានផុតកំណត់។ សូមឱ្យអតិថិជនបង្កើតកូដថ្មី។"
        : "មិនអាចអានកូដ QR នេះបានទេ។ សូមព្យាយាមម្តងទៀត។");
      setPreview(null);
    } finally {
      setIsScanning(false);
    }
  }

  async function handleConfirm() {
    if (!preview) return;
    setError("");
    setIsConfirming(true);
    try {
      setConfirmation(await confirmStamp(preview.scan_id));
      setPreview(null);
      setToken("");
    } catch (requestError) {
      setError(requestError instanceof ApiError && requestError.code === "QR_TOKEN_ALREADY_USED"
        ? "កូដនេះត្រូវបានប្រើរួចហើយ។"
        : "មិនអាចបន្ថែមត្រាបានទេ។ សូមព្យាយាមម្តងទៀត។");
    } finally {
      setIsConfirming(false);
    }
  }

  async function handleLogout() {
    setIsLoggingOut(true);
    try {
      await logoutStaff();
      router.replace("/staff/login");
    } catch {
      setError("មិនអាចចាកចេញពីគណនីបានទេ។");
      setIsLoggingOut(false);
    }
  }

  if (isLoading) {
    return <div className="flex min-h-screen items-center justify-center bg-[#f4f7f9]" role="status"><LoaderCircle aria-hidden="true" className="size-7 animate-spin text-emerald-700" /><span className="sr-only">កំពុងផ្ទុក…</span></div>;
  }

  if (!staff) {
    return <main className="flex min-h-screen items-center justify-center bg-[#f4f7f9] px-5"><div className="max-w-md rounded-3xl bg-white p-8 text-center shadow-xl ring-1 ring-black/5"><CircleAlert className="mx-auto size-10 text-red-600" /><h1 className="mt-4 text-xl font-semibold">មិនអាចបង្ហាញផ្ទាំងគ្រប់គ្រងបានទេ</h1><Link className="mt-6 inline-flex h-11 items-center rounded-xl bg-emerald-700 px-5 font-semibold text-white" href="/staff/login">ចូលគណនីម្តងទៀត</Link></div></main>;
  }

  const shownCard = confirmation?.card ?? preview?.card;
  const shownCustomer = confirmation?.customer ?? preview?.customer;

  return (
    <main className="min-h-screen bg-[#f4f7f9] text-slate-950">
      <div className="mx-auto grid min-h-screen max-w-[1500px] lg:grid-cols-[248px_1fr]">
        <aside className="hidden border-r border-slate-200/80 bg-white px-5 py-6 lg:flex lg:flex-col">
          <Link className="flex items-center gap-3 rounded-xl px-2 outline-none focus-visible:ring-4 focus-visible:ring-emerald-600/15" href="/">
            <span className="flex size-10 items-center justify-center rounded-xl bg-emerald-700 text-white shadow-sm"><Gift className="size-5" /></span>
            <span><strong className="block text-sm font-semibold">Digital Stamp</strong><small className="text-xs text-slate-500">Staff workspace</small></span>
          </Link>
          <nav aria-label="Staff navigation" className="mt-9 space-y-1.5">
            <a className="flex h-11 items-center gap-3 rounded-xl bg-emerald-50 px-3.5 text-sm font-semibold text-emerald-800 ring-1 ring-inset ring-emerald-700/8" href="#scan"><LayoutDashboard className="size-[18px]" />ទិដ្ឋភាពទូទៅ</a>
            <a className="flex h-11 items-center gap-3 rounded-xl px-3.5 text-sm font-medium text-slate-600 transition-[background-color,color] duration-150 hover:bg-slate-100 hover:text-slate-950" href="#scan"><ScanLine className="size-[18px]" />ស្កេន QR</a>
            <span className="flex h-11 cursor-not-allowed items-center gap-3 rounded-xl px-3.5 text-sm font-medium text-slate-400"><UsersRound className="size-[18px]" />អតិថិជន</span>
          </nav>
          <div className="mt-auto rounded-2xl bg-slate-50 p-3.5 ring-1 ring-slate-900/5">
            <div className="flex items-center gap-3"><span className="flex size-9 items-center justify-center rounded-xl bg-emerald-100 text-xs font-bold text-emerald-800">{initials(staff.name)}</span><span className="min-w-0"><strong className="block truncate text-sm font-semibold">{staff.name}</strong><small className="block truncate text-slate-500">{staff.email}</small></span></div>
            <button className="mt-3 flex h-9 w-full items-center justify-center gap-2 rounded-lg text-sm font-medium text-slate-600 transition-[background-color,color,transform] duration-150 hover:bg-white hover:text-slate-950 active:scale-[0.96]" disabled={isLoggingOut} onClick={() => void handleLogout()}><LogOut className="size-4" />ចាកចេញ</button>
          </div>
        </aside>

        <div className="min-w-0">
          <header className="sticky top-0 z-10 border-b border-slate-200/75 bg-[#f4f7f9]/88 backdrop-blur-xl">
            <div className="flex h-18 items-center justify-between px-5 sm:px-8 lg:px-10">
              <div><p className="text-xs font-semibold tracking-[0.14em] text-emerald-700 uppercase">ផ្ទាំងបុគ្គលិក</p><h1 className="mt-0.5 text-lg font-semibold tracking-tight">សួស្តី, {staff.name}</h1></div>
              <div className="flex items-center gap-2"><button aria-label="Notifications" className="flex size-10 items-center justify-center rounded-xl bg-white text-slate-600 shadow-sm ring-1 ring-slate-900/7 transition-[background-color,color,transform] duration-150 hover:bg-slate-50 hover:text-slate-950 active:scale-[0.96]"><Bell className="size-[18px]" /></button><button aria-label="Log out" className="flex size-10 items-center justify-center rounded-xl bg-white text-slate-600 shadow-sm ring-1 ring-slate-900/7 transition-[background-color,color,transform] duration-150 hover:bg-slate-50 hover:text-slate-950 active:scale-[0.96] lg:hidden" onClick={() => void handleLogout()}><LogOut className="size-[18px]" /></button></div>
            </div>
          </header>

          <div className="px-5 py-7 sm:px-8 lg:px-10 lg:py-9">
            <section className="grid gap-4 sm:grid-cols-3">
              {[{ label: "ត្រាថ្ងៃនេះ", icon: Check, tone: "bg-emerald-50 text-emerald-700" }, { label: "អតិថិជនថ្ងៃនេះ", icon: UsersRound, tone: "bg-sky-50 text-sky-700" }, { label: "រង្វាន់បានដូរ", icon: Gift, tone: "bg-amber-50 text-amber-700" }].map((item) => <div className="rounded-2xl bg-white p-5 shadow-[0_8px_30px_-20px_oklch(0.25_0.03_240/0.25)] ring-1 ring-slate-900/6" key={item.label}><div className="flex items-center justify-between"><span className={`flex size-10 items-center justify-center rounded-xl ${item.tone}`}><item.icon className="size-[18px]" /></span><span className="text-2xl font-semibold tabular-nums">—</span></div><p className="mt-4 text-sm text-slate-500">{item.label}</p></div>)}
            </section>

            <section className="mt-6 grid items-start gap-6 xl:grid-cols-[minmax(0,1.1fr)_minmax(340px,.9fr)]" id="scan">
              <div className="overflow-hidden rounded-3xl bg-white shadow-[0_24px_70px_-36px_oklch(0.25_0.05_160/0.3)] ring-1 ring-slate-900/6">
                <div className="border-b border-slate-100 px-6 py-5 sm:px-7"><div className="flex items-center gap-3"><span className="flex size-10 items-center justify-center rounded-xl bg-emerald-100 text-emerald-800"><QrCode className="size-5" /></span><div><h2 className="font-semibold">ស្កេនកូដអតិថិជន</h2><p className="mt-0.5 text-sm text-slate-500">បន្ថែមត្រាថ្មីដោយសុវត្ថិភាព</p></div></div></div>
                <div className="p-6 sm:p-7">
                  <div className="flex min-h-52 flex-col items-center justify-center rounded-2xl border border-dashed border-emerald-700/25 bg-emerald-50/45 px-6 text-center">
                    <span className="relative flex size-16 items-center justify-center rounded-2xl bg-white text-emerald-700 shadow-sm ring-1 ring-emerald-900/8"><ScanLine className="size-8" /><span className="absolute -right-1 -bottom-1 size-3 rounded-full bg-emerald-500 ring-4 ring-white" /></span>
                    <h3 className="mt-4 font-semibold">រួចរាល់សម្រាប់ស្កេន</h3><p className="mt-1 max-w-sm text-sm leading-6 text-slate-500">បញ្ចូលតម្លៃដែលបានអានពី QR របស់អតិថិជន ដើម្បីមើលកាតមុនបញ្ជាក់។</p>
                  </div>
                  <Link className="mt-5 flex h-12 w-full items-center justify-center gap-2 rounded-xl bg-emerald-700 px-5 font-semibold text-white shadow-[0_10px_24px_-12px_oklch(0.45_0.13_165/0.7)] outline-none transition-[background-color,box-shadow,transform] duration-150 hover:bg-emerald-800 focus-visible:ring-4 focus-visible:ring-emerald-600/20 active:scale-[0.96]" href="/staff/scan"><ScanLine className="size-5" />បើកកាមេរ៉ាស្កេន</Link>
                  <div className="my-5 flex items-center gap-3 text-xs text-slate-400"><span className="h-px flex-1 bg-slate-200" /><span>ឬបញ្ចូល token</span><span className="h-px flex-1 bg-slate-200" /></div>
                  <form onSubmit={handlePreview}><label className="sr-only" htmlFor="qr-token">លេខកូដ QR</label><div className="flex flex-col gap-3 sm:flex-row"><input autoComplete="off" className="h-11 min-w-0 flex-1 rounded-xl border border-slate-200 bg-white px-4 text-sm outline-none transition-[border-color,box-shadow] duration-150 placeholder:text-slate-400 focus:border-emerald-600 focus:ring-4 focus:ring-emerald-600/12" id="qr-token" onChange={(event) => { setToken(event.target.value); setError(""); }} placeholder="បញ្ចូល token ពី QR…" value={token} /><Button className="h-11 rounded-xl px-4 transition-[background-color,transform] duration-150 active:scale-[0.96]" disabled={isScanning} type="submit" variant="outline">{isScanning ? <LoaderCircle className="animate-spin" /> : <ArrowRight />}ពិនិត្យ</Button></div></form>
                  {error && <div aria-live="polite" className="mt-4 flex gap-2.5 rounded-xl bg-red-50 px-4 py-3 text-sm text-red-700 ring-1 ring-inset ring-red-700/10" role="alert"><CircleAlert className="mt-0.5 size-4 shrink-0" />{error}</div>}
                </div>
              </div>

              <div className="rounded-3xl bg-white p-6 shadow-[0_24px_70px_-36px_oklch(0.25_0.05_240/0.3)] ring-1 ring-slate-900/6 sm:p-7">
                {shownCustomer && shownCard ? <>
                  <div className="flex items-start justify-between gap-4"><div className="flex min-w-0 items-center gap-3.5"><span className="flex size-12 shrink-0 items-center justify-center rounded-2xl bg-sky-100 font-bold text-sky-800">{initials(shownCustomer.name)}</span><div className="min-w-0"><p className="truncate font-semibold">{shownCustomer.name}</p><p className="mt-0.5 text-sm text-slate-500">{shownCustomer.phone}</p></div></div>{confirmation ? <span className="inline-flex items-center gap-1.5 rounded-full bg-emerald-50 px-3 py-1.5 text-xs font-semibold text-emerald-700"><BadgeCheck className="size-3.5" />បានជោគជ័យ</span> : <span className="inline-flex items-center gap-1.5 rounded-full bg-amber-50 px-3 py-1.5 text-xs font-semibold text-amber-700"><Clock3 className="size-3.5" />រង់ចាំបញ្ជាក់</span>}</div>
                  <div className="mt-7 rounded-2xl bg-slate-50 p-5 ring-1 ring-slate-900/5"><div className="flex items-end justify-between"><div><p className="text-sm text-slate-500">ត្រាបច្ចុប្បន្ន</p><p className="mt-1 text-3xl font-semibold tabular-nums">{shownCard.stamp_count}<span className="text-base font-medium text-slate-400"> / {shownCard.required_stamps}</span></p></div><span className="text-sm font-semibold text-emerald-700">{progress(shownCard)}%</span></div><div aria-label={`${progress(shownCard)} percent complete`} aria-valuemax={100} aria-valuemin={0} aria-valuenow={progress(shownCard)} className="mt-4 h-2 overflow-hidden rounded-full bg-slate-200" role="progressbar"><div className="h-full rounded-full bg-emerald-600 transition-[width] duration-500" style={{ width: `${progress(shownCard)}%` }} /></div><p className="mt-3 flex items-center gap-2 text-sm text-slate-500"><Sparkles className="size-4 text-amber-500" />នៅសល់ {Math.max(0, shownCard.required_stamps - shownCard.stamp_count)} ត្រាទៀតដើម្បីទទួលរង្វាន់</p></div>
                  {confirmation ? <button className="mt-5 flex h-12 w-full items-center justify-center gap-2 rounded-xl bg-slate-950 font-semibold text-white transition-[background-color,transform] duration-150 hover:bg-slate-800 active:scale-[0.96]" onClick={() => setConfirmation(null)}>ស្កេនអតិថិជនបន្ទាប់<ArrowRight className="size-4" /></button> : <Button className="mt-5 h-12 w-full rounded-xl bg-emerald-700 font-semibold text-white transition-[background-color,transform] duration-150 hover:bg-emerald-800 active:scale-[0.96]" disabled={isConfirming} onClick={() => void handleConfirm()}>{isConfirming ? <LoaderCircle className="animate-spin" /> : <Check />}បញ្ជាក់ និងបន្ថែម 1 ត្រា</Button>}
                </> : <div className="flex min-h-[330px] flex-col items-center justify-center text-center"><span className="flex size-14 items-center justify-center rounded-2xl bg-slate-100 text-slate-400"><UserRound className="size-6" /></span><h2 className="mt-4 font-semibold">មិនទាន់មានអតិថិជន</h2><p className="mt-1 max-w-xs text-sm leading-6 text-slate-500">ព័ត៌មានអតិថិជន និងវឌ្ឍនភាពកាតនឹងបង្ហាញនៅទីនេះក្រោយពេលស្កេន។</p></div>}
              </div>
            </section>
          </div>
        </div>
      </div>
    </main>
  );
}
