"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";

import { LogoutButton } from "@/components/logout-button";
import { getCurrentUser, type AuthenticatedUser } from "@/lib/auth-client";

const navigationItems = [
  { href: "/", label: "Dashboard" },
  { href: "/invoices", label: "Invoices" },
  { href: "/payments", label: "Payments" },
  { href: "/reconciliations", label: "Reconciliations" },
  { href: "/customers", label: "Customers" },
  { href: "/ledger", label: "Ledger" },
];

export function Navigation() {
  const pathname = usePathname();
  const [currentUser, setCurrentUser] = useState<AuthenticatedUser | null>(null);

  useEffect(() => {
    getCurrentUser()
      .then(setCurrentUser)
      .catch(() => setCurrentUser(null));
  }, []);

  if (pathname === "/login") {
    return null;
  }

  const visibleNavigationItems = currentUser?.role === "finance_admin" ? [...navigationItems,{href:"/users",label:"Users"},{href:"/refunds",label:"Refunds"}] : currentUser?.role === "approver" ? [...navigationItems,{href:"/refunds",label:"Refunds"}] : navigationItems;

  return (
    <header className="top-navigation">
      <nav className="navigation-content">
        <Link className="brand" href="/">
          Ledgerflow
        </Link>

        <div className="navigation-links">
          {visibleNavigationItems.map((item) => {
            const isActive = pathname === item.href;

            return (
              <Link
                className={isActive ? "navigation-link active" : "navigation-link"}
                href={item.href}
                key={item.href}
              >
                {item.label}
              </Link>
            );
          })}
        </div>
        <LogoutButton />
      </nav>
    </header>
  );
}
