'use client';

import React, { useEffect, useState } from 'react';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { CreateTestDialog } from '@/components/CreateTestDialog';
import { EditTestDialog } from '@/components/EditTestDialog';
import { testsApi, Test } from '@/lib/api';
import { 
  Plus, 
  Edit, 
  Trash2, 
  Eye, 
  Clock, 
  Users, 
  FileText, 
  AlertCircle,
  Calendar
} from 'lucide-react';
import { formatDuration, formatDate } from '@/lib/utils';
import Link from 'next/link';

export default function ManageTestsPage() {
  const [tests, setTests] = useState<Test[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [showEditDialog, setShowEditDialog] = useState(false);
  const [selectedTest, setSelectedTest] = useState<Test | null>(null);

  useEffect(() => {
    fetchTests();
  }, []);

  const fetchTests = async () => {
    try {
      const response = await testsApi.getAll({ limit: 50 });
      setTests(response.data.data || []);
    } catch (error) {
      setError('Failed to load tests');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteTest = async (testId: number) => {
    if (!confirm('Are you sure you want to delete this test? This action cannot be undone.')) {
      return;
    }

    try {
      await testsApi.delete(testId);
      setTests(tests.filter(test => test.id !== testId));
    } catch (error) {
      alert('Failed to delete test');
    }
  };

  const handleTestCreated = (newTest: Test) => {
    setTests(prev => [newTest, ...prev]);
  };

  const handleEditTest = (test: Test) => {
    setSelectedTest(test);
    setShowEditDialog(true);
  };

  const handleTestUpdated = (updatedTest: Test) => {
    setTests(prev => prev.map(test =>
      test.id === updatedTest.id ? updatedTest : test
    ));
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
            <h1 className="text-2xl font-bold text-gray-900">Manage Tests</h1>
            <p className="text-gray-600">Create, edit, and manage your tests</p>
          </div>
          <Button
            className="flex items-center"
            onClick={() => setShowCreateDialog(true)}
          >
            <Plus className="h-4 w-4 mr-2" />
            Create New Test
          </Button>
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
              <h3 className="text-lg font-medium text-gray-900 mb-2">No Tests Created</h3>
              <p className="text-gray-600 mb-4">
                You haven't created any tests yet. Create your first test to get started.
              </p>
              <Button
                className="flex items-center mx-auto"
                onClick={() => setShowCreateDialog(true)}
              >
                <Plus className="h-4 w-4 mr-2" />
                Create Your First Test
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 gap-6">
            {tests.map((test) => (
              <Card key={test.id} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div>
                      <CardTitle className="text-lg">{test.title}</CardTitle>
                      <CardDescription>{test.description}</CardDescription>
                    </div>
                    <div className={`px-2 py-1 rounded-full text-xs font-medium ${
                      test.is_active 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-gray-100 text-gray-800'
                    }`}>
                      {test.is_active ? 'Active' : 'Inactive'}
                    </div>
                  </div>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div className="flex items-center text-gray-600">
                      <Clock className="h-4 w-4 mr-2" />
                      <span>{formatDuration(test.duration_minutes)}</span>
                    </div>
                    <div className="flex items-center text-gray-600">
                      <FileText className="h-4 w-4 mr-2" />
                      <span>{test.total_marks} marks</span>
                    </div>
                    <div className="flex items-center text-gray-600">
                      <Users className="h-4 w-4 mr-2" />
                      <span>Pass: {test.passing_marks}</span>
                    </div>
                    <div className="flex items-center text-gray-600">
                      <Calendar className="h-4 w-4 mr-2" />
                      <span>{formatDate(test.created_at)}</span>
                    </div>
                  </div>

                  {(test.start_time || test.end_time) && (
                    <div className="text-xs text-gray-500 space-y-1">
                      {test.start_time && (
                        <div>Available from: {formatDate(test.start_time)}</div>
                      )}
                      {test.end_time && (
                        <div>Available until: {formatDate(test.end_time)}</div>
                      )}
                    </div>
                  )}

                  <div className="flex items-center justify-between pt-4 border-t">
                    <div className="flex space-x-2">
                      <Link href={`/manage/tests/${test.id}`}>
                        <Button variant="outline" size="sm" className="flex items-center">
                          <Eye className="h-4 w-4 mr-1" />
                          View
                        </Button>
                      </Link>
                      <Button
                        variant="outline"
                        size="sm"
                        className="flex items-center"
                        onClick={() => handleEditTest(test)}
                      >
                        <Edit className="h-4 w-4 mr-1" />
                        Edit
                      </Button>
                      <Link href={`/manage/results?test=${test.id}`}>
                        <Button variant="outline" size="sm" className="flex items-center">
                          <FileText className="h-4 w-4 mr-1" />
                          Results
                        </Button>
                      </Link>
                    </div>
                    <Button 
                      variant="destructive" 
                      size="sm" 
                      onClick={() => handleDeleteTest(test.id)}
                      className="flex items-center"
                    >
                      <Trash2 className="h-4 w-4 mr-1" />
                      Delete
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {/* Create Test Dialog */}
        <CreateTestDialog
          isOpen={showCreateDialog}
          onClose={() => setShowCreateDialog(false)}
          onTestCreated={handleTestCreated}
        />

        {/* Edit Test Dialog */}
        <EditTestDialog
          isOpen={showEditDialog}
          onClose={() => {
            setShowEditDialog(false);
            setSelectedTest(null);
          }}
          test={selectedTest}
          onTestUpdated={handleTestUpdated}
        />
      </div>
    </Layout>
  );
}
