# GoCBT API Documentation

This document provides comprehensive documentation for the GoCBT REST API.

## 📋 Overview

The GoCBT API is a RESTful web service that provides endpoints for managing users, tests, questions, sessions, and results. All endpoints return JSON responses and use standard HTTP status codes.

### Base URL
```
http://localhost:8081/api/v1
```

### Authentication
Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### Response Format
All API responses follow this structure:
```json
{
  "success": true,
  "data": {},
  "message": "Success message",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error Response Format
```json
{
  "success": false,
  "error": "Error message",
  "code": "ERROR_CODE",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 🔐 Authentication Endpoints

### POST /auth/register
Register a new user account.

**Request Body:**
```json
{
  "username": "student1",
  "email": "student1@example.com",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe",
  "role": "student"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "student1",
      "email": "student1@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "student",
      "created_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

### POST /auth/login
Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "username": "student1",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "student1",
      "email": "student1@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "student"
    }
  }
}
```

### POST /auth/logout
Logout user (invalidate token).

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

## 👥 User Management Endpoints

### GET /users
Get list of users (Admin only).

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `role` (optional): Filter by role (admin, teacher, student)

**Response:**
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": 1,
        "username": "student1",
        "email": "student1@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "role": "student",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 25,
      "pages": 3
    }
  }
}
```

### GET /users/{id}
Get user by ID.

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "student1",
    "email": "student1@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "student",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### PUT /users/{id}
Update user information.

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Smith",
  "email": "john.smith@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "student1",
    "email": "john.smith@example.com",
    "first_name": "John",
    "last_name": "Smith",
    "role": "student",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

## 📝 Test Management Endpoints

### GET /tests
Get list of tests.

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page
- `status` (optional): Filter by status (active, inactive)

**Response:**
```json
{
  "success": true,
  "data": {
    "tests": [
      {
        "id": 1,
        "title": "Mathematics Quiz",
        "description": "Basic algebra and geometry",
        "instructions": "Answer all questions carefully",
        "duration_minutes": 60,
        "total_marks": 100,
        "passing_marks": 60,
        "start_time": "2024-01-15T09:00:00Z",
        "end_time": "2024-01-15T17:00:00Z",
        "created_by": 2,
        "created_at": "2024-01-14T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 5,
      "pages": 1
    }
  }
}
```

### GET /tests/available
Get available tests for current user (Student only).

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Mathematics Quiz",
      "description": "Basic algebra and geometry",
      "duration_minutes": 60,
      "total_marks": 100,
      "start_time": "2024-01-15T09:00:00Z",
      "end_time": "2024-01-15T17:00:00Z"
    }
  ]
}
```

### POST /tests
Create a new test (Teacher/Admin only).

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "Science Test",
  "description": "Physics and Chemistry basics",
  "instructions": "Read each question carefully",
  "duration_minutes": 90,
  "total_marks": 150,
  "passing_marks": 90,
  "start_time": "2024-01-16T09:00:00Z",
  "end_time": "2024-01-16T17:00:00Z"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 2,
    "title": "Science Test",
    "description": "Physics and Chemistry basics",
    "instructions": "Read each question carefully",
    "duration_minutes": 90,
    "total_marks": 150,
    "passing_marks": 90,
    "start_time": "2024-01-16T09:00:00Z",
    "end_time": "2024-01-16T17:00:00Z",
    "created_by": 2,
    "created_at": "2024-01-15T12:00:00Z"
  }
}
```

### GET /tests/{id}
Get test details by ID.

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Mathematics Quiz",
    "description": "Basic algebra and geometry",
    "instructions": "Answer all questions carefully",
    "duration_minutes": 60,
    "total_marks": 100,
    "passing_marks": 60,
    "start_time": "2024-01-15T09:00:00Z",
    "end_time": "2024-01-15T17:00:00Z",
    "questions": [
      {
        "id": 1,
        "question_text": "What is 2 + 2?",
        "question_type": "multiple_choice",
        "marks": 5,
        "options": [
          {"id": 1, "option_text": "3", "is_correct": false},
          {"id": 2, "option_text": "4", "is_correct": true},
          {"id": 3, "option_text": "5", "is_correct": false}
        ]
      }
    ]
  }
}
```

## ❓ Question Management Endpoints

### GET /questions
Get list of questions.

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `test_id` (optional): Filter by test ID
- `type` (optional): Filter by question type

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "test_id": 1,
      "question_text": "What is the capital of France?",
      "question_type": "multiple_choice",
      "marks": 5,
      "options": [
        {"id": 1, "option_text": "London", "is_correct": false},
        {"id": 2, "option_text": "Paris", "is_correct": true},
        {"id": 3, "option_text": "Berlin", "is_correct": false}
      ]
    }
  ]
}
```

### POST /questions
Create a new question (Teacher/Admin only).

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "test_id": 1,
  "question_text": "What is 5 × 6?",
  "question_type": "multiple_choice",
  "marks": 5,
  "options": [
    {"option_text": "25", "is_correct": false},
    {"option_text": "30", "is_correct": true},
    {"option_text": "35", "is_correct": false}
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 2,
    "test_id": 1,
    "question_text": "What is 5 × 6?",
    "question_type": "multiple_choice",
    "marks": 5,
    "options": [
      {"id": 4, "option_text": "25", "is_correct": false},
      {"id": 5, "option_text": "30", "is_correct": true},
      {"id": 6, "option_text": "35", "is_correct": false}
    ]
  }
}
```

## 🎯 Test Session Endpoints

