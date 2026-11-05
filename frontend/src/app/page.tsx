import { getDashboardData } from "@/lib/server-api";

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

export default async function DashboardPage() {
  const { invoices, journalEntries } = await getDashboardData();

  const paidInvoices = invoices.filter((invoice) => invoice.status === "PAID");
  const openInvoices = invoices.filter(
    (invoice) => invoice.status === "DRAFT" || invoice.status === "ISSUED",
  );

  const paidAmount = paidInvoices.reduce(
    (total, invoice) => total + invoice.amount,
    0,
  );

  const openAmount = openInvoices.reduce(
    (total, invoice) => total + invoice.amount,
    0,
  );

  return (
    <main className="dashboard">
      <section className="hero">
        <div>
          <p className="eyebrow">Finance Operations</p>
          <h1>Ledgerflow Dashboard</h1>
          <p className="subtitle">
            Monitor invoice, payment reconciliation, dan pembukuan Anda.
          </p>
        </div>

        <span className="status">
          <span className="status-dot" />
          API connected
        </span>
      </section>

      <section className="metrics">
        <article className="metric-card">
          <p>Total Invoice</p>
          <strong>{invoices.length}</strong>
        </article>

        <article className="metric-card">
          <p>Outstanding</p>
          <strong>{formatCurrency(openAmount, "IDR")}</strong>
          <span>{openInvoices.length} invoice belum lunas</span>
        </article>

        <article className="metric-card">
          <p>Payment Reconciled</p>
          <strong>{formatCurrency(paidAmount, "IDR")}</strong>
          <span>{paidInvoices.length} invoice lunas</span>
        </article>

        <article className="metric-card">
          <p>Journal Entries</p>
          <strong>{journalEntries.length}</strong>
          <span>Double-entry postings</span>
        </article>
      </section>

      <section className="content-grid">
        <article className="panel">
          <div className="panel-header">
            <div>
              <p className="eyebrow">Receivables</p>
              <h2>Recent Invoices</h2>
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
                    <th>Due Date</th>
                    <th>Status</th>
                    <th className="amount">Amount</th>
                  </tr>
                </thead>
                <tbody>
                  {invoices.map((invoice) => (
                    <tr key={invoice.id}>
                      <td>{invoice.number}</td>
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
        </article>

        <article className="panel">
          <div className="panel-header">
            <div>
              <p className="eyebrow">Accounting</p>
              <h2>Recent Ledger Activity</h2>
            </div>
          </div>

          {journalEntries.length === 0 ? (
            <p className="empty-state">Belum ada journal entry.</p>
          ) : (
            <div className="journal-list">
              {journalEntries.slice(0, 5).map((entry) => (
                <div className="journal-item" key={entry.id}>
                  <div>
                    <strong>{entry.description}</strong>
                    <span>{formatDate(entry.entry_date)}</span>
                  </div>
                  <div className="journal-lines">
                    {entry.lines.map((line) => (
                      <span key={`${entry.id}-${line.account_code}`}>
                        {line.account_code}:{" "}
                        {line.debit > 0
                          ? `Debit ${formatCurrency(line.debit, line.currency)}`
                          : `Credit ${formatCurrency(line.credit, line.currency)}`}
                      </span>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          )}
        </article>
      </section>
    </main>
  );
}