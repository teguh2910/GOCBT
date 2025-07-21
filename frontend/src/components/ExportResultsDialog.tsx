'use client';

import React, { useState } from 'react';
import { Button } from '@/components/ui/Button';
import { Test, TestResult } from '@/lib/api';
import { X, Download, FileText, BarChart3 } from 'lucide-react';
import { formatDate } from '@/lib/utils';

interface ExportResultsDialogProps {
  isOpen: boolean;
  onClose: () => void;
  test: Test;
  results: TestResult[];
}

export function ExportResultsDialog({ isOpen, onClose, test, results }: ExportResultsDialogProps) {
  const [exportFormat, setExportFormat] = useState<'csv' | 'detailed'>('csv');
  const [includeStatistics, setIncludeStatistics] = useState(true);
  const [exporting, setExporting] = useState(false);



  const exportBasicCSV = () => {
    const headers = [
      'Student ID',
      'Questions Answered',
      'Total Questions',
      'Marks Obtained',
      'Total Marks',
      'Percentage',
      'Grade',
      'Status',
      'Time Taken',
      'Completed At'
    ];

    const csvData = results.map(result => {
      const timeTaken = result.time_taken ? `${Math.floor(result.time_taken / 60)}:${(result.time_taken % 60).toString().padStart(2, '0')}` : 'N/A';

      return [
        `Student #${result.user_id}`,
        result.answered_questions,
        result.total_questions,
        result.marks_obtained,
        result.total_marks,
        `${result.percentage.toFixed(1)}%`,
        result.grade || 'F',
        result.is_passed ? 'Pass' : 'Fail',
        timeTaken,
        formatDate(result.completed_at)
      ];
    });

    return [headers, ...csvData];
  };

  const exportDetailedCSV = () => {
    // Calculate statistics
    const totalStudents = results.length;
    const completedResults = results.filter(r => r.completed_at);
    const passedStudents = results.filter(r => r.is_passed).length;
    const averageScore = results.reduce((sum, r) => sum + r.marks_obtained, 0) / totalStudents;
    const averagePercentage = results.reduce((sum, r) => sum + r.percentage, 0) / totalStudents;

    const statisticsData = [
      ['TEST STATISTICS'],
      ['Test Title', test.title],
      ['Test Description', test.description],
      ['Total Marks', test.total_marks],
      ['Passing Marks', test.passing_marks],
      ['Duration (minutes)', test.duration_minutes],
      [''],
      ['RESULTS SUMMARY'],
      ['Total Students', totalStudents],
      ['Completed', completedResults.length],
      ['Passed', passedStudents],
      ['Pass Rate', `${((passedStudents / totalStudents) * 100).toFixed(1)}%`],
      ['Average Score', averageScore.toFixed(1)],
      ['Average Percentage', `${averagePercentage.toFixed(1)}%`],
      [''],
      ['DETAILED RESULTS'],
      ...exportBasicCSV()
    ];

    return statisticsData;
  };

  const handleExport = () => {
    if (results.length === 0) {
      alert('No results to export');
      return;
    }

    setExporting(true);

    try {
      let csvData;
      let filename;

      if (exportFormat === 'csv') {
        csvData = exportBasicCSV();
        filename = `${test.title}_results_${new Date().toISOString().split('T')[0]}.csv`;
      } else {
        csvData = exportDetailedCSV();
        filename = `${test.title}_detailed_report_${new Date().toISOString().split('T')[0]}.csv`;
      }

      // Create CSV content
      const csvContent = csvData.map(row => 
        row.map(cell => `"${cell}"`).join(',')
      ).join('\n');

      // Create and download file
      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      const url = URL.createObjectURL(blob);
      
      link.setAttribute('href', url);
      link.setAttribute('download', filename);
      link.style.visibility = 'hidden';
      
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      
      URL.revokeObjectURL(url);
      onClose();
    } catch (error) {
      console.error('Export failed:', error);
      alert('Failed to export results. Please try again.');
    } finally {
      setExporting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full">
        <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Export Results</h2>
          <Button variant="ghost" size="sm" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>

        <div className="p-6 space-y-6">
          <div>
            <h3 className="text-sm font-medium text-gray-900 dark:text-white mb-3">Export Format</h3>
            <div className="space-y-3">
              <label className="flex items-center space-x-3 cursor-pointer">
                <input
                  type="radio"
                  name="exportFormat"
                  value="csv"
                  checked={exportFormat === 'csv'}
                  onChange={(e) => setExportFormat(e.target.value as 'csv')}
                  className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
                />
                <div className="flex items-center space-x-2">
                  <FileText className="h-4 w-4 text-gray-500" />
                  <div>
                    <div className="text-sm font-medium text-gray-900 dark:text-white">Basic CSV</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">Student results only</div>
                  </div>
                </div>
              </label>

              <label className="flex items-center space-x-3 cursor-pointer">
                <input
                  type="radio"
                  name="exportFormat"
                  value="detailed"
                  checked={exportFormat === 'detailed'}
                  onChange={(e) => setExportFormat(e.target.value as 'detailed')}
                  className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
                />
                <div className="flex items-center space-x-2">
                  <BarChart3 className="h-4 w-4 text-gray-500" />
                  <div>
                    <div className="text-sm font-medium text-gray-900 dark:text-white">Detailed Report</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">Includes test statistics and summary</div>
                  </div>
                </div>
              </label>
            </div>
          </div>

          <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
            <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-2">Export Summary</h4>
            <div className="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <div>Test: {test.title}</div>
              <div>Students: {results.length}</div>
              <div>Format: {exportFormat === 'csv' ? 'Basic CSV' : 'Detailed Report'}</div>
            </div>
          </div>
        </div>

        <div className="flex items-center justify-end space-x-4 p-6 border-t border-gray-200 dark:border-gray-700">
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={handleExport} disabled={exporting}>
            <Download className="h-4 w-4 mr-2" />
            {exporting ? 'Exporting...' : 'Export'}
          </Button>
        </div>
      </div>
    </div>
  );
}