### POST /sessions/start
Start a new test session.

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "test_id": 1
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "session_token": "sess_abc123def456",
    "test": {
      "id": 1,
      "title": "Mathematics Quiz",
      "duration_minutes": 60,
      "total_marks": 100
    },
    "started_at": "2024-01-15T14:00:00Z",
    "expires_at": "2024-01-15T15:00:00Z"
  }
}
```

### GET /sessions/{token}
Get current session details.

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": {
    "session_token": "sess_abc123def456",
    "test_id": 1,
    "user_id": 1,
    "status": "active",
    "started_at": "2024-01-15T14:00:00Z",
    "current_question": 3,
    "time_remaining": 3420,
    "questions": [
      {
        "id": 1,
        "question_text": "What is 2 + 2?",
        "question_type": "multiple_choice",
        "marks": 5,
        "options": [
          {"id": 1, "option_text": "3"},
          {"id": 2, "option_text": "4"},
          {"id": 3, "option_text": "5"}
        ]
      }
    ]
  }
}
```

### POST /sessions/{token}/submit-answer
Submit answer for a question.

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "question_id": 1,
  "selected_option_id": 2
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "question_id": 1,
    "submitted": true,
    "next_question_id": 2
  }
}
```

### POST /sessions/{token}/submit
Submit the entire test session.

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": {
    "session_token": "sess_abc123def456",
    "submitted_at": "2024-01-15T14:45:00Z",
    "result_id": 1,
    "score": 85.5,
    "total_marks": 100,
    "percentage": 85.5,
    "is_passed": true
  }
}
```

### GET /sessions/my
Get user's test sessions.

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "test_id": 1,
      "session_token": "sess_abc123def456",
      "status": "completed",
      "started_at": "2024-01-15T14:00:00Z",
      "submitted_at": "2024-01-15T14:45:00Z"
    }
  ]
}
```

## 📊 Results Endpoints

### GET /results/my
Get user's test results.

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "test_id": 1,
      "user_id": 1,
      "session_token": "sess_abc123def456",
      "score": 85.5,
      "total_marks": 100,
      "percentage": 85.5,
      "grade": "B",
      "is_passed": true,
      "completed_at": "2024-01-15T14:45:00Z"
    }
  ]
}
```

### GET /results/{id}
Get detailed result by ID.

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "test_id": 1,
    "user_id": 1,
    "score": 85.5,
    "total_marks": 100,
    "percentage": 85.5,
    "grade": "B",
    "is_passed": true,
    "completed_at": "2024-01-15T14:45:00Z",
    "answers": [
      {
        "question_id": 1,
        "selected_option_id": 2,
        "is_correct": true,
        "marks_awarded": 5
      }
    ]
  }
}
```

### GET /results/test/{test_id}
Get all results for a specific test (Teacher/Admin only).

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user": {
        "id": 1,
        "username": "student1",
        "first_name": "John",
        "last_name": "Doe"
      },
      "score": 85.5,
      "percentage": 85.5,
      "grade": "B",
      "is_passed": true,
      "completed_at": "2024-01-15T14:45:00Z"
    }
  ]
}
```

## 📈 Analytics Endpoints

### GET /analytics/dashboard
Get dashboard analytics (Teacher/Admin only).

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": {
    "total_tests": 5,
    "total_students": 25,
    "active_sessions": 3,
    "completed_tests": 45,
    "average_score": 78.5,
    "recent_activity": [
      {
        "type": "test_completed",
        "user": "student1",
        "test": "Mathematics Quiz",
        "timestamp": "2024-01-15T14:45:00Z"
      }
    ]
  }
}
```

## 🚨 Error Codes

| Code | Description |
|------|-------------|
| `INVALID_REQUEST` | Request body is malformed |
| `UNAUTHORIZED` | Authentication required |
| `FORBIDDEN` | Insufficient permissions |
| `NOT_FOUND` | Resource not found |
| `VALIDATION_ERROR` | Input validation failed |
| `DUPLICATE_ENTRY` | Resource already exists |
| `SESSION_EXPIRED` | Test session has expired |
| `TEST_NOT_AVAILABLE` | Test is not currently available |
| `INTERNAL_ERROR` | Server error |

## 📝 Rate Limiting

API endpoints are rate limited to prevent abuse:
- **Authentication endpoints**: 5 requests per minute
- **General endpoints**: 100 requests per minute
- **File upload endpoints**: 10 requests per minute

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642248000
```

## 🔒 Security Headers

All API responses include security headers:
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000
```

## 📚 SDKs and Examples

### JavaScript/Node.js Example
```javascript
const axios = require('axios');

const api = axios.create({
  baseURL: 'http://localhost:8081/api/v1',
  headers: {
    'Content-Type': 'application/json'
  }
});

// Login
const login = async (username, password) => {
  const response = await api.post('/auth/login', {
    username,
    password
  });

  // Set token for future requests
  api.defaults.headers.Authorization = `Bearer ${response.data.data.token}`;

  return response.data;
};

// Get available tests
const getAvailableTests = async () => {
  const response = await api.get('/tests/available');
  return response.data.data;
};
```

### Python Example
```python
import requests

class GoCBTAPI:
    def __init__(self, base_url='http://localhost:8081/api/v1'):
        self.base_url = base_url
        self.session = requests.Session()

    def login(self, username, password):
        response = self.session.post(f'{self.base_url}/auth/login', json={
            'username': username,
            'password': password
        })

        if response.status_code == 200:
            token = response.json()['data']['token']
            self.session.headers.update({
                'Authorization': f'Bearer {token}'
            })

        return response.json()

    def get_available_tests(self):
        response = self.session.get(f'{self.base_url}/tests/available')
        return response.json()['data']
```
