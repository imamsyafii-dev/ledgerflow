export type Invoice = {
  id: string;
  customer_name: string;
  number: string;
  amount: number;
  currency: string;
  status: "DRAFT" | "ISSUED" | "PAID" | "VOID";
  due_date: string;
};

export type Customer = {
  id: string;
  organization_id: string;
  name: string;
  email: string | null;
  created_at: string;
  updated_at: string;
};

export type CreateInvoiceInput = {
  customerID: string;
  number: string;
  amount: number;
  currency: string;
  dueDate: string;
};

export type JournalLine = {
  account_code: string;
  account_name: string;
  debit: number;
  credit: number;
  currency: string;
};

export type JournalEntry = {
  id: string;
  description: string;
  entry_date: string;
  lines: JournalLine[];
};

type ApiResponse<T> = {
  data: T;
};

const apiBaseURL = process.env.NEXT_PUBLIC_API_BASE_URL;
const organizationID = process.env.NEXT_PUBLIC_ORGANIZATION_ID;

function getConfiguration() {
  if (!apiBaseURL || !organizationID) {
    throw new Error(
      "NEXT_PUBLIC_API_BASE_URL dan NEXT_PUBLIC_ORGANIZATION_ID harus diisi",
    );
  }

  return {
    apiBaseURL,
    organizationID,
  };
}

async function getFromAPI<T>(path: string): Promise<T> {
  const { apiBaseURL } = getConfiguration();

  const response = await fetch(`${apiBaseURL}${path}`, {
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error(`API request failed with status ${response.status}`);
  }

  return response.json();
}

export async function getInvoices() {
  const { organizationID } = getConfiguration();

  const query = new URLSearchParams({
    organization_id: organizationID,
  });

  const response = await getFromAPI<ApiResponse<Invoice[]>>(
    `/api/v1/invoices?${query}`,
  );

  return response.data;
}

export async function createInvoice(input: CreateInvoiceInput) {
  const { apiBaseURL, organizationID } = getConfiguration();

  const response = await fetch(`${apiBaseURL}/api/v1/invoices`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      organization_id: organizationID,
      customer_id: input.customerID,
      number: input.number,
      amount: input.amount,
      currency: input.currency,
      due_date: input.dueDate,
    }),
  });

  if (!response.ok) {
    const payload = await response.json().catch(() => null);

    throw new Error(payload?.error ?? "Failed to create invoice");
  }

  const payload = await response.json() as ApiResponse<Invoice>;

  return payload.data;
}

export async function getCustomers() {
  const { organizationID } = getConfiguration();

  const query = new URLSearchParams({
    organization_id: organizationID,
  });

  const response = await getFromAPI<ApiResponse<Customer[]>>(
    `/api/v1/customers?${query}`,
  );

  return response.data;
}

export async function getJournalEntries() {
  const { organizationID } = getConfiguration();

  const query = new URLSearchParams({
    organization_id: organizationID,
  });

  const response = await getFromAPI<ApiResponse<JournalEntry[]>>(
    `/api/v1/ledger/journal-entries?${query}`,
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