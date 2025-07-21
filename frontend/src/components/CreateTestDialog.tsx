'use client';

import React, { useState } from 'react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { testsApi, Test } from '@/lib/api';
import { X, Calendar, Clock, FileText, Users } from 'lucide-react';

interface CreateTestDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onTestCreated: (test: Test) => void;
}

export function CreateTestDialog({ isOpen, onClose, onTestCreated }: CreateTestDialogProps) {
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    duration_minutes: 60,
    total_marks: 100,
    passing_marks: 50,
    instructions: '',
    start_time: '',
    end_time: '',
    is_active: true,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value, type } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'number' ? parseInt(value) || 0 : value,
    }));
  };

  const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: checked,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      // Prepare data for API
      const testData: any = {
        title: formData.title,
        description: formData.description,
        duration_minutes: formData.duration_minutes,
        total_marks: formData.total_marks,
        passing_marks: formData.passing_marks,
        instructions: formData.instructions,
        is_active: formData.is_active,
      };

      // Add optional date fields if provided
      if (formData.start_time) {
        testData.start_time = new Date(formData.start_time).toISOString();
      }
      if (formData.end_time) {
        testData.end_time = new Date(formData.end_time).toISOString();
      }

      const response = await testsApi.create(testData);
      const newTest = response.data.data;
      
      onTestCreated(newTest);
      onClose();
      
      // Reset form
      setFormData({
        title: '',
        description: '',
        duration_minutes: 60,
        total_marks: 100,
        passing_marks: 50,
        instructions: '',
        start_time: '',
        end_time: '',
        is_active: true,
      });
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create test. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Create New Test</h2>
          <Button variant="ghost" size="sm" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          {error && (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded">
              {error}
            </div>
          )}

          {/* Basic Information */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">Basic Information</h3>
            
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Test Title *
              </label>
              <Input
                id="title"
                name="title"
                type="text"
                value={formData.title}
                onChange={handleChange}
                required
                placeholder="Enter test title"
              />
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Description *
              </label>
              <textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleChange}
                required
                rows={3}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                placeholder="Enter test description"
              />
            </div>

            <div>
              <label htmlFor="instructions" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Instructions
              </label>
              <textarea
                id="instructions"
                name="instructions"
                value={formData.instructions}
                onChange={handleChange}
                rows={4}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                placeholder="Enter test instructions for students"
              />
            </div>
          </div>

          {/* Test Configuration */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">Test Configuration</h3>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label htmlFor="duration_minutes" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  <Clock className="h-4 w-4 inline mr-1" />
                  Duration (minutes) *
                </label>
                <Input
                  id="duration_minutes"
                  name="duration_minutes"
                  type="number"
                  value={formData.duration_minutes}
                  onChange={handleChange}
                  required
                  min="1"
                  placeholder="60"
                />
              </div>

              <div>
                <label htmlFor="total_marks" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  <FileText className="h-4 w-4 inline mr-1" />
                  Total Marks *
                </label>
                <Input
                  id="total_marks"
                  name="total_marks"
                  type="number"
                  value={formData.total_marks}
                  onChange={handleChange}
                  required
                  min="1"
                  placeholder="100"
                />
              </div>

              <div>
                <label htmlFor="passing_marks" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  <Users className="h-4 w-4 inline mr-1" />
                  Passing Marks *
                </label>
                <Input
                  id="passing_marks"
                  name="passing_marks"
                  type="number"
                  value={formData.passing_marks}
                  onChange={handleChange}
                  required
                  min="1"
                  max={formData.total_marks}
                  placeholder="50"
                />
              </div>

              <div className="flex items-center space-x-2">
                <input
                  id="is_active"
                  name="is_active"
                  type="checkbox"
                  checked={formData.is_active}
                  onChange={handleCheckboxChange}
                  className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                />
                <label htmlFor="is_active" className="text-sm font-medium text-gray-700 dark:text-gray-300">
                  Active (students can take this test)
                </label>
              </div>
            </div>
          </div>

          {/* Schedule (Optional) */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">Schedule (Optional)</h3>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label htmlFor="start_time" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  <Calendar className="h-4 w-4 inline mr-1" />
                  Available From
                </label>
                <Input
                  id="start_time"
                  name="start_time"
                  type="datetime-local"
                  value={formData.start_time}
                  onChange={handleChange}
                />
              </div>

              <div>
                <label htmlFor="end_time" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  <Calendar className="h-4 w-4 inline mr-1" />
                  Available Until
                </label>
                <Input
                  id="end_time"
                  name="end_time"
                  type="datetime-local"
                  value={formData.end_time}
                  onChange={handleChange}
                />
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="flex items-center justify-end space-x-4 pt-6 border-t border-gray-200 dark:border-gray-700">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? 'Creating...' : 'Create Test'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
