# GoCBT Development Guide

This guide provides comprehensive information for developers who want to contribute to or extend the GoCBT project.

## 🏗️ Architecture Overview

### System Architecture

GoCBT follows a modern three-tier architecture:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │    Database     │
│   (React)       │◄──►│     (Go)        │◄──►│  (SQLite/PG)    │
│                 │    │                 │    │                 │
│ • Next.js       │    │ • REST API      │    │ • User Data     │
│ • Tailwind CSS  │    │ • JWT Auth      │    │ • Test Data     │
│ • TypeScript    │    │ • Middleware    │    │ • Results       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Backend Architecture

```
cmd/server/main.go
├── Configuration Loading
├── Database Connection
├── Middleware Setup
├── Route Registration
└── Server Startup

internal/
├── api/           # HTTP handlers and routes
├── auth/          # Authentication and JWT
├── config/        # Configuration management
├── database/      # Database layer and repositories
├── middleware/    # HTTP middleware (CORS, Auth, etc.)
├── models/        # Data models and interfaces
├── services/      # Business logic layer
└── utils/         # Utility functions and validation
```

### Frontend Architecture

```
frontend/src/
├── app/           # Next.js app router pages
├── components/    # Reusable React components
├── contexts/      # React context providers
├── lib/           # Utility libraries and API client
└── styles/        # Global styles and Tailwind config
```

## 🛠️ Development Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- Git
- Docker (optional)

### Local Development

1. **Clone Repository**
```bash
git clone https://github.com/your-username/gocbt.git
cd gocbt
```

2. **Backend Setup**
```bash
# Install dependencies
go mod tidy

# Copy environment file
cp .env.example .env

# Run backend
go run cmd/server/main.go
```

3. **Frontend Setup**
```bash
# Navigate to frontend
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

### Development Tools

#### Recommended VS Code Extensions
- Go (official Go extension)
- ES7+ React/Redux/React-Native snippets
- Tailwind CSS IntelliSense
- Prettier - Code formatter
- GitLens

#### Go Tools
```bash
# Install useful Go tools
go install github.com/cosmtrek/air@latest          # Live reload
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Linting
go install github.com/swaggo/swag/cmd/swag@latest  # API documentation
```

## 📊 Database Schema

### Core Tables

#### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'teacher', 'student')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Tests Table
```sql
CREATE TABLE tests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    instructions TEXT,
    duration_minutes INTEGER NOT NULL,
    total_marks INTEGER NOT NULL,
    passing_marks INTEGER NOT NULL,
    start_time DATETIME,
    end_time DATETIME,
    created_by INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id)
);
```

#### Questions Table
```sql
CREATE TABLE questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    test_id INTEGER NOT NULL,
    question_text TEXT NOT NULL,
    question_type VARCHAR(20) NOT NULL CHECK (question_type IN ('multiple_choice', 'true_false', 'short_answer')),
    marks INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE
);
```

### Relationships

```
Users (1) ──── (N) Tests (created_by)
Tests (1) ──── (N) Questions
Tests (1) ──── (N) TestSessions
TestSessions (1) ──── (N) SessionAnswers
TestSessions (1) ──── (1) TestResults
```

## 🔧 API Development

### Adding New Endpoints

1. **Define the Handler**
```go
// internal/api/new_handler.go
func (h *NewHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
    // Validate input
    var req CreateResourceRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteErrorResponse(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Business logic
    resource, err := h.service.CreateResource(req)
    if err != nil {
        utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return response
    utils.WriteSuccessResponse(w, resource)
}
```

2. **Register Routes**
```go
// cmd/server/main.go
func setupRoutes(handlers...) *mux.Router {
    router := mux.NewRouter()
    
    // Add new routes
    router.HandleFunc("/api/v1/resources", handler.CreateResource).Methods("POST")
    router.HandleFunc("/api/v1/resources/{id}", handler.GetResource).Methods("GET")
    
    return router
}
```

3. **Add Service Layer**
```go
// internal/services/new_service.go
type NewService struct {
    repo NewRepository
}

func (s *NewService) CreateResource(req CreateResourceRequest) (*Resource, error) {
    // Validation
    if err := validateResource(req); err != nil {
        return nil, err
    }

    // Business logic
    resource := &Resource{
        Name: req.Name,
        // ... other fields
    }

    // Save to database
    return s.repo.Create(resource)
}
```

### Request/Response Patterns

#### Standard Request Structure
```go
type CreateResourceRequest struct {
    Name        string `json:"name" validate:"required,min=1,max=100"`
    Description string `json:"description" validate:"max=500"`
}
```

#### Standard Response Structure
```go
type APIResponse struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data,omitempty"`
    Error     string      `json:"error,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
}
```

## 🎨 Frontend Development

### Component Structure

```typescript
// components/ComponentName.tsx
import React from 'react';

interface ComponentNameProps {
  title: string;
  onAction?: () => void;
}

export const ComponentName: React.FC<ComponentNameProps> = ({
  title,
  onAction
}) => {
  return (
    <div className="component-container">
      <h2 className="text-xl font-semibold">{title}</h2>
      {onAction && (
        <button onClick={onAction} className="btn btn-primary">
          Action
        </button>
      )}
    </div>
  );
};
```

### State Management

#### Using React Context
```typescript
// contexts/AppContext.tsx
interface AppContextType {
  user: User | null;
  setUser: (user: User | null) => void;
}

export const AppContext = createContext<AppContextType | undefined>(undefined);

export const useApp = () => {
  const context = useContext(AppContext);
  if (!context) {
    throw new Error('useApp must be used within AppProvider');
  }
  return context;
};
```

### API Integration

```typescript
// lib/api.ts
export const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  timeout: 10000,
});

// Add request interceptor for auth
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

## 🧪 Testing

### Backend Testing

#### Unit Tests
```go
// internal/services/test_service_test.go
func TestCreateTest(t *testing.T) {
    // Setup
    mockRepo := &MockTestRepository{}
    service := NewTestService(mockRepo)

    // Test data
    req := CreateTestRequest{
        Title:           "Test Title",
        DurationMinutes: 60,
        TotalMarks:      100,
    }

    // Execute
    test, err := service.CreateTest(1, req)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "Test Title", test.Title)
    assert.Equal(t, 60, test.DurationMinutes)
}
```

#### Integration Tests
```go
func TestCreateTestIntegration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()

    // Create service with real repository
    repo := NewTestRepository(db)
    service := NewTestService(repo)

    // Test the full flow
    test, err := service.CreateTest(1, validRequest)
    
    assert.NoError(t, err)
    assert.NotZero(t, test.ID)
}
```

### Frontend Testing

#### Component Tests
```typescript
// components/__tests__/ComponentName.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { ComponentName } from '../ComponentName';

