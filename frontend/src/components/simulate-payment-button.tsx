"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { simulateMockPayment } from "@/lib/api";

type SimulatePaymentButtonProps = {
  invoiceID: string;
  amount: number;
  currency: string;
};

export function SimulatePaymentButton({
  invoiceID,
  amount,
  currency,
}: SimulatePaymentButtonProps) {
  const router = useRouter();

  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handlePayment() {
    setError("");
    setIsSubmitting(true);

    try {
      await simulateMockPayment({
        invoiceID,
        amount,
        currency,
      });

      router.refresh();
    } catch (error) {
      setError(
        error instanceof Error ? error.message : "Failed to simulate payment",
      );
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div className="invoice-action">
      <button
        className="payment-button"
        type="button"
        onClick={handlePayment}
        disabled={isSubmitting}
      >
        {isSubmitting ? "Processing..." : "Simulate Payment"}
      </button>

      {error && <span className="action-error">{error}</span>}
    </div>
  );
}