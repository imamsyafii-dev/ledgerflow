# Ledgerflow Roadmap

Dokumen ini berisi fitur yang belum tersedia dan arah pengembangan Ledgerflow setelah MVP saat ini.

## Status MVP saat ini

Sudah tersedia:

- Authentication dasar: login, logout, dan current user
- Data role: finance admin, approver, dan viewer
- Customer management
- Invoice management dan issue invoice
- Mock payment webhook
- Payment recording
- Reconciliation record
- Double-entry ledger dasar
- Dashboard finance
- Halaman customer, invoice, payment, reconciliation, dan ledger

## Prioritas pengembangan

### 1. User Management dan Role Authorization

Tujuan: finance admin dapat mengelola user dalam organisasinya dan sistem membatasi akses sesuai role.

- List users per organization
- Create user
- Activate dan deactivate user
- Ubah role user
- Halaman `/users`
- Hanya `finance_admin` yang dapat mengelola user
- `viewer` hanya dapat melihat data
- `approver` memiliki akses approval refund

### 2. Audit Log

Tujuan: setiap aksi penting dapat ditelusuri oleh finance team.

Aksi yang dicatat:

- User login
- Customer dibuat atau diubah
- Invoice dibuat dan diterbitkan
- Payment webhook diterima
- Payment direkonsiliasi
- Refund diajukan, disetujui, atau ditolak
- User dan role diubah

Data audit log:

- Actor atau user
- Organization
- Action
- Resource type
- Resource ID
- Request metadata
- Timestamp

### 3. Refund Approval Workflow

Tujuan: refund tidak dapat dilakukan tanpa approval.

```text
Finance Admin membuat request refund
        │
        ▼
Status: pending_approval
        │
        ▼
Approver menyetujui atau menolak
        │
        ▼
Status: approved / rejected
        │
        ▼
Payment atau refund record diperbarui
```

- Refund request
- Approval dan rejection
- Approval notes
- Role-based authorization
- Audit log untuk seluruh perubahan status

### 4. Payment Link dan Public Invoice Page

Tujuan: invoice dapat dibagikan kepada customer.

- Public invoice URL
- Invoice detail tanpa login
- Tombol simulasi pembayaran
- Expiry date untuk payment link
- Payment status pada halaman public
- Invoice tidak dapat dibayar ulang setelah `paid`

### 5. Background Jobs

Tujuan: proses asynchronous tidak bergantung pada request API.

- Reminder invoice jatuh tempo
- Retry payment processing yang gagal
- Retry webhook processing
- Scheduled overdue invoice check
- Worker terpisah dari API server

Tahap implementasi:

1. PostgreSQL-based job table
2. Worker Go sederhana
3. Redis queue
4. Retry policy dan dead-letter queue

### 6. CSV Bank Statement Import

Tujuan: melakukan rekonsiliasi pembayaran dari mutasi rekening bank.

- Upload CSV mutasi bank
- Validasi struktur file
- Preview data sebelum import
- Penyimpanan bank transaction
- Match transaction dengan payment atau invoice
- Manual reconciliation jika tidak ada kecocokan

### 7. Duplicate Transaction Detection

Tujuan: mencegah pencatatan pembayaran ganda.

Deteksi berdasarkan:

- External transaction ID
- Nominal pembayaran
- Reference number
- Waktu transaksi
- Customer atau invoice terkait

Status yang mungkin:

- Matched
- Unmatched
- Duplicate suspected
- Manually resolved

### 8. Finance Reports dan Export

- Invoice aging report
- Outstanding invoice report
- Payment report
- Reconciliation report
- General ledger report
- Export CSV
- Export PDF

### 9. API Documentation

- OpenAPI specification
- Swagger UI
- Contoh request dan response
- Dokumentasi authentication
- Dokumentasi webhook idempotency
- Postman collection

### 10. Testing

Backend:

- Unit test service layer
- Repository integration test
- Authentication middleware test
- Webhook idempotency test
- Reconciliation flow integration test
- Ledger balancing test

Frontend:

- Component test
- Form validation test
- Auth redirect test
- End-to-end flow: login → invoice → payment → reconciliation

### 11. Observability

- Structured JSON logging
- Request ID
- Error tracking
- Metrics endpoint
- API latency metrics
- Background job metrics
- Health check dan readiness check

### 12. Reliable Event Delivery

Tujuan: memastikan event bisnis tidak hilang ketika database atau worker mengalami kegagalan.

- Outbox table
- Transactional outbox pattern
- Event publisher worker
- Retry event delivery
- Idempotent event consumers

Contoh event:

- `invoice.issued`
- `payment.succeeded`
- `payment.reconciled`
- `refund.approved`

### 13. Deployment dan CI/CD

- Dockerfile backend
- Dockerfile frontend
- Production Docker Compose
- GitHub Actions untuk test dan build
- Environment-specific configuration
- Database migration pipeline
- Deployment guide

## Target arsitektur

```text
Next.js Dashboard
       │
       ▼
Go API
 ├── PostgreSQL
 ├── Redis
 ├── Worker
 ├── Outbox Publisher
 ├── Mock Payment Gateway
 └── Observability Stack
```
