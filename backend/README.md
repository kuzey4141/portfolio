# Portfolio Website - Buğra Kuzey Deveci

A modern full-stack portfolio website built with Go backend and React frontend, deployed on AWS.

## About Me

I'm Buğra Kuzey Deveci, a 22-year-old senior Management Information System student at Düzce University, based in Kocaeli, Turkey. I'm passionate about developing innovative digital solutions using modern web technologies, with expertise in full-stack development including React, TypeScript, Go, and PostgreSQL.

## Live Demo
- Website: http://3.78.181.203
- Admin Panel: http://3.78.181.203/admin

## Features

### Frontend (React + TypeScript)
- Responsive portfolio website
- Admin dashboard for content management
- Dynamic project showcase with image uploads
- Contact form with validation
- Dark/Light mode support

### Backend (Go + Gin Framework)
- RESTful API with JWT authentication
- PostgreSQL database with connection pooling
- Email notifications via Resend API
- CRUD operations for portfolio content
- Production-ready with static file serving

### Admin Panel
- Project management (add/edit/delete)
- Contact message handling
- Home and About page editing
- Secure authentication system

## Technology Stack

**Frontend:**
- React 18 with TypeScript
- SCSS for styling
- Material-UI components
- Lucide React icons

**Backend:**
- Go 1.21+ with Gin framework
- PostgreSQL database
- JWT authentication
- bcrypt password hashing
- Resend API for emails

**Deployment:**
- AWS EC2 (Ubuntu 22.04)
- AWS RDS (PostgreSQL)
- Nginx reverse proxy
- Environment-based configuration

## Local Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL

### Setup
```bash
# Clone repository
git clone <your-repo>
cd portfolio

# Backend setup
cd backend
go mod tidy
cp .env.example .env  # Configure your database and API keys

# Frontend setup
cd ../frontend
npm install

# Start development (single command)
cd ../backend
go run main.go  # Starts both backend and frontend
```

### Environment Variables
```
DB_HOST=your-postgres-host
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=portfolio
RESEND_API_KEY=your-resend-key
TO_EMAIL=your-email@domain.com
PORT=8081
GIN_MODE=release
```

## Production Deployment

### AWS Setup
1. Launch EC2 instance (Ubuntu 22.04)
2. Create RDS PostgreSQL database
3. Configure security groups (ports 22, 80, 8081)

### Deployment Steps
```bash
# Install dependencies
sudo apt update
sudo apt install nginx golang-go git -y

# Clone and build
git clone <your-repo>
cd portfolio/backend
go build -o portfolio-app .

# Configure nginx
sudo nano /etc/nginx/sites-available/default
# Add reverse proxy configuration

# Start services
export GIN_MODE=release
./portfolio-app > app.log 2>&1 &
sudo systemctl restart nginx
```

### Database Schema
The application creates these tables:
- `users` - Admin authentication
- `home` - Homepage content
- `about` - About page content
- `projects` - Portfolio projects
- `contact` - Contact form submissions

## API Endpoints

### Public Routes
- `GET /api/home` - Homepage data
- `GET /api/about` - About page data
- `GET /api/projects` - Project listings
- `POST /api/contact` - Submit contact form
- `POST /api/login` - Admin authentication

### Admin Routes (JWT Required)
- `GET|POST|PUT|DELETE /api/admin/projects` - Project management
- `GET|DELETE /api/admin/contact` - Contact management
- `PUT /api/admin/home` - Homepage updates
- `PUT /api/admin/about` - About page updates

## Admin Access
- Username: admin
- Password: admin123
- Access: http://your-domain/admin

## Project Structure
```
portfolio/
├── backend/
│   ├── main.go              # Application entry point
│   ├── server/              # Server utilities
│   ├── routes/              # API route handlers
│   ├── db/                  # Database connection
│   ├── auth/                # JWT authentication
│   ├── middleware/          # HTTP middleware
│   ├── home/                # Home page handlers
│   ├── about/               # About page handlers
│   ├── projects/            # Project handlers
│   ├── contact/             # Contact handlers
│   ├── user/                # User management
│   └── mail/                # Email service
└── frontend/
    ├── src/
    │   ├── components/      # React components
    │   ├── services/        # API services
    │   └── assets/          # Static assets
    └── public/
```

## License
MIT License
```