# Bank Sampah API

RESTful API pengelolaan bank sampah berbasis **Golang (Gin + GORM)** yang di-deploy di **Google Cloud Platform (Cloud Run + Cloud SQL)**.

## Stack Teknologi

| Komponen | Teknologi |
|---|---|
| Bahasa | Go 1.21 |
| Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL (Cloud SQL) |
| Auth | JWT (golang-jwt/jwt v5) |
| Deploy | GCP Cloud Run |
| Container | Docker |

## Struktur Folder

```
banksampah-api/
├── cmd/
│   └── api/
│       └── main.go              ← Entry point, jalankan dari sini
├── internal/
│   ├── config/
│   │   └── config.go            ← Load env, koneksi DB, setup router
│   ├── middleware/
│   │   └── jwt.go               ← AuthRequired, AdminOnly, OperatorOrAdmin
│   ├── models/
│   │   └── models.go            ← Semua struct + AutoMigrate
│   └── handlers/
│       ├── auth.go              ← Register + Login
│       ├── transaksi.go         ← CRUD transaksi + kalkulasi otomatis
│       └── handlers.go          ← Kecamatan, BSU, Nasabah, Sampah, Laporan
├── Dockerfile                   ← Build image untuk Cloud Run
├── .env.example                 ← Template environment variable
└── go.mod
```

## Cara Menjalankan Lokal

```bash
# 1. Clone repo
git clone https://github.com/yourusername/banksampah-api
cd banksampah-api

# 2. Buat file .env dari template
cp .env.example .env
# Edit .env sesuai konfigurasi PostgreSQL lokal kamu

# 3. Download dependency
go mod tidy

# 4. Jalankan
go run cmd/api/main.go

# Server jalan di http://localhost:8080
```

## Daftar Endpoint

### Auth (Public)
| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | /api/v1/auth/register | Daftar nasabah baru |
| POST | /api/v1/auth/login | Login, dapat token JWT |

### Kecamatan
| Method | Endpoint | Role |
|---|---|---|
| GET | /api/v1/kecamatan | Semua |
| POST | /api/v1/kecamatan | Admin |
| PUT | /api/v1/kecamatan/:id | Admin |
| DELETE | /api/v1/kecamatan/:id | Admin |

### Transaksi
| Method | Endpoint | Role |
|---|---|---|
| GET | /api/v1/transaksi/saya | Nasabah (transaksi sendiri) |
| POST | /api/v1/transaksi | Operator, Admin |
| GET | /api/v1/transaksi | Operator, Admin |
| PUT | /api/v1/transaksi/:id | Operator, Admin |
| DELETE | /api/v1/transaksi/:id | Operator, Admin |

### Laporan
| Method | Endpoint | Query Params | Role |
|---|---|---|---|
| GET | /api/v1/laporan | ?kecamatan_id=&from=&to= | Admin |

## Deploy ke GCP Cloud Run

```bash
# 1. Build & push image ke Artifact Registry
gcloud builds submit --tag gcr.io/PROJECT_ID/banksampah-api

# 2. Deploy ke Cloud Run
gcloud run deploy banksampah-api \
  --image gcr.io/PROJECT_ID/banksampah-api \
  --platform managed \
  --region asia-southeast2 \
  --set-env-vars DB_HOST=...,JWT_SECRET=...
```

## Contoh Request (Postman)

### Login
```json
POST /api/v1/auth/login
{
  "email": "budi@email.com",
  "password": "password123"
}
```

### Catat Transaksi
```json
POST /api/v1/transaksi
Authorization: Bearer <token_operator>

{
  "id_nasabah": 5,
  "id_sub_kategori": 1,
  "berat_kg": 2.5
}
```

Response (nilai_rupiah dihitung otomatis):
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "id_nasabah": 5,
    "id_sub_kategori": 1,
    "berat_kg": 2.5,
    "nilai_rupiah": 7500,
    "tanggal": "2025-05-24T..."
  }
}
```
