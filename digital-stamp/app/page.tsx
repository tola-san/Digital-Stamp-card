import Link from "next/link";

export default function Home() {
  return (
    <div className="flex min-h-screen bg-zinc-50 px-6 py-16 text-zinc-950">
      <main className="mx-auto flex w-full max-w-5xl flex-col justify-center">
        <span className="mb-5 w-fit rounded-full bg-sky-100 px-3 py-1 text-sm font-medium text-sky-800">
          Foundation ready
        </span>
        <h1 className="max-w-3xl text-4xl font-bold tracking-tight sm:text-6xl">
          A simpler way to reward every return visit.
        </h1>
        <p className="mt-6 max-w-2xl text-lg leading-8 text-zinc-600">
          Digital Stamp Card gives customers a mobile loyalty card and gives staff a secure way to issue stamps and redeem rewards.
        </p>

        <div className="mt-8 flex flex-wrap gap-3">
          <Link
            className="rounded-xl bg-sky-700 px-5 py-3 text-sm font-semibold text-white shadow-sm transition-[background-color,transform] duration-150 hover:bg-sky-800 active:scale-[0.96]"
            href="/register"
          >
            Create customer card
          </Link>
          <Link
            className="rounded-xl border border-zinc-200 bg-white px-5 py-3 text-sm font-semibold text-zinc-800 shadow-sm transition-[background-color,border-color,transform] duration-150 hover:border-zinc-300 hover:bg-zinc-100 active:scale-[0.96]"
            href="/login"
          >
            Customer login
          </Link>
        </div>

        <div className="mt-12 grid gap-4 sm:grid-cols-3">
          {[
            ["Customer portal", "Register, collect stamps, and follow reward progress."],
            ["Staff dashboard", "Manage customers and control every stamp issued."],
            ["Secure QR flow", "Use short-lived, single-use codes for each purchase."],
          ].map(([title, description]) => (
            <section key={title} className="rounded-2xl border border-zinc-200 bg-white p-6 shadow-sm">
              <h2 className="font-semibold">{title}</h2>
              <p className="mt-2 text-sm leading-6 text-zinc-600">{description}</p>
            </section>
          ))}
        </div>
      </main>
    </div>
  );
}
