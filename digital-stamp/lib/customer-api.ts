import { api } from "@/lib/api";


import type {
  ApiResponse,
  CurrentCustomerResult,
  CustomerCardResult,
  LoginCustomer,
  LoginCustomerResult,
  RegisterCustomer,
  RegisterCustomerResult,
  TransactionPage,
} from "@/types/customer";

/* Register a customer and create their browser session. */
export async function registerCustomer(
  input: RegisterCustomer,
): Promise<RegisterCustomerResult> {
  const response = await api.post<ApiResponse<RegisterCustomerResult>>(
    "/customers",
    input,
  );

  return response.data.data;
}

/* Log in an existing customer using their phone number. */
export async function loginCustomer(
  input: LoginCustomer,
): Promise<LoginCustomerResult> {
  const response = await api.post<ApiResponse<LoginCustomerResult>>(
    "/customer-sessions",
    input,
  );

  return response.data.data;
}

/* Load the customer connected to the HttpOnly session cookie. */
export async function getCurrentCustomer(): Promise<CurrentCustomerResult> {
  const response = await api.get<ApiResponse<CurrentCustomerResult>>(
    "/customers/me",
  );

  return response.data.data;
}

/* Load the current customer's digital stamp card. */
export async function getCustomerCard(): Promise<CustomerCardResult> {
  const response = await api.get<ApiResponse<CustomerCardResult>>(
    "/customers/me/card",
  );

  return response.data.data;
}

/* Load one cursor-based page of the customer's transaction history. */
export async function getCustomerTransactions(
  limit = 20,
  cursor?: string,
): Promise<TransactionPage> {
  const response = await api.get<ApiResponse<TransactionPage>>(
    "/customers/me/transactions",
    {
      params: {
        limit,
        ...(cursor ? { cursor } : {}),
      },
    },
  );

  return response.data.data;
}

/* Revoke the current customer session and clear its cookie. */
export async function logoutCustomer(): Promise<void> {
  await api.delete("/customer-sessions/current");
}
