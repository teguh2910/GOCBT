'use client';

import React, { useState } from 'react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { questionsApi, Question } from '@/lib/api';
import { X, Plus, Trash2 } from 'lucide-react';

interface CreateQuestionDialogProps {
  isOpen: boolean;
  onClose: () => void;
  testId: number;
  onQuestionCreated: (question: Question) => void;
}

export function CreateQuestionDialog({ isOpen, onClose, testId, onQuestionCreated }: CreateQuestionDialogProps) {
  const [formData, setFormData] = useState({
    question_text: '',
    question_type: 'multiple_choice' as 'multiple_choice' | 'true_false' | 'short_answer',
    marks: 10,
    order_index: 1,
  });
  const [options, setOptions] = useState([
    { option_text: '', is_correct: false, order_index: 1 },
    { option_text: '', is_correct: false, order_index: 2 },
  ]);
  const [answers, setAnswers] = useState([
    { answer_text: '', is_case_sensitive: false },
  ]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'number' ? parseInt(value) || 0 : value,
    }));
  };

  const handleOptionChange = (index: number, field: string, value: string | boolean) => {
    setOptions(prev => prev.map((option, i) => 
      i === index ? { ...option, [field]: value } : option
    ));
  };

  const handleAnswerChange = (index: number, field: string, value: string | boolean) => {
    setAnswers(prev => prev.map((answer, i) => 
      i === index ? { ...answer, [field]: value } : answer
    ));
  };

  const addOption = () => {
    setOptions(prev => [...prev, { 
      option_text: '', 
      is_correct: false, 
      order_index: prev.length + 1 
    }]);
  };

  const removeOption = (index: number) => {
    if (options.length > 2) {
      setOptions(prev => prev.filter((_, i) => i !== index));
    }
  };

  const addAnswer = () => {
    setAnswers(prev => [...prev, { answer_text: '', is_case_sensitive: false }]);
  };

  const removeAnswer = (index: number) => {
    if (answers.length > 1) {
      setAnswers(prev => prev.filter((_, i) => i !== index));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const questionData: any = {
        test_id: testId,
        question_text: formData.question_text,
        question_type: formData.question_type,
        marks: formData.marks,
        order_index: formData.order_index,
      };

      // Add options for multiple choice and true/false
      if (formData.question_type === 'multiple_choice' || formData.question_type === 'true_false') {
        const validOptions = options.filter(opt => opt.option_text.trim() !== '');
        if (validOptions.length < 2) {
          setError('Please provide at least 2 options');
          setLoading(false);
          return;
        }
        
        const correctOptions = validOptions.filter(opt => opt.is_correct);
        if (correctOptions.length === 0) {
          setError('Please mark at least one option as correct');
          setLoading(false);
          return;
        }

        questionData.options = validOptions;
      }

      // Add answers for short answer
      if (formData.question_type === 'short_answer') {
        const validAnswers = answers.filter(ans => ans.answer_text.trim() !== '');
        if (validAnswers.length === 0) {
          setError('Please provide at least one correct answer');
          setLoading(false);
          return;
        }
        questionData.answers = validAnswers;
      }

      const response = await questionsApi.create(questionData);
      const newQuestion = response.data.data;
      
      onQuestionCreated(newQuestion);
      onClose();
      
      // Reset form
      setFormData({
        question_text: '',
        question_type: 'multiple_choice',
        marks: 10,
        order_index: 1,
      });
      setOptions([
        { option_text: '', is_correct: false, order_index: 1 },
        { option_text: '', is_correct: false, order_index: 2 },
      ]);
      setAnswers([{ answer_text: '', is_case_sensitive: false }]);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create question. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Add Question</h2>
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

          {/* Question Details */}
          <div className="space-y-4">
            <div>
              <label htmlFor="question_text" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Question Text *
              </label>
              <textarea
                id="question_text"
                name="question_text"
                value={formData.question_text}
                onChange={handleChange}
                required
                rows={3}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                placeholder="Enter your question"
              />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label htmlFor="question_type" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Question Type *
                </label>
                <select
                  id="question_type"
                  name="question_type"
                  value={formData.question_type}
                  onChange={handleChange}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                >
                  <option value="multiple_choice">Multiple Choice</option>
                  <option value="true_false">True/False</option>
                  <option value="short_answer">Short Answer</option>
                </select>
              </div>

              <div>
                <label htmlFor="marks" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Marks *
                </label>
                <Input
                  id="marks"
                  name="marks"
                  type="number"
                  value={formData.marks}
                  onChange={handleChange}
                  required
                  min="1"
                />
              </div>

              <div>
                <label htmlFor="order_index" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Order
                </label>
                <Input
                  id="order_index"
                  name="order_index"
                  type="number"
                  value={formData.order_index}
                  onChange={handleChange}
                  min="1"
                />
              </div>
            </div>
          </div>

          {/* Options for Multiple Choice and True/False */}
          {(formData.question_type === 'multiple_choice' || formData.question_type === 'true_false') && (
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-medium text-gray-900 dark:text-white">Options</h3>
                {formData.question_type === 'multiple_choice' && (
                  <Button type="button" variant="outline" size="sm" onClick={addOption}>
                    <Plus className="h-4 w-4 mr-1" />
                    Add Option
                  </Button>
                )}
              </div>
              
              {options.map((option, index) => (
                <div key={index} className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    checked={option.is_correct}
                    onChange={(e) => handleOptionChange(index, 'is_correct', e.target.checked)}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  />
                  <Input
                    value={option.option_text}
                    onChange={(e) => handleOptionChange(index, 'option_text', e.target.value)}
                    placeholder={`Option ${index + 1}`}
                    className="flex-1"
                  />
                  {formData.question_type === 'multiple_choice' && options.length > 2 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => removeOption(index)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              ))}
            </div>
          )}

          {/* Answers for Short Answer */}
          {formData.question_type === 'short_answer' && (
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-medium text-gray-900 dark:text-white">Correct Answers</h3>
                <Button type="button" variant="outline" size="sm" onClick={addAnswer}>
                  <Plus className="h-4 w-4 mr-1" />
                  Add Answer
                </Button>
              </div>
              
              {answers.map((answer, index) => (
                <div key={index} className="flex items-center space-x-2">
                  <Input
                    value={answer.answer_text}
                    onChange={(e) => handleAnswerChange(index, 'answer_text', e.target.value)}
                    placeholder={`Correct answer ${index + 1}`}
                    className="flex-1"
                  />
                  <label className="flex items-center space-x-1">
                    <input
                      type="checkbox"
                      checked={answer.is_case_sensitive}
                      onChange={(e) => handleAnswerChange(index, 'is_case_sensitive', e.target.checked)}
                      className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                    />
                    <span className="text-sm text-gray-600 dark:text-gray-400">Case sensitive</span>
                  </label>
                  {answers.length > 1 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => removeAnswer(index)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              ))}
            </div>
          )}

          {/* Actions */}
          <div className="flex items-center justify-end space-x-4 pt-6 border-t border-gray-200 dark:border-gray-700">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? 'Adding...' : 'Add Question'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
