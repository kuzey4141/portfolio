package contact // contact paketi olarak tanımlandı

import (
	"context"          // Veritabanı işlemlerinde context için
	"fmt"              // Konsola yazdırma ve formatlama için
	"net/http"         // HTTP statü kodları için
	"strconv"          // String-integer dönüşümü için
	"github.com/gin-gonic/gin" // Gin framework kullanımı
	"github.com/jackc/pgx/v5"  // PostgreSQL bağlantı kütüphanesi
)

// Contact struct contact tablosunu temsil eder
type Contact struct {
	ID      int    `json:"id"`      // ID alanı, JSON'da "id" olarak gönderilir
	Email   string `json:"email"`   // Email alanı, JSON'da "email"
	Phone   string `json:"phone"`   // Telefon numarası, JSON'da "phone"
	Message string `json:"message"` // Mesaj içeriği, JSON'da "message"
}

var Conn *pgx.Conn // Global veritabanı bağlantısı tutulur

func SetDB(conn *pgx.Conn) {
	Conn = conn // Dışarıdan gelen bağlantı global Conn değişkenine atanır
}

func DeleteContact(c *gin.Context) {
	idStr := c.Param("id") // URL parametresinden ID alınır (/api/contact/:id)
	id, err := strconv.Atoi(idStr) // String ID integer'a çevrilir
	if err != nil {
		fmt.Println("ID dönüştürme hatası:", err) // Hata konsola yazılır
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"}) // 400 Bad Request dönülür
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM contact WHERE id=$1", id) // Silme sorgusu çalıştırılır
	if err != nil {
		fmt.Println("Silme hatası:", err) // Hata konsola yazılır
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silme işlemi başarısız"}) // 500 Internal Server Error dönülür
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Contact ID %d başarıyla silindi", id)}) // Başarı mesajı JSON olarak dönülür
}

func GetContacts(c *gin.Context) {
	rows, err := Conn.Query(context.Background(), "SELECT id, email, phone, message FROM contact") // Tüm contact kayıtları çekilir
	if err != nil {
		fmt.Println("Veri çekme hatası:", err) // Hata konsola yazılır
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veri alınamadı"}) // 500 dönülür
		return
	}
	defer rows.Close() // Fonksiyon sonunda kaynak kapatılır

	var contacts []Contact // Contact struct'larından oluşan liste
	for rows.Next() { // Satır satır döngü
		var cct Contact
		if err := rows.Scan(&cct.ID, &cct.Email, &cct.Phone, &cct.Message); err != nil {
			fmt.Println("Satır okunurken hata:", err) // Hata konsola yazılır
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Satır okunamadı"}) // 500 dönülür
			return
		}
		contacts = append(contacts, cct) // Listeye ekle
	}

	c.JSON(http.StatusOK, contacts) // JSON olarak tüm kayıtları dön
}

func UpdateContact(c *gin.Context) {
	var contact Contact
	if err := c.ShouldBindJSON(&contact); err != nil { // JSON verisi struct'a bind edilir
		fmt.Println("JSON çözümleme hatası:", err) // Hata konsola yazılır
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"}) // 400 dönülür
		return
	}

	_, err := Conn.Exec(context.Background(),
		"UPDATE contact SET email=$1, phone=$2, message=$3 WHERE id=$4",
		contact.Email, contact.Phone, contact.Message, contact.ID) // Güncelleme sorgusu çalıştırılır
	if err != nil {
		fmt.Println("Veritabanı güncelleme hatası:", err) // Hata konsola yazılır
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız"}) // 500 dönülür
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Contact ID %d başarıyla güncellendi", contact.ID)}) // Başarı mesajı dönülür
}

func CreateContact(c *gin.Context) {
	var contact Contact
	if err := c.ShouldBindJSON(&contact); err != nil { // JSON verisi struct'a bind edilir
		fmt.Println("JSON çözümleme hatası:", err) // Hata konsola yazılır
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"}) // 400 dönülür
		return
	}

	_, err := Conn.Exec(context.Background(),
		"INSERT INTO contact (email, phone, message) VALUES ($1, $2, $3)",
		contact.Email, contact.Phone, contact.Message) // Yeni kayıt ekleme sorgusu çalıştırılır
	if err != nil {
		fmt.Println("Veritabanı ekleme hatası:", err) // Hata konsola yazılır
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kayıt eklenemedi"}) // 500 dönülür
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Yeni iletişim kaydı başarıyla eklendi"}) // Başarı mesajı 201 ile dönülür
}