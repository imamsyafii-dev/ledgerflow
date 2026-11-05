"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { logout } from "@/lib/auth-client";

export function LogoutButton() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);

  async function handleLogout() {
    setIsLoading(true);

    try {
      await logout();
      router.replace("/login");
      router.refresh();
    } catch (error) {
      alert(error instanceof Error ? error.message : "Logout gagal.");
      setIsLoading(false);
    }
  }

  return (
    <button
      className="logout-button"
      type="button"
      onClick={handleLogout}
      disabled={isLoading}
    >
      {isLoading ? "Logging out..." : "Logout"}
    </button>
  );
}
