package home

import (
	"context"       // Veritabanı sorguları için context kullanımı
	"encoding/json" // JSON formatına dönüştürmek için
	"fmt"
	"net/http" // HTTP server ve istek/yanıt işlemleri için
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5" // PostgreSQL bağlantısı için pgx kütüphanesi
)

// Home struct, home tablosundaki verileri temsil eder
type Home struct {
	ID          int    `json:"id"`          // JSON çıktısında id olarak görünecek
	Title       string `json:"title"`       // Başlık alanı
	Description string `json:"description"` // İçerik alanı
}

// Conn, veritabanı bağlantısı için global değişken
var Conn *pgx.Conn

// SetDB fonksiyonu, main.go'dan veritabanı bağlantısını alır
func SetDB(conn *pgx.Conn) {
	Conn = conn
}

// DeleteHome belirli ID'ye göre bir home kaydını siler
func DeleteHome(w http.ResponseWriter, r *http.Request) {
	// URL'den ID'yi al (örnek: /api/home/delete/3)
	idStr := strings.TrimPrefix(r.URL.Path, "/api/home/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Geçersiz ID", http.StatusBadRequest)
		fmt.Println("ID dönüştürme hatası:", err)
		return
	}

	// SQL sorgusuyla sil
	_, err = Conn.Exec(context.Background(), "DELETE FROM home WHERE id=$1", id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Silme işlemi başarısız", http.StatusInternalServerError)
		fmt.Println("Silme hatası:", err)
		return
	}

	fmt.Fprintf(w, "Home ID %d başarıyla silindi", id)
}

// GetHomes, HTTP GET isteği geldiğinde home tablosundaki tüm verileri döner
func GetHomes(w http.ResponseWriter, r *http.Request) {
	// Veritabanından id, title, content sütunlarını seçiyoruz
	rows, err := Conn.Query(context.Background(), "SELECT id, title, description FROM home")
	if err != nil {
		fmt.Println(err)
		// Hata varsa HTTP 500 hatası döneriz
		http.Error(w, "Veri alınamadı", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // İşimiz bittiğinde veritabanı satırlarını kapatıyoruz

	var homes []Home // Boş bir slice oluşturuyoruz, verileri buraya koyacağız
	for rows.Next() {
		var h Home
		// Satırdan verileri Home structına dolduruyoruz
		if err := rows.Scan(&h.ID, &h.Title, &h.Description); err != nil {
			// Satır okunamazsa hata döneriz
			http.Error(w, "Satır okunamadı", http.StatusInternalServerError)
			return
		}
		homes = append(homes, h) // Slice'a ekliyoruz
	}

	// Yanıtın türü JSON olduğunu belirtiyoruz
	w.Header().Set("Content-Type", "application/json")
	// JSON formatına çevirip yanıt olarak gönderiyoruz
	json.NewEncoder(w).Encode(homes)
}

// UpdateHome bir home kaydını günceller
func UpdateHome(w http.ResponseWriter, r *http.Request) {
	var h Home
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil { // İstek gövdesindeki JSON Home struct'a decode edilir
		fmt.Println("JSON çözümleme hatası:", err)            // Hata varsa konsola yazdır
		http.Error(w, "Geçersiz veri", http.StatusBadRequest) // HTTP 400 hatası dön
		return
	}

	_, err := Conn.Exec(context.Background(), "UPDATE home SET title=$1, content=$2 WHERE id=$3", h.Title, h.Description, h.ID) // Veritabanında güncelleme yapılır
	if err != nil {
		fmt.Println("Veritabanı güncelleme hatası:", err)                     // Hata varsa konsola yazdır
		http.Error(w, "Güncelleme başarısız", http.StatusInternalServerError) // HTTP 500 hatası dön
		return
	}

	fmt.Fprintf(w, "Home ID %d başarıyla güncellendi", h.ID) // Başarılı mesaj HTTP cevabına yazılır
}

// CreateHome fonksiyonu, yeni bir home kaydı ekler
func CreateHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Sadece POST metoduna izin verilir
		http.Error(w, "Sadece POST isteği kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	var h Home                                                 // JSON'dan gelen verileri tutmak için Home struct'ı oluşturulur
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil { // JSON verisi çözülür
		fmt.Println("JSON çözümleme hatası:", err)
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	// SQL INSERT sorgusu çalıştırılır
	_, err := Conn.Exec(
		context.Background(),
		"INSERT INTO home (title, description) VALUES ($1, $2)", // Veritabanına yeni kayıt eklenir
		h.Title, h.Description, // Parametreler olarak title ve description verilir
	)
	if err != nil {
		fmt.Println("Veritabanı ekleme hatası:", err)
		http.Error(w, "Kayıt eklenemedi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)              // HTTP 201: başarıyla oluşturuldu
	fmt.Fprint(w, "Home kaydı başarıyla eklendi.") // Başarı mesajı gönderilir
}
