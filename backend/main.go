package main // Ana paket

import (
	"context"  // Veritabanı işlemleri için context paketi
	"fmt"      // Konsola yazdırmak için
	"log"      // Loglama için
	"net/http" // HTTP server için

	"portfolio/about"    // about paketini import ettik
	"portfolio/contact"  // contact paketini import ettik
	"portfolio/home"     // home paketini import ettik
	"portfolio/projects" // projects paketini import ettik
	"portfolio/user"     // user paketini import ettik

	"github.com/jackc/pgx/v5" // pgx PostgreSQL kütüphanesi
)

func main() {
	// Veritabanı bağlantı stringi (kendi bilgilerine göre düzenle)
	connStr := "postgres://postgres:6303523@localhost:5432/portfolio?sslmode=disable"

	// Veritabanına bağlan
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Veritabanına bağlanırken hata: %v", err) // Bağlantı hatası varsa programı durdur
	}
	defer conn.Close(context.Background()) // Program kapanınca bağlantıyı kapat

	// Her pakete veritabanı bağlantısını set et
	home.SetDB(conn)     // home paketine bağlantı
	about.SetDB(conn)    // about paketine bağlantı
	projects.SetDB(conn) // projects paketine bağlantı
	contact.SetDB(conn)  // contact paketine bağlantı
	user.SetDB(conn)     // user paketine bağlantı

	// HTTP endpointlerini tanımla (GET için)
	http.HandleFunc("/api/home", home.GetHomes)            // home verilerini çekmek için
	http.HandleFunc("/api/about", about.GetAbouts)         // about verilerini çekmek için
	http.HandleFunc("/api/projects", projects.GetProjects) // projects verilerini çekmek için
	http.HandleFunc("/api/contact", contact.GetContacts)   // contact verilerini çekmek için
	http.HandleFunc("/api/user", user.GetUsers)            // user verilerini çekmek için

	// HTTP endpointlerini tanımla (DELETE için)
	http.HandleFunc("/api/home/delete/", home.DeleteHome)            // belirli bir home kaydını silmek için
	http.HandleFunc("/api/about/delete/", about.DeleteAbout)         // belirli bir about kaydını silmek için
	http.HandleFunc("/api/projects/delete/", projects.DeleteProject) // belirli bir project kaydını silmek için
	http.HandleFunc("/api/contact/delete/", contact.DeleteContact)   // belirli bir contact kaydını silmek için
	http.HandleFunc("/api/user/delete/", user.DeleteUser)            // belirli bir user kaydını silmek için

	// HTTP endpointlerini tanımla (PUT için - güncelleme)
	http.HandleFunc("/api/home/update", home.UpdateHome)            // home güncelleme için PUT
	http.HandleFunc("/api/about/update", about.UpdateAbout)         // about güncelleme için PUT
	http.HandleFunc("/api/projects/update", projects.UpdateProject) // projects güncelleme için PUT
	http.HandleFunc("/api/contact/update", contact.UpdateContact)   // contact güncelleme için PUT
	http.HandleFunc("/api/user/update", user.UpdateUser)            // user güncelleme için PUT

	// POST endpointlerini bağla
	http.HandleFunc("/api/about/create", about.CreateAbout)
	http.HandleFunc("/api/contact/create", contact.CreateContact)
	http.HandleFunc("/api/home/create", home.CreateHome)
	http.HandleFunc("/api/projects/create", projects.CreateProject)
	http.HandleFunc("/api/users/create", user.CreateUser)

	// Sunucu başlatma
	fmt.Println("Sunucu 8080 portunda başladı...") // Konsola bilgi yazdır
	log.Fatal(http.ListenAndServe(":8080", nil))   // HTTP sunucusunu başlat, hata varsa logla ve kapat
}
