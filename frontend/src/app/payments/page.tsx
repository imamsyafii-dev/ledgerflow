import Link from "next/link";
import { getPayments } from "@/lib/api";

export const instant = false;

function formatCurrency(amount: number, currency: string) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency,
    maximumFractionDigits: 0,
  }).format(amount);
}

function formatDate(value: string | null) {
  if (!value) {
    return "-";
  }

  return new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));
}

export default async function PaymentsPage() {
  const payments = await getPayments();

  return (
    <main className="dashboard">
      <section className="hero">
        <div>
          <p className="eyebrow">Cash collection</p>
          <h1>Payments</h1>
          <p className="subtitle">
            Transaksi yang diterima dari payment gateway dan direkonsiliasi.
          </p>
        </div>

        <Link className="back-link" href="/">
          ← Dashboard
        </Link>
      </section>

      <section className="panel">
        {payments.length === 0 ? (
          <p className="empty-state">Belum ada payment.</p>
        ) : (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Payment ID</th>
                  <th>Invoice</th>
                  <th>Customer</th>
                  <th>Provider</th>
                  <th>Status</th>
                  <th>Paid At</th>
                  <th className="amount">Amount</th>
                </tr>
              </thead>
              <tbody>
                {payments.map((payment) => (
                  <tr key={payment.id}>
                    <td>{payment.provider_payment_id}</td>
                    <td>{payment.invoice_number}</td>
                    <td>{payment.customer_name}</td>
                    <td>{payment.provider}</td>
                    <td>
                      <span className="badge badge-paid">
                        {payment.status}
                      </span>
                    </td>
                    <td>{formatDate(payment.paid_at)}</td>
                    <td className="amount">
                      {formatCurrency(payment.amount, payment.currency)}
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