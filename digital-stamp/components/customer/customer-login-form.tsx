"use client";

import { FormEvent, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { ArrowRight, CircleAlert, LoaderCircle, Phone } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ApiError } from "@/lib/api";
import { loginCustomer } from "@/services/customer-service";

const internationalPhonePattern = /^\+[1-9][0-9]{7,14}$/;

function normalizePhone(phone: string) {
  return phone.replace(/[\s\-().]/g, "");
}

export function CustomerLoginForm() {
  const router = useRouter();
  const [phone, setPhone] = useState("");
  const [phoneError, setPhoneError] = useState("");
  const [formError, setFormError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPhoneError("");
    setFormError("");

    const normalizedPhone = normalizePhone(phone);
    if (!internationalPhonePattern.test(normalizedPhone)) {
      setPhoneError("សូមប្រើទម្រង់អន្តរជាតិ ឧទាហរណ៍ +855 12 345 678។");
      return;
    }

    setIsSubmitting(true);

    try {
      await loginCustomer({ phone: normalizedPhone });
      router.replace("/card");
    } catch (error) {
      if (error instanceof ApiError) {
        if (
          error.code === "INVALID_CREDENTIALS" ||
          error.code === "CUSTOMER_NOT_FOUND"
        ) {
          setPhoneError("រកមិនឃើញគណនីដែលប្រើលេខទូរសព្ទនេះទេ។");
        } else if (error.code === "INVALID_PHONE") {
          setPhoneError("លេខទូរសព្ទដែលបានបញ្ចូលមិនត្រឹមត្រូវទេ។");
        } else {
          setFormError("មិនអាចចូលគណនីបានទេ។ សូមព្យាយាមម្តងទៀត។");
        }
      } else {
        setFormError("មិនអាចភ្ជាប់ទៅសេវាកម្មបានទេ។ សូមព្យាយាមម្តងទៀត។");
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <Card className="w-full gap-0 rounded-3xl border-0 bg-white py-0 shadow-[0_24px_70px_-24px_oklch(0.32_0.08_235/0.35),0_2px_8px_oklch(0_0_0/0.06)] ring-1 ring-sky-950/8">
      <CardHeader className="gap-2 px-6 pt-7 pb-0 sm:px-9 sm:pt-9">
        <CardTitle className="text-2xl font-semibold tracking-tight text-zinc-950">
          សូមស្វាគមន៍មកវិញ
        </CardTitle>
        <CardDescription className="text-base leading-6 text-zinc-600">
          បញ្ចូលលេខទូរសព្ទដែលភ្ជាប់ជាមួយកាតត្រារបស់អ្នក។
        </CardDescription>
      </CardHeader>

      <CardContent className="px-6 pt-7 pb-7 sm:px-9 sm:pt-8 sm:pb-9">
        <form noValidate onSubmit={handleSubmit}>
          <label
            className="mb-2 block text-sm font-medium text-zinc-800"
            htmlFor="login-phone"
          >
            លេខទូរសព្ទ
          </label>
          <div className="relative">
            <Phone
              aria-hidden="true"
              className="pointer-events-none absolute top-1/2 left-3.5 size-5 -translate-y-1/2 text-zinc-400"
              strokeWidth={1.8}
            />
            <input
              aria-describedby={phoneError ? "login-phone-error" : "login-phone-hint"}
              aria-invalid={Boolean(phoneError)}
              autoComplete="tel"
              className="h-12 w-full rounded-xl border border-zinc-200 bg-white pr-4 pl-11 text-base text-zinc-950 outline-none transition-[border-color,box-shadow] duration-150 placeholder:text-zinc-400 hover:border-zinc-300 focus:border-sky-600 focus:ring-4 focus:ring-sky-600/12 aria-invalid:border-red-500 aria-invalid:ring-4 aria-invalid:ring-red-500/10"
              id="login-phone"
              inputMode="tel"
              name="phone"
              onChange={(event) => {
                setPhone(event.target.value);
                setPhoneError("");
              }}
              placeholder="+855 12 345 678"
              type="tel"
              value={phone}
            />
          </div>

          {phoneError ? (
            <p
              className="mt-2 flex items-center gap-1.5 text-sm text-red-600"
              id="login-phone-error"
            >
              <CircleAlert aria-hidden="true" className="size-4" />
              {phoneError}
            </p>
          ) : (
            <p className="mt-2 text-sm text-zinc-500" id="login-phone-hint">
              សូមប្រើលេខកូដប្រទេសដូចគ្នានឹងពេលចុះឈ្មោះ។
            </p>
          )}

          {formError && (
            <div
              aria-live="polite"
              className="mt-5 flex gap-2.5 rounded-xl bg-red-50 px-4 py-3 text-sm leading-5 text-red-700 ring-1 ring-inset ring-red-600/10"
              role="alert"
            >
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
              <>
                <LoaderCircle aria-hidden="true" className="size-5 animate-spin" />
                កំពុងបើកកាតរបស់អ្នក…
              </>
            ) : (
              <>
                បើកកាតត្រារបស់ខ្ញុំ
                <ArrowRight aria-hidden="true" className="size-5" strokeWidth={2} />
              </>
            )}
          </Button>

          <p className="mt-6 text-center text-sm text-zinc-600">
            អតិថិជនថ្មី?{" "}
            <Link
              className="font-semibold text-sky-700 underline-offset-4 hover:underline focus-visible:rounded-sm focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-sky-600/15"
              href="/register"
            >
              បង្កើតកាតត្រា
            </Link>
          </p>
        </form>
      </CardContent>
    </Card>
  );
}
