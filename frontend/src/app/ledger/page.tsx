import Link from "next/link";
import { getJournalEntries } from "@/lib/api";

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

export default async function LedgerPage() {
  const journalEntries = await getJournalEntries();

  return (
    <main className="dashboard">
      <section className="hero">
        <div>
          <p className="eyebrow">Accounting</p>
          <h1>General Ledger</h1>
          <p className="subtitle">
            Journal entry otomatis yang dibuat dari payment terrekonsiliasi.
          </p>
        </div>

        <Link className="back-link" href="/">
          ← Dashboard
        </Link>
      </section>

      <section className="ledger-entries">
        {journalEntries.length === 0 ? (
          <article className="panel">
            <p className="empty-state">Belum ada journal entry.</p>
          </article>
        ) : (
          journalEntries.map((entry) => (
            <article className="panel ledger-entry" key={entry.id}>
              <div className="panel-header">
                <div>
                  <p className="eyebrow">{entry.reference_type}</p>
                  <h2>{entry.description}</h2>
                </div>

                <span className="entry-date">
                  {formatDate(entry.entry_date)}
                </span>
              </div>

              <div className="table-wrapper">
                <table>
                  <thead>
                    <tr>
                      <th>Account</th>
                      <th className="amount">Debit</th>
                      <th className="amount">Credit</th>
                    </tr>
                  </thead>
                  <tbody>
                    {entry.lines.map((line) => (
                      <tr key={`${entry.id}-${line.account_code}`}>
                        <td>
                          {line.account_code} — {line.account_name}
                        </td>
                        <td className="amount">
                          {line.debit > 0
                            ? formatCurrency(line.debit, line.currency)
                            : "-"}
                        </td>
                        <td className="amount">
                          {line.credit > 0
                            ? formatCurrency(line.credit, line.currency)
                            : "-"}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </article>
          ))
        )}
      </section>
    </main>
  );
}
