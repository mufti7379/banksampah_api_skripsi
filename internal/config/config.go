package config

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yourusername/banksampah-api/internal/handlers"
	"github.com/yourusername/banksampah-api/internal/middleware"
)

// Config menyimpan semua konfigurasi aplikasi dari environment variable
type Config struct {
	Port     string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
	JWTSecret string
}

// Load membaca environment variable — di GCP disimpan di Secret Manager
func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASS", "password"),
		DBName:    getEnv("DB_NAME", "banksampah"),
		JWTSecret: getEnv("JWT_SECRET", "rahasia-ganti-ini"),
	}
}

// getEnv mengambil env variable, pakai fallback jika tidak ada
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// ConnectDB membuat koneksi ke PostgreSQL menggunakan GORM
func ConnectDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// SetupRouter mendaftarkan semua route API
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Handler per resource
	kecamatanH := handlers.NewKecamatanHandler(db)
	bsuH       := handlers.NewBSUHandler(db)
	authH      := handlers.NewAuthHandler(db)
	nasabahH   := handlers.NewNasabahHandler(db)
	transaksiH := handlers.NewTransaksiHandler(db)
	sampahH    := handlers.NewSampahHandler(db)
	laporanH   := handlers.NewLaporanHandler(db)

	// Health check — untuk GCP Cloud Run
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		// ── AUTH (public, tidak butuh token) ──────────────────────────
		auth := api.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login",    authH.Login)
		}

		// ── KECAMATAN (public read, admin write) ──────────────────────
		kec := api.Group("/kecamatan")
		{
			kec.GET("",     kecamatanH.GetAll)         // semua role bisa baca
			kec.GET("/:id", kecamatanH.GetByID)
			kec.POST("",    middleware.AdminOnly(), kecamatanH.Create)
			kec.PUT("/:id", middleware.AdminOnly(), kecamatanH.Update)
			kec.DELETE("/:id", middleware.AdminOnly(), kecamatanH.Delete)
		}

		// ── BSU ───────────────────────────────────────────────────────
		bsu := api.Group("/bsu", middleware.AuthRequired())
		{
			bsu.GET("",        bsuH.GetAll)
			bsu.GET("/:id",    bsuH.GetByID)
			bsu.POST("",       middleware.AdminOnly(), bsuH.Create)
			bsu.PUT("/:id",    middleware.AdminOnly(), bsuH.Update)
			bsu.DELETE("/:id", middleware.AdminOnly(), bsuH.Delete)
		}

		// ── NASABAH ───────────────────────────────────────────────────
		nasabah := api.Group("/nasabah", middleware.AuthRequired())
		{
			nasabah.GET("/profile", nasabahH.GetProfile)   // nasabah lihat diri sendiri
			nasabah.PUT("/profile", nasabahH.UpdateProfile)
			nasabah.GET("/saldo",   nasabahH.GetSaldo)
			nasabah.GET("",        middleware.OperatorOrAdmin(), nasabahH.GetAll)
			nasabah.GET("/:id",    middleware.OperatorOrAdmin(), nasabahH.GetByID)
		}

		// ── TRANSAKSI ─────────────────────────────────────────────────
		trx := api.Group("/transaksi", middleware.AuthRequired())
		{
			trx.GET("",      middleware.OperatorOrAdmin(), transaksiH.GetAll)
			trx.GET("/:id",  transaksiH.GetByID)
			trx.POST("",     middleware.OperatorOrAdmin(), transaksiH.Create) // operator input
			trx.PUT("/:id",  middleware.OperatorOrAdmin(), transaksiH.Update)
			trx.DELETE("/:id", middleware.OperatorOrAdmin(), transaksiH.Delete)
			// nasabah lihat transaksinya sendiri
			trx.GET("/saya", transaksiH.GetMine)
		}

		// ── MASTER DATA SAMPAH ────────────────────────────────────────
		sampah := api.Group("/sampah", middleware.AuthRequired())
		{
			sampah.GET("/jenis",        sampahH.GetJenis)
			sampah.POST("/jenis",       middleware.AdminOnly(), sampahH.CreateJenis)
			sampah.PUT("/jenis/:id",    middleware.AdminOnly(), sampahH.UpdateJenis)

			sampah.GET("/kategori",     sampahH.GetKategori)
			sampah.POST("/kategori",    middleware.AdminOnly(), sampahH.CreateKategori)
			sampah.PUT("/kategori/:id", middleware.AdminOnly(), sampahH.UpdateKategori)

			sampah.GET("/sub",          sampahH.GetSubKategori)
			sampah.POST("/sub",         middleware.AdminOnly(), sampahH.CreateSubKategori)
			sampah.PUT("/sub/:id",      middleware.OperatorOrAdmin(), sampahH.UpdateSubKategori)
		}

		// ── LAPORAN ───────────────────────────────────────────────────
		laporan := api.Group("/laporan", middleware.AuthRequired(), middleware.AdminOnly())
		{
			laporan.GET("", laporanH.GetLaporan) // ?kecamatan_id=&from=&to=
		}
	}

	return r
}
