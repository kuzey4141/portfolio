-- HOME tablosu
CREATE TABLE IF NOT EXISTS home (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT
);

INSERT INTO home (title, description) VALUES 
('Hoşgeldiniz', 'Bu benim portföy sitem.');

-- CONTACT tablosu
CREATE TABLE IF NOT EXISTS contact (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    message TEXT
);

INSERT INTO contact (email, phone, message) VALUES 
('kuzey@example.com', '555-1234', 'İletişim mesajı.');

-- PROJECTS tablosu
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    message TEXT
);

INSERT INTO projects (name, description, message) VALUES 
('Projex 1', 'İlk proje açıklaması.', 'Bu proje hakkında mesaj.');

-- ABOUT tablosu
CREATE TABLE IF NOT EXISTS about (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL
);

INSERT INTO about (content) VALUES 
('Bu benim hakkımda bölümüm.');

-- USERS tablosu (Admin paneli için)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE
);

INSERT INTO users (username, password, email) VALUES 
('admin', 'hashed_password_here', 'admin@example.com');
