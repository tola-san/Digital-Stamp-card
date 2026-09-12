import Link from "next/link";

export default function Home() {
  return (
    <div className="flex min-h-screen bg-zinc-50 px-6 py-16 text-zinc-950">
      <main className="mx-auto flex w-full max-w-5xl flex-col justify-center">
        <span className="mb-5 w-fit rounded-full bg-sky-100 px-3 py-1 text-sm font-medium text-sky-800">
          ប្រព័ន្ធមូលដ្ឋានរួចរាល់
        </span>
        <h1 className="max-w-3xl text-4xl font-bold tracking-tight sm:text-6xl">
          វិធីសាមញ្ញជាងមុន ដើម្បីផ្តល់រង្វាន់រាល់ពេលអតិថិជនត្រឡប់មកវិញ។
        </h1>
        <p className="mt-6 max-w-2xl text-lg leading-8 text-zinc-600">
          កាតត្រាឌីជីថលផ្តល់ជូនអតិថិជននូវកាតសមាជិកលើទូរសព្ទ និងផ្តល់ឱ្យបុគ្គលិកនូវវិធីសុវត្ថិភាពក្នុងការផ្តល់ត្រា និងប្តូរយករង្វាន់។
        </p>

        <div className="mt-8 flex flex-wrap gap-3">
          <Link
            className="rounded-xl bg-sky-700 px-5 py-3 text-sm font-semibold text-white shadow-sm transition-[background-color,transform] duration-150 hover:bg-sky-800 active:scale-[0.96]"
            href="/register"
          >
            បង្កើតកាតអតិថិជន
          </Link>
          <Link
            className="rounded-xl border border-zinc-200 bg-white px-5 py-3 text-sm font-semibold text-zinc-800 shadow-sm transition-[background-color,border-color,transform] duration-150 hover:border-zinc-300 hover:bg-zinc-100 active:scale-[0.96]"
            href="/login"
          >
            ចូលគណនីអតិថិជន
          </Link>
          <Link
            className="rounded-xl border border-sky-200 bg-sky-50 px-5 py-3 text-sm font-semibold text-sky-800 shadow-sm transition-[background-color,border-color,transform] duration-150 hover:border-sky-300 hover:bg-sky-100 active:scale-[0.96]"
            href="/staff/login"
          >
            ចូលគណនីបុគ្គលិក
          </Link>
        </div>

        <div className="mt-12 grid gap-4 sm:grid-cols-3">
          {[
            ["ផ្នែកអតិថិជន", "ចុះឈ្មោះ ប្រមូលត្រា និងតាមដានដំណើរឆ្ពោះទៅរករង្វាន់។"],
            ["ផ្ទាំងគ្រប់គ្រងបុគ្គលិក", "គ្រប់គ្រងអតិថិជន និងត្រួតពិនិត្យរាល់ត្រាដែលបានផ្តល់។"],
            ["ប្រព័ន្ធ QR សុវត្ថិភាព", "ប្រើកូដដែលមានសុពលភាពខ្លី និងប្រើបានតែម្តងសម្រាប់ការទិញនីមួយៗ។"],
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
