'use client';

import React, { createContext, useContext, useEffect, useState } from 'react';
import { User, authApi } from '@/lib/api';

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (username: string, password: string) => Promise<void>;
  register: (data: {
    username: string;
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    role: string;
  }) => Promise<void>;
  logout: () => void;
  loading: boolean;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check for existing token on mount
    const savedToken = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');

    if (savedToken && savedUser) {
      setToken(savedToken);
      setUser(JSON.parse(savedUser));
    }
    setLoading(false);
  }, []);

  const login = async (username: string, password: string) => {
    try {
      // Basic client-side validation
      if (!username.trim() || !password) {
        throw new Error('Username and password are required');
      }

      // Check for suspicious patterns
      if (username.includes('<') || username.includes('>') || username.includes('script')) {
        throw new Error('Invalid characters in username');
      }

      const response = await authApi.login(username.trim(), password);
      const { token: newToken, user: newUser } = response.data;

      // Validate token format (basic JWT structure check)
      if (!newToken || newToken.split('.').length !== 3) {
        throw new Error('Invalid authentication token');
      }

      setToken(newToken);
      setUser(newUser);
      localStorage.setItem('token', newToken);
      localStorage.setItem('user', JSON.stringify(newUser));
    } catch (error) {
      throw error;
    }
  };

  const register = async (data: {
    username: string;
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    role: string;
  }) => {
    try {
      // Sanitize and validate input data
      const sanitizedData = {
        username: data.username.trim(),
        email: data.email.trim().toLowerCase(),
        password: data.password,
        first_name: data.first_name.trim(),
        last_name: data.last_name.trim(),
        role: data.role
      };

      // Basic validation
      if (!sanitizedData.username || !sanitizedData.email || !sanitizedData.password ||
          !sanitizedData.first_name || !sanitizedData.last_name) {
        throw new Error('All fields are required');
      }

      // Check for suspicious patterns
      const textFields = [sanitizedData.username, sanitizedData.email,
                         sanitizedData.first_name, sanitizedData.last_name];
      for (const field of textFields) {
        if (field.includes('<') || field.includes('>') || field.includes('script')) {
          throw new Error('Invalid characters detected');
        }
      }

      const response = await authApi.register(sanitizedData);
      const { token: newToken, user: newUser } = response.data;

      // Validate token format
      if (!newToken || newToken.split('.').length !== 3) {
        throw new Error('Invalid authentication token');
      }

      setToken(newToken);
      setUser(newUser);
      localStorage.setItem('token', newToken);
      localStorage.setItem('user', JSON.stringify(newUser));
    } catch (error) {
      throw error;
    }
  };

  const logout = () => {
    setUser(null);
    setToken(null);
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  };

  const value = {
    user,
    token,
    login,
    register,
    logout,
    loading,
    isAuthenticated: !!user && !!token,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
