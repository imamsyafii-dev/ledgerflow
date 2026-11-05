export type Invoice = {
  id: string;
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