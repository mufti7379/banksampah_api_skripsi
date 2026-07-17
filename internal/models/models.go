package models

import (
	"time"

	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// WILAYAH
// ─────────────────────────────────────────────

type Kabupaten struct {
	ID             uint   `json:"id"              gorm:"primaryKey"`
	NamaKabupaten  string `json:"nama_kabupaten"  gorm:"not null"`
}

type Kecamatan struct {
	ID             uint      `json:"id"              gorm:"primaryKey"`
	IDKabupaten    uint      `json:"id_kabupaten"    gorm:"not null"`
	Kabupaten      Kabupaten `json:"kabupaten"       gorm:"foreignKey:IDKabupaten"`
	NamaKecamatan  string    `json:"nama_kecamatan"  gorm:"not null"`
	CreatedAt      time.Time `json:"created_at"`
}

// ─────────────────────────────────────────────
// BSU
// ─────────────────────────────────────────────

type BSU struct {
	ID              uint      `json:"id"               gorm:"primaryKey"`
	IDKecamatan     uint      `json:"id_kecamatan"     gorm:"not null"`
	Kecamatan       Kecamatan `json:"kecamatan"        gorm:"foreignKey:IDKecamatan"`
	NamaBSU         string    `json:"nama_bsu"         gorm:"not null"`
	Alamat          string    `json:"alamat"`
	Status          string    `json:"status"           gorm:"default:aktif"`
	PenanggungJawab string    `json:"penanggung_jawab"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
}

// ─────────────────────────────────────────────
// USER (nasabah + operator + admin dalam 1 tabel)
// ─────────────────────────────────────────────

type User struct {
	ID           uint      `json:"id"            gorm:"primaryKey"`
	BSUID        *uint     `json:"bsu_id"`                          // pointer agar bisa NULL untuk admin
	BSU          *BSU      `json:"bsu,omitempty" gorm:"foreignKey:BSUID"`
	Nama         string    `json:"nama"          gorm:"not null"`
	Email        string    `json:"email"         gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-"`                               // json:"-" agar tidak pernah muncul di response
	NIK          string    `json:"nik"`
	Role         string    `json:"role"          gorm:"default:nasabah"` // "nasabah" | "operator" | "admin"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ─────────────────────────────────────────────
// HIERARKI SAMPAH
// ─────────────────────────────────────────────

type JenisSampah struct {
	ID        uint      `json:"id"   gorm:"primaryKey"`
	Nama      string    `json:"nama" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

type KategoriSampah struct {
	ID             uint        `json:"id"              gorm:"primaryKey"`
	IDJenisSampah  uint        `json:"id_jenis_sampah" gorm:"not null"`
	JenisSampah    JenisSampah `json:"jenis_sampah"    gorm:"foreignKey:IDJenisSampah"`
	Nama           string      `json:"nama"            gorm:"not null"`
	CreatedAt      time.Time   `json:"created_at"`
}

type SubKategoriSampah struct {
	ID               uint           `json:"id"                gorm:"primaryKey"`
	IDKategoriSampah uint           `json:"id_kategori_sampah" gorm:"not null"`
	KategoriSampah   KategoriSampah `json:"kategori_sampah"    gorm:"foreignKey:IDKategoriSampah"`
	Nama             string         `json:"nama"               gorm:"not null"`
	HargaPerKg       float64        `json:"harga_per_kg"       gorm:"type:decimal(10,2);not null"` // DECIMAL, bukan TEXT
	CreatedAt        time.Time      `json:"created_at"`
}

// ─────────────────────────────────────────────
// TRANSAKSI
// ─────────────────────────────────────────────

type Transaksi struct {
	ID              uint              `json:"id"               gorm:"primaryKey"`
	IDNasabah       uint              `json:"id_nasabah"       gorm:"not null"`
	Nasabah         User              `json:"nasabah"          gorm:"foreignKey:IDNasabah"`
	IDSubKategori   uint              `json:"id_sub_kategori"  gorm:"not null"`
	SubKategori     SubKategoriSampah `json:"sub_kategori"     gorm:"foreignKey:IDSubKategori"`
	BeratKg         float64           `json:"berat_kg"         gorm:"type:decimal(12,2);not null"`
	NilaiRupiah     float64           `json:"nilai_rupiah"     gorm:"type:decimal(12,2)"` // dihitung otomatis
	Tanggal         time.Time         `json:"tanggal"`
	CreatedAt       time.Time         `json:"created_at"`
}

// ─────────────────────────────────────────────
// AUTO MIGRATE
// ─────────────────────────────────────────────

// AutoMigrate membuat semua tabel otomatis saat app pertama jalan
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&Kabupaten{},
		&Kecamatan{},
		&BSU{},
		&User{},
		&JenisSampah{},
		&KategoriSampah{},
		&SubKategoriSampah{},
		&Transaksi{},
	)
}
