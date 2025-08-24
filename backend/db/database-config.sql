
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE
);



-- HOME table - Main page content
CREATE TABLE IF NOT EXISTS home (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT
);



-- ABOUT table - About me section
CREATE TABLE IF NOT EXISTS about (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL
);



-- PROJECTS table - Portfolio projects with enhanced fields
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    message TEXT,
    image_url TEXT,
    technologies TEXT,
    github_url TEXT,
    demo_url TEXT
);



-- CONTACT table - Updated with all current fields
CREATE TABLE IF NOT EXISTS contact (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    message TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT false
);



-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_contact_created_at ON contact(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_contact_is_read ON contact(is_read);
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Set sequences to current values (if data already exists)
SELECT setval('users_id_seq', COALESCE((SELECT MAX(id) FROM users), 1));
SELECT setval('home_id_seq', COALESCE((SELECT MAX(id) FROM home), 1));
SELECT setval('about_id_seq', COALESCE((SELECT MAX(id) FROM about), 1));
SELECT setval('projects_id_seq', COALESCE((SELECT MAX(id) FROM projects), 1));
SELECT setval('contact_id_seq', COALESCE((SELECT MAX(id) FROM contact), 1));