package home

import (
	"context"       // Veritabanı işlemleri için context
	"fmt"           // Formatlı yazdırma için
	"net/http"      // HTTP statü kodları için
	"strconv"       // String to int dönüşümü için

	"github.com/gin-gonic/gin"    // Gin framework
	"github.com/jackc/pgx/v5"     // PostgreSQL bağlantısı için pgx kütüphanesi
)

// Home struct, home tablosundaki verileri temsil eder
type Home struct {
	ID          int    `json:"id"`          // JSON çıktısında id olarak görünür
	Title       string `json:"title"`       // Başlık alanı
	Description string `json:"description"` // Açıklama alanı
}

// Conn, veritabanı bağlantısı için global değişken
var Conn *pgx.Conn

// SetDB fonksiyonu, main.go'dan veritabanı bağlantısını alır
func SetDB(conn *pgx.Conn) {
	Conn = conn
}

// DeleteHome belirli ID'ye göre bir home kaydını siler
func DeleteHome(c *gin.Context) {
	idStr := c.Param("id") // URL parametresinden id'yi al (/api/home/:id şeklinde)
	id, err := strconv.Atoi(idStr) // string'i integer'a çevir
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"}) // Hatalı ID ise 400 dön
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM home WHERE id=$1", id) // Veritabanından silme sorgusu
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silme işlemi başarısız"}) // Silme başarısızsa 500 dön
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Home ID %d başarıyla silindi", id)}) // Başarı mesajı dön
}

// GetHomes, HTTP GET isteği geldiğinde home tablosundaki tüm verileri döner
func GetHomes(c *gin.Context) {
	rows, err := Conn.Query(context.Background(), "SELECT id, title, description FROM home") // Tüm home kayıtlarını seç
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veri alınamadı"}) // Veri çekilemezse 500 dön
		return
	}
	defer rows.Close() // İş bittiğinde bağlantıyı kapat

	var homes []Home              // Boş slice oluştur
	for rows.Next() {             // Satırlar arasında döngü
		var h Home
		if err := rows.Scan(&h.ID, &h.Title, &h.Description); err != nil { // Satırdan verileri al
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Satır okunamadı"}) // Hata olursa 500 dön
			return
		}
		homes = append(homes, h) // Slice'a ekle
	}

	c.JSON(http.StatusOK, homes) // JSON olarak tüm kayıtları dön
}

// UpdateHome bir home kaydını günceller
func UpdateHome(c *gin.Context) {
	var h Home
	if err := c.ShouldBindJSON(&h); err != nil { // JSON verisini struct'a bind et (decode et)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"}) // JSON hatalıysa 400 dön
		return
	}

	_, err := Conn.Exec(context.Background(), "UPDATE home SET title=$1, description=$2 WHERE id=$3", h.Title, h.Description, h.ID) // Güncelleme sorgusu
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız"}) // Başarısızsa 500 dön
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Home ID %d başarıyla güncellendi", h.ID)}) // Başarı mesajı dön
}

// CreateHome fonksiyonu, yeni bir home kaydı ekler
func CreateHome(c *gin.Context) {
	var h Home
	if err := c.ShouldBindJSON(&h); err != nil { // JSON verisini struct'a bind et
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"}) // Geçersiz veri ise 400 dön
		return
	}

	_, err := Conn.Exec(context.Background(), "INSERT INTO home (title, description) VALUES ($1, $2)", h.Title, h.Description) // Yeni kayıt ekle
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kayıt eklenemedi"}) // Eklenemezse 500 dön
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Home kaydı başarıyla eklendi"}) // Başarı mesajı dön (201)
}
