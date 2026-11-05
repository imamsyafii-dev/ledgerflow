# Ledgerflow

Ledgerflow adalah platform lokal untuk membantu tim finance mengelola invoice, menerima notifikasi pembayaran, melakukan rekonsiliasi, dan mencatat transaksi ke *double-entry ledger*.

Project ini dibuat sebagai portfolio project dengan domain payment dan finance operations.

## Problem yang diselesaikan

Bisnis dapat menerima pembayaran dari banyak channel. Tanpa proses yang terstruktur, tim finance harus mencocokkan invoice, transaksi pembayaran, dan mutasi bank secara manual.

Alur inti Ledgerflow saat ini:

1. Finance membuat customer dan invoice.
2. Invoice diterbitkan.
3. Mock payment gateway mengirim webhook pembayaran berhasil.
4. Sistem menyimpan payment dan webhook event.
5. Sistem menandai invoice sebagai `paid`.
6. Sistem membuat reconciliation record.
7. Sistem membuat journal entry double-entry.

## Fitur saat ini

- Multi-tenant organization
- Authentication berbasis JWT dalam HttpOnly cookie
- Login, logout, dan endpoint current user (`GET /api/v1/auth/me`)
- Data role user: `finance_admin`, `approver`, dan `viewer`
- Daftar user organization untuk `finance_admin`
- Customer management
- Invoice creation dan invoice issuing
- Mock payment gateway webhook
- Webhook idempotency menggunakan event identifier
- Payment recording
- Automatic invoice reconciliation
- Double-entry ledger sederhana
- Dashboard finance
- Halaman customer, invoice, payment, reconciliation, dan ledger
- PostgreSQL melalui Docker Compose
- Go backend dengan Gin
- Next.js dashboard frontend

> Role disimpan di user dan JWT. Endpoint daftar user dibatasi untuk `finance_admin`; pembuatan user, perubahan role, dan perubahan status user masih masuk roadmap.

## Arsitektur

```text
Next.js Dashboard
       │
       │ HTTP + HttpOnly Cookie
       ▼
Go API (Gin)
 ├── PostgreSQL
 │   ├── organizations
 │   ├── users
 │   ├── customers
 │   ├── invoices
 │   ├── payments
 │   ├── webhook_events
 │   ├── reconciliations
 │   └── ledger tables
 │
 └── Mock Payment Gateway Webhook
```

## Tech stack

| Layer | Teknologi |
| --- | --- |
| Backend API | Go, Gin |
| Database | PostgreSQL 16 |
| Database driver | pgx |
| Authentication | JWT, HttpOnly Cookie, bcrypt |
| Frontend | Next.js 16, React 19, TypeScript |
| Local infrastructure | Docker Compose |

## Struktur repository

```text
ledgerflow/
├── backend/
│   ├── cmd/
│   │   ├── api/              # Entry point Go API
│   │   └── seed-admin/       # Command untuk membuat admin awal
│   ├── internal/
│   │   ├── auth/
│   │   ├── config/
│   │   ├── database/
│   │   ├── http/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   └── migrations/
├── frontend/
│   └── src/
│       ├── app/
│       ├── components/
│       └── lib/
├── docs/
│   └── ROADMAP.md
├── docker-compose.yml
└── .env.example
```

## Menjalankan project secara lokal

### Prerequisites

- Go
- Node.js dan npm
- Docker Desktop

### 1. Siapkan environment variable

Dari root repository:

```bash
cp .env.example .env
```

Untuk frontend:

```bash
cp frontend/.env.example frontend/.env.local
```

### 2. Jalankan PostgreSQL

```bash
docker compose up -d
docker compose ps
```

### 3. Jalankan database migration

Jalankan pada database baru, sesuai urutan nama file:

```bash
for file in backend/migrations/*.sql; do
  docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U ledgerflow -d ledgerflow < "$file"
done
```

> Migration saat ini belum memiliki migration runner yang mencatat file yang sudah dieksekusi. Jangan ulangi perintah ini pada database yang sudah memiliki migration tersebut.

### 4. Jalankan backend

```bash
cd backend
go run ./cmd/api
```

API berjalan pada `http://localhost:8080`.

```bash
curl http://localhost:8080/health
```

### 5. Buat organization dan admin awal

Buat organization melalui public bootstrap endpoint. Simpan nilai `data.id` dari response.

```bash
curl -X POST http://localhost:8080/api/v1/organizations \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ledgerflow Demo",
    "slug": "ledgerflow-demo"
  }'
```

Lalu, dari folder `backend`, jalankan:

```bash
export ADMIN_ORGANIZATION_ID="organization-id-di-sini"
export ADMIN_NAME="Finance Admin"
export ADMIN_EMAIL="admin@ledgerflow.local"
export ADMIN_PASSWORD="password-yang-kuat"

go run ./cmd/seed-admin
```

### 6. Jalankan frontend

Di terminal lain:

```bash
cd frontend
npm install
npm run dev
```

Buka `http://localhost:3000`.

## API utama

| Method | Endpoint | Keterangan |
| --- | --- | --- |
| POST | `/api/v1/auth/login` | Login dan membuat session cookie |
| POST | `/api/v1/auth/logout` | Menghapus session cookie |
| GET | `/api/v1/auth/me` | Mengambil user yang sedang login |
| POST | `/api/v1/organizations` | Bootstrap organization baru, tanpa JWT |
| GET | `/api/v1/users` | Menampilkan user organisasi; hanya `finance_admin` |
| POST | `/api/v1/customers` | Membuat customer |
| GET | `/api/v1/customers` | Menampilkan customer organisasi aktif |
| POST | `/api/v1/invoices` | Membuat invoice |
| GET | `/api/v1/invoices` | Menampilkan invoice organisasi aktif |
| POST | `/api/v1/invoices/:invoice_id/issue` | Menerbitkan invoice draft |
| GET | `/api/v1/payments` | Menampilkan payments |
| GET | `/api/v1/reconciliations` | Menampilkan hasil rekonsiliasi |
| GET | `/api/v1/ledger/journal-entries` | Menampilkan journal entries |
| POST | `/webhooks/mock-gateway` | Simulasi webhook payment gateway |

## Catatan keamanan saat ini

- Endpoint API utama dilindungi JWT middleware.
- `POST /api/v1/organizations` bersifat public untuk bootstrap lokal; sebelum production sebaiknya ditambahkan rate limiting atau dibuat sebagai alur registrasi yang membuat admin pertama sekaligus.
- Token disimpan dalam HttpOnly cookie agar tidak dapat dibaca JavaScript di browser.
- Data customer dan invoice dibatasi berdasarkan `organization_id` dari JWT.
- Mock webhook masih merupakan endpoint publik untuk kebutuhan simulasi lokal.
- Ganti `JWT_SECRET` sebelum memakai project di luar local development.

## Development commands

Backend:

```bash
cd backend
gofmt -w .
go test ./...
```

Frontend:

```bash
cd frontend
npm run build
npm run lint
```

## Roadmap

Lihat [docs/ROADMAP.md](docs/ROADMAP.md) untuk fitur pengembangan berikutnya.
