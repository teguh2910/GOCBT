// Frontend validation utilities

export interface ValidationResult {
  isValid: boolean;
  errors: string[];
}

// Email validation
export const validateEmail = (email: string): boolean => {
  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return emailRegex.test(email.trim());
};

// Username validation
export const validateUsername = (username: string): boolean => {
  const trimmed = username.trim();
  if (trimmed.length < 3 || trimmed.length > 50) {
    return false;
  }
  const usernameRegex = /^[a-zA-Z0-9_]+$/;
  return usernameRegex.test(trimmed);
};

// Password strength validation
export const validatePassword = (password: string): ValidationResult => {
  const errors: string[] = [];
  
  if (password.length < 8) {
    errors.push('Password must be at least 8 characters long');
  }
  
  if (password.length > 128) {
    errors.push('Password must be no more than 128 characters long');
  }
  
  if (!/[A-Z]/.test(password)) {
    errors.push('Password must contain at least one uppercase letter');
  }
  
  if (!/[a-z]/.test(password)) {
    errors.push('Password must contain at least one lowercase letter');
  }
  
  if (!/[0-9]/.test(password)) {
    errors.push('Password must contain at least one digit');
  }
  
  if (!/[^a-zA-Z0-9]/.test(password)) {
    errors.push('Password must contain at least one special character');
  }
  
  return {
    isValid: errors.length === 0,
    errors
  };
};

// Text length validation
export const validateTextLength = (text: string, minLength: number, maxLength: number): boolean => {
  const length = text.trim().length;
  return length >= minLength && length <= maxLength;
};

// Required field validation
export const validateRequired = (value: string): boolean => {
  return value.trim().length > 0;
};

// Number validation
export const validateNumber = (value: string, min?: number, max?: number): boolean => {
  const num = parseFloat(value);
  if (isNaN(num)) {
    return false;
  }
  
  if (min !== undefined && num < min) {
    return false;
  }
  
  if (max !== undefined && num > max) {
    return false;
  }
  
  return true;
};

// Integer validation
export const validateInteger = (value: string, min?: number, max?: number): boolean => {
  const num = parseInt(value, 10);
  if (isNaN(num) || !Number.isInteger(num)) {
    return false;
  }
  
  if (min !== undefined && num < min) {
    return false;
  }
  
  if (max !== undefined && num > max) {
    return false;
  }
  
  return true;
};

// Sanitize HTML to prevent XSS
export const sanitizeHTML = (input: string): string => {
  const div = document.createElement('div');
  div.textContent = input;
  return div.innerHTML;
};

// Validate and sanitize form data
export const sanitizeFormData = (data: Record<string, any>): Record<string, any> => {
  const sanitized: Record<string, any> = {};
  
  for (const [key, value] of Object.entries(data)) {
    if (typeof value === 'string') {
      // Trim whitespace and sanitize HTML
      sanitized[key] = sanitizeHTML(value.trim());
    } else {
      sanitized[key] = value;
    }
  }
  
  return sanitized;
};

// Check for potentially dangerous patterns
export const containsDangerousPatterns = (input: string): boolean => {
  const dangerousPatterns = [
    /<script/i,
    /javascript:/i,
    /on\w+\s*=/i,
    /data:text\/html/i,
    /vbscript:/i,
    /<iframe/i,
    /<object/i,
    /<embed/i,
    /<link/i,
    /<meta/i,
  ];
  
  return dangerousPatterns.some(pattern => pattern.test(input));
};

// Validate file upload
export const validateFile = (file: File, allowedTypes: string[], maxSize: number): ValidationResult => {
  const errors: string[] = [];
  
  // Check file type
  if (!allowedTypes.includes(file.type)) {
    errors.push(`File type ${file.type} is not allowed`);
  }
  
  // Check file size
  if (file.size > maxSize) {
    errors.push(`File size must be less than ${Math.round(maxSize / 1024 / 1024)}MB`);
  }
  
  // Check file name for dangerous patterns
  if (containsDangerousPatterns(file.name)) {
    errors.push('File name contains invalid characters');
  }
  
  return {
    isValid: errors.length === 0,
    errors
  };
};

// Validate URL
export const validateURL = (url: string): boolean => {
  try {
    const urlObj = new URL(url);
    return ['http:', 'https:'].includes(urlObj.protocol);
  } catch {
    return false;
  }
};

// Form validation helper
export const validateForm = (
  data: Record<string, any>,
  rules: Record<string, (value: any) => ValidationResult | boolean>
): ValidationResult => {
  const errors: string[] = [];
  
  for (const [field, rule] of Object.entries(rules)) {
    const value = data[field];
    const result = rule(value);
    
    if (typeof result === 'boolean') {
      if (!result) {
        errors.push(`${field} is invalid`);
      }
    } else {
      if (!result.isValid) {
        errors.push(...result.errors.map(error => `${field}: ${error}`));
      }
    }
  }
  
  return {
    isValid: errors.length === 0,
    errors
  };
};

// Rate limiting helper for client-side
export class ClientRateLimit {
  private attempts: Map<string, number[]> = new Map();
  
  constructor(private maxAttempts: number, private windowMs: number) {}
  
  isAllowed(key: string): boolean {
    const now = Date.now();
    const attempts = this.attempts.get(key) || [];
    
    // Remove old attempts outside the window
    const validAttempts = attempts.filter(time => now - time < this.windowMs);
    
    if (validAttempts.length >= this.maxAttempts) {
      return false;
    }
    
    // Add current attempt
    validAttempts.push(now);
    this.attempts.set(key, validAttempts);
    
    return true;
  }
  
  reset(key: string): void {
    this.attempts.delete(key);
  }
}
