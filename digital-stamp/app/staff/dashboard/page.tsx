import type { Metadata } from "next";

import { StaffDashboard } from "@/components/staff/staff-dashboard";

export const metadata: Metadata = {
  title: "Staff dashboard | Digital Stamp Card",
  description: "Manage Digital Stamp Card shop activity.",
};

export default function StaffDashboardPage() {
  return <StaffDashboard />;
}
