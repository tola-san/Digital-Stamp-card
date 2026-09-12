import { api } from "@/lib/api";
import type {
  LoginStaff,
  StaffApiResponse,
  StaffLoginResult,
} from "@/types/staff";

/* Log in a staff member and let Gin create the HttpOnly session cookie. */
export async function loginStaff(input: LoginStaff): Promise<StaffLoginResult> {
  const response = await api.post<StaffApiResponse<StaffLoginResult>>(
    "/staff-sessions",
    input,
  );

  return response.data.data;
}

/* Load the staff member connected to the current session cookie. */
export async function getCurrentStaff(): Promise<StaffLoginResult> {
  const response = await api.get<StaffApiResponse<StaffLoginResult>>("/staff/me");

  return response.data.data;
}

/* Revoke the current staff session and clear its cookie. */
export async function logoutStaff(): Promise<void> {
  await api.delete("/staff-sessions/current");
}
