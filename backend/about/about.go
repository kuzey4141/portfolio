package about // about paketi tanımlanır, bu dosyadaki tüm kodlar bu pakete aittir

import (
	"context"       // Veritabanı işlemleri için context kullanılır
	"encoding/json" // JSON formatı ile veri çözümleme ve oluşturma işlemleri için
	"fmt"           // Terminale hata veya bilgi mesajı yazmak için
	"net/http"      // HTTP sunucusu işlemleri için
	"strconv"       // String ifadeleri integer'a çevirmek için (örneğin ID)
	"strings"       // String parçalama ve işleme işlemleri için

	"github.com/jackc/pgx/v5" // PostgreSQL veritabanı ile iletişim kurmak için kullanılan pgx kütüphanesi
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
func GetAbouts(w http.ResponseWriter, r *http.Request) {
	rows, err := Conn.Query(context.Background(), "SELECT id, content FROM about") // SQL sorgusu çalıştırılır
	if err != nil {
		fmt.Println(err)                                                // Hata varsa terminale yazdırılır
		http.Error(w, "Veri alınamadı", http.StatusInternalServerError) // HTTP 500 hatası döndürülür
		return
	}
	defer rows.Close() // Sorgu tamamlandıktan sonra kaynak serbest bırakılır

	var abouts []About // Gelen verileri tutmak için boş bir slice tanımlanır
	for rows.Next() {  // Her satır için döngü
		var a About                                          // Geçici About nesnesi
		if err := rows.Scan(&a.ID, &a.Content); err != nil { // Satırdaki veriler struct'a atanır
			fmt.Println(err)                                                 // Hata varsa terminale yazdır
			http.Error(w, "Satır okunamadı", http.StatusInternalServerError) // HTTP 500 hatası gönder
			return
		}
		abouts = append(abouts, a) // Struct slice'a eklenir
	}

	w.Header().Set("Content-Type", "application/json") // Yanıtın içeriği JSON olduğunu belirtir
	json.NewEncoder(w).Encode(abouts)                  // abouts dizisi JSON formatında client’a gönderilir
}

// DeleteAbout fonksiyonu, belirtilen ID'ye sahip about kaydını siler
func DeleteAbout(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/") // URL'yi "/" karakterine göre parçalar (örneğin: /api/about/delete/2 → ["", "api", "about", "2"])
	if len(parts) != 4 {                    // ID parametresi yoksa veya eksikse
		http.Error(w, "Geçersiz istek", http.StatusBadRequest) // HTTP 400 hatası gönder
		return
	}

	idStr := parts[3]              // ID değeri URL'nin 4. parçasıdır
	id, err := strconv.Atoi(idStr) // String olarak gelen ID, integer’a çevrilir
	if err != nil {
		fmt.Println(err)                                    // Hata yazdırılır
		http.Error(w, "Geçersiz ID", http.StatusBadRequest) // HTTP 400 hatası gönderilir
		return
	}

	_, err = Conn.Exec(context.Background(), "DELETE FROM about WHERE id=$1", id) // SQL silme sorgusu çalıştırılır
	if err != nil {
		fmt.Println(err)                                                        // Hata yazdırılır
		http.Error(w, "Silme işlemi başarısız", http.StatusInternalServerError) // HTTP 500 hatası gönderilir
		return
	}

	w.WriteHeader(http.StatusOK)            // HTTP 200 yanıtı gönderilir
	w.Write([]byte("About kaydı silindi.")) // Silme işlemi başarı mesajı olarak client'a yazılır
}

// UpdateAbout fonksiyonu, belirtilen ID’ye sahip about kaydının içeriğini günceller
func UpdateAbout(w http.ResponseWriter, r *http.Request) {
	var a About // Güncellenecek veriyi tutacak struct

	if err := json.NewDecoder(r.Body).Decode(&a); err != nil { // İstekten gelen JSON gövdesi çözümlenir ve a değişkenine atanır
		fmt.Println("JSON çözümleme hatası:", err)            // JSON parsing hatası varsa yazdırılır
		http.Error(w, "Geçersiz veri", http.StatusBadRequest) // HTTP 400 hatası döndürülür
		return
	}

	_, err := Conn.Exec(context.Background(), "UPDATE about SET content=$1 WHERE id=$2", a.Content, a.ID) // SQL güncelleme sorgusu çalıştırılır
	if err != nil {
		fmt.Println("Veritabanı güncelleme hatası:", err)                     // Hata varsa terminale yaz
		http.Error(w, "Güncelleme başarısız", http.StatusInternalServerError) // HTTP 500 döndür
		return
	}

	fmt.Fprintf(w, "About ID %d başarıyla güncellendi", a.ID) // Başarı mesajı client'a gönderilir
}

// CreateAbout fonksiyonu, yeni bir about kaydı ekler
func CreateAbout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Sadece POST metodu kabul edilir
		http.Error(w, "Sadece POST istekleri kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	var a About                                                // Gövdeden gelen veriyi tutmak için About struct'ı oluştur
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil { // JSON çözümlemesi yapılır
		fmt.Println("JSON çözümleme hatası:", err)
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	// SQL INSERT sorgusu çalıştırılır
	_, err := Conn.Exec(context.Background(), "INSERT INTO about (content) VALUES ($1)", a.Content)
	if err != nil {
		fmt.Println("Veritabanı ekleme hatası:", err)
		http.Error(w, "Kayıt eklenemedi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)                    // HTTP 201 yanıtı gönder
	fmt.Fprint(w, "Yeni About kaydı başarıyla eklendi.") // Başarı mesajı yazdır
}
