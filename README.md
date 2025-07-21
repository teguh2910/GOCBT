# GoCBT - Computer-Based Test Application

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

A modern, secure, and lightweight Computer-Based Test (CBT) application built with Go and React for schools and educational institutions. GoCBT provides a comprehensive platform for creating, managing, and conducting online examinations with real-time monitoring and automated scoring.

## ✨ Features

### 🎓 For Students
- **Secure Authentication**: JWT-based login with role-based access control
- **Intuitive Test Interface**: Clean, responsive design for optimal test-taking experience
- **Real-time Progress**: Live progress tracking and time remaining indicators
- **Auto-save**: Automatic answer saving to prevent data loss
- **Results Dashboard**: Immediate access to test results and performance analytics

### 👨‍🏫 For Teachers
- **Test Creation**: Easy-to-use interface for creating tests with multiple question types
- **Question Bank**: Comprehensive question management with reusable question pools
- **Real-time Monitoring**: Live monitoring of student progress during tests
- **Flexible Scheduling**: Set test availability windows and time limits
- **Detailed Analytics**: Comprehensive reports on student performance and test statistics

### 🔧 For Administrators
- **User Management**: Complete user administration with role assignments
- **System Monitoring**: Dashboard for system health and usage statistics
- **Security Controls**: Advanced security features and audit logging
- **Bulk Operations**: Import/export functionality for users and test data

## 🚀 Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gorilla Mux for routing
- **Database**: SQLite (development), PostgreSQL (production)
- **Authentication**: JWT with bcrypt password hashing
- **Security**: Comprehensive input validation, XSS prevention, SQL injection protection

### Frontend
- **Framework**: React 18 with Next.js 14
- **Styling**: Tailwind CSS with custom components
- **State Management**: React Context API
- **HTTP Client**: Axios with security interceptors
- **UI Components**: Custom component library with dark mode support

### DevOps & Deployment
- **Containerization**: Docker with multi-stage builds
- **Database Migrations**: Custom migration system
- **Environment Configuration**: Environment-based configuration
- **Security**: Rate limiting, security headers, CORS protection

## 📁 Project Structure

```
gocbt/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/                     # HTTP handlers and REST endpoints
│   ├── auth/                    # JWT authentication and middleware
│   ├── config/                  # Configuration management
│   ├── database/                # Database connection and repositories
│   ├── middleware/              # Security and utility middleware
│   ├── models/                  # Data models and interfaces
│   ├── services/                # Business logic layer
│   └── utils/                   # Utility functions and validation
├── frontend/                    # React/Next.js frontend application
│   ├── src/
│   │   ├── app/                 # Next.js app router pages
│   │   ├── components/          # Reusable React components
│   │   ├── contexts/            # React context providers
│   │   ├── lib/                 # Utility libraries and API client
│   │   └── styles/              # Global styles and Tailwind config
│   ├── public/                  # Static assets
│   └── package.json             # Frontend dependencies
├── migrations/                  # Database migration files
├── docs/                        # Comprehensive documentation
├── examples/                    # Example data and configurations
├── docker-compose.yml           # Docker development setup
├── Dockerfile                   # Production Docker image
├── go.mod                       # Go module dependencies
└── README.md                    # This file
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Node.js 18+** - [Download Node.js](https://nodejs.org/)
- **Git** - [Download Git](https://git-scm.com/)

### Option 1: Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/your-username/gocbt.git
cd gocbt

# Start with Docker Compose
docker-compose up -d

# Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

### Option 2: Local Development

```bash
# Clone the repository
git clone https://github.com/your-username/gocbt.git
cd gocbt

# Backend setup
go mod tidy
go run cmd/server/main.go

# Frontend setup (in a new terminal)
cd frontend
npm install
npm run dev

# Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8081
```

### Default Login Credentials

After setup, you can login with these default accounts:

| Role | Username | Password |
|------|----------|----------|
| Admin | admin | Admin123! |
| Teacher | teacher1 | Teacher123! |
| Student | student1 | Student123! |

## 📖 Documentation

- **[Setup Guide](docs/SETUP.md)** - Detailed installation and configuration
- **[API Documentation](docs/API.md)** - Complete REST API reference
- **[User Guide](docs/USER_GUIDE.md)** - How to use GoCBT for students and teachers
- **[Development Guide](docs/DEVELOPMENT.md)** - Contributing and development setup
- **[Security Guide](SECURITY.md)** - Security features and best practices

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the root directory:

```env
# Server Configuration
APP_ENV=development
SERVER_HOST=localhost
SERVER_PORT=8081

# Database Configuration
DB_DRIVER=sqlite
DB_FILEPATH=./gocbt.db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION=24h

# CORS Configuration
CORS_ORIGINS=http://localhost:3000

# Frontend Configuration
NEXT_PUBLIC_API_URL=http://localhost:8081/api/v1
```

### Database Setup

The application supports both SQLite and PostgreSQL:

#### SQLite (Development)
```env
DB_DRIVER=sqlite
DB_FILEPATH=./gocbt.db
```

#### PostgreSQL (Production)
```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=gocbt
DB_PASSWORD=your-password
DB_NAME=gocbt
DB_SSLMODE=disable
```

## 🧪 Testing

```bash
# Run backend tests
go test ./...

# Run frontend tests
cd frontend
npm test

# Run security tests
go test ./internal/utils -v
```

## 🚢 Deployment

### Docker Production Deployment

```bash
# Build production image
docker build -t gocbt:latest .

# Run with environment variables
docker run -d \
  -p 8080:8080 \
  -e APP_ENV=production \
  -e JWT_SECRET=your-production-secret \
  -e DB_DRIVER=postgres \
  -e DB_HOST=your-db-host \
  gocbt:latest
```

### Manual Deployment

```bash
# Build backend
go build -o gocbt cmd/server/main.go

# Build frontend
cd frontend
npm run build

# Deploy built files to your server
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](docs/CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Check our [docs](docs/) directory
- **Issues**: Report bugs on [GitHub Issues](https://github.com/your-username/gocbt/issues)
- **Discussions**: Join our [GitHub Discussions](https://github.com/your-username/gocbt/discussions)

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/) and [React](https://reactjs.org/)
- UI components inspired by modern design systems
- Security best practices from [OWASP](https://owasp.org/)

---

**Made with ❤️ for education**
