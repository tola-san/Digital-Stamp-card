import type { Customer, StampCard, StampTransaction } from "@/types/domain";

/* Request body for POST /customers. */
export type RegisterCustomer = {
  name: string;
  phone: string;
};

/* Request body for POST /customer-sessions. */
export type LoginCustomer = {
  phone: string;
};

export type RegisterCustomerResult = {
  customer: Customer;
  card: StampCard;
};

export type LoginCustomerResult = {
  customer: Customer;
};

export type CurrentCustomerResult = {
  customer: Customer;
};

export type CustomerCardResult = {
  card: StampCard;
};

export type TransactionPage = {
  transactions: StampTransaction[];
  next_cursor: string | null;
};

/*
 * Standard response returned by the Gin API.
 *
 * Example:
 * {
 *   "data": {
 *     ...
 *   }
 * }
 */
export type ApiResponse<T> = {
  data: T;
};
