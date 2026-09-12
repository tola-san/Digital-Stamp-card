"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import {
  ArrowRight,
  CircleAlert,
  Eye,
  EyeOff,
  LoaderCircle,
  LockKeyhole,
  Mail,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ApiError } from "@/lib/api";
import { loginStaff } from "@/services/staff-service";

export function StaffLoginForm() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [formError, setFormError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFormError("");

    if (!email.trim() || !password) {
      setFormError("សូមបញ្ចូលអ៊ីមែល និងពាក្យសម្ងាត់។");
      return;
    }

    setIsSubmitting(true);

    try {
      await loginStaff({ email: email.trim().toLowerCase(), password });
      router.replace("/staff/dashboard");
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) {
        setFormError("អ៊ីមែល ឬពាក្យសម្ងាត់មិនត្រឹមត្រូវទេ។");
      } else {
        setFormError(
          error instanceof ApiError
            ? "មិនអាចចូលគណនីបានទេ។ សូមព្យាយាមម្តងទៀត។"
            : "មិនអាចភ្ជាប់ទៅសេវាកម្មបានទេ។ សូមព្យាយាមម្តងទៀត។",
        );
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <Card className="w-full gap-0 rounded-3xl border-0 bg-white py-0 shadow-[0_24px_70px_-24px_oklch(0.32_0.08_235/0.35),0_2px_8px_oklch(0_0_0/0.06)] ring-1 ring-sky-950/8">
      <CardHeader className="gap-2 px-6 pt-7 pb-0 sm:px-9 sm:pt-9">
        <span className="mb-2 flex size-11 items-center justify-center rounded-2xl bg-sky-100 text-sky-800 ring-1 ring-inset ring-sky-900/8">
          <LockKeyhole aria-hidden="true" className="size-5" strokeWidth={1.8} />
        </span>
        <CardTitle className="text-2xl font-semibold tracking-tight text-zinc-950">
          ចូលគណនីបុគ្គលិក
        </CardTitle>
        <CardDescription className="text-base leading-6 text-zinc-600">
          ប្រើគណនីបុគ្គលិកដែលបានបង្កើតរួច។
        </CardDescription>
      </CardHeader>

      <CardContent className="px-6 pt-7 pb-7 sm:px-9 sm:pt-8 sm:pb-9">
        <form noValidate onSubmit={handleSubmit}>
          <label className="mb-2 block text-sm font-medium text-zinc-800" htmlFor="staff-email">
            អាសយដ្ឋានអ៊ីមែល
          </label>
          <div className="relative">
            <Mail aria-hidden="true" className="pointer-events-none absolute top-1/2 left-3.5 size-5 -translate-y-1/2 text-zinc-400" strokeWidth={1.8} />
            <input
              autoComplete="email"
              className="h-12 w-full rounded-xl border border-zinc-200 bg-white pr-4 pl-11 text-base text-zinc-950 outline-none transition-[border-color,box-shadow] duration-150 placeholder:text-zinc-400 hover:border-zinc-300 focus:border-sky-600 focus:ring-4 focus:ring-sky-600/12"
              id="staff-email"
              name="email"
              onChange={(event) => {
                setEmail(event.target.value);
                setFormError("");
              }}
              placeholder="owner@example.com"
              required
              type="email"
              value={email}
            />
          </div>

          <label className="mt-5 mb-2 block text-sm font-medium text-zinc-800" htmlFor="staff-password">
            ពាក្យសម្ងាត់
          </label>
          <div className="relative">
            <LockKeyhole aria-hidden="true" className="pointer-events-none absolute top-1/2 left-3.5 size-5 -translate-y-1/2 text-zinc-400" strokeWidth={1.8} />
            <input
              autoComplete="current-password"
              className="h-12 w-full rounded-xl border border-zinc-200 bg-white pr-12 pl-11 text-base text-zinc-950 outline-none transition-[border-color,box-shadow] duration-150 placeholder:text-zinc-400 hover:border-zinc-300 focus:border-sky-600 focus:ring-4 focus:ring-sky-600/12"
              id="staff-password"
              minLength={8}
              name="password"
              onChange={(event) => {
                setPassword(event.target.value);
                setFormError("");
              }}
              placeholder="យ៉ាងតិច 8 តួអក្សរ"
              required
              type={showPassword ? "text" : "password"}
              value={password}
            />
            <button
              aria-label={showPassword ? "លាក់ពាក្យសម្ងាត់" : "បង្ហាញពាក្យសម្ងាត់"}
              className="absolute top-1/2 right-2 flex size-9 -translate-y-1/2 items-center justify-center rounded-lg text-zinc-500 outline-none transition-[background-color,color,transform] duration-150 hover:bg-zinc-100 hover:text-zinc-800 focus-visible:ring-4 focus-visible:ring-sky-600/15 active:scale-[0.96]"
              onClick={() => setShowPassword((current) => !current)}
              type="button"
            >
              {showPassword ? <EyeOff aria-hidden="true" className="size-5" /> : <Eye aria-hidden="true" className="size-5" />}
            </button>
          </div>

          {formError && (
            <div aria-live="polite" className="mt-5 flex gap-2.5 rounded-xl bg-red-50 px-4 py-3 text-sm leading-5 text-red-700 ring-1 ring-inset ring-red-600/10" role="alert">
              <CircleAlert aria-hidden="true" className="mt-0.5 size-4 shrink-0" />
              {formError}
            </div>
          )}

          <Button
            className="mt-7 h-12 w-full rounded-xl bg-sky-700 px-5 text-base font-semibold text-white shadow-[0_10px_24px_-12px_oklch(0.5_0.15_240/0.8)] transition-[background-color,box-shadow,transform] duration-150 hover:bg-sky-800 active:scale-[0.96] active:translate-y-0"
            disabled={isSubmitting}
            type="submit"
          >
            {isSubmitting ? (
              <><LoaderCircle aria-hidden="true" className="size-5 animate-spin" />កំពុងចូលគណនី…</>
            ) : (
              <>បើកផ្ទាំងគ្រប់គ្រង<ArrowRight aria-hidden="true" className="size-5" strokeWidth={2} /></>
            )}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
