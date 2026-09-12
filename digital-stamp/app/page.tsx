import Link from "next/link";
import {
  ArrowRight,
  Check,
  Coffee,
  Gift,
  QrCode,
  ShieldCheck,
  Smartphone,
  UserRoundPlus,
} from "lucide-react";

const trustItems = [
  {
    icon: Smartphone,
    title: "មិនចាំបាច់ទាញយកកម្មវិធី",
    description: "បើកកាតត្រាបានភ្លាមៗតាមកម្មវិធីរុករក។",
  },
  {
    icon: ShieldCheck,
    title: "រក្សាទុកដោយសុវត្ថិភាព",
    description: "ត្រារបស់អ្នកមិនបាត់បង់ដូចកាតក្រដាសទេ។",
  },
  {
    icon: Gift,
    title: "ដឹងពេលរង្វាន់រួចរាល់",
    description: "មើលចំនួនត្រា និងរង្វាន់បានគ្រប់ពេល។",
  },
];

const steps = [
  {
    icon: UserRoundPlus,
    number: "០១",
    title: "ចុះឈ្មោះ",
    description: "បញ្ចូលឈ្មោះ និងលេខទូរសព្ទ ដើម្បីបង្កើតកាតឥតគិតថ្លៃ។",
  },
  {
    icon: QrCode,
    number: "០២",
    title: "ប្រមូលត្រា",
    description: "ស្កេន QR របស់ហាងរាល់ពេលទិញ ដើម្បីទទួលបានត្រាថ្មី។",
  },
  {
    icon: Gift,
    number: "០៣",
    title: "ទទួលរង្វាន់",
    description: "ពេលត្រាគ្រប់ចំនួន សូមប្តូរយករង្វាន់នៅហាង។",
  },
];

