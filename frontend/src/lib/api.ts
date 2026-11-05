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

export type User = {
  id: string;
  organization_id: string;
  name: string;
  email: string;
  role: "finance_admin" | "approver" | "viewer";
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export type CreateUserInput = {
  name: string;
  email: string;
  password: string;
  role: User["role"];
};

export type CreateInvoiceInput = {
  customerID: string;
  number: string;
  amount: number;
  currency: string;
  dueDate: string;
};

export type CreateCustomerInput = {
  name: string;
  email: string;
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
  organization_id: string;
  reference_type: string;
  reference_id: string;
  description: string;
  entry_date: string;
  lines: JournalLine[];
};

export type Payment = {
  id: string;
  invoice_id: string;
  invoice_number: string;
  customer_name: string;
  provider: string;
  provider_payment_id: string;
  amount: number;
  currency: string;
  status: "PENDING" | "SUCCEEDED" | "FAILED" | "REFUNDED";
  paid_at: string | null;
  created_at: string;
};

export type Reconciliation = {
  id: string;
  invoice_id: string;
  invoice_number: string;
  customer_name: string;
  payment_id: string;
  provider_payment_id: string;
  amount: number;
  status: "MATCHED" | "MANUAL_REVIEW" | "REVERSED";
  reconciled_at: string;
};

export type RefundRequest = { 
  id: string; 
  payment_id: 
  string; invoice_number: 
  string; customer_name: 
  string; amount: number; 
  currency: string; 
  requested_by: string; 
  requested_by_name: 
  string; reason: 
  string; status: "PENDING" | "APPROVED" | "REJECTED"; 
  decided_by: string | null; 
  decision_note: string | null; 
  decided_at: string | null; 
  executed_by: string | null;
  executed_at: string | null;
  provider_refund_id: string | null;
  completed_at: string | null;
  failure_reason: string | null;
  created_at: string 
};
export async function decideRefund(
  id: string, 
  status: "APPROVED" | "REJECTED", 
  note: string
) { 
  const response = await fetch(`${getAPIBaseURL()}/api/v1/refund-requests/${id}/decision`, { method: "PATCH", credentials: "include", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ status, note }) }); 
  if (!response.ok) { 
    const p = await response.json().catch(() => null); 
    throw new Error(p?.error ?? "Failed to decide refund") 
  } 
  return response.json() 
}
export async function executeRefund(refundID: string) {
  const response = await fetch(
    `${getAPIBaseURL()}/api/v1/refund-requests/${refundID}/execute`,
    {
      method: "POST",
      credentials: "include",
    },
  );

  if (!response.ok) {
    const payload = await response.json().catch(() => null);

    throw new Error(payload?.error ?? "Failed to execute refund");
  }

  return response.json();
}
export async function requestRefund(
  paymentID: string, 
  reason: string
) { 
  const r = await fetch(`${getAPIBaseURL()}/api/v1/payments/${paymentID}/refund-requests`, { method: "POST", credentials: "include", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ reason }) }); 
  if (!r.ok) { 
    const p = await r.json().catch(() => null); 
    throw new Error(p?.error ?? "Failed to request refund") 
  } 
  return r.json() 
}

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

function getAPIBaseURL() {
  if (!apiBaseURL) {
    throw new Error("NEXT_PUBLIC_API_BASE_URL harus diisi");
  }

  return apiBaseURL;
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
    credentials: "include",
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

export async function issueInvoice(invoiceID: string) {
  const { apiBaseURL, organizationID } = getConfiguration();

  const response = await fetch(
    `${apiBaseURL}/api/v1/invoices/${invoiceID}/issue`,
    {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        organization_id: organizationID,
      }),
    },
  );

  if (!response.ok) {
    const payload = await response.json().catch(() => null);

    throw new Error(payload?.error ?? "Failed to issue invoice");
  }

  const payload = await response.json() as ApiResponse<Invoice>;

  return payload.data;
}

export type SimulatePaymentInput = {
  invoiceID: string;
  amount: number;
  currency: string;
};

export async function simulateMockPayment(input: SimulatePaymentInput) {
  const { apiBaseURL, organizationID } = getConfiguration();
  const uniqueID = crypto.randomUUID();

  const response = await fetch(`${apiBaseURL}/webhooks/mock-gateway`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      event_id: `evt_mock_${uniqueID}`,
      event_type: "payment.succeeded",
      payment_id: `pay_mock_${uniqueID}`,
      organization_id: organizationID,
      invoice_id: input.invoiceID,
      amount: input.amount,
      currency: input.currency,
    }),
  });

  if (!response.ok) {
    const payload = await response.json().catch(() => null);

    throw new Error(payload?.error ?? "Failed to simulate payment");
  }

  return response.json();
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

export async function createCustomer(input: CreateCustomerInput) {
  const { apiBaseURL, organizationID } = getConfiguration();

  const response = await fetch(`${apiBaseURL}/api/v1/customers`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      organization_id: organizationID,
      name: input.name,
      email: input.email,
    }),
  });

  if (!response.ok) {
    const payload = await response.json().catch(() => null);

    throw new Error(payload?.error ?? "Failed to create customer");
  }

  const payload = await response.json() as ApiResponse<Customer>;

  return payload.data;
}

export async function createUser(input: CreateUserInput) {
  const response = await fetch(`${getAPIBaseURL()}/api/v1/users`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      name: input.name,
      email: input.email,
      password: input.password,
      role: input.role,
    }),
  });

  if (!response.ok) {
    const payload = await response.json().catch(() => null);
    throw new Error(payload?.error ?? "Failed to create user");
  }

  const payload = (await response.json()) as ApiResponse<User>;
  return payload.data;
}

export async function updateUserRole(userID: string, role: User["role"]) {
  const response = await fetch(
    `${getAPIBaseURL()}/api/v1/users/${userID}/role`,
    {
      method: "PATCH",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ role }),
    },
  );

  if (!response.ok) {
    const payload = await response.json().catch(() => null);
    throw new Error(payload?.error ?? "Failed to update user role");
  }

  const payload = (await response.json()) as ApiResponse<User>;
  return payload.data;
}

export async function updateUserStatus(userID: string, isActive: boolean) {
  const response = await fetch(
    `${getAPIBaseURL()}/api/v1/users/${userID}/status`,
    {
      method: "PATCH",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ is_active: isActive }),
    },
  );

  if (!response.ok) {
    const payload = await response.json().catch(() => null);
    throw new Error(payload?.error ?? "Failed to update user status");
  }

  const payload = (await response.json()) as ApiResponse<User>;
  return payload.data;
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

export async function getPayments() {
  const { organizationID } = getConfiguration();

  const query = new URLSearchParams({
    organization_id: organizationID,
  });

  const response = await getFromAPI<ApiResponse<Payment[]>>(
    `/api/v1/payments?${query}`,
  );

  return response.data;
}

export async function getReconciliations() {
  const { organizationID } = getConfiguration();

  const query = new URLSearchParams({
    organization_id: organizationID,
  });

  const response = await getFromAPI<ApiResponse<Reconciliation[]>>(
    `/api/v1/reconciliations?${query}`,
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
