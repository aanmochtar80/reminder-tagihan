# Reminder Tagihan Web App

Aplikasi web modern untuk manajemen pelanggan dan pengiriman pengingat tagihan otomatis melalui WhatsApp, dibangun dengan Golang (Gin), SQLite, dan TailwindCSS.

## Fitur Utama
- **Dashboard SaaS-like**: Statistik tagihan bulanan, jumlah pelanggan.
- **WhatsApp Web Multi-Device**: Terintegrasi langsung dengan `whatsmeow` tanpa perlu API pihak ke-3.
- **Manajemen Pelanggan & Tagihan**: CRUD pelanggan dan invoice dengan perhitungan jatuh tempo otomatis.
- **Cron Scheduler**: Pengingat tagihan H-3, Hari H, dan telat bayar H+1 dikirim otomatis via WhatsApp.
- **Riwayat & Aktivitas**: Log aktivitas user dan riwayat pengiriman WA.
- **Sangat Ringan**: Menggunakan SQLite dan di-compile menjadi *single binary*.

---

## 🚀 Cara Install & Run (Local Development)

### Prasyarat
- Go 1.21+
- Node.js & npm (Hanya untuk build TailwindCSS, jika diperlukan)

### Langkah-langkah
1. Clone / buka direktori proyek ini.
2. Install dependensi Go:
   ```bash
   go mod tidy
   ```
3. Copy file `.env.example` menjadi `.env` dan sesuaikan nilainya:
   ```bash
   cp .env.example .env
   ```
4. Build Tailwind CSS (Opsional jika sudah ada `style.css` terbaru):
   ```bash
   npm install
   npm run build:css
   ```
5. Jalankan aplikasi:
   ```bash
   go run ./cmd/app/main.go
   ```
6. Buka browser: `http://localhost:8080`.
   - **Username**: `admin`
   - **Password**: `admin`

---

## 🛠️ Cara Build (Production Binary)
Untuk menjalankan aplikasi tanpa memerlukan Go compiler di target server:
```bash
# Build untuk Linux (jika VPS Anda Linux)
CGO_ENABLED=1 GOOS=linux go build -o reminder_tagihan ./cmd/app
```
Setelah di-build, Anda hanya perlu membawa file `reminder_tagihan`, folder `web`, dan file `.env` ke server.

---

## 📱 Cara Setup WhatsApp
1. Login ke aplikasi sebagai admin.
2. Buka menu **WhatsApp Gateway** dari sidebar.
3. Buka aplikasi WhatsApp di HP Anda > **Perangkat Tertaut** > **Tautkan Perangkat**.
4. Arahkan kamera HP ke layar monitor untuk scan **QR Code**.
5. Tunggu hingga status berubah menjadi **Connected & Logged In**.
6. Sistem sekarang siap mengirimkan pesan otomatis melalui Cron Jobs. *HP Anda tidak perlu terus menyala/online setelah terhubung.*

---

## 💾 Cara Backup Database
Aplikasi ini menggunakan SQLite, sehingga seluruh data tersimpan dalam 1 folder `database`.
1. Hentikan aplikasi sejenak.
2. Copy folder `database` (di dalamnya terdapat `app.db` dan `whatsapp.db`) ke tempat aman atau ke flashdisk/cloud.
3. **Restore**: Timpa folder `database` yang ada dengan folder backup Anda.

---

## ☁️ Cara Deploy ke VPS (Menggunakan Docker)

Cara termudah dan paling direkomendasikan untuk VPS (Ubuntu/Debian) adalah menggunakan **Docker Compose**.

1. Install Docker dan Docker Compose di VPS Anda.
2. Upload seluruh source code ini ke VPS Anda.
3. Masuk ke direktori proyek di VPS, jalankan:
   ```bash
   docker-compose up -d --build
   ```
4. Aplikasi akan berjalan di port `8080` dan berjalan di background (detach mode).
5. (Opsional) Set up Nginx sebagai Reverse Proxy untuk diarahkan ke `localhost:8080` dengan SSL Let's Encrypt.
