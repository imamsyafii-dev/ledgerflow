import Link from "next/link";
import { redirect } from "next/navigation";

import { getCurrentUser, getUsers } from "@/lib/server-api";

export const instant = false;

function formatDate(value: string) {
  return new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
  }).format(new Date(value));
}

function formatRole(role: string) {
  return role
    .split("_")
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

export default async function UsersPage() {
  const currentUser = await getCurrentUser();

  if (currentUser.role !== "finance_admin") {
    redirect("/");
  }

  const users = await getUsers();

  return (
    <main className="dashboard">
      <section className="hero">
        <div>
          <p className="eyebrow">Organization settings</p>
          <h1>Users</h1>
          <p className="subtitle">
            Daftar user yang memiliki akses ke organization Anda.
          </p>
        </div>

        <Link className="back-link" href="/">
          ← Dashboard
        </Link>
      </section>

      <section className="panel">
        <div className="panel-header">
          <div>
            <p className="eyebrow">Access control</p>
            <h2>Organization Users</h2>
          </div>
        </div>

        {users.length === 0 ? (
          <p className="empty-state">Belum ada user.</p>
        ) : (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Status</th>
                  <th>Created At</th>
                </tr>
              </thead>
              <tbody>
                {users.map((user) => (
                  <tr key={user.id}>
                    <td>{user.name}</td>
                    <td>{user.email}</td>
                    <td>{formatRole(user.role)}</td>
                    <td>
                      <span
                        className={
                          user.is_active ? "badge badge-paid" : "badge badge-void"
                        }
                      >
                        {user.is_active ? "Active" : "Inactive"}
                      </span>
                    </td>
                    <td>{formatDate(user.created_at)}</td>
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
