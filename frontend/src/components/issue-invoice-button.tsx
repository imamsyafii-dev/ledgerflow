"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { issueInvoice } from "@/lib/api";

type IssueInvoiceButtonProps = {
  invoiceID: string;
};

export function IssueInvoiceButton({
  invoiceID,
}: IssueInvoiceButtonProps) {
  const router = useRouter();

  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleIssue() {
    setError("");
    setIsSubmitting(true);

    try {
      await issueInvoice(invoiceID);
      router.refresh();
    } catch (error) {
      setError(
        error instanceof Error ? error.message : "Failed to issue invoice",
      );
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div className="invoice-action">
      <button
        className="secondary-button"
        type="button"
        onClick={handleIssue}
        disabled={isSubmitting}
      >
        {isSubmitting ? "Issuing..." : "Issue"}
      </button>

      {error && <span className="action-error">{error}</span>}
    </div>
  );
}