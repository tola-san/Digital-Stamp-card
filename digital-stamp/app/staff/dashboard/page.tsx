import type { Metadata } from "next";

import { StaffDashboard } from "@/components/staff/staff-dashboard";

export const metadata: Metadata = {
  title: "ផ្ទាំងគ្រប់គ្រងបុគ្គលិក | កាតត្រាឌីជីថល",
  description: "គ្រប់គ្រងសកម្មភាពកាតត្រាឌីជីថលរបស់ហាង។",
};

export default function StaffDashboardPage() {
  return <StaffDashboard />;
}
