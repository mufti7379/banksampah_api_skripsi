package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yourusername/banksampah-api/internal/models"
)

type TransaksiHandler struct{ db *gorm.DB }

func NewTransaksiHandler(db *gorm.DB) *TransaksiHandler {
	return &TransaksiHandler{db: db}
}

type TransaksiInput struct {
	IDNasabah     uint      `json:"id_nasabah"      binding:"required"`
	IDSubKategori uint      `json:"id_sub_kategori" binding:"required"`
	BeratKg       float64   `json:"berat_kg"        binding:"required,gt=0"`
	Tanggal       time.Time `json:"tanggal"`
}

// POST /api/v1/transaksi
// Hanya operator atau admin yang bisa akses (sudah di-guard di router)
func (h *TransaksiHandler) Create(c *gin.Context) {
	var input TransaksiInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// ── VALIDASI 1: pastikan id_nasabah adalah user dengan role "nasabah" ──
	// Ini kunci keamanan sistem — bukan hanya cek ID, tapi juga cek role-nya
	var nasabah models.User
	result := h.db.Where("id = ? AND role = ?", input.IDNasabah, "nasabah").First(&nasabah)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id_nasabah tidak valid atau bukan role nasabah",
		})
		return
	}

	// ── VALIDASI 2: pastikan sub_kategori ada ──────────────────────────────
	var subKategori models.SubKategoriSampah
	if err := h.db.First(&subKategori, input.IDSubKategori).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "sub_kategori tidak ditemukan"})
		return
	}

	// ── HITUNG nilai rupiah OTOMATIS ───────────────────────────────────────
	// Rumus: nilai_rupiah = berat_kg × harga_per_kg
	// Kalkulasi di backend — client tidak perlu kirim harga
	nilaiRupiah := input.BeratKg * subKategori.HargaPerKg

	// Set tanggal default ke sekarang jika tidak dikirim
	if input.Tanggal.IsZero() {
		input.Tanggal = time.Now()
	}

	transaksi := models.Transaksi{
		IDNasabah:     input.IDNasabah,
		IDSubKategori: input.IDSubKategori,
		BeratKg:       input.BeratKg,
		NilaiRupiah:   nilaiRupiah, // hasil kalkulasi otomatis
		Tanggal:       input.Tanggal,
	}

	if err := h.db.Create(&transaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal simpan transaksi"})
		return
	}

	// Load relasi untuk response yang lengkap
	h.db.Preload("Nasabah").Preload("SubKategori.KategoriSampah.JenisSampah").
		First(&transaksi, transaksi.ID)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Transaksi berhasil dicatat",
		"data":    transaksi,
	})
}

// GET /api/v1/transaksi/saya
// Nasabah lihat transaksinya sendiri
func (h *TransaksiHandler) GetMine(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var transaksi []models.Transaksi
	h.db.Preload("SubKategori").
		Where("id_nasabah = ?", userID).
		Order("tanggal DESC").
		Find(&transaksi)

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": transaksi})
}

// GET /api/v1/transaksi — untuk operator & admin
func (h *TransaksiHandler) GetAll(c *gin.Context) {
	var transaksi []models.Transaksi
	h.db.Preload("Nasabah").Preload("SubKategori").
		Order("tanggal DESC").Find(&transaksi)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": transaksi})
}

// GET /api/v1/transaksi/:id
func (h *TransaksiHandler) GetByID(c *gin.Context) {
	var transaksi models.Transaksi
	if err := h.db.Preload("Nasabah").Preload("SubKategori.KategoriSampah.JenisSampah").
		First(&transaksi, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaksi tidak ditemukan"})
		return
	}

	// Nasabah hanya boleh lihat transaksinya sendiri
	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")
	if role == "nasabah" && transaksi.IDNasabah != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Tidak punya akses ke transaksi ini"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": transaksi})
}

// PUT /api/v1/transaksi/:id
func (h *TransaksiHandler) Update(c *gin.Context) {
	var transaksi models.Transaksi
	if err := h.db.First(&transaksi, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaksi tidak ditemukan"})
		return
	}

	var input TransaksiInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Hitung ulang nilai rupiah jika berat atau sub-kategori berubah
	var subKategori models.SubKategoriSampah
	h.db.First(&subKategori, input.IDSubKategori)

	transaksi.BeratKg       = input.BeratKg
	transaksi.IDSubKategori = input.IDSubKategori
	transaksi.NilaiRupiah   = input.BeratKg * subKategori.HargaPerKg

	h.db.Save(&transaksi)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": transaksi})
}

// DELETE /api/v1/transaksi/:id
func (h *TransaksiHandler) Delete(c *gin.Context) {
	if err := h.db.Delete(&models.Transaksi{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaksi tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Transaksi dihapus"})
}
