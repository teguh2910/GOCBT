'use client';

import React, { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { CreateQuestionDialog } from '@/components/CreateQuestionDialog';
import { EditTestDialog } from '@/components/EditTestDialog';
import { testsApi, Test, Question } from '@/lib/api';
import { 
  ArrowLeft,
  Plus,
  Edit,
  Trash2,
  Clock,
  FileText,
  Users,
  AlertCircle,
  CheckCircle,
  XCircle
} from 'lucide-react';
import { formatDuration, formatDate } from '@/lib/utils';

export default function TestDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const testId = parseInt(params.id as string);

  const [test, setTest] = useState<Test | null>(null);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateQuestionDialog, setShowCreateQuestionDialog] = useState(false);
  const [showEditTestDialog, setShowEditTestDialog] = useState(false);

  useEffect(() => {
    if (testId) {
      fetchTestData();
    }
  }, [testId]);

  const fetchTestData = async () => {
    try {
      const [testRes, questionsRes] = await Promise.all([
        testsApi.getById(testId),
        testsApi.getQuestions(testId),
      ]);

      setTest(testRes.data.data);
      setQuestions(questionsRes.data.data || []);
    } catch (error) {
      setError('Failed to load test data');
    } finally {
      setLoading(false);
    }
  };

  const handleQuestionCreated = (newQuestion: Question) => {
    setQuestions(prev => [...prev, newQuestion]);
  };

  const handleTestUpdated = (updatedTest: Test) => {
    setTest(updatedTest);
  };

  const handleDeleteQuestion = async (questionId: number) => {
    if (!confirm('Are you sure you want to delete this question? This action cannot be undone.')) {
      return;
    }

    try {
      // Note: This would need a delete question API endpoint
      // await questionsApi.delete(questionId);
      setQuestions(prev => prev.filter(q => q.id !== questionId));
      alert('Question deleted successfully');
    } catch (error) {
      alert('Failed to delete question');
    }
  };

  if (loading) {
    return (
      <Layout>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </Layout>
    );
  }

  if (error || !test) {
    return (
      <Layout>
        <div className="max-w-2xl mx-auto">
          <Card>
            <CardContent className="p-6 text-center">
              <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">Test Not Found</h3>
              <p className="text-gray-600 dark:text-gray-400 mb-4">
                {error || 'The test you are looking for could not be found.'}
              </p>
              <Button onClick={() => router.back()}>
                <ArrowLeft className="h-4 w-4 mr-2" />
                Go Back
              </Button>
            </CardContent>
          </Card>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="max-w-6xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <Button variant="outline" onClick={() => router.back()}>
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Tests
            </Button>
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{test.title}</h1>
              <p className="text-gray-600 dark:text-gray-400">{test.description}</p>
            </div>
          </div>
          <div className="flex items-center space-x-3">
            <Button
              variant="outline"
              onClick={() => setShowEditTestDialog(true)}
            >
              <Edit className="h-4 w-4 mr-2" />
              Edit Test
            </Button>
            <div className={`px-3 py-1 rounded-full text-sm font-medium ${
              test.is_active
                ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'
            }`}>
              {test.is_active ? 'Active' : 'Inactive'}
            </div>
          </div>
        </div>

        {/* Test Overview */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="p-2 bg-blue-100 dark:bg-blue-900 rounded-lg">
                  <Clock className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Duration</p>
                  <p className="text-2xl font-bold text-gray-900 dark:text-white">{formatDuration(test.duration_minutes)}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="p-2 bg-green-100 dark:bg-green-900 rounded-lg">
                  <FileText className="h-6 w-6 text-green-600 dark:text-green-400" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Total Marks</p>
                  <p className="text-2xl font-bold text-gray-900 dark:text-white">{test.total_marks}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="p-2 bg-purple-100 dark:bg-purple-900 rounded-lg">
                  <Users className="h-6 w-6 text-purple-600 dark:text-purple-400" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Passing Marks</p>
                  <p className="text-2xl font-bold text-gray-900 dark:text-white">{test.passing_marks}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="p-2 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
                  <FileText className="h-6 w-6 text-yellow-600 dark:text-yellow-400" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Questions</p>
                  <p className="text-2xl font-bold text-gray-900 dark:text-white">{questions.length}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Test Details */}
        <Card>
          <CardHeader>
            <CardTitle>Test Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 className="font-medium text-gray-900 dark:text-white mb-2">Instructions</h4>
                <p className="text-gray-600 dark:text-gray-400 whitespace-pre-wrap">
                  {test.instructions || 'No instructions provided'}
                </p>
              </div>
              <div className="space-y-3">
                <div>
                  <span className="font-medium text-gray-900 dark:text-white">Created:</span>
                  <span className="ml-2 text-gray-600 dark:text-gray-400">{formatDate(test.created_at)}</span>
                </div>
                {test.start_time && (
                  <div>
                    <span className="font-medium text-gray-900 dark:text-white">Available from:</span>
                    <span className="ml-2 text-gray-600 dark:text-gray-400">{formatDate(test.start_time)}</span>
                  </div>
                )}
                {test.end_time && (
                  <div>
                    <span className="font-medium text-gray-900 dark:text-white">Available until:</span>
                    <span className="ml-2 text-gray-600 dark:text-gray-400">{formatDate(test.end_time)}</span>
                  </div>
                )}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Questions */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Questions ({questions.length})</CardTitle>
                <CardDescription>Manage questions for this test</CardDescription>
              </div>
              <Button onClick={() => setShowCreateQuestionDialog(true)}>
                <Plus className="h-4 w-4 mr-2" />
                Add Question
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            {questions.length === 0 ? (
              <div className="text-center py-8">
                <FileText className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">No Questions Yet</h3>
                <p className="text-gray-600 dark:text-gray-400 mb-4">
                  Add questions to make this test available to students.
                </p>
                <Button onClick={() => setShowCreateQuestionDialog(true)}>
                  <Plus className="h-4 w-4 mr-2" />
                  Add Your First Question
                </Button>
              </div>
            ) : (
              <div className="space-y-4">
                {questions.map((question, index) => (
                  <div key={question.id} className="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center space-x-2 mb-2">
                          <span className="text-sm font-medium text-gray-500 dark:text-gray-400">
                            Question {index + 1}
                          </span>
                          <span className="px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-xs rounded">
                            {question.question_type.replace('_', ' ')}
                          </span>
                          <span className="text-sm text-gray-500 dark:text-gray-400">
                            {question.marks} marks
                          </span>
                        </div>
                        <p className="text-gray-900 dark:text-white mb-3">{question.question_text}</p>
                        
                        {/* Show options for multiple choice and true/false */}
                        {question.options && question.options.length > 0 && (
                          <div className="space-y-1">
                            {question.options.map((option, optIndex) => (
                              <div key={option.id} className="flex items-center space-x-2 text-sm">
                                {option.is_correct ? (
                                  <CheckCircle className="h-4 w-4 text-green-500" />
                                ) : (
                                  <XCircle className="h-4 w-4 text-gray-300" />
                                )}
                                <span className={option.is_correct ? 'text-green-700 dark:text-green-400 font-medium' : 'text-gray-600 dark:text-gray-400'}>
                                  {String.fromCharCode(65 + optIndex)}. {option.option_text}
                                </span>
                              </div>
                            ))}
                          </div>
                        )}

                        {/* Show answers for short answer */}
                        {question.correct_answers && question.correct_answers.length > 0 && (
                          <div className="space-y-1">
                            <p className="text-sm font-medium text-gray-700 dark:text-gray-300">Correct answers:</p>
                            {question.correct_answers.map((answer, ansIndex) => (
                              <div key={answer.id} className="flex items-center space-x-2 text-sm">
                                <CheckCircle className="h-4 w-4 text-green-500" />
                                <span className="text-green-700 dark:text-green-400">
                                  {answer.answer_text}
                                  {answer.is_case_sensitive && (
                                    <span className="ml-1 text-xs text-gray-500">(case sensitive)</span>
                                  )}
                                </span>
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                      
                      <div className="flex space-x-2 ml-4">
                        <Button variant="outline" size="sm">
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button 
                          variant="destructive" 
                          size="sm"
                          onClick={() => handleDeleteQuestion(question.id)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Create Question Dialog */}
        <CreateQuestionDialog
          isOpen={showCreateQuestionDialog}
          onClose={() => setShowCreateQuestionDialog(false)}
          testId={testId}
          onQuestionCreated={handleQuestionCreated}
        />

        {/* Edit Test Dialog */}
        <EditTestDialog
          isOpen={showEditTestDialog}
          onClose={() => setShowEditTestDialog(false)}
          test={test}
          onTestUpdated={handleTestUpdated}
        />
      </div>
    </Layout>
  );
}
