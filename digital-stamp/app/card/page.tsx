import type { Metadata } from "next";

import { CustomerCardView } from "@/components/customer/customer-card-view";

export const metadata: Metadata = {
  title: "កាតត្រារបស់ខ្ញុំ | កាតត្រាឌីជីថល",
  description: "មើលត្រាបច្ចុប្បន្ន និងដំណើរឆ្ពោះទៅរករង្វាន់របស់អ្នក។",
};

export default function CardPage() {
  return <CustomerCardView />;
}
