import axios, { AxiosError } from "axios";

export const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api";

/* Error response returned by the Gin API. */
type ApiErrorBody = {
  error?: {
    code?: string;
    message?: string;
  };
};

export class ApiError extends Error {
  constructor(
    message: string,
    public readonly status: number,
    public readonly code = "API_ERROR",
  ) {
    super(message);
    this.name = "ApiError";
  }
}

/* Shared HTTP client for calls from interactive Client Components. */
export const api = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10_000,
});

api.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiErrorBody>) => {
    const status = error.response?.status ?? 500;
    const message =
      error.response?.data?.error?.message ??
      error.message ??
      "The request could not be completed.";
    const code = error.response?.data?.error?.code ?? "API_ERROR";

    return Promise.reject(new ApiError(message, status, code));
  },
);
