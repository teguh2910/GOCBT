'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { sessionsApi, testsApi, TestSession, Test } from '@/lib/api';
import { 
  Clock, 
  Play, 
  CheckCircle, 
  AlertCircle, 
  Eye, 
  RefreshCw,
  Calendar,
  Activity
} from 'lucide-react';
import { formatTime, formatDate, getStatusColor, formatStatus } from '@/lib/utils';

export default function SessionsPage() {
  const [sessions, setSessions] = useState<TestSession[]>([]);
  const [tests, setTests] = useState<Test[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchSessions();
  }, []);

  const fetchSessions = async () => {
    try {
      const [sessionsRes, testsRes] = await Promise.all([
        sessionsApi.getMy(),
        testsApi.getAll({ limit: 100 }) // Get all tests to match with sessions
      ]);

      setSessions(sessionsRes.data.data || []);
      setTests(testsRes.data.data || []);
    } catch (error) {
      setError('Failed to load sessions');
    } finally {
      setLoading(false);
    }
  };

  const getTestTitle = (testId: number) => {
    const test = tests.find(t => t.id === testId);
    return test ? test.title : `Test #${testId}`;
  };

  const getTimeRemaining = (session: TestSession) => {
    if (session.status !== 'in_progress' || !session.expires_at) {
      return null;
    }

    const now = new Date().getTime();
    const expiresAt = new Date(session.expires_at).getTime();
    const remaining = Math.max(0, Math.floor((expiresAt - now) / 1000));
    
    return remaining;
  };

  const activeSessions = sessions.filter(s => s.status === 'in_progress');
  const completedSessions = sessions.filter(s => s.status === 'submitted' || s.status === 'completed');
  const notStartedSessions = sessions.filter(s => s.status === 'not_started');

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
            <h1 className="text-2xl font-bold text-gray-900">My Test Sessions</h1>
            <p className="text-gray-600">View and manage your test sessions</p>
          </div>
          <Button variant="outline" onClick={fetchSessions} className="flex items-center">
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
        </div>

        {error && (
          <Card>
            <CardContent className="p-6 text-center">
              <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">Error</h3>
              <p className="text-gray-600">{error}</p>
              <Button onClick={fetchSessions} className="mt-4">
                Try Again
              </Button>
            </CardContent>
          </Card>
        )}

        {/* Statistics */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="p-2 bg-blue-100 rounded-lg">
                  <Activity className="h-6 w-6 text-blue-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Active Sessions</p>
                  <p className="text-2xl font-bold text-gray-900">{activeSessions.length}</p>
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
                  <p className="text-sm font-medium text-gray-600">Completed</p>
                  <p className="text-2xl font-bold text-gray-900">{completedSessions.length}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="p-2 bg-yellow-100 rounded-lg">
                  <Clock className="h-6 w-6 text-yellow-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Pending</p>
                  <p className="text-2xl font-bold text-gray-900">{notStartedSessions.length}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Active Sessions */}
        {activeSessions.length > 0 && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Play className="h-5 w-5 mr-2 text-blue-600" />
                Active Sessions
              </CardTitle>
              <CardDescription>Tests you are currently taking</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {activeSessions.map((session) => {
                  const timeRemaining = getTimeRemaining(session);
                  return (
                    <div key={session.id} className="flex items-center justify-between p-4 border rounded-lg bg-blue-50 border-blue-200">
                      <div>
                        <h4 className="font-medium text-gray-900 dark:text-white">{getTestTitle(session.test_id)}</h4>
                        <p className="text-sm text-gray-600 dark:text-gray-400">Session #{session.id}</p>
                        <div className="flex items-center mt-1 text-xs text-gray-500 dark:text-gray-400">
                          <Calendar className="h-3 w-3 mr-1" />
                          <span>Started: {session.started_at ? formatDate(session.started_at) : 'Not started'}</span>
                        </div>
                      </div>
                      <div className="text-right">
                        {timeRemaining !== null && (
                          <div className={`text-lg font-mono font-bold ${
                            timeRemaining < 300 ? 'text-red-600' : 'text-blue-600'
                          }`}>
                            {formatTime(timeRemaining)}
                          </div>
                        )}
                        <div className="text-sm text-gray-600 mb-2">
                          Question {session.current_question_index + 1}
                        </div>
                        <Link href={`/sessions/${session.session_token}`}>
                          <Button size="sm" className="bg-blue-600 hover:bg-blue-700">
                            Continue Test
                          </Button>
                        </Link>
                      </div>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>
        )}

        {/* All Sessions */}
        <Card>
          <CardHeader>
            <CardTitle>All Sessions</CardTitle>
            <CardDescription>Complete history of your test sessions</CardDescription>
          </CardHeader>
          <CardContent>
            {sessions.length === 0 ? (
              <div className="text-center py-8">
                <Activity className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">No Sessions</h3>
                <p className="text-gray-600 mb-4">You haven't started any test sessions yet.</p>
                <Link href="/tests">
                  <Button>Browse Available Tests</Button>
                </Link>
              </div>
            ) : (
              <div className="space-y-4">
                {sessions.map((session) => {
                  const timeRemaining = getTimeRemaining(session);
                  return (
                    <div key={session.id} className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50">
                      <div className="flex-1">
                        <div className="flex items-center justify-between">
                          <h4 className="font-medium text-gray-900 dark:text-white">{getTestTitle(session.test_id)}</h4>
                          <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(session.status)}`}>
                            {formatStatus(session.status)}
                          </span>
                        </div>

                        <div className="mt-2 grid grid-cols-2 md:grid-cols-4 gap-4 text-sm text-gray-600 dark:text-gray-400">
                          <div>
                            <span className="font-medium">Session:</span> #{session.id}
                          </div>
                          <div>
                            <span className="font-medium">Progress:</span> Question {session.current_question_index + 1}
                          </div>
                          <div>
                            <span className="font-medium">Started:</span> {
                              session.started_at ? formatDate(session.started_at) : 'Not started'
                            }
                          </div>
                          <div>
                            <span className="font-medium">Expires:</span> {formatDate(session.expires_at)}
                          </div>
                        </div>

                        {timeRemaining !== null && (
                          <div className="mt-2">
                            <span className="text-sm font-medium text-gray-600">Time Remaining: </span>
                            <span className={`font-mono text-sm ${
                              timeRemaining < 300 ? 'text-red-600 font-bold' : 'text-blue-600'
                            }`}>
                              {formatTime(timeRemaining)}
                            </span>
                          </div>
                        )}
                      </div>

                      <div className="ml-4 flex space-x-2">
                        {session.status === 'in_progress' && (
                          <Link href={`/sessions/${session.session_token}`}>
                            <Button size="sm">
                              <Play className="h-3 w-3 mr-1" />
                              Continue
                            </Button>
                          </Link>
                        )}
                        {session.status === 'not_started' && (
                          <Link href={`/sessions/${session.session_token}`}>
                            <Button size="sm" variant="outline">
                              <Play className="h-3 w-3 mr-1" />
                              Start
                            </Button>
                          </Link>
                        )}
                        {(session.status === 'submitted' || session.status === 'completed') && (
                          <Link href={`/results/session/${session.id}`}>
                            <Button size="sm" variant="outline">
                              <Eye className="h-3 w-3 mr-1" />
                              View Results
                            </Button>
                          </Link>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </Layout>
  );
}
