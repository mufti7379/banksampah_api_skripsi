package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yourusername/banksampah-api/internal/models"
)

// ─────────────────────────────────────────────
// RESPONSE HELPER
// ─────────────────────────────────────────────

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}
func created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": data})
}
func notFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": msg})
}
func badRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": msg})
}

// ─────────────────────────────────────────────
// KECAMATAN
// ─────────────────────────────────────────────

type KecamatanHandler struct{ db *gorm.DB }
func NewKecamatanHandler(db *gorm.DB) *KecamatanHandler { return &KecamatanHandler{db: db} }

func (h *KecamatanHandler) GetAll(c *gin.Context) {
	var data []models.Kecamatan
	h.db.Preload("Kabupaten").Find(&data)
	ok(c, data)
}
func (h *KecamatanHandler) GetByID(c *gin.Context) {
	var data models.Kecamatan
	if err := h.db.Preload("Kabupaten").First(&data, c.Param("id")).Error; err != nil {
		notFound(c, "Kecamatan tidak ditemukan"); return
	}
	ok(c, data)
}
func (h *KecamatanHandler) Create(c *gin.Context) {
	var data models.Kecamatan
	if err := c.ShouldBindJSON(&data); err != nil { badRequest(c, err.Error()); return }
	h.db.Create(&data)
	created(c, data)
}
func (h *KecamatanHandler) Update(c *gin.Context) {
	var data models.Kecamatan
	if err := h.db.First(&data, c.Param("id")).Error; err != nil { notFound(c, "Kecamatan tidak ditemukan"); return }
	c.ShouldBindJSON(&data)
	h.db.Save(&data)
	ok(c, data)
}
func (h *KecamatanHandler) Delete(c *gin.Context) {
	h.db.Delete(&models.Kecamatan{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Kecamatan dihapus"})
}

// ─────────────────────────────────────────────
// BSU
// ─────────────────────────────────────────────

type BSUHandler struct{ db *gorm.DB }
func NewBSUHandler(db *gorm.DB) *BSUHandler { return &BSUHandler{db: db} }

func (h *BSUHandler) GetAll(c *gin.Context) {
	var data []models.BSU
	h.db.Preload("Kecamatan.Kabupaten").Find(&data)
	ok(c, data)
}
func (h *BSUHandler) GetByID(c *gin.Context) {
	var data models.BSU
	if err := h.db.Preload("Kecamatan.Kabupaten").First(&data, c.Param("id")).Error; err != nil {
		notFound(c, "BSU tidak ditemukan"); return
	}
	ok(c, data)
}
func (h *BSUHandler) Create(c *gin.Context) {
	var data models.BSU
	if err := c.ShouldBindJSON(&data); err != nil { badRequest(c, err.Error()); return }
	h.db.Create(&data)
	created(c, data)
}
func (h *BSUHandler) Update(c *gin.Context) {
	var data models.BSU
	if err := h.db.First(&data, c.Param("id")).Error; err != nil { notFound(c, "BSU tidak ditemukan"); return }
	c.ShouldBindJSON(&data)
	h.db.Save(&data)
	ok(c, data)
}
func (h *BSUHandler) Delete(c *gin.Context) {
	h.db.Delete(&models.BSU{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "BSU dihapus"})
}

// ─────────────────────────────────────────────
// NASABAH
// ─────────────────────────────────────────────

type NasabahHandler struct{ db *gorm.DB }
func NewNasabahHandler(db *gorm.DB) *NasabahHandler { return &NasabahHandler{db: db} }

func (h *NasabahHandler) GetAll(c *gin.Context) {
	var data []models.User
	h.db.Where("role = ?", "nasabah").Find(&data)
	ok(c, data)
}
func (h *NasabahHandler) GetByID(c *gin.Context) {
	var data models.User
	if err := h.db.Where("id = ? AND role = ?", c.Param("id"), "nasabah").First(&data).Error; err != nil {
		notFound(c, "Nasabah tidak ditemukan"); return
	}
	ok(c, data)
}
func (h *NasabahHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var data models.User
	h.db.First(&data, userID)
	ok(c, data)
}
func (h *NasabahHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var data models.User
	h.db.First(&data, userID)
	c.ShouldBindJSON(&data)
	h.db.Save(&data)
	ok(c, data)
}
func (h *NasabahHandler) GetSaldo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var total float64
	h.db.Model(&models.Transaksi{}).
		Where("id_nasabah = ?", userID).
		Select("COALESCE(SUM(nilai_rupiah), 0)").Scan(&total)
	ok(c, gin.H{"saldo": total})
}

// ─────────────────────────────────────────────
// SAMPAH (jenis, kategori, sub-kategori)
// ─────────────────────────────────────────────

type SampahHandler struct{ db *gorm.DB }
func NewSampahHandler(db *gorm.DB) *SampahHandler { return &SampahHandler{db: db} }

func (h *SampahHandler) GetJenis(c *gin.Context) {
	var data []models.JenisSampah; h.db.Find(&data); ok(c, data)
}
func (h *SampahHandler) CreateJenis(c *gin.Context) {
	var data models.JenisSampah
	if err := c.ShouldBindJSON(&data); err != nil { badRequest(c, err.Error()); return }
	h.db.Create(&data); created(c, data)
}
func (h *SampahHandler) UpdateJenis(c *gin.Context) {
	var data models.JenisSampah
	h.db.First(&data, c.Param("id")); c.ShouldBindJSON(&data); h.db.Save(&data); ok(c, data)
}
func (h *SampahHandler) GetKategori(c *gin.Context) {
	var data []models.KategoriSampah
	h.db.Preload("JenisSampah").Find(&data); ok(c, data)
}
func (h *SampahHandler) CreateKategori(c *gin.Context) {
	var data models.KategoriSampah
	if err := c.ShouldBindJSON(&data); err != nil { badRequest(c, err.Error()); return }
	h.db.Create(&data); created(c, data)
}
func (h *SampahHandler) UpdateKategori(c *gin.Context) {
	var data models.KategoriSampah
	h.db.First(&data, c.Param("id")); c.ShouldBindJSON(&data); h.db.Save(&data); ok(c, data)
}
func (h *SampahHandler) GetSubKategori(c *gin.Context) {
	var data []models.SubKategoriSampah
	h.db.Preload("KategoriSampah.JenisSampah").Find(&data); ok(c, data)
}
func (h *SampahHandler) CreateSubKategori(c *gin.Context) {
	var data models.SubKategoriSampah
	if err := c.ShouldBindJSON(&data); err != nil { badRequest(c, err.Error()); return }
	h.db.Create(&data); created(c, data)
}
func (h *SampahHandler) UpdateSubKategori(c *gin.Context) {
	var data models.SubKategoriSampah
	h.db.First(&data, c.Param("id")); c.ShouldBindJSON(&data); h.db.Save(&data); ok(c, data)
}

// ─────────────────────────────────────────────
// LAPORAN
// ─────────────────────────────────────────────

type LaporanHandler struct{ db *gorm.DB }
func NewLaporanHandler(db *gorm.DB) *LaporanHandler { return &LaporanHandler{db: db} }

// GET /api/v1/laporan?kecamatan_id=1&from=2025-01-01&to=2025-12-31
func (h *LaporanHandler) GetLaporan(c *gin.Context) {
	kecamatanID := c.Query("kecamatan_id")
	from        := c.Query("from")
	to          := c.Query("to")

	query := h.db.Model(&models.Transaksi{}).
		Joins("JOIN users ON users.id = transaksi.id_nasabah").
		Joins("JOIN bsus ON bsus.id = users.bsu_id")

	// Filter opsional berdasarkan kecamatan
	if kecamatanID != "" {
		query = query.Where("bsus.id_kecamatan = ?", kecamatanID)
	}
	if from != "" { query = query.Where("transaksi.tanggal >= ?", from) }
	if to   != "" { query = query.Where("transaksi.tanggal <= ?", to)   }

	var totalBerat, totalNilai float64
	var jumlahTransaksi int64
	query.Count(&jumlahTransaksi)
	query.Select("COALESCE(SUM(berat_kg), 0)").Scan(&totalBerat)
	query.Select("COALESCE(SUM(nilai_rupiah), 0)").Scan(&totalNilai)

	ok(c, gin.H{
		"jumlah_transaksi": jumlahTransaksi,
		"total_berat_kg":   totalBerat,
		"total_nilai_rupiah": totalNilai,
		"filter": gin.H{
			"kecamatan_id": kecamatanID,
			"from": from,
			"to":   to,
		},
	})
}
