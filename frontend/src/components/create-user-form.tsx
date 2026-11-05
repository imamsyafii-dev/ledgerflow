"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

import { createUser, type User } from "@/lib/api";

export function CreateUserForm() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<User["role"]>("viewer");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setIsSubmitting(true);

    try {
      await createUser({ name, email, password, role });
      setName("");
      setEmail("");
      setPassword("");
      setRole("viewer");
      router.refresh();
    } catch (error) {
      setError(error instanceof Error ? error.message : "Failed to create user");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form className="user-form" onSubmit={handleSubmit}>
      <div className="form-field">
        <label htmlFor="user-name">Name</label>
        <input
          id="user-name"
          value={name}
          onChange={(event) => setName(event.target.value)}
          placeholder="Finance Team Member"
          required
        />
      </div>

      <div className="form-field">
        <label htmlFor="user-email">Email</label>
        <input
          id="user-email"
          type="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          placeholder="member@ledgerflow.local"
          required
        />
      </div>

      <div className="form-field">
        <label htmlFor="user-password">Temporary Password</label>
        <input
          id="user-password"
          type="password"
          minLength={8}
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          placeholder="At least 8 characters"
          required
        />
      </div>

      <div className="form-field">
        <label htmlFor="user-role">Role</label>
        <select
          id="user-role"
          value={role}
          onChange={(event) => setRole(event.target.value as User["role"])}
        >
          <option value="viewer">Viewer</option>
          <option value="approver">Approver</option>
          <option value="finance_admin">Finance Admin</option>
        </select>
      </div>

      <div className="form-actions">
        <button className="primary-button" type="submit" disabled={isSubmitting}>
          {isSubmitting ? "Creating..." : "Create User"}
        </button>
      </div>

      {error && <p className="form-error">{error}</p>}
    </form>
  );
}
