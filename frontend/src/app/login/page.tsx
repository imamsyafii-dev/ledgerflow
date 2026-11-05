import { LoginForm } from "@/components/login-form";

export default function LoginPage() {
  return (
    <main className="login-page">
      <section className="login-card">
        <p className="eyebrow">Ledgerflow</p>
        <h1>Welcome back</h1>
        <p className="subtitle">
          Sign in to access your finance operations dashboard.
        </p>

        <LoginForm />
      </section>
    </main>
  );
}