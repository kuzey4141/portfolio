package about // about paketi tanımlanır, bu dosyadaki tüm kodlar bu pakete aittir

import (
	"context"  // Veritabanı işlemleri için context kullanılır
	"fmt"      // Terminale hata veya bilgi mesajı yazmak için
	"net/http" // HTTP statü kodları için
	"strconv"  // String ifadeleri integer'a çevirmek için (örneğin ID)

	"github.com/gin-gonic/gin" // Gin framework'ü kullanmak için
	"github.com/jackc/pgx/v5"  // PostgreSQL veritabanı ile iletişim kurmak için pgx kütüphanesi
)

// About struct'ı, veritabanındaki "about" tablosundaki bir satırı temsil eder
type About struct {
	ID      int    `json:"id"`      // JSON çıktısında "id" alanı olarak gösterilir
	Content string `json:"content"` // JSON çıktısında "content" alanı olarak gösterilir
}

var Conn *pgx.Conn // Veritabanı bağlantısını global olarak tutacak değişken

// SetDB fonksiyonu, dışarıdan alınan veritabanı bağlantısını bu pakette kullanılmak üzere ayarlar
func SetDB(conn *pgx.Conn) {
	Conn = conn // Gelen bağlantı global değişkene atanır
}

// GetAbouts fonksiyonu, veritabanından tüm "about" kayıtlarını çeker ve JSON olarak döner
func GetAbouts(c *gin.Context) {
	rows, err := Conn.Query(context.Background(), "SELECT id, content FROM about") // SQL sorgusu çalıştırılır
	if err != nil {
		fmt.Println(err)                                                         // Hata varsa terminale yazdırılır
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veri alınamadı"}) // HTTP 500 hatası JSON ile dönülür
		return
	}
	defer rows.Close() // Sorgu tamamlandıktan sonra kaynak serbest bırakılır

	var abouts []About // Gelen verileri tutmak için boş bir slice tanımlanır
	for rows.Next() {  // Her satır için döngü
		var a About                                          // Geçici About nesnesi
		if err := rows.Scan(&a.ID, &a.Content); err != nil { // Satırdaki veriler struct'a atanır
			fmt.Println(err)                                                          // Hata varsa terminale yazdır
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Satır okunamadı"}) // Hata varsa JSON dön
			return
		}
		abouts = append(abouts, a) // Struct slice'a eklenir
	}

	c.JSON(http.StatusOK, abouts) // JSON olarak tüm kayıtları dön
}

// DeleteAbout fonksiyonu, belirtilen ID'ye sahip about kaydını siler
func DeleteAbout(c *gin.Context) {
	idStr := c.Param("id")         // URL parametresinden ID alınır (/api/about/:id şeklinde)
	id, err := strconv.Atoi(idStr) // String ID integer'a çevrilir
	if err != nil {
		fmt.Println(err)                                             // Hata yazdırılır
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"}) // Geçersiz ID ise 400 dönülür
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM about WHERE id=$1", id) // Silme sorgusu çalıştırılır
	if err != nil {
		fmt.Println(err)                                                                 // Hata yazdırılır
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silme işlemi başarısız"}) // Başarısızsa 500 dön
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "About kaydı silindi."}) // Başarı mesajı JSON olarak dönülür
}

// UpdateAbout fonksiyonu, belirtilen ID’ye sahip about kaydının içeriğini günceller
func UpdateAbout(c *gin.Context) {
	var a About
	if err := c.ShouldBindJSON(&a); err != nil { // JSON verisi struct'a bind edilir
		fmt.Println("JSON çözümleme hatası:", err)                     // Hata varsa yazdırılır
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"}) // Hatalı veri ise 400 dönülür
		return
	}

	_, err := Conn.Exec(context.Background(), "UPDATE about SET content=$1 WHERE id=$2", a.Content, a.ID) // Güncelleme sorgusu çalıştırılır
	if err != nil {
		fmt.Println("Veritabanı güncelleme hatası:", err)                              // Hata yazdırılır
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız"}) // Başarısızsa 500 dönülür
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("About ID %d başarıyla güncellendi", a.ID)}) // Başarı mesajı JSON ile dönülür
}

// CreateAbout fonksiyonu, yeni bir about kaydı ekler
func CreateAbout(c *gin.Context) {
	var a About
	if err := c.ShouldBindJSON(&a); err != nil { // JSON verisi struct'a bind edilir
		fmt.Println("JSON çözümleme hatası:", err)                     // Hata varsa yazdırılır
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"}) // Hatalı veri ise 400 dönülür
		return
	}

	_, err := Conn.Exec(context.Background(), "INSERT INTO about (content) VALUES ($1)", a.Content) // Yeni kayıt ekleme sorgusu çalıştırılır
	if err != nil {
		fmt.Println("Veritabanı ekleme hatası:", err)                              // Hata yazdırılır
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kayıt eklenemedi"}) // Eklenemezse 500 dönülür
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Yeni About kaydı başarıyla eklendi"}) // Başarı mesajı 201 ile dönülür
}
