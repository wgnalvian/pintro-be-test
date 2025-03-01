# Golang Server

## Cara Menjalankan Aplikasi

1. **Clone repository**
   ```sh
   git clone https://github.com/username/repository.git
   cd repository
   ```
2. **Buat file `.env`** dan tambahkan konfigurasi berikut:
   ```env
   MIDTRANS_KEY=1234
   other...
   ```
3. **Jalankan aplikasi**
   ```sh
   go mod tidy
   go run main.go
   ```
4. **Akses server** 

---

## Arsitektur Aplikasi
Aplikasi ini menggunakan **Golang** dengan struktur berikut:
```
/repository
│── /config      # Konfigurasi database dan environment
│── /controller # Handler untuk request API
│── /entity     # Data Object
│── /route      # Routing API
│── /service    # Logika bisnis
│── /utils      # Util
│── /dto        # Data transfer object
│── main.go      # Entry point aplikasi
```
- **Gin** framework golang
- **Mongo** koneksion mongodb
- **Midtrans API** untuk pembayaran online

---

## Endpoint API

| Method | Endpoint          | Deskripsi                     |
|--------|------------------|-----------------------------|
| POST   | `/register`      | Registrasi user baru       |
| POST   | `/login`         | Autentikasi user           |
| GET    | `/transactions`  | Melihat daftar transaksi   |
| POST   | `/token`         | Generate token and Top up saldo               |
| POST   | `/transfer`      | Transfer saldo ke user lain|
| GET    | `/user`          | Get user data               |

---
