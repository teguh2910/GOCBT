-- GoCBT Sample Data
-- This file contains sample data for demonstration and testing purposes

-- Insert sample users
INSERT INTO users (username, email, password_hash, first_name, last_name, role, created_at) VALUES
-- Admin user
('admin', 'admin@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'System', 'Administrator', 'admin', datetime('now')),

-- Teachers
('teacher1', 'teacher1@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Sarah', 'Johnson', 'teacher', datetime('now')),
('teacher2', 'teacher2@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Michael', 'Brown', 'teacher', datetime('now')),
('teacher3', 'teacher3@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Emily', 'Davis', 'teacher', datetime('now')),

-- Students
('student1', 'student1@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'John', 'Smith', 'student', datetime('now')),
('student2', 'student2@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Emma', 'Wilson', 'student', datetime('now')),
('student3', 'student3@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'James', 'Taylor', 'student', datetime('now')),
('student4', 'student4@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Olivia', 'Anderson', 'student', datetime('now')),
('student5', 'student5@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'William', 'Thomas', 'student', datetime('now')),
('student6', 'student6@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Sophia', 'Jackson', 'student', datetime('now')),
('student7', 'student7@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Benjamin', 'White', 'student', datetime('now')),
('student8', 'student8@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Isabella', 'Harris', 'student', datetime('now')),
('student9', 'student9@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Lucas', 'Martin', 'student', datetime('now')),
('student10', 'student10@gocbt.example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Mia', 'Garcia', 'student', datetime('now'));

-- Insert sample tests
INSERT INTO tests (title, description, instructions, duration_minutes, total_marks, passing_marks, start_time, end_time, created_by, created_at) VALUES
('Mathematics Quiz - Algebra', 'Basic algebra concepts and problem solving', 'Read each question carefully. Show your work where applicable. You have 60 minutes to complete this test.', 60, 100, 60, datetime('now', '+1 hour'), datetime('now', '+7 days'), 2, datetime('now')),

('Science Test - Physics', 'Fundamental physics principles and applications', 'Answer all questions. Use the formula sheet provided. Calculator is allowed for numerical problems.', 90, 150, 90, datetime('now', '+2 hours'), datetime('now', '+5 days'), 2, datetime('now')),

('English Literature Quiz', 'Reading comprehension and literary analysis', 'Read the passages carefully before answering. Support your answers with evidence from the text.', 45, 75, 45, datetime('now', '+30 minutes'), datetime('now', '+3 days'), 3, datetime('now')),

('History Test - World War II', 'Major events and consequences of World War II', 'Answer all questions in complete sentences. Provide specific dates and examples where possible.', 75, 120, 72, datetime('now', '+1 day'), datetime('now', '+10 days'), 3, datetime('now')),

('Computer Science - Programming Basics', 'Introduction to programming concepts', 'Code examples are provided. Trace through the code carefully before answering.', 50, 80, 48, datetime('now', '+3 hours'), datetime('now', '+6 days'), 4, datetime('now'));

-- Insert sample questions for Mathematics Quiz
INSERT INTO questions (test_id, question_text, question_type, marks, created_at) VALUES
-- Math Test Questions
(1, 'What is the value of x in the equation 2x + 5 = 13?', 'multiple_choice', 5, datetime('now')),
(1, 'Simplify the expression: 3(x + 4) - 2x', 'multiple_choice', 5, datetime('now')),
(1, 'If y = 2x + 3, what is the value of y when x = 4?', 'multiple_choice', 5, datetime('now')),
(1, 'Solve for x: x² - 5x + 6 = 0', 'multiple_choice', 10, datetime('now')),
(1, 'The slope of a line passing through points (2, 3) and (4, 7) is:', 'multiple_choice', 5, datetime('now')),

-- Physics Test Questions
(2, 'What is the SI unit of force?', 'multiple_choice', 5, datetime('now')),
(2, 'According to Newton''s first law, an object at rest will:', 'multiple_choice', 5, datetime('now')),
(2, 'The acceleration due to gravity on Earth is approximately:', 'multiple_choice', 5, datetime('now')),
(2, 'If a car travels 100 meters in 10 seconds, its average speed is:', 'multiple_choice', 10, datetime('now')),
(2, 'The formula for kinetic energy is:', 'multiple_choice', 10, datetime('now')),

-- English Literature Questions
(3, 'In Shakespeare''s Romeo and Juliet, the two families are:', 'multiple_choice', 5, datetime('now')),
(3, 'What literary device is used in "The wind whispered through the trees"?', 'multiple_choice', 5, datetime('now')),
(3, 'The main theme of George Orwell''s "1984" is:', 'multiple_choice', 10, datetime('now')),

-- History Questions
(4, 'World War II began in which year?', 'multiple_choice', 5, datetime('now')),
(4, 'The D-Day invasion took place in:', 'multiple_choice', 10, datetime('now')),
(4, 'Which country was NOT part of the Axis powers?', 'multiple_choice', 5, datetime('now')),

-- Computer Science Questions
(5, 'Which of the following is a programming language?', 'multiple_choice', 5, datetime('now')),
(5, 'What does "HTML" stand for?', 'multiple_choice', 5, datetime('now')),
(5, 'In programming, what is a variable?', 'multiple_choice', 10, datetime('now'));

-- Insert sample question options
INSERT INTO question_options (question_id, option_text, is_correct, created_at) VALUES
-- Math Question 1 options
(1, '3', false, datetime('now')),
(1, '4', true, datetime('now')),
(1, '5', false, datetime('now')),
(1, '6', false, datetime('now')),

-- Math Question 2 options
(2, 'x + 12', true, datetime('now')),
(2, '5x + 12', false, datetime('now')),
(2, 'x + 6', false, datetime('now')),
(2, '3x + 2', false, datetime('now')),

-- Math Question 3 options
(3, '8', false, datetime('now')),
(3, '10', false, datetime('now')),
(3, '11', true, datetime('now')),
(3, '12', false, datetime('now')),

-- Math Question 4 options
(4, 'x = 2, x = 3', true, datetime('now')),
(4, 'x = 1, x = 6', false, datetime('now')),
(4, 'x = -2, x = -3', false, datetime('now')),
(4, 'No real solutions', false, datetime('now')),

-- Math Question 5 options
(5, '1', false, datetime('now')),
(5, '2', true, datetime('now')),
(5, '3', false, datetime('now')),
(5, '4', false, datetime('now')),

-- Physics Question 1 options
(6, 'Joule', false, datetime('now')),
(6, 'Newton', true, datetime('now')),
(6, 'Watt', false, datetime('now')),
(6, 'Pascal', false, datetime('now')),

-- Physics Question 2 options
(7, 'Remain at rest unless acted upon by a force', true, datetime('now')),
(7, 'Accelerate continuously', false, datetime('now')),
(7, 'Move at constant velocity', false, datetime('now')),
(7, 'Decelerate gradually', false, datetime('now')),

-- Physics Question 3 options
(8, '9.8 m/s²', true, datetime('now')),
(8, '10 m/s²', false, datetime('now')),
(8, '9.6 m/s²', false, datetime('now')),
(8, '8.9 m/s²', false, datetime('now')),

-- Physics Question 4 options
(9, '5 m/s', false, datetime('now')),
(9, '10 m/s', true, datetime('now')),
(9, '15 m/s', false, datetime('now')),
(9, '20 m/s', false, datetime('now')),

-- Physics Question 5 options
(10, 'KE = mv²', false, datetime('now')),
(10, 'KE = ½mv²', true, datetime('now')),
(10, 'KE = 2mv²', false, datetime('now')),
(10, 'KE = m²v', false, datetime('now')),

-- English Question 1 options
(11, 'Montagues and Capulets', true, datetime('now')),
(11, 'Smiths and Jones', false, datetime('now')),
(11, 'Hatfields and McCoys', false, datetime('now')),
(11, 'Lancasters and Yorks', false, datetime('now')),

-- English Question 2 options
(12, 'Metaphor', false, datetime('now')),
(12, 'Personification', true, datetime('now')),
(12, 'Simile', false, datetime('now')),
(12, 'Alliteration', false, datetime('now')),

-- English Question 3 options
(13, 'Love and romance', false, datetime('now')),
(13, 'Totalitarianism and surveillance', true, datetime('now')),
(13, 'Adventure and exploration', false, datetime('now')),
(13, 'Comedy and humor', false, datetime('now')),

-- History Question 1 options
(14, '1938', false, datetime('now')),
(14, '1939', true, datetime('now')),
(14, '1940', false, datetime('now')),
(14, '1941', false, datetime('now')),

-- History Question 2 options
(15, 'Normandy, France', true, datetime('now')),
(15, 'Sicily, Italy', false, datetime('now')),
(15, 'Calais, France', false, datetime('now')),
(15, 'Dover, England', false, datetime('now')),

-- History Question 3 options
(16, 'Germany', false, datetime('now')),
(16, 'Italy', false, datetime('now')),
(16, 'Japan', false, datetime('now')),
(16, 'Soviet Union', true, datetime('now')),

-- Computer Science Question 1 options
(17, 'Microsoft Word', false, datetime('now')),
(17, 'Python', true, datetime('now')),
(17, 'Adobe Photoshop', false, datetime('now')),
(17, 'Google Chrome', false, datetime('now')),

-- Computer Science Question 2 options
(18, 'Home Text Markup Language', false, datetime('now')),
(18, 'HyperText Markup Language', true, datetime('now')),
(18, 'High Tech Modern Language', false, datetime('now')),
(18, 'Hyperlink and Text Markup Language', false, datetime('now')),

-- Computer Science Question 3 options
(19, 'A fixed value that cannot change', false, datetime('now')),
(19, 'A storage location with a name and value that can change', true, datetime('now')),
(19, 'A type of computer virus', false, datetime('now')),
(19, 'A programming error', false, datetime('now'));

-- Insert sample test sessions (some completed, some active)
INSERT INTO test_sessions (test_id, user_id, session_token, status, started_at, submitted_at, created_at) VALUES
(1, 5, 'sess_math_john_001', 'completed', datetime('now', '-2 days'), datetime('now', '-2 days', '+45 minutes'), datetime('now', '-2 days')),
(1, 6, 'sess_math_emma_001', 'completed', datetime('now', '-1 day'), datetime('now', '-1 day', '+50 minutes'), datetime('now', '-1 day')),
(2, 5, 'sess_physics_john_001', 'completed', datetime('now', '-3 hours'), datetime('now', '-2 hours'), datetime('now', '-3 hours')),
(3, 7, 'sess_english_james_001', 'active', datetime('now', '-30 minutes'), NULL, datetime('now', '-30 minutes')),
(1, 8, 'sess_math_olivia_001', 'completed', datetime('now', '-4 hours'), datetime('now', '-3 hours'), datetime('now', '-4 hours'));

-- Insert sample test results
INSERT INTO test_results (test_id, user_id, session_token, score, total_marks, percentage, grade, is_passed, completed_at, created_at) VALUES
(1, 5, 'sess_math_john_001', 85, 100, 85.0, 'B', true, datetime('now', '-2 days', '+45 minutes'), datetime('now', '-2 days', '+45 minutes')),
(1, 6, 'sess_math_emma_001', 92, 100, 92.0, 'A', true, datetime('now', '-1 day', '+50 minutes'), datetime('now', '-1 day', '+50 minutes')),
(2, 5, 'sess_physics_john_001', 78, 150, 52.0, 'C', false, datetime('now', '-2 hours'), datetime('now', '-2 hours')),
(1, 8, 'sess_math_olivia_001', 95, 100, 95.0, 'A', true, datetime('now', '-3 hours'), datetime('now', '-3 hours'));

-- Insert sample session answers
INSERT INTO session_answers (session_id, question_id, selected_option_id, answer_text, is_correct, marks_awarded, answered_at) VALUES
-- John's Math Test answers
(1, 1, 2, NULL, true, 5, datetime('now', '-2 days', '+5 minutes')),
(1, 2, 1, NULL, true, 5, datetime('now', '-2 days', '+10 minutes')),
(1, 3, 3, NULL, true, 5, datetime('now', '-2 days', '+15 minutes')),
(1, 4, 1, NULL, true, 10, datetime('now', '-2 days', '+25 minutes')),
(1, 5, 2, NULL, true, 5, datetime('now', '-2 days', '+30 minutes')),

-- Emma's Math Test answers
(2, 1, 2, NULL, true, 5, datetime('now', '-1 day', '+3 minutes')),
(2, 2, 1, NULL, true, 5, datetime('now', '-1 day', '+8 minutes')),
(2, 3, 3, NULL, true, 5, datetime('now', '-1 day', '+12 minutes')),
(2, 4, 1, NULL, true, 10, datetime('now', '-1 day', '+20 minutes')),
(2, 5, 2, NULL, true, 5, datetime('now', '-1 day', '+25 minutes'));

-- Note: Password for all sample users is "password123" (hashed with bcrypt)
-- Default login credentials:
-- Admin: admin / password123
-- Teachers: teacher1, teacher2, teacher3 / password123
-- Students: student1-student10 / password123
