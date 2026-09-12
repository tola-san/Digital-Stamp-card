"use client";

import { FormEvent, useState } from "react";
import Link from "next/link";
import {
  ArrowRight,
  Check,
  CircleAlert,
  LoaderCircle,
  Phone,
  UserRound,
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
import { registerCustomer } from "@/services/customer-service";
import type { RegisterCustomerResult } from "@/types/customer";

type FieldErrors = {
  name?: string;
  phone?: string;
};

const internationalPhonePattern = /^\+[1-9][0-9]{7,14}$/;

function normalizePhone(phone: string) {
  return phone.replace(/[\s\-().]/g, "");
}

export function CustomerRegistrationForm() {

  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({});
  const [formError, setFormError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [registration, setRegistration] =
    useState<RegisterCustomerResult | null>(null);

  function validateForm() {
    const errors: FieldErrors = {};
    const trimmedName = name.trim();
    const normalizedPhone = normalizePhone(phone);

    if (trimmedName.length < 2 || trimmedName.length > 100) {
      errors.name = "Enter a name between 2 and 100 characters.";
    }

    if (!internationalPhonePattern.test(normalizedPhone)) {
      errors.phone = "Use international format, for example +855 12 345 678.";
    }

    setFieldErrors(errors);
    return Object.keys(errors).length === 0;
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFormError("");

    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);

    try {
      const result = await registerCustomer({
        name: name.trim(),
        phone: normalizePhone(phone),
      });

      setRegistration(result);
    } catch (error) {
      if (error instanceof ApiError) {
        if (error.code === "PHONE_ALREADY_REGISTERED") {
          setFieldErrors({ phone: "This phone number is already registered." });
        } else if (error.code === "INVALID_NAME") {
          setFieldErrors({ name: error.message });
        } else if (error.code === "INVALID_PHONE") {
          setFieldErrors({ phone: error.message });
        } else {
          setFormError(error.message);
        }
      } else {
        setFormError("We could not connect to the service. Please try again.");
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  if (registration) {
    return (
      <Card className="w-full gap-0 rounded-xl border-0 bg-white py-0 shadow-[0_24px_70px_-24px_oklch(0.32_0.08_235/0.35),0_2px_8px_oklch(0_0_0/0.06)] ring-1 ring-sky-950/8">
        <CardContent className="flex flex-col items-center px-6 py-10 text-center sm:px-10 sm:py-12">
          <div className="flex size-16 items-center justify-center rounded-2xl bg-sky-600 text-white shadow-[0_12px_24px_-10px_oklch(0.58_0.16_240/0.7)]">
            <Check aria-hidden="true" className="size-8" strokeWidth={2.5} />
          </div>
          <p className="mt-6 text-sm font-semibold tracking-wide text-sky-700 uppercase">
            Registration complete
          </p>
          <h2 className="mt-2 text-3xl font-semibold tracking-tight text-zinc-950">
            Welcome, {registration.customer.name}
          </h2>
          <p className="mt-3 max-w-sm text-base leading-7 text-zinc-600">
            Your digital stamp card is ready. Your session has been securely
            saved on this device.
          </p>

          <div className="mt-8 flex w-full items-center justify-between rounded-2xl bg-sky-50 p-5 ring-1 ring-inset ring-sky-900/8">
            <div className="text-left">
              <p className="text-sm font-medium text-sky-700">Starting balance</p>
              <p className="mt-1 text-sm text-sky-800/70">
                {registration.customer.phone}
              </p>
            </div>
            <div className="text-right">
              <span className="text-3xl font-semibold tabular-nums text-sky-700">
                {registration.card.stamp_count}
              </span>
              <p className="text-xs font-medium text-sky-800/70">stamps</p>
            </div>
          </div>

          <Link
            className="mt-6 flex h-12 w-full items-center justify-center gap-2 rounded-xl bg-sky-700 px-5 text-base font-semibold text-white shadow-[0_10px_24px_-12px_oklch(0.5_0.15_240/0.8)] transition-[background-color,box-shadow,transform] duration-150 hover:bg-sky-800 active:scale-[0.96]"
            href="/card"
          >
            View my stamp card
            <ArrowRight aria-hidden="true" className="size-5" strokeWidth={2} />
          </Link>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="w-full gap-0 rounded-xl border-0 bg-white py-0 shadow-[0_24px_70px_-24px_oklch(0.32_0.08_235/0.35),0_2px_8px_oklch(0_0_0/0.06)] ring-1 ring-sky-950/8">
      <CardHeader className="gap-2 px-6 pt-7 pb-0 sm:px-9 sm:pt-9">
        <CardTitle className="text-2xl font-semibold tracking-tight text-zinc-950">
          Create your account
        </CardTitle>
        <CardDescription className="text-base leading-6 text-zinc-600">
          Enter your details to receive your digital stamp card.
        </CardDescription>
      </CardHeader>

      <CardContent className="px-6 pt-7 pb-7 sm:px-9 sm:pt-8 sm:pb-9">
        <form noValidate onSubmit={handleSubmit}>
          <div className="space-y-5">
            <div>
              <label
                className="mb-2 block text-sm font-medium text-zinc-800"
                htmlFor="customer-name"
              >
                Full name
              </label>
              <div className="relative">
                <UserRound
                  aria-hidden="true"
                  className="pointer-events-none absolute top-1/2 left-3.5 size-5 -translate-y-1/2 text-zinc-400"
                  strokeWidth={1.8}
                />
                <input
                  aria-describedby={fieldErrors.name ? "name-error" : undefined}
                  aria-invalid={Boolean(fieldErrors.name)}
                  autoComplete="name"
                  className="h-12 w-full rounded-xl border border-zinc-200 bg-white pr-4 pl-11 text-base text-zinc-950 outline-none transition-[border-color,box-shadow] duration-150 placeholder:text-zinc-400 hover:border-zinc-300 focus:border-sky-600 focus:ring-4 focus:ring-sky-600/12 aria-invalid:border-red-500 aria-invalid:ring-4 aria-invalid:ring-red-500/10"
                  id="customer-name"
                  maxLength={100}
                  name="name"
                  onChange={(event) => {
                    setName(event.target.value);
                    setFieldErrors((errors) => ({ ...errors, name: undefined }));
                  }}
                  placeholder="Sokha Chan"
                  type="text"
                  value={name}
                />
              </div>
              {fieldErrors.name && (
                <p className="mt-2 flex items-center gap-1.5 text-sm text-red-600" id="name-error">
                  <CircleAlert aria-hidden="true" className="size-4" />
                  {fieldErrors.name}
                </p>
              )}
            </div>

            <div>
              <label
                className="mb-2 block text-sm font-medium text-zinc-800"
                htmlFor="customer-phone"
              >
                Phone number
              </label>
              <div className="relative">
                <Phone
                  aria-hidden="true"
                  className="pointer-events-none absolute top-1/2 left-3.5 size-5 -translate-y-1/2 text-zinc-400"
                  strokeWidth={1.8}
                />
                <input
                  aria-describedby={fieldErrors.phone ? "phone-error" : "phone-hint"}
                  aria-invalid={Boolean(fieldErrors.phone)}
                  autoComplete="tel"
                  className="h-12 w-full rounded-xl border border-zinc-200 bg-white pr-4 pl-11 text-base text-zinc-950 outline-none transition-[border-color,box-shadow] duration-150 placeholder:text-zinc-400 hover:border-zinc-300 focus:border-sky-600 focus:ring-4 focus:ring-sky-600/12 aria-invalid:border-red-500 aria-invalid:ring-4 aria-invalid:ring-red-500/10"
                  id="customer-phone"
                  inputMode="tel"
                  name="phone"
                  onChange={(event) => {
                    setPhone(event.target.value);
                    setFieldErrors((errors) => ({ ...errors, phone: undefined }));
                  }}
                  placeholder="+855 12 345 678"
                  type="tel"
                  value={phone}
                />
              </div>
              {fieldErrors.phone ? (
                <p className="mt-2 flex items-center gap-1.5 text-sm text-red-600" id="phone-error">
                  <CircleAlert aria-hidden="true" className="size-4" />
                  {fieldErrors.phone}
                </p>
              ) : (
                <p className="mt-2 text-sm text-zinc-500" id="phone-hint">
                  Include your country code. You will use this number to return.
                </p>
              )}
            </div>
          </div>

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
                Creating your card…
              </>
            ) : (
              <>
                Create my stamp card
                <ArrowRight aria-hidden="true" className="size-5" strokeWidth={2} />
              </>
            )}
          </Button>

          <p className="mt-4 text-center text-xs leading-5 text-zinc-500">
            By continuing, you agree to keep your phone number associated with
            this loyalty account.
          </p>

          <p className="mt-4 text-center text-sm text-zinc-600">
            Already registered?{" "}
            <Link
              className="font-semibold text-sky-700 underline-offset-4 hover:underline focus-visible:rounded-sm focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-sky-600/15"
              href="/login"
            >
              Open your card
            </Link>
          </p>
        </form>
      </CardContent>
    </Card>
  );
}
