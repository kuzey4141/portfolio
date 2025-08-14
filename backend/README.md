# Portfolio Projesi

## Backend API
- Go + Gin framework
- PostgreSQL veritabanı
- JWT kimlik doğrulama
- bcrypt şifre hashleme

## Özellikler
- **Public routes:** `/api/*` (portfolio içeriği - herkese açık)
- **Admin routes:** `/api/admin/*` (içerik yönetimi - sadece admin)
- JWT token kimlik doğrulama sistemi
- Güvenli şifre saklama (bcrypt hash)
- Middleware ile otomatik güvenlik kontrolü

## Kurulum
1. PostgreSQL kurulumu yapın
2. `backend/db/database-config.sql` dosyasını çalıştırın
3. Backend klasöründe `go run main.go` komutu ile sunucuyu başlatın

## Admin Girişi
- **Kullanıcı adı:** admin
- **Şifre:** admin123

## API Endpoint'leri

### Herkese Açık
- `GET /api/home` - Ana sayfa bilgileri
- `GET /api/about` - Hakkımda bilgileri
- `GET /api/projects` - Proje listesi
- `POST /api/contact` - İletişim formu
- `POST /api/login` - Giriş yapma

### Sadece Admin (Token Gerekli)
- `GET /api/admin/contact` - Gelen mesajları görüntüle
- `POST /api/admin/projects` - Yeni proje ekle
- `PUT /api/admin/home` - Ana sayfa güncelle
- `DELETE /api/admin/users/:id` - Kullanıcı sil

## Güvenlik
- Şifreler bcrypt ile hashlenmiş olarak saklanır
- JWT token'lar 24 saat geçerlidir
- Admin endpoint'leri middleware ile korunur
- SQL injection koruması vardır

## Teknolojiler
- **Backend:** Go, Gin Framework
- **Veritabanı:** PostgreSQL
- **Güvenlik:** JWT, bcrypt
- **API:** RESTful API yapısı