export default function Home() {
  return (
    <main className="min-h-screen overflow-hidden bg-[#f3f8fc] text-zinc-950">
      <header className="relative z-20 border-b border-sky-950/8 bg-white/80 backdrop-blur-xl">
        <div className="mx-auto flex h-18 w-full max-w-7xl items-center justify-between px-5 sm:px-8 lg:px-10">
          <Link
            className="group flex items-center gap-2.5 rounded-xl outline-none focus-visible:ring-4 focus-visible:ring-sky-600/15"
            href="/"
          >
            <span className="flex size-10 items-center justify-center rounded-xl bg-sky-700 text-white shadow-sm transition-transform duration-150 group-hover:-rotate-3">
              <Coffee aria-hidden="true" className="size-5" strokeWidth={2} />
            </span>
            <span className="font-semibold text-zinc-950">កាតត្រាឌីជីថល</span>
          </Link>

          <Link
            className="rounded-xl px-3.5 py-2 text-sm font-semibold text-zinc-600 outline-none transition-[background-color,color,transform] duration-150 hover:bg-sky-50 hover:text-sky-800 focus-visible:ring-4 focus-visible:ring-sky-600/15 active:scale-[0.96]"
            href="/staff/login"
          >
            ចូលបុគ្គលិក
          </Link>
        </div>
      </header>

      <section className="relative">
        <div
          aria-hidden="true"
          className="absolute -top-52 -right-44 size-[34rem] rounded-full bg-sky-200/55 blur-3xl"
        />
        <div
          aria-hidden="true"
          className="absolute -bottom-48 -left-48 size-[30rem] rounded-full bg-amber-100/70 blur-3xl"
        />

        <div className="relative mx-auto grid w-full max-w-7xl items-center gap-14 px-5 py-16 sm:px-8 sm:py-20 lg:grid-cols-[1.02fr_0.98fr] lg:gap-20 lg:px-10 lg:py-24">
          <div className="max-w-2xl">
            <div className="inline-flex items-center gap-2 rounded-full bg-white/85 px-3.5 py-2 text-sm font-semibold text-sky-800 shadow-sm ring-1 ring-sky-950/8 backdrop-blur-sm">
              <Coffee aria-hidden="true" className="size-4" strokeWidth={2} />
              កាតត្រាឌីជីថលសម្រាប់ហាងកាហ្វេ
            </div>

            <h1 className="mt-7 text-4xl leading-[1.35] font-semibold text-balance sm:text-5xl sm:leading-[1.3] lg:text-[3.65rem]">
              ទិញកាហ្វេ ប្រមូលត្រា និងទទួល
              <span className="text-sky-700"> រង្វាន់ពិសេស</span>
            </h1>
            <p className="mt-6 max-w-xl text-base leading-8 text-zinc-600 sm:text-lg sm:leading-9">
              រក្សាត្រារបស់អ្នកនៅលើទូរសព្ទ ដោយមិនចាំបាច់យកកាតក្រដាស។
              ចុះឈ្មោះបានឆាប់រហ័ស ហើយត្រឡប់មកប្រើវិញដោយលេខទូរសព្ទ។
            </p>

            <div className="mt-8 flex flex-col gap-3 sm:flex-row">
              <Link
                className="inline-flex min-h-12 items-center justify-center gap-2 rounded-xl bg-sky-700 px-5 py-3 text-base font-semibold text-white shadow-[0_12px_28px_-14px_oklch(0.5_0.15_240/0.85)] outline-none transition-[background-color,box-shadow,transform] duration-150 hover:bg-sky-800 hover:shadow-[0_16px_34px_-16px_oklch(0.42_0.15_240/0.9)] focus-visible:ring-4 focus-visible:ring-sky-600/20 active:scale-[0.96]"
                href="/register"
              >
                បង្កើតកាតឥតគិតថ្លៃ
                <ArrowRight aria-hidden="true" className="size-5" strokeWidth={2} />
              </Link>
              <Link
                className="inline-flex min-h-12 items-center justify-center rounded-xl border border-zinc-200 bg-white px-5 py-3 text-base font-semibold text-zinc-800 shadow-sm outline-none transition-[background-color,border-color,transform] duration-150 hover:border-sky-200 hover:bg-sky-50 focus-visible:ring-4 focus-visible:ring-sky-600/15 active:scale-[0.96]"
                href="/login"
              >
                ខ្ញុំមានកាតរួចហើយ
              </Link>
            </div>

            <div className="mt-7 flex flex-wrap gap-x-6 gap-y-2 text-sm text-zinc-600">
              <span className="flex items-center gap-2">
                <Check aria-hidden="true" className="size-4 text-sky-700" strokeWidth={2.5} />
                ប្រើតែលេខទូរសព្ទ
              </span>
              <span className="flex items-center gap-2">
                <Check aria-hidden="true" className="size-4 text-sky-700" strokeWidth={2.5} />
                មិនគិតថ្លៃ
              </span>
            </div>
          </div>

          <div className="relative mx-auto w-full max-w-xl lg:mx-0">
            <div
              aria-hidden="true"
              className="absolute inset-x-12 -bottom-7 h-24 rounded-full bg-sky-900/20 blur-2xl"
            />
            <div className="relative overflow-hidden rounded-[2rem] bg-sky-800 p-6 text-white shadow-[0_32px_90px_-34px_oklch(0.3_0.12_240/0.85),0_4px_14px_oklch(0_0_0/0.1)] ring-1 ring-white/15 sm:p-8">
              <div
                aria-hidden="true"
                className="absolute -top-24 -right-20 size-64 rounded-full bg-sky-400/25 blur-3xl"
              />
              <div
                aria-hidden="true"
                className="absolute -bottom-32 -left-24 size-64 rounded-full bg-sky-950/35 blur-3xl"
              />

              <div className="relative flex items-start justify-between gap-6">
                <div>
                  <div className="flex items-center gap-2 text-sky-100">
                    <Coffee aria-hidden="true" className="size-5" strokeWidth={1.8} />
                    <span className="text-sm font-semibold">រង្វាន់កាហ្វេ</span>
                  </div>
                  <p className="mt-3 text-xl font-semibold sm:text-2xl">កាតរបស់អ្នក</p>
                  <p className="mt-1 text-sm text-sky-100/70">សមាជិកហាងកាហ្វេ</p>
                </div>
                <div className="rounded-2xl bg-white/10 px-4 py-3 text-right ring-1 ring-white/15 backdrop-blur-sm">
                  <p className="text-2xl font-semibold tabular-nums sm:text-3xl">6/10</p>
                  <p className="mt-0.5 text-xs text-sky-100/70">ត្រាបានប្រមូល</p>
                </div>
              </div>

              <div className="relative mt-8 grid grid-cols-5 gap-3 sm:gap-4">
                {Array.from({ length: 10 }, (_, index) => {
                  const isCollected = index < 6;

                  return (
                    <div
                      aria-label={`ត្រាទី ${index + 1}: ${isCollected ? "បានប្រមូល" : "មិនទាន់មាន"}`}
                      className={`flex aspect-square items-center justify-center rounded-full ring-1 ${
                        isCollected
                          ? "bg-amber-300 text-sky-950 shadow-[inset_0_1px_0_oklch(1_0_0/0.55),0_8px_18px_-12px_oklch(0.8_0.15_80/0.9)] ring-amber-100/70"
                          : "bg-white/8 text-white/35 ring-white/18"
                      }`}
                      key={index}
                    >
                      <Coffee
                        aria-hidden="true"
                        className="size-6 sm:size-7"
                        strokeWidth={isCollected ? 2 : 1.5}
                      />
                    </div>
                  );
                })}
              </div>

              <div className="relative mt-8">
                <div className="flex items-center justify-between text-xs font-medium text-sky-100/80">
                  <span>បានបំពេញ 60%</span>
                  <span>នៅសល់ 4 ត្រា</span>
                </div>
                <div className="mt-2 h-2 overflow-hidden rounded-full bg-black/20 ring-1 ring-white/10">
                  <div className="h-full w-3/5 rounded-full bg-amber-300" />
                </div>
              </div>

              <div className="relative mt-6 flex items-center gap-3 rounded-2xl bg-white/10 p-4 ring-1 ring-white/12">
                <span className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-amber-300 text-sky-950">
                  <Gift aria-hidden="true" className="size-5" strokeWidth={2} />
                </span>
                <div>
                  <p className="font-semibold">ជិតទទួលបានរង្វាន់ហើយ!</p>
                  <p className="mt-0.5 text-sm text-sky-100/70">ប្រមូលតែ 4 ត្រាទៀតប៉ុណ្ណោះ។</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="border-y border-sky-950/8 bg-white/75" aria-label="អត្ថប្រយោជន៍">
        <div className="mx-auto grid w-full max-w-7xl divide-y divide-sky-950/8 px-5 sm:px-8 md:grid-cols-3 md:divide-x md:divide-y-0 lg:px-10">
          {trustItems.map(({ icon: Icon, title, description }) => (
            <article className="flex gap-4 py-7 md:px-7 md:first:pl-0 md:last:pr-0" key={title}>
              <span className="flex size-11 shrink-0 items-center justify-center rounded-2xl bg-sky-100 text-sky-800 ring-1 ring-inset ring-sky-900/8">
                <Icon aria-hidden="true" className="size-5" strokeWidth={1.8} />
              </span>
              <div>
                <h2 className="font-semibold text-zinc-900">{title}</h2>
                <p className="mt-1 text-sm leading-6 text-zinc-600">{description}</p>
              </div>
            </article>
          ))}
        </div>
      </section>

      <section className="bg-white px-5 py-18 sm:px-8 sm:py-24 lg:px-10" aria-labelledby="how-it-works">
        <div className="mx-auto w-full max-w-7xl">
          <div className="max-w-2xl">
            <p className="text-sm font-semibold text-sky-700">ងាយស្រួល និងឆាប់រហ័ស</p>
            <h2 id="how-it-works" className="mt-3 text-3xl leading-[1.4] font-semibold text-zinc-950 sm:text-4xl">
              ចាប់ផ្តើមប្រើក្នុងបីជំហាន
            </h2>
            <p className="mt-4 text-base leading-8 text-zinc-600">
              មិនត្រូវការកម្មវិធីបន្ថែម ឬកាតក្រដាសទៀតទេ។
            </p>
          </div>

          <div className="mt-10 grid gap-5 md:grid-cols-3">
            {steps.map(({ icon: Icon, number, title, description }) => (
              <article
                className="relative overflow-hidden rounded-3xl bg-[#f3f8fc] p-7 ring-1 ring-sky-950/8 transition-[box-shadow,transform] duration-150 hover:-translate-y-1 hover:shadow-[0_20px_50px_-28px_oklch(0.35_0.1_240/0.35)]"
                key={number}
              >
                <span className="absolute top-5 right-6 text-3xl font-semibold text-sky-900/10">{number}</span>
                <span className="flex size-12 items-center justify-center rounded-2xl bg-sky-700 text-white shadow-sm">
                  <Icon aria-hidden="true" className="size-6" strokeWidth={1.8} />
                </span>
                <h3 className="mt-6 text-xl font-semibold text-zinc-950">{title}</h3>
                <p className="mt-2 text-sm leading-7 text-zinc-600">{description}</p>
              </article>
            ))}
          </div>

          <div className="mt-12 flex flex-col items-start justify-between gap-5 rounded-3xl bg-sky-800 p-6 text-white shadow-[0_24px_70px_-30px_oklch(0.3_0.12_240/0.65)] sm:flex-row sm:items-center sm:p-8">
            <div>
              <p className="text-xl font-semibold">រួចរាល់ដើម្បីប្រមូលត្រាដំបូងរបស់អ្នកហើយឬនៅ?</p>
              <p className="mt-2 text-sm leading-6 text-sky-100/75">បង្កើតកាតដោយឥតគិតថ្លៃក្នុងពេលតែប៉ុន្មានវិនាទី។</p>
            </div>
            <Link
              className="inline-flex min-h-12 shrink-0 items-center justify-center gap-2 rounded-xl bg-white px-5 py-3 font-semibold text-sky-800 outline-none transition-[background-color,transform] duration-150 hover:bg-sky-50 focus-visible:ring-4 focus-visible:ring-white/30 active:scale-[0.96]"
              href="/register"
            >
              បង្កើតកាតឥឡូវនេះ
              <ArrowRight aria-hidden="true" className="size-5" strokeWidth={2} />
            </Link>
          </div>
        </div>
      </section>
    </main>
  );
}
