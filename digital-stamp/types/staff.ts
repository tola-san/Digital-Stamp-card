import type { Staff } from "@/types/domain";

/* Request body for POST /staff-sessions. */
export type LoginStaff = {
  email: string;
  password: string;
};

export type StaffLoginResult = {
  staff: Staff;
};

/* Standard successful response returned by the Gin API. */
export type StaffApiResponse<T> = {
  data: T;
};
