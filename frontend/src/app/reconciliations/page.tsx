import Link from "next/link";
import { getReconciliations } from "@/lib/api";

export const instant = false;

function formatCurrency(amount: number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(amount);
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));
}

export default async function ReconciliationsPage() {
  const reconciliations = await getReconciliations();

  return (
    <main className="dashboard">
      <section className="hero">
        <div>
          <p className="eyebrow">Finance controls</p>
          <h1>Reconciliations</h1>
          <p className="subtitle">
            Hasil pencocokan otomatis antara payment dan invoice.
          </p>
        </div>

        <Link className="back-link" href="/">
          ← Dashboard
        </Link>
      </section>

      <section className="panel">
        {reconciliations.length === 0 ? (
          <p className="empty-state">Belum ada reconciliation.</p>
        ) : (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Invoice</th>
                  <th>Customer</th>
                  <th>Payment ID</th>
                  <th>Status</th>
                  <th>Reconciled At</th>
                  <th className="amount">Amount</th>
                </tr>
              </thead>
              <tbody>
                {reconciliations.map((reconciliation) => (
                  <tr key={reconciliation.id}>
                    <td>{reconciliation.invoice_number}</td>
                    <td>{reconciliation.customer_name}</td>
                    <td>{reconciliation.provider_payment_id}</td>
                    <td>
                      <span className="badge badge-paid">
                        {reconciliation.status}
                      </span>
                    </td>
                    <td>{formatDate(reconciliation.reconciled_at)}</td>
                    <td className="amount">
                      {formatCurrency(reconciliation.amount)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </main>
  );
}