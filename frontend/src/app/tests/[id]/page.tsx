'use client';

import React, { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { 
  testsApi, 
  sessionsApi, 
  Test, 
  Question, 
  TestSession, 
  UserAnswer 
} from '@/lib/api';
import { 
  Clock, 
  CheckCircle, 
  AlertCircle, 
  ArrowLeft, 
  ArrowRight,
  Flag,
  Send
} from 'lucide-react';
import { formatTime, formatDuration } from '@/lib/utils';

export default function TestPage() {
  const params = useParams();
  const router = useRouter();
  const testId = parseInt(params.id as string);

  const [test, setTest] = useState<Test | null>(null);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [session, setSession] = useState<TestSession | null>(null);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [answers, setAnswers] = useState<Record<number, UserAnswer>>({});
  const [timeRemaining, setTimeRemaining] = useState<number>(0);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  // Current answer state
  const [selectedOption, setSelectedOption] = useState<number | null>(null);
  const [textAnswer, setTextAnswer] = useState('');

  useEffect(() => {
    fetchTestData();
  }, [testId]);

  useEffect(() => {
    if (session && session.status === 'in_progress') {
      const timer = setInterval(() => {
        setTimeRemaining(prev => {
          if (prev <= 1) {
            handleSubmitTest();
            return 0;
          }
          return prev - 1;
        });
      }, 1000);

      return () => clearInterval(timer);
    }
  }, [session]);

  useEffect(() => {
    // Load current answer when question changes
    const currentQuestion = questions[currentQuestionIndex];
    if (currentQuestion && answers[currentQuestion.id]) {
      const answer = answers[currentQuestion.id];
      if (answer.selected_option_id) {
        setSelectedOption(answer.selected_option_id);
      }
      if (answer.answer_text) {
        setTextAnswer(answer.answer_text);
      }
    } else {
      setSelectedOption(null);
      setTextAnswer('');
    }
  }, [currentQuestionIndex, questions, answers]);

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

  const startTest = async () => {
    try {
      const response = await sessionsApi.start(testId);
      const newSession = response.data.data;
      setSession(newSession);
      setTimeRemaining(newSession.remaining_time_seconds || 0);
    } catch (error) {
      setError('Failed to start test');
    }
  };

  const submitAnswer = async () => {
    if (!session || !questions[currentQuestionIndex]) return;

    const currentQuestion = questions[currentQuestionIndex];
    const answerData: any = {
      question_id: currentQuestion.id,
    };

    if (currentQuestion.question_type === 'multiple_choice' || currentQuestion.question_type === 'true_false') {
      if (selectedOption) {
        answerData.selected_option_id = selectedOption;
      }
    } else {
      answerData.answer_text = textAnswer;
    }

    try {
      const response = await sessionsApi.submitAnswer(session.session_token, answerData);
      const answer = response.data.data;
      
      setAnswers(prev => ({
        ...prev,
        [currentQuestion.id]: answer,
      }));

      // Update progress
      await sessionsApi.updateProgress(session.session_token, currentQuestionIndex);
    } catch (error) {
      console.error('Failed to submit answer:', error);
    }
  };

  const nextQuestion = async () => {
    await submitAnswer();
    if (currentQuestionIndex < questions.length - 1) {
      setCurrentQuestionIndex(currentQuestionIndex + 1);
    }
  };

  const previousQuestion = () => {
    if (currentQuestionIndex > 0) {
      setCurrentQuestionIndex(currentQuestionIndex - 1);
    }
  };

  const handleSubmitTest = async () => {
    if (!session) return;

    setSubmitting(true);
    try {
      // Submit current answer first
      await submitAnswer();
      
      // Submit the test
      await sessionsApi.submit(session.session_token);
      
      // Redirect to results
      router.push(`/results/session/${session.id}`);
    } catch (error) {
      setError('Failed to submit test');
      setSubmitting(false);
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

  if (error) {
    return (
      <Layout>
        <Card>
          <CardContent className="p-6 text-center">
            <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">Error</h3>
            <p className="text-gray-600">{error}</p>
            <Button onClick={() => router.back()} className="mt-4">
              Go Back
            </Button>
          </CardContent>
        </Card>
      </Layout>
    );
  }

  if (!test) {
    return (
      <Layout>
        <Card>
          <CardContent className="p-6 text-center">
            <h3 className="text-lg font-medium text-gray-900">Test not found</h3>
          </CardContent>
        </Card>
      </Layout>
    );
  }

  // Test not started yet
  if (!session) {
    return (
      <Layout>
        <div className="max-w-2xl mx-auto">
          <Card>
            <CardHeader>
              <CardTitle>{test.title}</CardTitle>
              <CardDescription>{test.description}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="grid grid-cols-2 gap-4">
                <div className="text-center p-4 bg-gray-50 rounded-lg">
                  <Clock className="h-8 w-8 text-blue-600 mx-auto mb-2" />
                  <p className="text-sm text-gray-600">Duration</p>
                  <p className="font-medium">{formatDuration(test.duration_minutes)}</p>
                </div>
                <div className="text-center p-4 bg-gray-50 rounded-lg">
                  <CheckCircle className="h-8 w-8 text-green-600 mx-auto mb-2" />
                  <p className="text-sm text-gray-600">Total Marks</p>
                  <p className="font-medium">{test.total_marks}</p>
                </div>
              </div>

              {test.instructions && (
                <div>
                  <h4 className="font-medium mb-2">Instructions:</h4>
                  <p className="text-gray-600 whitespace-pre-wrap">{test.instructions}</p>
                </div>
              )}

              <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                <div className="flex">
                  <AlertCircle className="h-5 w-5 text-yellow-600 mt-0.5" />
                  <div className="ml-3">
                    <h4 className="text-sm font-medium text-yellow-800">Important Notes:</h4>
                    <ul className="mt-2 text-sm text-yellow-700 list-disc list-inside space-y-1">
                      <li>Once started, the timer cannot be paused</li>
                      <li>Your answers are automatically saved</li>
                      <li>You can navigate between questions freely</li>
                      <li>Submit the test before time runs out</li>
                    </ul>
                  </div>
                </div>
              </div>

              <div className="flex space-x-4">
                <Button variant="outline" onClick={() => router.back()}>
                  Cancel
                </Button>
                <Button onClick={startTest} className="flex-1">
                  Start Test
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </Layout>
    );
  }

  const currentQuestion = questions[currentQuestionIndex];
  const progress = ((currentQuestionIndex + 1) / questions.length) * 100;

  return (
    <Layout>
      <div className="max-w-4xl mx-auto">
        {/* Header with timer and progress */}
        <div className="bg-white border-b sticky top-0 z-10 p-4 mb-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-xl font-bold">{test.title}</h1>
              <p className="text-sm text-gray-600">
                Question {currentQuestionIndex + 1} of {questions.length}
              </p>
            </div>
            <div className="flex items-center space-x-4">
              <div className="text-right">
                <p className="text-sm text-gray-600">Time Remaining</p>
                <p className={`font-mono text-lg font-bold ${
                  timeRemaining < 300 ? 'text-red-600' : 'text-gray-900'
                }`}>
                  {formatTime(timeRemaining)}
                </p>
              </div>
              <Button 
                onClick={handleSubmitTest} 
                disabled={submitting}
                variant="destructive"
                size="sm"
              >
                <Send className="h-4 w-4 mr-2" />
                {submitting ? 'Submitting...' : 'Submit Test'}
              </Button>
            </div>
          </div>
          
          {/* Progress bar */}
          <div className="mt-4">
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div 
                className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                style={{ width: `${progress}%` }}
              />
            </div>
          </div>
        </div>

        {/* Question content */}
        {currentQuestion && (
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">
                {currentQuestion.question_text}
              </CardTitle>
              <CardDescription>
                {currentQuestion.marks} mark{currentQuestion.marks !== 1 ? 's' : ''}
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Multiple Choice / True False */}
              {(currentQuestion.question_type === 'multiple_choice' || 
                currentQuestion.question_type === 'true_false') && (
                <div className="space-y-3">
                  {currentQuestion.options?.map((option) => (
                    <label
                      key={option.id}
                      className={`flex items-center p-4 border rounded-lg cursor-pointer transition-colors ${
                        selectedOption === option.id
                          ? 'border-blue-500 bg-blue-50'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}
                    >
                      <input
                        type="radio"
                        name="answer"
                        value={option.id}
                        checked={selectedOption === option.id}
                        onChange={() => setSelectedOption(option.id)}
                        className="sr-only"
                      />
                      <div className={`w-4 h-4 rounded-full border-2 mr-3 ${
                        selectedOption === option.id
                          ? 'border-blue-500 bg-blue-500'
                          : 'border-gray-300'
                      }`}>
                        {selectedOption === option.id && (
                          <div className="w-2 h-2 bg-white rounded-full mx-auto mt-0.5" />
                        )}
                      </div>
                      <span>{option.option_text}</span>
                    </label>
                  ))}
                </div>
              )}

              {/* Short Answer */}
              {currentQuestion.question_type === 'short_answer' && (
                <div>
                  <Input
                    value={textAnswer}
                    onChange={(e) => setTextAnswer(e.target.value)}
                    placeholder="Enter your answer..."
                    className="w-full"
                  />
                </div>
              )}

              {/* Navigation */}
              <div className="flex items-center justify-between pt-6 border-t">
                <Button
                  variant="outline"
                  onClick={previousQuestion}
                  disabled={currentQuestionIndex === 0}
                >
                  <ArrowLeft className="h-4 w-4 mr-2" />
                  Previous
                </Button>

                <div className="text-sm text-gray-500">
                  {answers[currentQuestion.id] && (
                    <span className="flex items-center text-green-600">
                      <CheckCircle className="h-4 w-4 mr-1" />
                      Answered
                    </span>
                  )}
                </div>

                {currentQuestionIndex < questions.length - 1 ? (
                  <Button onClick={nextQuestion}>
                    Next
                    <ArrowRight className="h-4 w-4 ml-2" />
                  </Button>
                ) : (
                  <Button onClick={handleSubmitTest} disabled={submitting}>
                    <Send className="h-4 w-4 mr-2" />
                    {submitting ? 'Submitting...' : 'Submit Test'}
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        )}

        {/* Question navigation */}
        <Card className="mt-6">
          <CardHeader>
            <CardTitle className="text-base">Question Navigation</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-10 gap-2">
              {questions.map((question, index) => (
                <button
                  key={question.id}
                  onClick={() => setCurrentQuestionIndex(index)}
                  className={`w-10 h-10 rounded-lg text-sm font-medium transition-colors ${
                    index === currentQuestionIndex
                      ? 'bg-blue-600 text-white'
                      : answers[question.id]
                      ? 'bg-green-100 text-green-800 border border-green-300'
                      : 'bg-gray-100 text-gray-600 border border-gray-300 hover:bg-gray-200'
                  }`}
                >
                  {index + 1}
                </button>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </Layout>
  );
}
