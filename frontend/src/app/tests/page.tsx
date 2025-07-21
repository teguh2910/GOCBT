'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { testsApi, resultsApi, Test } from '@/lib/api';
import { Clock, FileText, Users, Calendar, AlertCircle } from 'lucide-react';
import { formatDuration, formatDate } from '@/lib/utils';

export default function TestsPage() {
  const [tests, setTests] = useState<Test[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchTests();
  }, []);

  const fetchTests = async () => {
    try {
      const [testsRes, resultsRes] = await Promise.all([
        testsApi.getAvailable(),
        resultsApi.getMy(),
      ]);

      const tests = testsRes.data.data || [];
      const results = resultsRes.data.data || [];

      // Filter out tests that the student has already completed
      const completedTestIds = results.map(result => result.test_id);
      const availableTests = tests.filter(test => !completedTestIds.includes(test.id));

      setTests(availableTests);
    } catch (error) {
      setError('Failed to load tests');
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

  return (
    <Layout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Available Tests</h1>
          <p className="text-gray-600">Choose a test to begin your assessment</p>
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

        {tests.length === 0 && !error ? (
          <Card>
            <CardContent className="p-6 text-center">
              <FileText className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">No Tests Available</h3>
              <p className="text-gray-600 dark:text-gray-400 mb-4">
                There are currently no new tests available for you to take.
              </p>
              <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
                <p className="text-sm text-blue-700 dark:text-blue-400 mb-3">
                  💡 <strong>Note:</strong> Tests disappear from this list once you complete them.
                  You can view your completed tests and results in the <strong>"My Results"</strong> section.
                </p>
                <Link href="/results">
                  <Button variant="outline" size="sm">
                    View My Results
                  </Button>
                </Link>
              </div>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {tests.map((test) => (
              <Card key={test.id} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <CardTitle className="text-lg">{test.title}</CardTitle>
                  <CardDescription>{test.description}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div className="flex items-center text-gray-600">
                      <Clock className="h-4 w-4 mr-2" />
                      <span>{formatDuration(test.duration_minutes)}</span>
                    </div>
                    <div className="flex items-center text-gray-600">
                      <FileText className="h-4 w-4 mr-2" />
                      <span>{test.total_marks} marks</span>
                    </div>
                  </div>

                  <div className="text-sm text-gray-600">
                    <div className="flex items-center">
                      <Users className="h-4 w-4 mr-2" />
                      <span>Passing: {test.passing_marks} marks</span>
                    </div>
                  </div>

                  {(test.start_time || test.end_time) && (
                    <div className="text-xs text-gray-500 space-y-1">
                      {test.start_time && (
                        <div className="flex items-center">
                          <Calendar className="h-3 w-3 mr-1" />
                          <span>Starts: {formatDate(test.start_time)}</span>
                        </div>
                      )}
                      {test.end_time && (
                        <div className="flex items-center">
                          <Calendar className="h-3 w-3 mr-1" />
                          <span>Ends: {formatDate(test.end_time)}</span>
                        </div>
                      )}
                    </div>
                  )}

                  {test.instructions && (
                    <div className="text-sm text-gray-600">
                      <p className="font-medium mb-1">Instructions:</p>
                      <p className="text-xs line-clamp-3">{test.instructions}</p>
                    </div>
                  )}

                  <div className="pt-4 border-t">
                    <Link href={`/tests/${test.id}`}>
                      <Button className="w-full">
                        Start Test
                      </Button>
                    </Link>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </Layout>
  );
}
