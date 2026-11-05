"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { createCustomer } from "@/lib/api";

export function CreateCustomerForm() {
  const router = useRouter();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setError("");
    setIsSubmitting(true);

    try {
      await createCustomer({
        name,
        email,
      });

      setName("");
      setEmail("");

      router.refresh();
    } catch (error) {
      setError(
        error instanceof Error
          ? error.message
          : "Failed to create customer",
      );
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form className="customer-form" onSubmit={handleSubmit}>
      <div className="form-field">
        <label htmlFor="customer-name">Customer Name</label>
        <input
          id="customer-name"
          value={name}
          onChange={(event) => setName(event.target.value)}
          placeholder="Budi Santoso"
          required
        />
      </div>

      <div className="form-field">
        <label htmlFor="customer-email">Email</label>
        <input
          id="customer-email"
          type="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          placeholder="budi@example.com"
        />
      </div>

      <div className="form-actions">
        <button
          className="primary-button"
          type="submit"
          disabled={isSubmitting}
        >
          {isSubmitting ? "Creating..." : "Create Customer"}
        </button>
      </div>

      {error && <p className="form-error">{error}</p>}
    </form>
  );
}