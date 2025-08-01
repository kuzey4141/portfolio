-- HOME tablosu
CREATE TABLE home (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT
);

INSERT INTO home (title, description) VALUES 
('Hoşgeldiniz', 'Bu benim portföy sitem.');

-- CONTACT tablosu
CREATE TABLE contact (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    Message string `json:"message"`
);

INSERT INTO contact (email, phone, message) VALUES 
('kuzey@example.com', '555-1234', 'İletişim mesajı.');

CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    message TEXT
);

INSERT INTO projects (name, description, message) VALUES 
('Projex 1', 'İlk proje açıklaması.', 'Bu proje hakkında mesaj.');


-- ABOUT tablosu
CREATE TABLE about (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL
);

INSERT INTO about (content) VALUES 
('Bu benim hakkımda bölümüm.');

-- USERS tablosu (Admin paneli için)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE
);

INSERT INTO users (username, password, email) VALUES 
('admin', 'hashed_password_here', 'admin@example.com');
