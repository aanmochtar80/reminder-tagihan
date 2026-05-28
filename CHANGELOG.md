# Changelog

Semua perubahan dan pencapaian pada proyek **Reminder Tagihan** akan didokumentasikan di file ini.

## [28 Mei 2026] - Pembaruan Fitur Besar & Perbaikan Krusial

### 🚀 Fitur Baru (Features)
- **Manajemen Tagihan (Invoices)**: Menambahkan fungsionalitas Edit dan Delete (Hapus) pada tagihan dengan konfirmasi berlapis.
- **Manajemen Pelanggan (Customers)**: Menambahkan fungsionalitas Edit dan Delete (Hapus) secara penuh. Penghapusan pelanggan kini akan secara otomatis ikut menghapus seluruh tagihan yang berelasi dengannya.
- **Keamanan (Security)**: Menambahkan fitur dan formulir khusus bagi Admin untuk mengubah *password* dengan aman (terenkripsi *Bcrypt*) melalui halaman Pengaturan.
- **Template WhatsApp**: Menambahkan template informasi rekening pembayaran (BNI, BCA, Jago, DANA, OVO).
- **Template WhatsApp**: Mendesain ulang struktur pesan pengingat tagihan dengan tambahan *emoji* dan format tebal (*bold*) agar lebih menarik dan interaktif, baik di tingkat sistem maupun UI.
- **Lokalisasi**: Menyesuaikan setelan *timezone* pada `docker-compose.yml` menjadi `Asia/Makassar` (WITA) agar sinkron dengan domisili server klien.

### 🎨 Pembaruan Tampilan (UI/UX)
- **Aksi Modern**: Merombak total tombol "Aksi" pada tabel Tagihan dan Pelanggan dari tombol blok teks kaku menjadi kumpulan Ikon SVG interaktif (*Heroicons*) yang minimalis dan elegan.
- **Pop-up Modal**: Mengintegrasikan sistem *Modal* (layar tumpang tindih) yang cantik dan responsif berbasis Alpine.js untuk fitur Edit Pelanggan dan Edit Tagihan, membuang transisi halaman tradisional.
- **Format Mata Uang**: Memastikan seluruh angka di aplikasi (Dashboard, Tabel, dan Input) terpisahkan oleh format ribuan Rupiah secara rapi.

### 🐛 Perbaikan Kutu (Bug Fixes)
- **Crash Format Uang**: Memperbaiki kelemahan pada fungsi logika `formatRupiah` yang sempat menyebabkan *panic/crash* saat bertemu data teks kosong atau angka 0. Kini mampu menangani multi-tipe (string, int, float) secara absolut.
- **Dashboard Statis**: Memperbaiki angka *Dashboard* yang sebelumnya hanya *placeholder* 0. Kini seluruh angka (Tagihan Bulan Ini, Jatuh Tempo Hari Ini, Telat Bayar) dihitung 100% dinamis dari kumpulan data asli di *database*.
- **Isu Zona Waktu (Timezone Miscalculation)**: Memperbaiki logika perbandingan tanggal di SQLite yang menyebabkan tagihan jatuh tempo terlempar ke daftar *Overdue* lebih cepat dari seharusnya. Diselesaikan dengan menyamakan parameter *query* menggunakan waktu absolut (UTC).
- **Polling WhatsApp Gateway**: Memperbaiki *bug* antarmuka di mana kode QR akan *stuck/loading* abadi sesaat setelah pengguna melakukan *logout*. Diselesaikan dengan memicu fungsi inisialisasi ulang (*Re-initialize Client*) tepat setelah pemutusan sesi.
