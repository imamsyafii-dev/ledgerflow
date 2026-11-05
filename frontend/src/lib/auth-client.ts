export type AuthenticatedUser = {
  id: string;
  organization_id: string;
  name: string;
  email: string;
  role: "finance_admin" | "approver" | "viewer";
  is_active: boolean;
};

type ApiResponse<T> = {
  data: T;
};

function getAPIBaseURL() {
  const apiBaseURL = process.env.NEXT_PUBLIC_API_BASE_URL;

  if (!apiBaseURL) {
    throw new Error("NEXT_PUBLIC_API_BASE_URL harus diisi");
  }

  return apiBaseURL;
}

export async function login(email: string, password: string) {
  const response = await fetch(`${getAPIBaseURL()}/api/v1/auth/login`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email,
      password,
    }),
  });

  if (!response.ok) {
    const payload = await response.json().catch(() => null);

    throw new Error(payload?.error ?? "Failed to login");
  }

  const payload = await response.json() as ApiResponse<AuthenticatedUser>;

  return payload.data;
}

export async function logout() {
  const response = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/auth/logout`,
    {
      method: "POST",
      credentials: "include",
    },
  );

  if (!response.ok && response.status !== 204) {
    throw new Error("Logout gagal. Silakan coba lagi.");
  }
}