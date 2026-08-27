# 📚 Student Management API

RESTful API untuk pengelolaan data mahasiswa menggunakan **Go** dan **Fiber (v2)**.

- **Base URL**: `http://localhost:3000/api/v1`

---

## 📌 Daftar Endpoint

| Method | Endpoint | Query / Param | Status Code | Deskripsi |
| :--- | :--- | :--- | :--- | :--- |
| **GET** | `/health` | - | `200` | Cek status kesehatan API |
| **GET** | `/students` | `page`, `limit`, `search`, `sort`, `order`, `is_active` | `200`, `400` | Ambil daftar mahasiswa + paginasi |
| **GET** | `/students/:id` | `id` | `200`, `400`, `404` | Ambil detail mahasiswa |
| **POST** | `/students` | - | `201`, `400`, `409`, `415`, `422` | Tambah mahasiswa baru |
| **PUT** | `/students/:id` | `id` | `200`, `400`, `404`, `409`, `415`, `422` | Ganti seluruh data mahasiswa |
| **PATCH** | `/students/:id` | `id` | `200`, `400`, `404`, `409`, `415`, `422` | Perbarui sebagian data mahasiswa |
| **DELETE** | `/students/:id` | `id` | `204`, `400`, `404` | Hapus data mahasiswa |

---

## 📩 Format Respons API

Seluruh respons API dibungkus dalam format standar berikut:

### 1. Respons Sukses (*Success Response*)
```json
{
  "success": true,
  "message": "Pesan sukses",
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 3,
    "total_pages": 1
  }
}
```

### 2. Respons Error (*Error Response*)
```json
{
  "success": false,
  "message": "Pesan error",
  "errors": {
    "field": "Keterangan detail error"
  }
}
```

---

## ⚠️ Penjelasan HTTP Status Code

| Status Code | Tipe | Keterangan |
| :--- | :--- | :--- |
| `200 OK` | Sukses | Request berhasil diproses. |
| `201 Created` | Sukses | Data mahasiswa baru berhasil dibuat. |
| `204 No Content` | Sukses | Data mahasiswa berhasil dihapus. |
| `400 Bad Request` | Error | ID, query, atau body JSON tidak valid. |
| `404 Not Found` | Error | Data mahasiswa / endpoint tidak ditemukan. |
| `409 Conflict` | Error | NIM sudah digunakan oleh mahasiswa lain. |
| `415 Unsupported Media` | Error | Header `Content-Type` bukan `application/json`. |
| `422 Unprocessable` | Error | Validasi konten gagal (field kosong, grade di luar 0-100). |

---

## 🚀 Cara Menjalankan

```bash
# Jalankan server
go run .
```

---

## 🤖 Dokumentasi Bantuan AI

* **Sumber / Alat AI**: ChatGPT (OpenAI GPT-4)
* **Tujuan Penggunaan**: Membantu percepatan perancangan backend Go, pembuatan query SQL native PostgreSQL, implementasi Repository Pattern, penyiapan middleware Fiber v2, serta standardisasi penanganan respons dan error API.
* **Rincian Bantuan per Komponen Teknis**:
  * **1. Konfigurasi & Koneksi Database PostgreSQL**:
    * **Modul Env (`config/env.go`)**: Menyusun pembaca file `.env` dengan helper `GetEnv` dan `GetEnvInt` untuk membaca kredensial database (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`, `DB_MAX_CONNS`).
    * **Connection Pool (`database/postgres.go`)**: Mengintegrasikan driver `github.com/jackc/pgx/v5/pgxpool` via fungsi `NewPool()`. Mengatur konfigurasi pool (`MaxConns`, `MinConns=2`, `MaxConnLifetime=1h`, `MaxConnIdleTime=30m`) serta mekanismenya melakukan verifikasi awal dengan `pool.Ping()` (timeout 5 detik).
  * **2. Skema Migrasi SQL & Data Pengujian**:
    * **Skema DDL (`migration/001_create_students.sql`)**: Merancang skema tabel `students` (`id SERIAL PRIMARY KEY`, `nim VARCHAR(30)`, `name VARCHAR(100)`, `grade DOUBLE PRECISION`, `is_active BOOLEAN`, `created_at TIMESTAMPTZ`).
    * **Integrity Constraints**: Menambahkan check constraint `CHECK (grade >= 0 AND grade <= 100)` dan unique index case-insensitive `CREATE UNIQUE INDEX ON students (LOWER(nim))`.
    * **Data Seeding (`migration/002_seed_students.sql`)**: Menyiapkan query DML `INSERT INTO students` untuk mengisi data awal mahasiswa pengujian.
  * **3. Layer Akses Data (Repository Pattern - `repository/student_repository.go`)**:
    * **Model & Struct DTO**: Mendefinisikan struct `Student`, `CreateStudentRequest`, `ReplaceStudentRequest`, dan `PatchStudentRequest` (dengan field pointer `*string`/`*float64` untuk memperbolehkan partial update).
    * **Paginasi, Filtering, & Searching (`List`)**: Menyusun query dinamis berbasis parameter `$1, $2, ...` menggunakan pencarian case-insensitive `LOWER(name) LIKE LOWER($1)`, filter `is_active = $2`, klausa `ORDER BY` dinamis dengan whitelist kolom (`id`, `nim`, `name`, `grade`, `is_active`), serta klausa `LIMIT` & `OFFSET`.
    * **Operasi CRUD (`Create`, `FindByID`, `Replace`, `Patch`, `Delete`)**:
      * Memanfaatkan klausul `RETURNING id, nim, name, grade, is_active` pada query `INSERT` dan `UPDATE` untuk mengembalikan data hasil secara *atomic*.
      * Membangun dynamic query `UPDATE` pada metode `Patch` berbasis ketersediaan pointer non-nil.
      * Menggunakan `pgx.ErrNoRows` untuk mendeteksi data yang tidak ditemukan pada `FindByID`, `Replace`, dan `Patch`.
  * **4. Layer Handler & Routing API Go Fiber (`handler.go` & `main.go`)**:
    * **Routing Group**: Menyusun rute API v1 (`/api/v1/students`) dengan penanganan endpoint `GET /`, `GET /:id`, `POST /`, `PUT /:id`, `PATCH /:id`, `DELETE /:id`, serta health check `/health` (dengan timeout ping 2 detik).
    * **Middleware Integration**: Memasang middleware Fiber `requestid.New()`, `logger.New()`, dan `cors.New()`, serta kustom middleware `requireJSON` untuk validasi header `Content-Type: application/json`.
  * **5. Standardisasi Respons & Error Handling (`helper.go`)**:
    * **Struktur JSON Uniform**: Menyiapkan fungsi helper `sendSuccess()` (format `{ success: true, message, data, meta }`) dan `sendError()` (format `{ success: false, message, errors }`).
    * **HTTP Status Code Mapping**:
      * `200 OK`: Pengambilan/pembaruan data sukses.
      * `201 Created`: Pembuatan data mahasiswa baru sukses.
      * `204 No Content`: Penghapusan data mahasiswa sukses.
      * `400 Bad Request`: Parameter ID atau payload JSON tidak valid.
      * `404 Not Found`: Data mahasiswa atau endpoint tidak ditemukan.
      * `409 Conflict`: Poin penanganan duplikasi NIM terdeteksi.
      * `415 Unsupported Media Type`: Request body tanpa header `application/json`.
      * `422 Unprocessable Entity`: Kegagalan validasi konten (misal: field wajib kosong, grade di luar range 0–100).
      * `503 Service Unavailable`: Database tidak dapat dihubungi saat health check.
