import type { Customer, Staff, StampCard } from "@/types/domain";

/* Request body for POST /staff-sessions. */
export type LoginStaff = {
  email: string;
  password: string;
};

export type StaffLoginResult = {
  staff: Staff;
};

export type StampScanPreview = {
  scan_id: string;
  customer: Customer;
  card: StampCard;
  expires_at: string;
};

export type StampConfirmation = {
  customer: Customer;
  card: StampCard;
  reward_available: boolean;
};

/* Standard successful response returned by the Gin API. */
export type StaffApiResponse<T> = {
  data: T;
};
