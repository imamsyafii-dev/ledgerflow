import Link from "next/link";
import { redirect } from "next/navigation";

import { ExecuteRefundButton } from "@/components/execute-refund-button";
import { RefundDecisionActions } from "@/components/refund-decision-actions";
import {
  getCurrentUser,
  getRefundRequests,
} from "@/lib/server-api";

export const instant = false;

function formatCurrency(amount: number, currency: string) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency,
    maximumFractionDigits: 0,
  }).format(amount);
}

function getStatusClass(status: string) {
  if (status === "COMPLETED") {
    return "badge badge-paid";
  }

  if (status === "REJECTED" || status === "FAILED") {
    return "badge badge-void";
  }

  return "badge badge-issued";
}

export default async function RefundsPage() {
  const user = await getCurrentUser();

  if (user.role === "viewer") {
    redirect("/");
  }

  const refunds = await getRefundRequests();

  return (
    <main className="dashboard">
      <section className="hero">
        <div>
          <p className="eyebrow">Approval workflow</p>
          <h1>Refund Requests</h1>
          <p className="subtitle">
            Kelola approval dan eksekusi refund payment.
          </p>
        </div>

        <Link className="back-link" href="/">
          ← Dashboard
        </Link>
      </section>

      <section className="panel">
        {refunds.length === 0 ? (
          <p className="empty-state">Belum ada refund request.</p>
        ) : (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Invoice</th>
                  <th>Customer</th>
                  <th className="amount">Amount</th>
                  <th>Requested By</th>
                  <th>Reason</th>
                  <th>Status</th>
                  <th>Provider Refund ID</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody>
                {refunds.map((refund) => (
                  <tr key={refund.id}>
                    <td>{refund.invoice_number}</td>
                    <td>{refund.customer_name}</td>
                    <td className="amount">
                      {formatCurrency(refund.amount, refund.currency)}
                    </td>
                    <td>{refund.requested_by_name}</td>
                    <td>{refund.reason}</td>
                    <td>
                      <span className={getStatusClass(refund.status)}>
                        {refund.status}
                      </span>
                    </td>
                    <td>{refund.provider_refund_id ?? "-"}</td>
                    <td>
                      {user.role === "approver" &&
                      refund.status === "PENDING" ? (
                        <RefundDecisionActions refund={refund} />
                      ) : null}

                      {user.role === "finance_admin" &&
                      refund.status === "APPROVED" ? (
                        <ExecuteRefundButton refundID={refund.id} />
                      ) : null}

                      {refund.status !== "PENDING" &&
                      refund.status !== "APPROVED" ? (
                        <span className="action-hint">No action</span>
                      ) : null}
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