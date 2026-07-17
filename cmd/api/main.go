package main

import (
	"log"

	"github.com/yourusername/banksampah-api/internal/config"
	"github.com/yourusername/banksampah-api/internal/models"
)

func main() {
	// 1. Load konfigurasi dari environment variable
	cfg := config.Load()

	// 2. Koneksi ke database PostgreSQL (Cloud SQL)
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	// 3. Auto-migrate: buat tabel otomatis dari struct model
	models.AutoMigrate(db)

	// 4. Setup router Gin + jalankan server
	r := config.SetupRouter(db)
	log.Printf("Server jalan di port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
