import type { Metadata } from "next";

import { StaffQRScanner } from "@/components/staff/staff-qr-scanner";

export const metadata: Metadata = {
  title: "ស្កេន QR | កាតត្រាឌីជីថល",
  description: "ស្កេនកូដ QR របស់អតិថិជន និងបន្ថែមត្រា។",
};

export default function StaffScanPage() {
  return <StaffQRScanner />;
}
