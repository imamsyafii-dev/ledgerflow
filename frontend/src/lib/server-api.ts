import "server-only";

import { cookies } from "next/headers";
import type {
  Customer,
  Invoice,
  JournalEntry,
  Payment,
  Reconciliation,
} from "@/lib/api";

type ApiResponse<T> = {
  data: T;
};

function getAPIBaseURL() {
  const apiBaseURL = process.env.API_BASE_URL;

  if (!apiBaseURL) {
    throw new Error("API_BASE_URL harus diisi");
  }

  return apiBaseURL;
}

async function getFromAPI<T>(path: string): Promise<T> {
  const cookieStore = await cookies();

  const response = await fetch(`${getAPIBaseURL()}${path}`, {
    cache: "no-store",
    headers: {
      Cookie: cookieStore.toString(),
    },
  });

  if (!response.ok) {
    throw new Error(`API request failed with status ${response.status}`);
  }

  return response.json();
}

export async function getInvoices() {
  const response = await getFromAPI<ApiResponse<Invoice[]>>(
    "/api/v1/invoices",
  );

  return response.data;
}

export async function getCustomers() {
  const response = await getFromAPI<ApiResponse<Customer[]>>(
    "/api/v1/customers",
  );

  return response.data;
}

export async function getPayments() {
  const response = await getFromAPI<ApiResponse<Payment[]>>(
    "/api/v1/payments",
  );

  return response.data;
}

export async function getReconciliations() {
  const response = await getFromAPI<ApiResponse<Reconciliation[]>>(
    "/api/v1/reconciliations",
  );

  return response.data;
}

export async function getJournalEntries() {
  const response = await getFromAPI<ApiResponse<JournalEntry[]>>(
    "/api/v1/ledger/journal-entries",
  );

  return response.data;
}

export async function getDashboardData() {
  const [invoices, journalEntries] = await Promise.all([
    getInvoices(),
    getJournalEntries(),
  ]);

  return {
    invoices,
    journalEntries,
  };
}