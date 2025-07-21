import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081/api/v1';

// Create axios instance with security measures
export const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000, // Prevent hanging requests
  headers: {
    'Content-Type': 'application/json',
  },
  // Security configurations
  withCredentials: false, // Prevent CSRF if not using cookies
  maxRedirects: 0, // Prevent redirect attacks
});

// Request interceptor to add auth token and security headers
api.interceptors.request.use(
  (config) => {
    // Add auth token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Add security headers
    config.headers['X-Requested-With'] = 'XMLHttpRequest';

    // Validate URL to prevent SSRF
    if (config.url && !config.url.startsWith('/') && !config.url.startsWith(API_BASE_URL)) {
      throw new Error('Invalid request URL');
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle auth errors and security
api.interceptors.response.use(
  (response) => {
    // Validate response content type for security
    const contentType = response.headers['content-type'];
    if (contentType && !contentType.includes('application/json')) {
      console.warn('Unexpected content type:', contentType);
    }
    return response;
  },
  (error) => {
    // Handle authentication errors
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
      return Promise.reject(new Error('Authentication required'));
    }

    // Handle rate limiting
    if (error.response?.status === 429) {
      return Promise.reject(new Error('Too many requests. Please try again later.'));
    }

    // Handle server errors
    if (error.response?.status >= 500) {
      return Promise.reject(new Error('Server error. Please try again later.'));
    }

    // Sanitize error messages to prevent XSS
    if (error.response?.data?.message) {
      const sanitizedMessage = error.response.data.message.replace(/<[^>]*>/g, '');
      error.response.data.message = sanitizedMessage;
    }

    return Promise.reject(error);
  }
);

// Types
export interface User {
  id: number;
  username: string;
  email: string;
  first_name: string;
  last_name: string;
  role: 'student' | 'teacher' | 'admin';
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Test {
  id: number;
  title: string;
  description: string;
  created_by: number;
  duration_minutes: number;
  total_marks: number;
  passing_marks: number;
  instructions: string;
  is_active: boolean;
  start_time?: string;
  end_time?: string;
  created_at: string;
  updated_at: string;
}

export interface Question {
  id: number;
  test_id: number;
  question_text: string;
  question_type: 'multiple_choice' | 'true_false' | 'short_answer';
  marks: number;
  order_index: number;
  options?: QuestionOption[];
  correct_answers?: CorrectAnswer[];
  created_at: string;
  updated_at: string;
}

export interface QuestionOption {
  id: number;
  question_id: number;
  option_text: string;
  is_correct: boolean;
  order_index: number;
  created_at: string;
}

export interface CorrectAnswer {
  id: number;
  question_id: number;
  answer_text: string;
  is_case_sensitive: boolean;
  created_at: string;
}

export interface TestSession {
  id: number;
  test_id: number;
  user_id: number;
  session_token: string;
  status: 'not_started' | 'in_progress' | 'completed' | 'submitted' | 'expired';
  started_at?: string;
  submitted_at?: string;
  expires_at: string;
  time_remaining?: number;
  current_question_index: number;
  created_at: string;
  updated_at: string;
  remaining_time_seconds?: number;
}

export interface UserAnswer {
  id: number;
  session_id: number;
  question_id: number;
  answer_text?: string;
  selected_option_id?: number;
  is_correct?: boolean;
  marks_awarded: number;
  answered_at: string;
}

export interface TestResult {
  id: number;
  session_id: number;
  test_id: number;
  user_id: number;
  total_questions: number;
  answered_questions: number;
  correct_answers: number;
  total_marks: number;
  marks_obtained: number;
  percentage: number;
  grade?: string;
  is_passed: boolean;
  time_taken?: number;
  completed_at: string;
}

export interface TestStatistics {
  test_id: number;
  total_attempts: number;
  completed_attempts: number;
  passed_attempts: number;
  average_score: number;
  highest_score: number;
  lowest_score: number;
  average_time_taken: number;
}

// API Response wrapper
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: string;
  message?: string;
}

// Auth API
export const authApi = {
  login: (username: string, password: string) =>
    api.post<{ token: string; user: User }>('/auth/login', { username, password }),
  
  register: (data: {
    username: string;
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    role: string;
  }) => api.post<{ token: string; user: User }>('/auth/register', data),
  
  getProfile: () => api.get<User>('/auth/profile'),
  
  refreshToken: () => api.post<{ token: string }>('/auth/refresh'),
};

// Tests API
export const testsApi = {
  getAvailable: () => api.get<ApiResponse<Test[]>>('/tests/available'),
  getAll: (params?: { creator?: number; limit?: number; offset?: number }) =>
    api.get<ApiResponse<Test[]>>('/tests', { params }),
  getById: (id: number) => api.get<ApiResponse<Test>>(`/tests/${id}`),
  create: (data: Partial<Test>) => api.post<ApiResponse<Test>>('/tests', data),
  update: (id: number, data: Partial<Test>) => api.put<ApiResponse<Test>>(`/tests/${id}`, data),
  delete: (id: number) => api.delete(`/tests/${id}`),
  getQuestions: (id: number) => api.get<ApiResponse<Question[]>>(`/tests/${id}/questions`),
};

// Questions API
export const questionsApi = {
  create: (data: {
    test_id: number;
    question_text: string;
    question_type: string;
    marks: number;
    order_index: number;
    options?: Array<{ option_text: string; is_correct: boolean; order_index: number }>;
    answers?: Array<{ answer_text: string; is_case_sensitive: boolean }>;
  }) => api.post<ApiResponse<Question>>('/questions', data),
  getById: (id: number) => api.get<ApiResponse<Question>>(`/questions/${id}`),
  update: (id: number, data: Partial<Question>) => api.put<ApiResponse<Question>>(`/questions/${id}`, data),
  delete: (id: number) => api.delete(`/questions/${id}`),
};

// Sessions API
export const sessionsApi = {
  start: (test_id: number) => api.post<ApiResponse<TestSession>>('/sessions/start', { test_id }),
  get: (token: string) => api.get<ApiResponse<TestSession>>(`/sessions/${token}`),
  submitAnswer: (token: string, data: {
    question_id: number;
    answer_text?: string;
    selected_option_id?: number;
  }) => api.post<ApiResponse<UserAnswer>>(`/sessions/${token}/answers`, data),
  getAnswers: (token: string) => api.get<ApiResponse<UserAnswer[]>>(`/sessions/${token}/answers`),
  submit: (token: string) => api.post<ApiResponse<TestSession>>(`/sessions/${token}/submit`),
  updateProgress: (token: string, current_question_index: number) =>
    api.put(`/sessions/${token}/progress`, { current_question_index }),
  getMy: () => api.get<ApiResponse<TestSession[]>>('/sessions/my'),
};

// Results API
export const resultsApi = {
  getMy: () => api.get<ApiResponse<TestResult[]>>('/results/my'),
  getById: (id: number) => api.get<ApiResponse<TestResult>>(`/results/${id}`),
  getBySession: (sessionId: number) => api.get<ApiResponse<TestResult>>(`/results/session/${sessionId}`),
  getTestResults: (testId: number) => api.get<ApiResponse<TestResult[]>>(`/results/test/${testId}`),
  getTestStatistics: (testId: number) => api.get<ApiResponse<TestStatistics>>(`/results/test/${testId}/statistics`),
};
