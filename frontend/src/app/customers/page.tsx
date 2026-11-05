import Link from "next/link";
import { getCustomers } from "@/lib/api";

export const instant = false;

function formatDate(value: string) {
  return new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
  }).format(new Date(value));
}

export default async function CustomersPage() {
  const customers = await getCustomers();

  return (
    <main className="dashboard">
      <section className="hero">
        <div>
          <p className="eyebrow">Directory</p>
          <h1>Customers</h1>
          <p className="subtitle">
            Daftar customer yang terhubung dengan invoice Anda.
          </p>
        </div>

        <Link className="back-link" href="/">
          ← Dashboard
        </Link>
      </section>

      <article className="panel">
        {customers.length === 0 ? (
          <p className="empty-state">Belum ada customer.</p>
        ) : (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Email</th>
                  <th>Created At</th>
                </tr>
              </thead>
              <tbody>
                {customers.map((customer) => (
                  <tr key={customer.id}>
                    <td>{customer.name}</td>
                    <td>{customer.email ?? "-"}</td>
                    <td>{formatDate(customer.created_at)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </article>
    </main>
  );
}