describe('ComponentName', () => {
  it('renders title correctly', () => {
    render(<ComponentName title="Test Title" />);
    expect(screen.getByText('Test Title')).toBeInTheDocument();
  });

  it('calls onAction when button is clicked', () => {
    const mockAction = jest.fn();
    render(<ComponentName title="Test" onAction={mockAction} />);
    
    fireEvent.click(screen.getByText('Action'));
    expect(mockAction).toHaveBeenCalledTimes(1);
  });
});
```

### Running Tests

```bash
# Backend tests
go test ./...
go test -cover ./...

# Frontend tests
cd frontend
npm test
npm run test:coverage
```

## 🚀 Deployment

### Docker Deployment

#### Multi-stage Dockerfile
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gocbt cmd/server/main.go

# Production stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gocbt .
EXPOSE 8080
CMD ["./gocbt"]
```

#### Docker Compose
```yaml
version: '3.8'
services:
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_DRIVER=postgres
      - DB_HOST=db
    depends_on:
      - db

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=gocbt
      - POSTGRES_USER=gocbt
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## 📝 Code Style Guidelines

### Go Code Style

#### Naming Conventions
- Use camelCase for variables and functions
- Use PascalCase for exported functions and types
- Use ALL_CAPS for constants
- Use descriptive names

#### Error Handling
```go
// Good
func CreateUser(req CreateUserRequest) (*User, error) {
    if err := validateUser(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    user, err := repo.Create(req)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}
```

### TypeScript Code Style

#### Interface Definitions
```typescript
// Use PascalCase for interfaces
interface UserProfile {
  id: number;
  username: string;
  email: string;
  createdAt: Date;
}

// Use descriptive prop names
interface ButtonProps {
  variant: 'primary' | 'secondary' | 'danger';
  size: 'sm' | 'md' | 'lg';
  isLoading?: boolean;
  onClick: () => void;
}
```

## 🤝 Contributing

### Pull Request Process

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Make your changes**
4. **Add tests for new functionality**
5. **Ensure all tests pass**
6. **Update documentation**
7. **Submit a pull request**

### Commit Message Format

```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Example:
```
feat(auth): add password reset functionality

- Add password reset endpoint
- Implement email notification
- Add frontend reset form

Closes #123
```

### Code Review Checklist

- [ ] Code follows style guidelines
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No security vulnerabilities
- [ ] Performance considerations addressed
- [ ] Error handling is appropriate

## 🔍 Debugging

### Backend Debugging

#### Using Delve Debugger
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug cmd/server/main.go
```

#### Logging
```go
import "log/slog"

// Structured logging
slog.Info("User created", 
    "user_id", user.ID, 
    "username", user.Username)

slog.Error("Database error", 
    "error", err, 
    "operation", "create_user")
```

### Frontend Debugging

#### React Developer Tools
- Install React Developer Tools browser extension
- Use Components and Profiler tabs
- Inspect component state and props

#### Console Debugging
```typescript
// Use console.group for organized logging
console.group('API Call');
console.log('Request:', request);
console.log('Response:', response);
console.groupEnd();

// Use console.table for arrays/objects
console.table(users);
```

## 📚 Additional Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [React Documentation](https://reactjs.org/docs/)
- [Next.js Documentation](https://nextjs.org/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)

### Tools and Libraries
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Testify](https://github.com/stretchr/testify) - Testing toolkit
- [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)

### Best Practices
- [Effective Go](https://golang.org/doc/effective_go.html)
- [React Best Practices](https://reactjs.org/docs/thinking-in-react.html)
- [API Design Guidelines](https://github.com/microsoft/api-guidelines)
- [Security Best Practices](https://owasp.org/www-project-top-ten/)
