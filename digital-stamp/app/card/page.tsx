import type { Metadata } from "next";

import { CustomerCardView } from "@/components/customer/customer-card-view";

export const metadata: Metadata = {
  title: "My stamp card | Digital Stamp Card",
  description: "View your current stamps and reward progress.",
};

export default function CardPage() {
  return <CustomerCardView />;
}
