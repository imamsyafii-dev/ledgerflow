"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const navigationItems = [
  { href: "/", label: "Dashboard" },
  { href: "/invoices", label: "Invoices" },
  { href: "/customers", label: "Customers" },
  { href: "/ledger", label: "Ledger" },
];

export function Navigation() {
  const pathname = usePathname();

  return (
    <header className="top-navigation">
      <nav className="navigation-content">
        <Link className="brand" href="/">
          Ledgerflow
        </Link>

        <div className="navigation-links">
          {navigationItems.map((item) => {
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
      </nav>
    </header>
  );
}
