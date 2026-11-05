import Link from "next/link";
import { CreateInvoiceForm } from "@/components/create-invoice-form";
import { getCustomers, getInvoices } from "@/lib/api";

export const instant = false;

function formatCurrency(amount: number, currency: string) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency,
    maximumFractionDigits: 0,
  }).format(amount);
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
  }).format(new Date(value));
}

export default async function InvoicesPage() {
  const [invoices, customers] = await Promise.all([
    getInvoices(),
    getCustomers(),
  ]);

  return (
    <main className="dashboard">
      <section className="hero">
        <div>
          <p className="eyebrow">Receivables</p>
          <h1>Invoices</h1>
          <p className="subtitle">
            Buat dan monitor invoice untuk customer Anda.
          </p>
        </div>

        <Link className="back-link" href="/">
          ← Dashboard
        </Link>
      </section>

      <section className="panel form-panel">
        <div className="panel-header">
          <div>
            <p className="eyebrow">New invoice</p>
            <h2>Create Invoice</h2>
          </div>
        </div>

        <CreateInvoiceForm customers={customers} />
      </section>

      <section className="panel">
        <div className="panel-header">
          <div>
            <p className="eyebrow">Receivables</p>
            <h2>All Invoices</h2>
          </div>
        </div>

        {invoices.length === 0 ? (
          <p className="empty-state">Belum ada invoice.</p>
        ) : (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Invoice</th>
                  <th>Customer</th>
                  <th>Due Date</th>
                  <th>Status</th>
                  <th className="amount">Amount</th>
                </tr>
              </thead>
              <tbody>
                {invoices.map((invoice) => (
                  <tr key={invoice.id}>
                    <td>{invoice.number}</td>
                    <td>{invoice.customer_id}</td>
                    <td>{formatDate(invoice.due_date)}</td>
                    <td>
                      <span className={`badge badge-${invoice.status.toLowerCase()}`}>
                        {invoice.status}
                      </span>
                    </td>
                    <td className="amount">
                      {formatCurrency(invoice.amount, invoice.currency)}
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