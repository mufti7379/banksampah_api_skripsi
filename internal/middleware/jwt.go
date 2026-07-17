package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	BSUID  *uint  `json:"bsu_id"`
	jwt.RegisteredClaims
}

// AuthRequired memastikan request punya token JWT yang valid
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token tidak ada"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token tidak valid atau sudah expired"})
			c.Abort()
			return
		}

		// Simpan info user ke context — bisa diambil di handler
		c.Set("user_id", claims.UserID)
		c.Set("role",    claims.Role)
		c.Set("bsu_id",  claims.BSUID)
		c.Next()
	}
}

// AdminOnly hanya izinkan role "admin"
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pastikan AuthRequired sudah jalan duluan
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Hanya admin yang boleh akses"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// OperatorOrAdmin izinkan role "operator" atau "admin"
func OperatorOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Tidak terautentikasi"})
			c.Abort()
			return
		}
		if role != "operator" && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Hanya operator atau admin yang boleh akses"})
			c.Abort()
			return
		}
		c.Next()
	}
}
