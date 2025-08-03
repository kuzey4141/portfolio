package user // user paketi

import (
	"context"            // DB işlemleri için context
	"fmt"                // Konsola yazdırmak için
	"strconv"            // String - int dönüşümü için
	"github.com/gin-gonic/gin" // Gin framework
	"github.com/jackc/pgx/v5"  // PostgreSQL kütüphanesi
)

// User struct, users tablosundaki verileri temsil eder
type User struct {
	ID       int    `json:"id"`                 // JSON'da id alanı
	Username string `json:"username"`           // JSON'da kullanıcı adı
	Password string `json:"password,omitempty"` // JSON'da şifre (gönderilmezse boş gösterilir)
	Email    string `json:"email"`              // JSON'da email
}

var Conn *pgx.Conn // Global veritabanı bağlantısı

func SetDB(conn *pgx.Conn) {
	Conn = conn // Bağlantıyı set et
}

// DeleteUser silme işlemi
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")              // URL'den id parametresi al
	id, err := strconv.Atoi(idStr)      // String to int dönüştür
	if err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz ID"}) // Hatalı ID ise JSON ile hata döndür
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id) // DB'den sil
	if err != nil {
		c.JSON(500, gin.H{"error": "Silme işlemi başarısız"}) // DB hatası ise JSON hata
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("User ID %d başarıyla silindi", id)}) // Başarı mesajı JSON
}

// GetUsers kullanıcı listesini döner (şifre hariç)
func GetUsers(c *gin.Context) {
	rows, err := Conn.Query(context.Background(), "SELECT id, username, email FROM users") // Kullanıcıları çek
	if err != nil {
		c.JSON(500, gin.H{"error": "Veri alınamadı"}) // Hata varsa JSON dön
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			c.JSON(500, gin.H{"error": "Satır okunamadı"}) // Satır okuma hatası
			return
		}
		users = append(users, u)
	}

	c.JSON(200, users) // JSON olarak kullanıcı listesini döndür
}

// UpdateUser kullanıcı güncelleme işlemi
func UpdateUser(c *gin.Context) {
	var u User
	if err := c.BindJSON(&u); err != nil { // JSON gövdesini User struct'a bağla
		c.JSON(400, gin.H{"error": "Geçersiz veri"}) // Hatalı veri ise JSON döndür
		return
	}

	if u.Password != "" {
		_, err := Conn.Exec(context.Background(), "UPDATE users SET username=$1, email=$2, password=$3 WHERE id=$4", u.Username, u.Email, u.Password, u.ID) // Şifre dahil güncelle
		if err != nil {
			c.JSON(500, gin.H{"error": "Güncelleme başarısız"})
			return
		}
	} else {
		_, err := Conn.Exec(context.Background(), "UPDATE users SET username=$1, email=$2 WHERE id=$3", u.Username, u.Email, u.ID) // Şifresiz güncelle
		if err != nil {
			c.JSON(500, gin.H{"error": "Güncelleme başarısız"})
			return
		}
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("User ID %d başarıyla güncellendi", u.ID)}) // Başarı mesajı
}

// CreateUser yeni kullanıcı ekleme
func CreateUser(c *gin.Context) {
	var u User
	if err := c.BindJSON(&u); err != nil { // JSON'u User struct'a bağla
		c.JSON(400, gin.H{"error": "Geçersiz veri"})
		return
	}

	// Şifre hashleme işlemi burada yapılmalı, şimdilik düz kayıt

	_, err := Conn.Exec(context.Background(), "INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", u.Username, u.Password, u.Email) // DB ekleme
	if err != nil {
		c.JSON(500, gin.H{"error": "Kullanıcı eklenemedi"})
		return
	}

	c.JSON(201, gin.H{"message": "Kullanıcı başarıyla oluşturuldu"}) // Başarı mesajı
}