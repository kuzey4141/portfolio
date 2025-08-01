package user // user paketi tanımlandı

import (
	"context"       // veritabanı işlemleri için context paketi
	"encoding/json" // JSON encoding/decoding için
	"fmt"           // formatlı IO için
	"net/http"      // HTTP istek ve yanıt işlemleri için
	"strconv"       // string-int dönüşümü için
	"strings"       // string işlemleri için

	"github.com/jackc/pgx/v5" // PostgreSQL için pgx kütüphanesi
)

// User struct'ı, user tablosundaki sütunları temsil eder
type User struct {
	ID       int    `json:"id"`                 // kullanıcı ID'si, JSON'da "id"
	Username string `json:"username"`           // kullanıcı adı, JSON'da "username"
	Password string `json:"password,omitempty"` // şifre, JSON'da "password", boşsa gösterilmez
	Email    string `json:"email"`              // email adresi, JSON'da "email"
}

var Conn *pgx.Conn // global veritabanı bağlantısı tutar

// SetDB fonksiyonu, dışarıdan gelen veritabanı bağlantısını Conn değişkenine atar
func SetDB(conn *pgx.Conn) {
	Conn = conn // Conn global değişkene atandı
}

// DeleteUser fonksiyonu, verilen ID'ye göre kullanıcı siler
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/user/delete/") // URL'den ID çıkarılır
	id, err := strconv.Atoi(idStr)                               // ID string'ten integer'a çevrilir
	if err != nil {
		http.Error(w, "Geçersiz ID", http.StatusBadRequest) // hatalı ID için 400 döner
		fmt.Println("ID dönüştürme hatası:", err)           // hata konsola yazılır
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id) // veritabanından silme işlemi
	if err != nil {
		fmt.Println("Silme hatası:", err)                                       // hata konsola yazılır
		http.Error(w, "Silme işlemi başarısız", http.StatusInternalServerError) // 500 hata döner
		return
	}

	fmt.Fprintf(w, "User ID %d başarıyla silindi", id) // başarı mesajı yazılır
}

// GetUsers fonksiyonu, kullanıcı listesini döner (şifre hariç)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := Conn.Query(context.Background(), "SELECT id, username, email FROM users") // sadece id, username, email sorgulanır
	if err != nil {
		fmt.Println("Veri çekme hatası:", err)                          // hata konsola yazılır
		http.Error(w, "Veri alınamadı", http.StatusInternalServerError) // 500 hata döner
		return
	}
	defer rows.Close() // fonksiyon sonunda veritabanı satırları kapatılır

	var users []User  // boş user slice oluşturulur
	for rows.Next() { // satır satır dönülür
		var u User                                                      // her satır için User struct
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil { // satır verileri struct'a atanır
			fmt.Println("Satır okunurken hata:", err)                        // hata varsa yazdırılır
			http.Error(w, "Satır okunamadı", http.StatusInternalServerError) // 500 hata döner
			return
		}
		users = append(users, u) // kullanıcı slice'ına ekle
	}

	w.Header().Set("Content-Type", "application/json") // içerik tipi JSON olarak ayarlanır
	json.NewEncoder(w).Encode(users)                   // users slice'ı JSON olarak yazılır
}

// UpdateUser fonksiyonu, kullanıcı bilgilerini günceller (şifre isteğe bağlı)
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil { // gelen JSON User struct'a decode edilir
		fmt.Println("JSON çözümleme hatası:", err)            // hata varsa yazdırılır
		http.Error(w, "Geçersiz veri", http.StatusBadRequest) // 400 hata döner
		return
	}

	if u.Password != "" {
		// şifre dolu ise şifre dahil güncelleme yap
		_, err := Conn.Exec(context.Background(), "UPDATE users SET username=$1, email=$2, password=$3 WHERE id=$4", u.Username, u.Email, u.Password, u.ID)
		if err != nil {
			fmt.Println("Veritabanı güncelleme hatası:", err)                     // hata yazdırılır
			http.Error(w, "Güncelleme başarısız", http.StatusInternalServerError) // 500 hata döner
			return
		}
	} else {
		// şifre boş ise sadece username ve email güncellenir
		_, err := Conn.Exec(context.Background(), "UPDATE users SET username=$1, email=$2 WHERE id=$3", u.Username, u.Email, u.ID)
		if err != nil {
			fmt.Println("Veritabanı güncelleme hatası:", err)                     // hata yazdırılır
			http.Error(w, "Güncelleme başarısız", http.StatusInternalServerError) // 500 hata döner
			return
		}
	}

	fmt.Fprintf(w, "User ID %d başarıyla güncellendi", u.ID) // başarı mesajı yazılır
}

// CreateUser yeni kullanıcı oluşturur
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Sadece POST kabul edilir
		http.Error(w, "Sadece POST isteği kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	var u User
	// İstekten gelen JSON'u User struct'ına dönüştür
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println("JSON çözümleme hatası:", err)
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	// Burada normalde şifre hashlenmeli (örneğin bcrypt ile)
	// Şifre hashleme işlemi yapılmadan kayıt önerilmez

	// Veritabanına yeni kullanıcı ekle
	_, err := Conn.Exec(
		context.Background(),
		"INSERT INTO users (username, password, email) VALUES ($1, $2, $3)",
		u.Username, u.Password, u.Email,
	)
	if err != nil {
		fmt.Println("Veritabanı ekleme hatası:", err)
		http.Error(w, "Kullanıcı eklenemedi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)                // Başarı kodu 201
	fmt.Fprint(w, "Kullanıcı başarıyla oluşturuldu") // Başarı mesajı
}
