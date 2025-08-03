package projects // projects paketi

import (
	"context"            // DB işlemleri için context
	"fmt"                // Konsola yazdırmak için
	"strconv"            // String - int dönüşümleri için
	"github.com/gin-gonic/gin" // Gin framework

	"github.com/jackc/pgx/v5"  // PostgreSQL kütüphanesi
)

// Project struct, projects tablosundaki verileri temsil eder
type Project struct {
	ID          int    `json:"id"`          // JSON'da id olarak gösterilir
	Name        string `json:"name"`        // JSON'da name olarak gösterilir
	Description string `json:"description"` // JSON'da description olarak gösterilir
	Url         string `json:"url"`         // JSON'da url olarak gösterilir
}

var Conn *pgx.Conn // Global veritabanı bağlantısı

func SetDB(conn *pgx.Conn) {
	Conn = conn // Global bağlantıyı set et
}

// DeleteProject Gin handler: Silme işlemi için
func DeleteProject(c *gin.Context) {
	idStr := c.Param("id")                     // URL'den :id parametresini al
	id, err := strconv.Atoi(idStr)             // String'i int'e çevir
	if err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz ID"}) // Hatalı ID için JSON hata döndür
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM projects WHERE id=$1", id) // Silme sorgusu
	if err != nil {
		c.JSON(500, gin.H{"error": "Silme işlemi başarısız"}) // DB hatası varsa JSON hata
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Project ID %d başarıyla silindi", id)}) // Başarı mesajı JSON olarak
}

// GetProjects tüm projeleri JSON olarak döner
func GetProjects(c *gin.Context) {
	rows, err := Conn.Query(context.Background(), "SELECT id, name, description, url FROM projects") // Tüm projeleri çek
	if err != nil {
		c.JSON(500, gin.H{"error": "Veri alınamadı"}) // Hata varsa JSON dön
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Url); err != nil {
			c.JSON(500, gin.H{"error": "Satır okunamadı"}) // Satır okuma hatası
			return
		}
		projects = append(projects, p)
	}

	c.JSON(200, projects) // Başarıyla projeleri JSON olarak gönder
}

// UpdateProject güncelleme işlemi için Gin handler
func UpdateProject(c *gin.Context) {
	idStr := c.Param("id")                     // URL'den id al
	id, err := strconv.Atoi(idStr)             // String'i int yap
	if err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz ID"}) // Hatalı ID ise hata dön
		return
	}

	var p Project
	if err := c.BindJSON(&p); err != nil {    // JSON'dan Project struct'a bağla
		c.JSON(400, gin.H{"error": "Geçersiz veri"}) // JSON parse hatası
		return
	}

	sql := `UPDATE projects SET name=$1, description=$2, url=$3 WHERE id=$4`
	_, err = Conn.Exec(context.Background(), sql, p.Name, p.Description, p.Url, id) // DB güncelleme
	if err != nil {
		c.JSON(500, gin.H{"error": "Güncelleme başarısız"}) // DB hatası
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Project ID %d başarıyla güncellendi", id)}) // Başarı mesajı
}

// CreateProject yeni proje eklemek için Gin handler
func CreateProject(c *gin.Context) {
	var p Project
	if err := c.BindJSON(&p); err != nil {    // JSON'u Project struct'a bağla
		c.JSON(400, gin.H{"error": "Geçersiz veri"}) // JSON parse hatası
		return
	}

	_, err := Conn.Exec(context.Background(),
		"INSERT INTO projects (name, description, url) VALUES ($1, $2, $3)",
		p.Name, p.Description, p.Url) // DB'ye kayıt ekle
	if err != nil {
		c.JSON(500, gin.H{"error": "Kayıt eklenemedi"}) // Hata varsa JSON dön
		return
	}

	c.JSON(201, gin.H{"message": "Proje kaydı başarıyla eklendi"}) // Başarı mesajı
}