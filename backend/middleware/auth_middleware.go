package middleware

import (
	"net/http"
	"portfolio/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT token doğrulama middleware'i
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization header'ını al
		authHeader := c.GetHeader("Authorization")

		// Header yoksa hata döndür
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header gerekli",
			})
			c.Abort() // İsteği durdur, sonraki handler'lara gitme
			return
		}

		// "Bearer token123456" formatında gelir
		// "Bearer " kısmını çıkar, sadece token'ı al
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Token boşsa hata
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token bulunamadı",
			})
			c.Abort()
			return
		}

		// Token'ı doğrula (auth package'daki fonksiyonu kullan)
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Geçersiz veya süresi dolmuş token",
			})
			c.Abort()
			return
		}

		// Token geçerliyse, kullanıcı bilgilerini context'e ekle
		// Böylece sonraki handler'lar kullanıcı bilgilerine erişebilir
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// Token geçerli, sonraki handler'a devam et
		c.Next()
	}
}

// GetCurrentUser context'den kullanıcı bilgilerini al (helper function)
func GetCurrentUser(c *gin.Context) (userID int, username string, exists bool) {
	userIDInterface, exists1 := c.Get("user_id")
	usernameInterface, exists2 := c.Get("username")

	if !exists1 || !exists2 {
		return 0, "", false
	}

	userID, ok1 := userIDInterface.(int)
	username, ok2 := usernameInterface.(string)

	if !ok1 || !ok2 {
		return 0, "", false
	}

	return userID, username, true
}
