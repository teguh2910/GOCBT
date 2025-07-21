'use client';

import React, { useEffect, useState } from 'react';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { resultsApi, testsApi, TestResult, Test } from '@/lib/api';
import { 
  BarChart3, 
  Clock, 
  CheckCircle, 
  XCircle, 
  AlertCircle,
  TrendingUp,
  Award
} from 'lucide-react';
import { formatDate, formatTime, getGradeColor } from '@/lib/utils';

export default function ResultsPage() {
  const [results, setResults] = useState<TestResult[]>([]);
  const [tests, setTests] = useState<Test[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchResults();
  }, []);

  const fetchResults = async () => {
    try {
      const [resultsRes, testsRes] = await Promise.all([
        resultsApi.getMy(),
        testsApi.getAll({ limit: 100 }) // Get all tests to match with results
      ]);

      setResults(resultsRes.data.data || []);
      setTests(testsRes.data.data || []);
    } catch (error) {
      setError('Failed to load results');
    } finally {
      setLoading(false);
    }
  };

  const getTestTitle = (testId: number) => {
    const test = tests.find(t => t.id === testId);
    return test ? test.title : `Test #${testId}`;
  };

  const calculateStats = () => {
    if (results.length === 0) return null;

    const totalTests = results.length;
    const passedTests = results.filter(r => r.is_passed).length;
    const averageScore = results.reduce((sum, r) => sum + r.percentage, 0) / totalTests;
    const highestScore = Math.max(...results.map(r => r.percentage));

    return {
      totalTests,
      passedTests,
      passRate: (passedTests / totalTests) * 100,
      averageScore,
      highestScore,
    };
  };

  const stats = calculateStats();

  if (loading) {
    return (
      <Layout>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">My Test Results</h1>
          <p className="text-gray-600">View your test performance and progress</p>
        </div>

        {error && (
          <Card>
            <CardContent className="p-6 text-center">
              <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">Error</h3>
              <p className="text-gray-600">{error}</p>
              <Button onClick={fetchResults} className="mt-4">
                Try Again
              </Button>
            </CardContent>
          </Card>
        )}

        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center">
                  <div className="p-2 bg-blue-100 rounded-lg">
                    <BarChart3 className="h-6 w-6 text-blue-600" />
                  </div>
                  <div className="ml-4">
                    <p className="text-sm font-medium text-gray-600">Total Tests</p>
                    <p className="text-2xl font-bold text-gray-900">{stats.totalTests}</p>
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
                    <p className="text-sm font-medium text-gray-600">Passed</p>
                    <p className="text-2xl font-bold text-gray-900">{stats.passedTests}</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center">
                  <div className="p-2 bg-purple-100 rounded-lg">
                    <TrendingUp className="h-6 w-6 text-purple-600" />
                  </div>
                  <div className="ml-4">
                    <p className="text-sm font-medium text-gray-600">Average Score</p>
                    <p className="text-2xl font-bold text-gray-900">{stats.averageScore.toFixed(1)}%</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center">
                  <div className="p-2 bg-yellow-100 rounded-lg">
                    <Award className="h-6 w-6 text-yellow-600" />
                  </div>
                  <div className="ml-4">
                    <p className="text-sm font-medium text-gray-600">Best Score</p>
                    <p className="text-2xl font-bold text-gray-900">{stats.highestScore.toFixed(1)}%</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {results.length === 0 && !error ? (
          <Card>
            <CardContent className="p-6 text-center">
              <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">No Results Yet</h3>
              <p className="text-gray-600">You haven't completed any tests yet. Take a test to see your results here.</p>
            </CardContent>
          </Card>
        ) : (
          <Card>
            <CardHeader>
              <CardTitle>Test Results</CardTitle>
              <CardDescription>Your recent test performance</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {results.map((result) => (
                  <div key={result.id} className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50">
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <div>
                          <h4 className="font-medium text-gray-900 dark:text-white">{getTestTitle(result.test_id)}</h4>
                          <p className="text-xs text-gray-500 dark:text-gray-400">Result #{result.id}</p>
                        </div>
                        <div className="flex items-center space-x-2">
                          {result.is_passed ? (
                            <CheckCircle className="h-5 w-5 text-green-500" />
                          ) : (
                            <XCircle className="h-5 w-5 text-red-500" />
                          )}
                          <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                            result.is_passed ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                          }`}>
                            {result.is_passed ? 'Passed' : 'Failed'}
                          </span>
                        </div>
                      </div>
                      
                      <div className="mt-2 grid grid-cols-2 md:grid-cols-4 gap-4 text-sm text-gray-600">
                        <div>
                          <span className="font-medium">Score:</span> {result.percentage.toFixed(1)}%
                        </div>
                        <div>
                          <span className="font-medium">Grade:</span> 
                          <span className={`ml-1 font-bold ${getGradeColor(result.grade || 'F')}`}>
                            {result.grade}
                          </span>
                        </div>
                        <div>
                          <span className="font-medium">Marks:</span> {result.marks_obtained}/{result.total_marks}
                        </div>
                        <div>
                          <span className="font-medium">Questions:</span> {result.answered_questions}/{result.total_questions}
                        </div>
                      </div>

                      <div className="mt-2 flex items-center space-x-4 text-xs text-gray-500">
                        <div className="flex items-center">
                          <Clock className="h-3 w-3 mr-1" />
                          <span>
                            {result.time_taken ? formatTime(result.time_taken) : 'N/A'}
                          </span>
                        </div>
                        <div>
                          <span>Completed: {formatDate(result.completed_at)}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </Layout>
  );
}
