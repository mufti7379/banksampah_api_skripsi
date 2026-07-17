package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/yourusername/banksampah-api/internal/middleware"
	"github.com/yourusername/banksampah-api/internal/models"
)

type AuthHandler struct{ db *gorm.DB }

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// ─── REGISTER ────────────────────────────────────────────────────────────────

type RegisterInput struct {
	BSUID uint   `json:"bsu_id" binding:"required"`
	Nama  string `json:"nama"   binding:"required"`
	Email string `json:"email"  binding:"required,email"`
	Pass  string `json:"password" binding:"required,min=6"`
	NIK   string `json:"nik"`
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Cek email sudah dipakai
	var existing models.User
	if result := h.db.Where("email = ?", input.Email).First(&existing); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Email sudah terdaftar"})
		return
	}

	// Hash password — password asli tidak pernah disimpan
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Pass), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal proses password"})
		return
	}

	user := models.User{
		BSUID:        &input.BSUID,
		Nama:         input.Nama,
		Email:        input.Email,
		PasswordHash: string(hash),
		NIK:          input.NIK,
		Role:         "nasabah", // nasabah selalu daftar sebagai role nasabah
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal daftar"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data":    user,
	})
}

// ─── LOGIN ───────────────────────────────────────────────────────────────────

type LoginInput struct {
	Email string `json:"email"    binding:"required,email"`
	Pass  string `json:"password" binding:"required"`
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Cari user berdasarkan email
	var user models.User
	if err := h.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Email atau password salah"})
		return
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Pass)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Email atau password salah"})
		return
	}

	// Buat token JWT — expired 24 jam
	claims := &middleware.Claims{
		UserID: user.ID,
		Role:   user.Role,
		BSUID:  user.BSUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal buat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  tokenStr,
		"user": gin.H{
			"id":    user.ID,
			"nama":  user.Nama,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
