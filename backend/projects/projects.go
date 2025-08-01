package projects // Bu dosya projects paketine ait

import (
	"context"       // Veritabanı işlemlerinde context kullanımı için
	"encoding/json" // JSON formatına dönüştürme ve çözümleme için
	"fmt"           // Konsola yazı yazdırmak için
	"net/http"      // HTTP istek/cevap işlemek için
	"strconv"       // String-int dönüşümü için
	"strings"       // String işlemleri (örneğin URL'den id ayıklamak) için

	"github.com/jackc/pgx/v5" // PostgreSQL veritabanı sürücüsü
)

// Project struct, projects tablosundaki verileri temsil eder
type Project struct {
	ID          int    `json:"id"`          // ID alanı, JSON'da "id" olarak adlandırılır
	Name        string `json:"name"`        // Proje adı, JSON'da "name"
	Description string `json:"description"` // Proje açıklaması, JSON'da "description"
	Url         string `json:"url"`         // Proje bağlantı adresi, JSON'da "url"
}

var Conn *pgx.Conn // Veritabanı bağlantısını tutan global değişken

// SetDB, dışarıdan alınan bağlantıyı paket içinde kullanmak için atar
func SetDB(conn *pgx.Conn) {
	Conn = conn // Global bağlantı değişkenine değer atanır
}

// DeleteProject, belirli bir ID'ye sahip projeyi siler
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/projects/delete/") // URL'den id kısmını çıkar
	id, err := strconv.Atoi(idStr)                                   // id'yi string'den int'e çevir
	if err != nil {
		fmt.Println(err)                                    // Hata varsa konsola yaz
		fmt.Println("ID dönüştürme hatası:", err)           // Ek bilgi yaz
		http.Error(w, "Geçersiz ID", http.StatusBadRequest) // HTTP 400 hatası gönder
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM projects WHERE id=$1", id) // SQL sorgusu ile silme işlemi yap
	if err != nil {
		fmt.Println(err)                                                        // Hata varsa konsola yaz
		fmt.Println("Silme hatası:", err)                                       // Ek hata bilgisi
		http.Error(w, "Silme işlemi başarısız", http.StatusInternalServerError) // HTTP 500 hatası gönder
		return
	}

	fmt.Fprintf(w, "Project ID %d başarıyla silindi", id) // Başarı mesajı döndür
}

// GetProjects tüm projeleri JSON olarak döndürür
func GetProjects(w http.ResponseWriter, r *http.Request) {
	rows, err := Conn.Query(context.Background(), "SELECT id, name, description, url FROM projects") // SQL sorgusu ile tüm projeleri çek
	if err != nil {
		fmt.Println(err)                                                // Hata yazdır
		fmt.Println("Query hatası:", err)                               // Ek hata bilgisi
		http.Error(w, "Veri alınamadı", http.StatusInternalServerError) // HTTP 500 hatası gönder
		return
	}
	defer rows.Close() // Sorgu sonuçlarını kapat

	var projects []Project // Proje dizisi tanımla
	for rows.Next() {      // Her satır için
		var p Project                                                             // Yeni Project nesnesi oluştur
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Url); err != nil { // Satır verilerini struct'a ata
			fmt.Println("Satır okunurken hata:", err)                        // Hata varsa yaz
			http.Error(w, "Satır okunamadı", http.StatusInternalServerError) // HTTP 500 hatası gönder
			return
		}
		projects = append(projects, p) // Okunan projeyi diziye ekle
	}

	w.Header().Set("Content-Type", "application/json") // İçerik türünü JSON olarak ayarla
	json.NewEncoder(w).Encode(projects)                // Proje listesini JSON olarak HTTP yanıtına yaz
}

// UpdateProject belirli ID'li bir projeyi günceller
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut { // PUT değilse
		http.Error(w, "Sadece PUT isteği kabul edilir", http.StatusMethodNotAllowed) // HTTP 405 döndür
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/projects/update/") // URL'den id'yi al
	id, err := strconv.Atoi(idStr)                                   // String'den int'e çevir
	if err != nil {
		fmt.Println("ID dönüştürme hatası:", err)           // Hata varsa yazdır
		http.Error(w, "Geçersiz ID", http.StatusBadRequest) // HTTP 400 döndür
		return
	}

	var p Project                                              // Güncellenecek proje nesnesi
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil { // JSON'dan proje nesnesini çıkar
		fmt.Println("JSON çözümleme hatası:", err)            // Hata varsa yazdır
		http.Error(w, "Geçersiz veri", http.StatusBadRequest) // HTTP 400 döndür
		return
	}

	sql := `UPDATE projects SET name=$1, description=$2, url=$3 WHERE id=$4`        // Güncelleme sorgusu
	_, err = Conn.Exec(context.Background(), sql, p.Name, p.Description, p.Url, id) // Sorguyu çalıştır
	if err != nil {
		fmt.Println("Veritabanı güncelleme hatası:", err)                     // Hata varsa yazdır
		http.Error(w, "Güncelleme başarısız", http.StatusInternalServerError) // HTTP 500 döndür
		return
	}

	fmt.Fprintf(w, "Project ID %d başarıyla güncellendi", id) // Başarı mesajı döndür
}

// CreateProject fonksiyonu, yeni bir proje kaydı ekler
func CreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Yalnızca POST isteklerine izin verilir
		http.Error(w, "Sadece POST isteği kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	var p Project                                              // İstekten gelecek JSON verisini tutacak Project struct'ı
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil { // JSON gövdesi çözümlenir
		fmt.Println("JSON çözümleme hatası:", err)
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	// SQL INSERT sorgusu: gelen name, description ve url alanları veritabanına eklenir
	_, err := Conn.Exec(
		context.Background(),
		"INSERT INTO projects (name, description, url) VALUES ($1, $2, $3)",
		p.Name, p.Description, p.Url, // Sırasıyla struct'taki veriler sorguya aktarılır
	)
	if err != nil {
		fmt.Println("Veritabanı ekleme hatası:", err)
		http.Error(w, "Kayıt eklenemedi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)               // HTTP 201: başarıyla oluşturuldu
	fmt.Fprint(w, "Proje kaydı başarıyla eklendi.") // Kullanıcıya mesaj döndürülür
}
