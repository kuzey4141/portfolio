// backend/db/connection.go
package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5" // PostgreSQL kütüphanesi
)

var Conn *pgx.Conn // Global veritabanı bağlantısı

// ConnectDB fonksiyonu bağlantıyı kurar
func ConnectDB() {
	connStr := "postgres://postgres:6303523@localhost:5432/portfolio?sslmode=disable" // Bağlantı stringi

	var err error
	Conn, err = pgx.Connect(context.Background(), connStr) // Veritabanına bağlan
	if err != nil {
		log.Fatalf("Veritabanına bağlanırken hata: %v", err) // Hata varsa logla ve çık
	}
}
