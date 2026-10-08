# SIAKAD Mini - UTS Praktikum Pemrograman Backend Lanjut

RESTful API akademik sederhana menggunakan Go, Fiber, PostgreSQL, JWT, bcrypt, RBAC, validation, pagination, filtering, search, dan transaction.

## Fitur
- 10 endpoint sesuai soal UTS.
- Role `admin` dan `mahasiswa`.
- Password bcrypt.
- JWT Bearer token.
- Soft delete mahasiswa.
- Pagination, filter prodi/angkatan, search NIM/nama, sort nama/IPK.
- Daftar mata kuliah dengan terisi dan sisa kuota.
- KRS dengan batas SKS berdasarkan IPK.
- Transaction + row locking saat mengambil mata kuliah.
- Error response JSON seragam.
- Request ID dan rate limit login sederhana.

## Struktur
```text
siakad-mini/
├── app/model
├── app/repository
├── app/service
├── config
├── database
├── helper
├── middleware
├── route
├── migrations
├── cmd/seed
└── main.go
```

## Menjalankan
1. Buat database PostgreSQL bernama `siakad_mini`.
2. Salin `.env.example` menjadi `.env` dan sesuaikan.
3. Jalankan migration:
   `psql -U postgres -d siakad_mini -f migrations/001_init.sql`
4. Jalankan seeder:
   `go run ./cmd/seed`
5. Jalankan API:
   `go run .`
6. Base URL: `http://localhost:3000`.

Admin seed:
- email: `admin@siakad.local`
- password: `admin12345`

Mahasiswa seed:
- email: `mhs01@siakad.local` sampai `mhs20@siakad.local`
- password awal masing-masing = NIM.

Contoh NIM pertama: `202600000001`.

## Endpoint
| Method | Endpoint | Akses |
|---|---|---|
| POST | /api/v1/auth/login | publik |
| GET | /api/v1/auth/me | semua role |
| GET | /api/v1/students | admin |
| POST | /api/v1/students | admin |
| GET | /api/v1/students/:id | admin / mahasiswa pemilik |
| PUT | /api/v1/students/:id | admin |
| DELETE | /api/v1/students/:id | admin |
| GET | /api/v1/courses | semua role |
| POST | /api/v1/enrollments | mahasiswa |
| DELETE | /api/v1/enrollments/:id | mahasiswa pemilik |

## Catatan
`204 No Content` sengaja tidak mengembalikan body, sesuai spesifikasi UTS.
