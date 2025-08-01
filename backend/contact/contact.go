package contact // contact paketi olarak tanımlandı

import (
	"context"       // Veritabanı işlemlerinde context için
	"encoding/json" // JSON verisi kodlama/çözme için
	"fmt"           // Konsola yazdırma ve formatlama için
	"net/http"      // HTTP istek ve cevap işlemleri için
	"strconv"       // String-integer dönüşümü için
	"strings"       // String işlemleri için

	"github.com/jackc/pgx/v5" // PostgreSQL bağlantı kütüphanesi
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

func DeleteContact(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/contact/delete/") // URL'den ID parametresi alınır
	id, err := strconv.Atoi(idStr)                                  // ID stringden integer'a çevrilir
	if err != nil {
		fmt.Println("ID dönüştürme hatası:", err)           // Hata varsa konsola yazdır
		http.Error(w, "Geçersiz ID", http.StatusBadRequest) // HTTP 400 hatası dön
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM contact WHERE id=$1", id) // Veritabanından ID ile kayıt silinir
	if err != nil {
		fmt.Println("Silme hatası:", err)                                       // Hata varsa konsola yazdır
		http.Error(w, "Silme işlemi başarısız", http.StatusInternalServerError) // HTTP 500 hatası dön
		return
	}

	fmt.Fprintf(w, "Contact ID %d başarıyla silindi", id) // Başarılı mesajı HTTP cevabına yazılır
}

func GetContacts(w http.ResponseWriter, r *http.Request) {
	rows, err := Conn.Query(context.Background(), "SELECT id, email, phone, message FROM contact") // Tüm contact kayıtları çekilir
	if err != nil {
		fmt.Println("Veri çekme hatası:", err)                          // Hata varsa konsola yazdır
		http.Error(w, "Veri alınamadı", http.StatusInternalServerError) // HTTP 500 hatası dön
		return
	}
	defer rows.Close() // Fonksiyon sonunda satırları kapat

	var contacts []Contact // Contact struct'larından oluşan boş liste oluşturulur
	for rows.Next() {      // Satır satır veriler okunur
		var c Contact
		if err := rows.Scan(&c.ID, &c.Email, &c.Phone, &c.Message); err != nil { // Satırdaki değerler struct'a aktarılır
			fmt.Println("Satır okunurken hata:", err)                        // Hata varsa konsola yazdır
			http.Error(w, "Satır okunamadı", http.StatusInternalServerError) // HTTP 500 hatası dön
			return
		}
		contacts = append(contacts, c) // Okunan contact listeye eklenir
	}

	w.Header().Set("Content-Type", "application/json") // Response header JSON olarak ayarlanır
	json.NewEncoder(w).Encode(contacts)                // contacts slice'ı JSON olarak HTTP cevabına yazılır
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {
	var c Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil { // HTTP isteği gövdesindeki JSON decode edilir Contact struct'a
		fmt.Println("JSON çözümleme hatası:", err)            // Hata varsa konsola yazdır
		http.Error(w, "Geçersiz veri", http.StatusBadRequest) // HTTP 400 hatası dön
		return
	}

	_, err := Conn.Exec(context.Background(), "UPDATE contact SET email=$1, phone=$2, message=$3 WHERE id=$4", c.Email, c.Phone, c.Message, c.ID) // Veritabanında güncelleme yapılır
	if err != nil {
		fmt.Println("Veritabanı güncelleme hatası:", err)                     // Hata varsa konsola yazdır
		http.Error(w, "Güncelleme başarısız", http.StatusInternalServerError) // HTTP 500 hatası dön
		return
	}

	fmt.Fprintf(w, "Contact ID %d başarıyla güncellendi", c.ID) // Başarılı mesaj HTTP cevabına yazılır
}

// CreateContact fonksiyonu, yeni bir contact kaydı oluşturur
func CreateContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Sadece POST metoduna izin verilir
		http.Error(w, "Sadece POST istekleri kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	var c Contact                                              // JSON'dan gelen veriyi tutacak struct
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil { // Gövdedeki JSON verisi çözümlenir
		fmt.Println("JSON çözümleme hatası:", err)
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	// INSERT SQL sorgusu çalıştırılır
	_, err := Conn.Exec(
		context.Background(),
		"INSERT INTO contact (email, phone, message) VALUES ($1, $2, $3)",
		c.Email, c.Phone, c.Message, // Veriler sorguya parametre olarak geçilir
	)
	if err != nil {
		fmt.Println("Veritabanı ekleme hatası:", err)
		http.Error(w, "Kayıt eklenemedi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)                       // HTTP 201: Başarıyla oluşturuldu
	fmt.Fprint(w, "Yeni iletişim kaydı başarıyla eklendi.") // Başarı mesajı client’a gönderilir
}
