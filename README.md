# 🥨 Keripik Pangsit API

Sistem manajemen penjualan Keripik Pangsit dengan fitur produk, transaksi, manajemen poin customer, dan penukaran poin.

## 🏗️ Arsitektur

```
keripik-pangsit/
├── config/           # Konfigurasi database (PostgreSQL & Redis)
├── controller/       # HTTP handler layer
├── helper/           # Helper function
├── repository/       # Data access layer (DB + Redis cache)
├── routes/           # Routing Gin
├── models/           # Domain models & request structs
├── migrations/       # SQL schema & seed data
├── main.go
├── .env.example
├── ERD.svg
└── keripik-pangsit.postman_collection.json
```

## ⚙️ Setup

### 1. Prasyarat
- Go 1.21+
- PostgreSQL 13+
- Redis 6+

### 2. Clone & Install
```bash
cp .env.example .env
# Sesuaikan konfigurasi di .env

go mod tidy
```

### 3. Buat Database & Jalankan Migrasi
```bash
createdb keripik_pangsit
psql -U postgres -d keripik_pangsit -f migrations/schema.sql
```

### 4. Jalankan Server
```bash
go run main.go
```

Server berjalan di `http://localhost:8080`

---

## 📡 API Endpoints

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/api/v1/health` | Health check |
| **Products** | | |
| POST | `/api/v1/products` | Tambah produk baru |
| GET | `/api/v1/products?tanggal=YYYY-MM-DD` | Lihat produk per tanggal |
| **Customers** | | |
| GET | `/api/v1/customers` | Lihat semua customer + poin |
| **Transactions** | | |
| POST | `/api/v1/transactions` | Buat transaksi baru |
| GET | `/api/v1/transactions?start_date=&end_date=` | Lihat transaksi per periode |
| GET | `/api/v1/transactions/summary?start_date=&end_date=` | Summary transaksi (income, bestseller, dll) |
| **Redeem Poin** | | |
| POST | `/api/v1/redeem` | Tukar poin dengan produk |
| GET | `/api/v1/redeem` | Lihat riwayat penukaran poin |

---

## 📦 Produk

| Ukuran | Harga |
|--------|-------|
| Small | Rp 10.000 |
| Medium | Rp 25.000 |
| Large | Rp 35.000 |

**Rasa:** Jagung Bakar, Rumput Laut, Original, Jagung Manis, Keju Asin, Keju Manis, Pedas

---

## 🎯 Sistem Poin

- **Earn:** Setiap kelipatan Rp 1.000 = 1 poin
- **Redeem:** Small = 200 poin | Medium = 300 poin | Large = 500 poin

---

## 🔧 Contoh Request

### Buat Transaksi
```json
POST /api/v1/transactions
{
  "nama_customer": "Budi",
  "product_id": "p1000000-0000-0000-0000-000000000001",
  "quantity": 2
}
```

### Tukar Poin
```json
POST /api/v1/redeem
{
  "nama_customer": "Budi",
  "ukuran": "Small",
  "product_id": "p1000000-0000-0000-0000-000000000001"
}
```

### Summary Transaksi
```
GET /api/v1/transactions/summary?start_date=2025-10-01&end_date=2025-12-31
```

---

## 🗄️ Teknologi

- **Language:** Go (Gin framework)
- **Database:** PostgreSQL (via sqlx + lib/pq)
- **Cache:** Redis (go-redis)
- **ID:** UUID (google/uuid)
