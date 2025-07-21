'use client';

import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { ExportResultsDialog } from '@/components/ExportResultsDialog';
import { resultsApi, testsApi, TestResult, Test, TestStatistics } from '@/lib/api';
import { 
  BarChart3, 
  Users, 
  TrendingUp, 
  Award, 
  Clock, 
  CheckCircle, 
  XCircle, 
  AlertCircle,
  Download,
  Filter
} from 'lucide-react';
import { formatDate, formatTime, getGradeColor } from '@/lib/utils';

export default function ManageResultsPage() {
  const searchParams = useSearchParams();
  const testId = searchParams.get('test');

  const [results, setResults] = useState<TestResult[]>([]);
  const [tests, setTests] = useState<Test[]>([]);
  const [selectedTest, setSelectedTest] = useState<Test | null>(null);
  const [statistics, setStatistics] = useState<TestStatistics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showExportDialog, setShowExportDialog] = useState(false);

  useEffect(() => {
    fetchTests();
  }, []);

  useEffect(() => {
    if (testId && tests.length > 0) {
      const test = tests.find(t => t.id === parseInt(testId));
      if (test) {
        setSelectedTest(test);
        fetchTestResults(parseInt(testId));
        fetchTestStatistics(parseInt(testId));
      }
    }
  }, [testId, tests]);

  const fetchTests = async () => {
    try {
      const response = await testsApi.getAll({ limit: 100 });
      setTests(response.data.data || []);
    } catch (error) {
      setError('Failed to load tests');
    } finally {
      setLoading(false);
    }
  };

  const fetchTestResults = async (testId: number) => {
    try {
      const response = await resultsApi.getTestResults(testId);
      setResults(response.data.data || []);
    } catch (error) {
      setError('Failed to load test results');
    }
  };

  const fetchTestStatistics = async (testId: number) => {
    try {
      const response = await resultsApi.getTestStatistics(testId);
      setStatistics(response.data.data);
    } catch (error) {
      console.error('Failed to load test statistics:', error);
    }
  };

  const handleTestSelect = (test: Test) => {
    setSelectedTest(test);
    setResults([]);
    setStatistics(null);
    fetchTestResults(test.id);
    fetchTestStatistics(test.id);
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

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Test Results</h1>
            <p className="text-gray-600">View and analyze student performance</p>
          </div>
          {selectedTest && (
            <Button
              variant="outline"
              className="flex items-center"
              onClick={() => setShowExportDialog(true)}
              disabled={results.length === 0}
            >
              <Download className="h-4 w-4 mr-2" />
              Export Results
            </Button>
          )}
        </div>

        {error && (
          <Card>
            <CardContent className="p-6 text-center">
              <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">Error</h3>
              <p className="text-gray-600">{error}</p>
              <Button onClick={fetchTests} className="mt-4">
                Try Again
              </Button>
            </CardContent>
          </Card>
        )}

        {/* Test Selection */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <Filter className="h-5 w-5 mr-2" />
              Select Test
            </CardTitle>
            <CardDescription>Choose a test to view its results and statistics</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {tests.map((test) => (
                <button
                  key={test.id}
                  onClick={() => handleTestSelect(test)}
                  className={`p-4 border rounded-lg text-left transition-colors ${
                    selectedTest?.id === test.id
                      ? 'border-blue-500 bg-blue-50'
                      : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
                  }`}
                >
                  <h4 className="font-medium text-gray-900">{test.title}</h4>
                  <p className="text-sm text-gray-600 mt-1">{test.description}</p>
                  <div className="flex items-center mt-2 text-xs text-gray-500">
                    <Clock className="h-3 w-3 mr-1" />
                    <span>{test.duration_minutes} min</span>
                    <span className="mx-2">•</span>
                    <span>{test.total_marks} marks</span>
                  </div>
                </button>
              ))}
            </div>
          </CardContent>
        </Card>

        {selectedTest && (
          <>
            {/* Statistics */}
            {statistics && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <Card>
                  <CardContent className="p-6">
                    <div className="flex items-center">
                      <div className="p-2 bg-blue-100 rounded-lg">
                        <Users className="h-6 w-6 text-blue-600" />
                      </div>
                      <div className="ml-4">
                        <p className="text-sm font-medium text-gray-600">Total Attempts</p>
                        <p className="text-2xl font-bold text-gray-900">{statistics.total_attempts}</p>
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
                        <p className="text-sm font-medium text-gray-600">Pass Rate</p>
                        <p className="text-2xl font-bold text-gray-900">
                          {statistics.total_attempts > 0 
                            ? ((statistics.passed_attempts / statistics.total_attempts) * 100).toFixed(1)
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
                        <TrendingUp className="h-6 w-6 text-purple-600" />
                      </div>
                      <div className="ml-4">
                        <p className="text-sm font-medium text-gray-600">Average Score</p>
                        <p className="text-2xl font-bold text-gray-900">{statistics.average_score.toFixed(1)}%</p>
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
                        <p className="text-sm font-medium text-gray-600">Highest Score</p>
                        <p className="text-2xl font-bold text-gray-900">{statistics.highest_score.toFixed(1)}%</p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* Results Table */}
            <Card>
              <CardHeader>
                <CardTitle>Student Results - {selectedTest.title}</CardTitle>
                <CardDescription>Individual student performance data</CardDescription>
              </CardHeader>
              <CardContent>
                {results.length === 0 ? (
                  <div className="text-center py-8">
                    <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 mb-2">No Results Yet</h3>
                    <p className="text-gray-600">No students have completed this test yet.</p>
                  </div>
                ) : (
                  <div className="overflow-x-auto">
                    <table className="w-full">
                      <thead>
                        <tr className="border-b">
                          <th className="text-left py-3 px-4 font-medium text-gray-900">Student</th>
                          <th className="text-left py-3 px-4 font-medium text-gray-900">Score</th>
                          <th className="text-left py-3 px-4 font-medium text-gray-900">Grade</th>
                          <th className="text-left py-3 px-4 font-medium text-gray-900">Status</th>
                          <th className="text-left py-3 px-4 font-medium text-gray-900">Time</th>
                          <th className="text-left py-3 px-4 font-medium text-gray-900">Completed</th>
                        </tr>
                      </thead>
                      <tbody>
                        {results.map((result) => (
                          <tr key={result.id} className="border-b hover:bg-gray-50">
                            <td className="py-3 px-4">
                              <div>
                                <div className="font-medium text-gray-900">Student #{result.user_id}</div>
                                <div className="text-sm text-gray-600">
                                  {result.answered_questions}/{result.total_questions} questions
                                </div>
                              </div>
                            </td>
                            <td className="py-3 px-4">
                              <div>
                                <div className="font-medium">{result.percentage.toFixed(1)}%</div>
                                <div className="text-sm text-gray-600">
                                  {result.marks_obtained}/{result.total_marks} marks
                                </div>
                              </div>
                            </td>
                            <td className="py-3 px-4">
                              <span className={`font-bold ${getGradeColor(result.grade || 'F')}`}>
                                {result.grade}
                              </span>
                            </td>
                            <td className="py-3 px-4">
                              <div className="flex items-center">
                                {result.is_passed ? (
                                  <CheckCircle className="h-4 w-4 text-green-500 mr-2" />
                                ) : (
                                  <XCircle className="h-4 w-4 text-red-500 mr-2" />
                                )}
                                <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                                  result.is_passed 
                                    ? 'bg-green-100 text-green-800' 
                                    : 'bg-red-100 text-red-800'
                                }`}>
                                  {result.is_passed ? 'Passed' : 'Failed'}
                                </span>
                              </div>
                            </td>
                            <td className="py-3 px-4 text-sm text-gray-600">
                              {result.time_taken ? formatTime(result.time_taken) : 'N/A'}
                            </td>
                            <td className="py-3 px-4 text-sm text-gray-600">
                              {formatDate(result.completed_at)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </CardContent>
            </Card>
          </>
        )}

        {/* Export Results Dialog */}
        {selectedTest && (
          <ExportResultsDialog
            isOpen={showExportDialog}
            onClose={() => setShowExportDialog(false)}
            test={selectedTest}
            results={results}
          />
        )}
      </div>
    </Layout>
  );
}
