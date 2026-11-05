"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { updateUserRole, updateUserStatus, type User } from "@/lib/api";

type UserActionsProps = {
  user: User;
  isCurrentUser: boolean;
};

export function UserActions({ user, isCurrentUser }: UserActionsProps) {
  const router = useRouter();
  const [role, setRole] = useState<User["role"]>(user.role);
  const [error, setError] = useState("");
  const [isSavingRole, setIsSavingRole] = useState(false);
  const [isSavingStatus, setIsSavingStatus] = useState(false);

  if (isCurrentUser) {
    return <span className="action-hint">Current user</span>;
  }

  async function handleRoleUpdate() {
    setError("");
    setIsSavingRole(true);

    try {
      await updateUserRole(user.id, role);
      router.refresh();
    } catch (error) {
      setError(error instanceof Error ? error.message : "Failed to update role");
    } finally {
      setIsSavingRole(false);
    }
  }

  async function handleStatusUpdate() {
    setError("");
    setIsSavingStatus(true);

    try {
      await updateUserStatus(user.id, !user.is_active);
      router.refresh();
    } catch (error) {
      setError(error instanceof Error ? error.message : "Failed to update status");
    } finally {
      setIsSavingStatus(false);
    }
  }

  return (
    <div className="user-actions">
      <div className="user-role-action">
        <select
          aria-label={`Role for ${user.name}`}
          value={role}
          onChange={(event) => setRole(event.target.value as User["role"])}
          disabled={isSavingRole || isSavingStatus}
        >
          <option value="viewer">Viewer</option>
          <option value="approver">Approver</option>
          <option value="finance_admin">Finance Admin</option>
        </select>
        <button
          className="secondary-button"
          type="button"
          onClick={handleRoleUpdate}
          disabled={role === user.role || isSavingRole || isSavingStatus}
        >
          {isSavingRole ? "Saving..." : "Save role"}
        </button>
      </div>

      <button
        className={user.is_active ? "danger-button" : "secondary-button"}
        type="button"
        onClick={handleStatusUpdate}
        disabled={isSavingRole || isSavingStatus}
      >
        {isSavingStatus
          ? "Saving..."
          : user.is_active
            ? "Deactivate"
            : "Activate"}
      </button>

      {error && <p className="action-error">{error}</p>}
    </div>
  );
}
