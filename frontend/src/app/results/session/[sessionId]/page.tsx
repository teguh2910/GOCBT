'use client';

import React, { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { resultsApi, TestResult } from '@/lib/api';
import { 
  ArrowLeft,
  Award,
  Clock,
  CheckCircle,
  XCircle,
  BarChart3,
  FileText,
  Calendar,
  AlertCircle
} from 'lucide-react';
import { formatDate, formatTime, getGradeColor } from '@/lib/utils';

export default function SessionResultPage() {
  const params = useParams();
  const router = useRouter();
  const sessionId = parseInt(params.sessionId as string);

  const [result, setResult] = useState<TestResult | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (sessionId) {
      fetchResult();
    }
  }, [sessionId]);

  const fetchResult = async () => {
    try {
      const response = await resultsApi.getBySession(sessionId);
      setResult(response.data.data);
    } catch (error) {
      setError('Failed to load test result');
    } finally {
      setLoading(false);
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

  if (error || !result) {
    return (
      <Layout>
        <div className="max-w-2xl mx-auto">
          <Card>
            <CardContent className="p-6 text-center">
              <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">Result Not Found</h3>
              <p className="text-gray-600 mb-4">
                {error || 'The test result you are looking for could not be found.'}
              </p>
              <div className="flex space-x-4 justify-center">
                <Button variant="outline" onClick={() => router.back()}>
                  <ArrowLeft className="h-4 w-4 mr-2" />
                  Go Back
                </Button>
                <Button onClick={fetchResult}>
                  Try Again
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <Button variant="outline" onClick={() => router.back()}>
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back
            </Button>
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Test Result</h1>
              <p className="text-gray-600">Session #{result.session_id}</p>
            </div>
          </div>
          <div className="text-right">
            <div className={`text-3xl font-bold ${getGradeColor(result.grade || 'F')}`}>
              {result.grade}
            </div>
            <div className="text-sm text-gray-600">Grade</div>
          </div>
        </div>

        {/* Overall Result */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              {result.is_passed ? (
                <CheckCircle className="h-6 w-6 text-green-600 mr-2" />
              ) : (
                <XCircle className="h-6 w-6 text-red-600 mr-2" />
              )}
              {result.is_passed ? 'Congratulations! You Passed' : 'Test Not Passed'}
            </CardTitle>
            <CardDescription>
              {result.is_passed 
                ? 'You have successfully completed this test with a passing score.'
                : 'You did not achieve the minimum passing score for this test.'
              }
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Final Score:</span>
                  <span className="text-2xl font-bold text-gray-900">
                    {result.percentage.toFixed(1)}%
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Marks Obtained:</span>
                  <span className="font-medium">
                    {result.marks_obtained} / {result.total_marks}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Questions Answered:</span>
                  <span className="font-medium">
                    {result.answered_questions} / {result.total_questions}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Correct Answers:</span>
                  <span className="font-medium text-green-600">
                    {result.correct_answers}
                  </span>
                </div>
              </div>
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Time Taken:</span>
                  <span className="font-medium">
                    {result.time_taken ? formatTime(result.time_taken) : 'N/A'}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Completed:</span>
                  <span className="font-medium">
                    {formatDate(result.completed_at)}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Test ID:</span>
                  <span className="font-medium">#{result.test_id}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Result ID:</span>
                  <span className="font-medium">#{result.id}</span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Performance Breakdown */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="p-2 bg-blue-100 rounded-lg">
                  <BarChart3 className="h-6 w-6 text-blue-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Score</p>
                  <p className="text-2xl font-bold text-gray-900">{result.percentage.toFixed(1)}%</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="p-2 bg-green-100 rounded-lg">
                  <CheckCircle className="h-6 w-6 text-green-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Accuracy</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {result.answered_questions > 0 
                      ? ((result.correct_answers / result.answered_questions) * 100).toFixed(1)
                      : 0}%
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="p-2 bg-purple-100 rounded-lg">
                  <Clock className="h-6 w-6 text-purple-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Time Efficiency</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {result.time_taken ? formatTime(result.time_taken) : 'N/A'}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Progress Bar */}
        <Card>
          <CardHeader>
            <CardTitle>Performance Overview</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between text-sm text-gray-600 mb-2">
                  <span>Overall Score</span>
                  <span>{result.percentage.toFixed(1)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div 
                    className={`h-3 rounded-full ${
                      result.is_passed ? 'bg-green-500' : 'bg-red-500'
                    }`}
                    style={{ width: `${Math.min(result.percentage, 100)}%` }}
                  />
                </div>
              </div>

              <div>
                <div className="flex justify-between text-sm text-gray-600 mb-2">
                  <span>Questions Answered</span>
                  <span>{result.answered_questions} / {result.total_questions}</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    className="bg-blue-500 h-2 rounded-full"
                    style={{ 
                      width: `${(result.answered_questions / result.total_questions) * 100}%` 
                    }}
                  />
                </div>
              </div>

              <div>
                <div className="flex justify-between text-sm text-gray-600 mb-2">
                  <span>Correct Answers</span>
                  <span>{result.correct_answers} / {result.answered_questions}</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    className="bg-green-500 h-2 rounded-full"
                    style={{ 
                      width: result.answered_questions > 0 
                        ? `${(result.correct_answers / result.answered_questions) * 100}%` 
                        : '0%'
                    }}
                  />
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Actions */}
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="font-medium text-gray-900">What's Next?</h3>
                <p className="text-sm text-gray-600 mt-1">
                  {result.is_passed 
                    ? 'Great job! You can view more tests or check your overall progress.'
                    : 'Don\'t worry! Review the material and try again when you\'re ready.'
                  }
                </p>
              </div>
              <div className="flex space-x-3">
                <Button variant="outline" onClick={() => router.push('/results')}>
                  <BarChart3 className="h-4 w-4 mr-2" />
                  All Results
                </Button>
                <Button onClick={() => router.push('/tests')}>
                  <FileText className="h-4 w-4 mr-2" />
                  Browse Tests
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </Layout>
  );
}
