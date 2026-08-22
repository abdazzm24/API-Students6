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
