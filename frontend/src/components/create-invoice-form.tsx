"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import {
  createInvoice,
  type Customer,
} from "@/lib/api";

type CreateInvoiceFormProps = {
  customers: Customer[];
};

export function CreateInvoiceForm({
  customers,
}: CreateInvoiceFormProps) {
  const router = useRouter();

  const [customerID, setCustomerID] = useState("");
  const [number, setNumber] = useState("");
  const [amount, setAmount] = useState("");
  const [currency, setCurrency] = useState("IDR");
  const [dueDate, setDueDate] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setError("");
    setIsSubmitting(true);

    try {
      await createInvoice({
        customerID,
        number,
        amount: Number(amount),
        currency,
        dueDate,
      });

      setCustomerID("");
      setNumber("");
      setAmount("");
      setCurrency("IDR");
      setDueDate("");

      router.refresh();
    } catch (error) {
      setError(
        error instanceof Error
          ? error.message
          : "Failed to create invoice",
      );
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form className="invoice-form" onSubmit={handleSubmit}>
      <div className="form-field">
        <label htmlFor="customer">Customer</label>
        <select
          id="customer"
          value={customerID}
          onChange={(event) => setCustomerID(event.target.value)}
          required
        >
          <option value="">Select customer</option>
          {customers.map((customer) => (
            <option key={customer.id} value={customer.id}>
              {customer.name}
            </option>
          ))}
        </select>
      </div>

      <div className="form-field">
        <label htmlFor="number">Invoice Number</label>
        <input
          id="number"
          value={number}
          onChange={(event) => setNumber(event.target.value)}
          placeholder="INV-2026-0002"
          required
        />
      </div>

      <div className="form-field">
        <label htmlFor="amount">Amount</label>
        <input
          id="amount"
          type="number"
          min="1"
          value={amount}
          onChange={(event) => setAmount(event.target.value)}
          placeholder="150000"
          required
        />
      </div>

      <div className="form-field">
        <label htmlFor="currency">Currency</label>
        <select
          id="currency"
          value={currency}
          onChange={(event) => setCurrency(event.target.value)}
          required
        >
          <option value="IDR">IDR</option>
        </select>
      </div>

      <div className="form-field">
        <label htmlFor="due-date">Due Date</label>
        <input
          id="due-date"
          type="date"
          value={dueDate}
          onChange={(event) => setDueDate(event.target.value)}
          required
        />
      </div>

      <div className="form-actions">
        <button
          className="primary-button"
          type="submit"
          disabled={isSubmitting || customers.length === 0}
        >
          {isSubmitting ? "Creating..." : "Create Invoice"}
        </button>

        {customers.length === 0 && (
          <span className="form-hint">Create a customer first.</span>
        )}
      </div>

      {error && <p className="form-error">{error}</p>}
    </form>
  );
}