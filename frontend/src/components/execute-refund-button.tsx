"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { executeRefund } from "@/lib/api";

type ExecuteRefundButtonProps = {
  refundID: string;
};

export function ExecuteRefundButton({
  refundID,
}: ExecuteRefundButtonProps) {
  const router = useRouter();
  const [isExecuting, setIsExecuting] = useState(false);
  const [error, setError] = useState("");

  async function handleExecute() {
    const confirmed = window.confirm(
      "Execute refund? This action will mark the payment as refunded.",
    );

    if (!confirmed) {
      return;
    }

    setError("");
    setIsExecuting(true);

    try {
      await executeRefund(refundID);
      router.refresh();
    } catch (error) {
      setError(
        error instanceof Error
          ? error.message
          : "Failed to execute refund",
      );
    } finally {
      setIsExecuting(false);
    }
  }

  return (
    <div className="refund-execution-action">
      <button
        className="secondary-button"
        type="button"
        onClick={handleExecute}
        disabled={isExecuting}
      >
        {isExecuting ? "Executing..." : "Execute Refund"}
      </button>

      {error && <p className="action-error">{error}</p>}
    </div>
  );
}