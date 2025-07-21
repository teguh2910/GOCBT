# GoCBT Setup Guide

This guide provides detailed instructions for setting up GoCBT in different environments.

## 📋 Prerequisites

### System Requirements

- **Operating System**: Windows 10+, macOS 10.15+, or Linux (Ubuntu 18.04+)
- **Memory**: Minimum 4GB RAM (8GB recommended)
- **Storage**: At least 2GB free space
- **Network**: Internet connection for downloading dependencies

### Required Software

#### Go (Backend)
- **Version**: Go 1.21 or higher
- **Download**: [https://golang.org/dl/](https://golang.org/dl/)
- **Verification**: Run `go version` in terminal

#### Node.js (Frontend)
- **Version**: Node.js 18 or higher
- **Download**: [https://nodejs.org/](https://nodejs.org/)
- **Verification**: Run `node --version` and `npm --version`

#### Git
- **Download**: [https://git-scm.com/](https://git-scm.com/)
- **Verification**: Run `git --version`

#### Optional: Docker
- **Docker Desktop**: [https://www.docker.com/products/docker-desktop](https://www.docker.com/products/docker-desktop)
- **Verification**: Run `docker --version` and `docker-compose --version`

## 🐳 Docker Setup (Recommended)

Docker provides the easiest way to get GoCBT running with all dependencies.

### Step 1: Clone Repository

```bash
git clone https://github.com/your-username/gocbt.git
cd gocbt
```

### Step 2: Environment Configuration

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

### Step 3: Start Services

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Step 4: Access Application

- **Frontend**: [http://localhost:3000](http://localhost:3000)
- **Backend API**: [http://localhost:8081](http://localhost:8081)

## 💻 Local Development Setup

For development with hot reloading and debugging capabilities.

### Step 1: Clone and Setup Backend

```bash
# Clone repository
git clone https://github.com/your-username/gocbt.git
cd gocbt

# Install Go dependencies
go mod tidy

# Create environment file
cp .env.example .env
# Edit .env with your configuration
```

### Step 2: Database Setup

#### SQLite (Default)
```bash
# SQLite database will be created automatically
# No additional setup required
```

#### PostgreSQL (Optional)
```bash
# Install PostgreSQL
# Create database
createdb gocbt

# Update .env file
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=your-username
DB_PASSWORD=your-password
DB_NAME=gocbt
DB_SSLMODE=disable
```

### Step 3: Start Backend

```bash
# Run backend server
go run cmd/server/main.go

# Or with live reloading (install air first)
go install github.com/cosmtrek/air@latest
air
```

### Step 4: Setup Frontend

```bash
# Open new terminal and navigate to frontend
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

### Step 5: Access Application

- **Frontend**: [http://localhost:3000](http://localhost:3000)
- **Backend API**: [http://localhost:8081](http://localhost:8081)

## 🔧 Configuration Options

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `APP_ENV` | Application environment | `development` | No |
| `SERVER_HOST` | Server host | `localhost` | No |
| `SERVER_PORT` | Server port | `8081` | No |
| `DB_DRIVER` | Database driver (sqlite/postgres) | `sqlite` | No |
| `DB_FILEPATH` | SQLite database file path | `./gocbt.db` | SQLite only |
| `DB_HOST` | PostgreSQL host | `localhost` | PostgreSQL only |
| `DB_PORT` | PostgreSQL port | `5432` | PostgreSQL only |
| `DB_USER` | PostgreSQL username | - | PostgreSQL only |
| `DB_PASSWORD` | PostgreSQL password | - | PostgreSQL only |
| `DB_NAME` | PostgreSQL database name | - | PostgreSQL only |
| `JWT_SECRET` | JWT signing secret | - | Yes |
| `JWT_EXPIRATION` | JWT token expiration | `24h` | No |
| `CORS_ORIGINS` | Allowed CORS origins | `*` | No |

### Database Configuration

#### SQLite Configuration
```env
DB_DRIVER=sqlite
DB_FILEPATH=./data/gocbt.db
```

#### PostgreSQL Configuration
```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=gocbt
DB_PASSWORD=secure-password
DB_NAME=gocbt
DB_SSLMODE=require
```

### Frontend Configuration

Create `frontend/.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8081/api/v1
NEXT_PUBLIC_APP_NAME=GoCBT
NEXT_PUBLIC_APP_VERSION=1.0.0
```

## 🗄️ Database Migration

### Automatic Migration

The application automatically runs migrations on startup.

### Manual Migration

```bash
# Run specific migration
go run cmd/server/main.go -migrate

# Reset database (development only)
rm gocbt.db
go run cmd/server/main.go
```

## 🧪 Testing Setup

### Backend Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/auth
```

### Frontend Tests

```bash
cd frontend

# Run tests
npm test

# Run tests with coverage
npm run test:coverage

# Run E2E tests
npm run test:e2e
```

## 🚀 Production Deployment

### Environment Setup

```env
APP_ENV=production
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
JWT_SECRET=your-very-secure-production-secret
DB_DRIVER=postgres
DB_HOST=your-production-db-host
DB_USER=your-production-db-user
DB_PASSWORD=your-production-db-password
DB_NAME=gocbt
DB_SSLMODE=require
CORS_ORIGINS=https://your-domain.com
```

### Build and Deploy

```bash
# Build backend
go build -o gocbt cmd/server/main.go

# Build frontend
cd frontend
npm run build

# Deploy to your server
```

### Docker Production

```bash
# Build production image
docker build -t gocbt:latest .

# Run production container
docker run -d \
  -p 8080:8080 \
  --env-file .env.production \
  gocbt:latest
```

## 🔍 Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Find process using port
lsof -i :8081  # macOS/Linux
netstat -ano | findstr :8081  # Windows

# Kill process
kill -9 <PID>  # macOS/Linux
taskkill /PID <PID> /F  # Windows
```

#### Database Connection Issues
```bash
# Check database status
# For PostgreSQL
pg_isready -h localhost -p 5432

# Check database permissions
# Ensure user has CREATE, SELECT, INSERT, UPDATE, DELETE permissions
```

#### Frontend Build Issues
```bash
# Clear npm cache
npm cache clean --force

# Delete node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Getting Help

1. Check the [FAQ](FAQ.md)
2. Search [GitHub Issues](https://github.com/your-username/gocbt/issues)
3. Create a new issue with:
   - Operating system and version
   - Go and Node.js versions
   - Error messages and logs
   - Steps to reproduce

## 📚 Next Steps

After successful setup:

1. Read the [User Guide](USER_GUIDE.md)
2. Explore the [API Documentation](API.md)
3. Check out [Development Guide](DEVELOPMENT.md) for contributing
4. Review [Security Guide](../SECURITY.md) for production deployment
