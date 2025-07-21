'use client';

import React, { useEffect, useState } from 'react';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { sessionsApi, TestSession } from '@/lib/api';
import { 
  Clock, 
  User, 
  Play, 
  Pause, 
  CheckCircle, 
  AlertCircle, 
  RefreshCw,
  Eye,
  Activity
} from 'lucide-react';
import { formatTime, formatDate, getStatusColor, formatStatus } from '@/lib/utils';

export default function ManageSessionsPage() {
  const [sessions, setSessions] = useState<TestSession[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    fetchSessions();
    
    // Set up auto-refresh every 30 seconds
    const interval = setInterval(fetchSessions, 30000);
    return () => clearInterval(interval);
  }, []);

  const fetchSessions = async (showRefreshing = false) => {
    if (showRefreshing) setRefreshing(true);
    
    try {
      // Note: This would need a new API endpoint to get all sessions for teachers
      // For now, we'll simulate with the user's sessions
      const response = await sessionsApi.getMy();
      setSessions(response.data.data || []);
    } catch (error) {
      setError('Failed to load sessions');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  const handleRefresh = () => {
    fetchSessions(true);
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'not_started':
        return <Clock className="h-4 w-4" />;
      case 'in_progress':
        return <Play className="h-4 w-4" />;
      case 'completed':
      case 'submitted':
        return <CheckCircle className="h-4 w-4" />;
      case 'expired':
        return <AlertCircle className="h-4 w-4" />;
      default:
        return <Clock className="h-4 w-4" />;
    }
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
  const expiredSessions = sessions.filter(s => s.status === 'expired');

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
            <h1 className="text-2xl font-bold text-gray-900">Active Sessions</h1>
            <p className="text-gray-600">Monitor ongoing and recent test sessions</p>
          </div>
          <Button 
            variant="outline" 
            onClick={handleRefresh}
            disabled={refreshing}
            className="flex items-center"
          >
            <RefreshCw className={`h-4 w-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
            {refreshing ? 'Refreshing...' : 'Refresh'}
          </Button>
        </div>

        {error && (
          <Card>
            <CardContent className="p-6 text-center">
              <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">Error</h3>
              <p className="text-gray-600">{error}</p>
              <Button onClick={() => fetchSessions()} className="mt-4">
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
                <div className="p-2 bg-red-100 rounded-lg">
                  <AlertCircle className="h-6 w-6 text-red-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Expired</p>
                  <p className="text-2xl font-bold text-gray-900">{expiredSessions.length}</p>
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
                Active Sessions ({activeSessions.length})
              </CardTitle>
              <CardDescription>Students currently taking tests</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {activeSessions.map((session) => {
                  const timeRemaining = getTimeRemaining(session);
                  return (
                    <div key={session.id} className="flex items-center justify-between p-4 border rounded-lg bg-blue-50 border-blue-200">
                      <div className="flex items-center space-x-4">
                        <div className="p-2 bg-blue-100 rounded-lg">
                          <User className="h-5 w-5 text-blue-600" />
                        </div>
                        <div>
                          <h4 className="font-medium text-gray-900">Student #{session.user_id}</h4>
                          <p className="text-sm text-gray-600">Test ID: {session.test_id}</p>
                          <div className="flex items-center mt-1 text-xs text-gray-500">
                            <Clock className="h-3 w-3 mr-1" />
                            <span>Started: {session.started_at ? formatDate(session.started_at) : 'Not started'}</span>
                          </div>
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
                        <div className="text-sm text-gray-600">
                          Question {session.current_question_index + 1}
                        </div>
                        <Button variant="outline" size="sm" className="mt-2">
                          <Eye className="h-3 w-3 mr-1" />
                          Monitor
                        </Button>
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
            <CardDescription>Complete list of test sessions</CardDescription>
          </CardHeader>
          <CardContent>
            {sessions.length === 0 ? (
              <div className="text-center py-8">
                <Activity className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">No Sessions</h3>
                <p className="text-gray-600">No test sessions have been started yet.</p>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b">
                      <th className="text-left py-3 px-4 font-medium text-gray-900">Student</th>
                      <th className="text-left py-3 px-4 font-medium text-gray-900">Test</th>
                      <th className="text-left py-3 px-4 font-medium text-gray-900">Status</th>
                      <th className="text-left py-3 px-4 font-medium text-gray-900">Progress</th>
                      <th className="text-left py-3 px-4 font-medium text-gray-900">Started</th>
                      <th className="text-left py-3 px-4 font-medium text-gray-900">Time Remaining</th>
                      <th className="text-left py-3 px-4 font-medium text-gray-900">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {sessions.map((session) => {
                      const timeRemaining = getTimeRemaining(session);
                      return (
                        <tr key={session.id} className="border-b hover:bg-gray-50">
                          <td className="py-3 px-4">
                            <div className="flex items-center">
                              <User className="h-4 w-4 text-gray-400 mr-2" />
                              <span className="font-medium">Student #{session.user_id}</span>
                            </div>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-gray-900">Test #{session.test_id}</span>
                          </td>
                          <td className="py-3 px-4">
                            <div className="flex items-center">
                              {getStatusIcon(session.status)}
                              <span className={`ml-2 px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(session.status)}`}>
                                {formatStatus(session.status)}
                              </span>
                            </div>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-sm text-gray-600">
                              Question {session.current_question_index + 1}
                            </span>
                          </td>
                          <td className="py-3 px-4 text-sm text-gray-600">
                            {session.started_at ? formatDate(session.started_at) : 'Not started'}
                          </td>
                          <td className="py-3 px-4">
                            {timeRemaining !== null ? (
                              <span className={`font-mono text-sm ${
                                timeRemaining < 300 ? 'text-red-600 font-bold' : 'text-gray-600'
                              }`}>
                                {formatTime(timeRemaining)}
                              </span>
                            ) : (
                              <span className="text-sm text-gray-400">-</span>
                            )}
                          </td>
                          <td className="py-3 px-4">
                            <Button variant="outline" size="sm">
                              <Eye className="h-3 w-3 mr-1" />
                              View
                            </Button>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </Layout>
  );
}
