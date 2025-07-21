'use client';

import React, { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { testsApi, sessionsApi, resultsApi, TestResult, TestSession, Test } from '@/lib/api';
import { 
  BookOpen, 
  Clock, 
  TrendingUp, 
  Users, 
  FileText, 
  CheckCircle,
  AlertCircle,
  BarChart3
} from 'lucide-react';
import Link from 'next/link';
import { formatTime, formatDate, getGradeColor, getStatusColor, formatStatus } from '@/lib/utils';

export default function DashboardPage() {
  const { user } = useAuth();
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    availableTests: 0,
    activeSessions: 0,
    completedTests: 0,
    averageScore: 0,
  });
  const [recentTests, setRecentTests] = useState<Test[]>([]);
  const [recentSessions, setRecentSessions] = useState<TestSession[]>([]);
  const [recentResults, setRecentResults] = useState<TestResult[]>([]);
  const [allTests, setAllTests] = useState<Test[]>([]);

  useEffect(() => {
    fetchDashboardData();
  }, [user]);

  const fetchDashboardData = async () => {
    try {
      if (user?.role === 'student') {
        // Student dashboard data
        const [testsRes, sessionsRes, resultsRes, allTestsRes] = await Promise.all([
          testsApi.getAvailable(),
          sessionsApi.getMy(),
          resultsApi.getMy(),
          testsApi.getAll({ limit: 100 }), // Get all tests for name resolution
        ]);

        const tests = testsRes.data.data || [];
        const sessions = sessionsRes.data.data || [];
        const results = resultsRes.data.data || [];
        const allTestsData = allTestsRes.data.data || [];

        // Filter out tests that the student has already completed
        const completedTestIds = results.map(result => result.test_id);
        const availableTests = tests.filter(test => !completedTestIds.includes(test.id));

        setRecentTests(availableTests.slice(0, 5));
        setRecentSessions(sessions.slice(0, 5));
        setRecentResults(results.slice(0, 5));
        setAllTests(allTestsData);

        const activeSessions = sessions.filter(s => s.status === 'in_progress').length;
        const averageScore = results.length > 0 
          ? results.reduce((sum, r) => sum + r.percentage, 0) / results.length 
          : 0;

        setStats({
          availableTests: availableTests.length,
          activeSessions,
          completedTests: results.length,
          averageScore,
        });
      } else {
        // Teacher/Admin dashboard data
        const testsRes = await testsApi.getAll({ limit: 10 });
        const tests = testsRes.data.data || [];
        setRecentTests(tests);

        setStats({
          availableTests: tests.length,
          activeSessions: 0,
          completedTests: 0,
          averageScore: 0,
        });
      }
    } catch (error) {
      console.error('Failed to fetch dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  const getTestTitle = (testId: number) => {
    const test = allTests.find(t => t.id === testId);
    return test ? test.title : `Test #${testId}`;
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
        {/* Header */}
        <div>
          <h1 className="text-2xl font-bold text-gray-900">
            Welcome back, {user?.first_name}!
          </h1>
          <p className="text-gray-600">
            {user?.role === 'student' 
              ? "Here's your learning progress and available tests."
              : "Manage your tests and monitor student progress."
            }
          </p>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="p-2 bg-blue-100 rounded-lg">
                  <FileText className="h-6 w-6 text-blue-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">
                    {user?.role === 'student' ? 'Available Tests' : 'Total Tests'}
                  </p>
                  <p className="text-2xl font-bold text-gray-900">{stats.availableTests}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          {user?.role === 'student' && (
            <>
              <Card>
                <CardContent className="p-6">
                  <div className="flex items-center">
                    <div className="p-2 bg-yellow-100 rounded-lg">
                      <Clock className="h-6 w-6 text-yellow-600" />
                    </div>
                    <div className="ml-4">
                      <p className="text-sm font-medium text-gray-600">Active Sessions</p>
                      <p className="text-2xl font-bold text-gray-900">{stats.activeSessions}</p>
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
                      <p className="text-sm font-medium text-gray-600">Completed Tests</p>
                      <p className="text-2xl font-bold text-gray-900">{stats.completedTests}</p>
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
                      <p className="text-2xl font-bold text-gray-900">
                        {stats.averageScore.toFixed(1)}%
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </>
          )}
        </div>

        {/* Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Available Tests / Recent Tests */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <BookOpen className="h-5 w-5 mr-2" />
                {user?.role === 'student' ? 'Available Tests' : 'Recent Tests'}
              </CardTitle>
              <CardDescription>
                {user?.role === 'student' 
                  ? 'Tests you can take right now'
                  : 'Recently created tests'
                }
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {recentTests.length > 0 ? (
                  recentTests.map((test) => (
                    <div key={test.id} className="flex items-center justify-between p-4 border rounded-lg">
                      <div>
                        <h4 className="font-medium">{test.title}</h4>
                        <p className="text-sm text-gray-600">{test.description}</p>
                        <div className="flex items-center mt-2 text-xs text-gray-500">
                          <Clock className="h-3 w-3 mr-1" />
                          {test.duration_minutes} min
                          <span className="mx-2">•</span>
                          {test.total_marks} marks
                        </div>
                      </div>
                      {user?.role === 'student' && (
                        <Link href={`/tests/${test.id}`}>
                          <Button size="sm">Start Test</Button>
                        </Link>
                      )}
                    </div>
                  ))
                ) : (
                  <div className="text-center py-4">
                    <p className="text-gray-500 dark:text-gray-400 mb-2">
                      {user?.role === 'student' ? 'No new tests available' : 'No tests created yet'}
                    </p>
                    {user?.role === 'student' && (
                      <p className="text-xs text-gray-400 dark:text-gray-500">
                        Completed tests are in "My Results"
                      </p>
                    )}
                  </div>
                )}
              </div>
              {recentTests.length > 0 && (
                <div className="mt-4">
                  <Link href={user?.role === 'student' ? '/tests' : '/manage/tests'}>
                    <Button variant="outline" className="w-full">
                      View All Tests
                    </Button>
                  </Link>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Recent Activity */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <BarChart3 className="h-5 w-5 mr-2" />
                {user?.role === 'student' ? 'Recent Results' : 'Recent Activity'}
              </CardTitle>
              <CardDescription>
                {user?.role === 'student' 
                  ? 'Your latest test results'
                  : 'Latest system activity'
                }
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {user?.role === 'student' && recentResults.length > 0 ? (
                  recentResults.map((result) => (
                    <div key={result.id} className="flex items-center justify-between p-4 border rounded-lg">
                      <div>
                        <h4 className="font-medium text-gray-900 dark:text-white">{getTestTitle(result.test_id)}</h4>
                        <p className="text-sm text-gray-600 dark:text-gray-400">
                          Score: {result.percentage.toFixed(1)}%
                        </p>
                        <div className="flex items-center mt-2">
                          <span className={`text-xs font-medium ${getGradeColor(result.grade || 'F')}`}>
                            Grade: {result.grade}
                          </span>
                          <span className="mx-2 text-gray-300 dark:text-gray-600">•</span>
                          <span className="text-xs text-gray-500 dark:text-gray-400">
                            {formatDate(result.completed_at)}
                          </span>
                        </div>
                      </div>
                      <div className={`px-2 py-1 rounded-full text-xs font-medium ${
                        result.is_passed
                          ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                          : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                      }`}>
                        {result.is_passed ? 'Passed' : 'Failed'}
                      </div>
                    </div>
                  ))
                ) : (
                  <p className="text-gray-500 text-center py-4">
                    {user?.role === 'student' ? 'No test results yet' : 'No recent activity'}
                  </p>
                )}
              </div>
              {user?.role === 'student' && recentResults.length > 0 && (
                <div className="mt-4">
                  <Link href="/results">
                    <Button variant="outline" className="w-full">
                      View All Results
                    </Button>
                  </Link>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </Layout>
  );
}
