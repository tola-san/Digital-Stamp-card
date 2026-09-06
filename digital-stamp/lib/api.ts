
import axios, {AxiosError} from "axios";
import { error } from "console";

export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api";
 


// Defined type of ApiErrorbody
type ApiErrorBody = {
  error?: {
    code?: string
    message?: string
  }
}


export class ApiError extends Error {

  constructor(
    message: string,
    public readonly status: number,
    public readonly code = "API_ERROR"
  ) {
    super(message);
    this.name = "ApiError";
  }
}


//  create handling method POST by consume axios HTTP

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

    const code =
      error.response?.data?.error?.code ??
      "API_ERROR";

    return Promise.reject(
      new ApiError(message, status, code)
    );
  }
);

// export async function apiFetch<T>(path: string, init: RequestInit = {}): Promise<T> {
//   const response = await fetch(`${API_URL}${path.startsWith("/") ? path : `/${path}`}`, {
//     ...init,
//     credentials: "include",
//     headers: {
//       ...(init.body ? { "Content-Type": "application/json" } : {}),
//       ...init.headers,
//     },
//   })

//   if (!response.ok) {
//     const body = (await response.json().catch(() => ({}))) as ApiErrorBody
//     throw new ApiError(
//       body.error?.message ?? "The request could not be completed.",
//       response.status,
//       body.error?.code,
//     )
//   }

//   if (response.status === 204) {
//     return undefined as T
//   }

//   return response.json() as Promise<T>
// }
