# GoCBT Frequently Asked Questions (FAQ)

## 🎓 For Students

### Getting Started

**Q: How do I log into GoCBT?**
A: Use the username and password provided by your teacher or administrator. Navigate to the login page and enter your credentials.

**Q: I forgot my password. What should I do?**
A: Contact your teacher or system administrator to reset your password. For security reasons, only they can reset student passwords.

**Q: Can I change my password?**
A: Yes, you can change your password from your profile settings after logging in.

### Taking Tests

**Q: Can I pause a test and continue later?**
A: No, once you start a test, you must complete it in one session. However, your answers are automatically saved as you progress.

**Q: What happens if my internet connection is lost during a test?**
A: Don't panic! Your answers are automatically saved. Simply refresh the page and log back in to continue where you left off.

**Q: Can I go back to previous questions?**
A: Yes, you can navigate between questions using the "Previous" and "Next" buttons, as long as the test allows it.

**Q: What happens when time runs out?**
A: The test will automatically submit when time expires. Make sure to manage your time effectively.

**Q: Can I retake a test if I'm not satisfied with my score?**
A: This depends on your teacher's settings. Most tests allow only one attempt to maintain academic integrity.

### Technical Issues

**Q: Which browsers are supported?**
A: GoCBT works best with modern browsers: Chrome 90+, Firefox 88+, Safari 14+, and Edge 90+. Internet Explorer is not supported.

**Q: The test won't load. What should I do?**
A: Try these steps:
1. Refresh the page
2. Clear your browser cache
3. Check your internet connection
4. Contact your teacher if the problem persists

**Q: Can I use my mobile phone or tablet?**
A: While GoCBT is responsive and works on mobile devices, we recommend using a computer or laptop for the best experience.

## 👨‍🏫 For Teachers

### Test Creation

**Q: How many questions can I add to a test?**
A: There's no strict limit, but we recommend 20-50 questions for optimal performance and student experience.

**Q: Can I import questions from other sources?**
A: Currently, questions must be entered manually through the interface. Bulk import features may be added in future versions.

**Q: Can I reuse questions across different tests?**
A: Yes, you can copy questions from existing tests when creating new ones.

**Q: How do I set different time limits for different students?**
A: You can modify individual student settings from the test monitoring dashboard during active sessions.

### Student Management

**Q: How do I add students to my class?**
A: Students must be created by an administrator. Once created, you can assign them to your tests.

**Q: Can I see students' progress in real-time?**
A: Yes, the monitoring dashboard shows live progress for all students currently taking tests.

**Q: How do I handle technical issues during a test?**
A: You can extend time, reset sessions, or provide additional attempts from the monitoring dashboard.

### Results and Grading

**Q: Are tests graded automatically?**
A: Multiple choice and true/false questions are graded automatically. Short answer questions may require manual review.

**Q: Can I modify grades after a test is submitted?**
A: Yes, you can adjust individual question scores from the results review page.

**Q: How do I export results?**
A: Use the "Export" button on the results page to download data in CSV or PDF format.

## 🔧 For Administrators

### System Setup

**Q: What are the minimum system requirements?**
A: 
- Server: 2GB RAM, 10GB storage, Linux/Windows/macOS
- Database: SQLite (development) or PostgreSQL (production)
- Clients: Modern web browsers

**Q: Can GoCBT handle multiple concurrent users?**
A: Yes, GoCBT is designed to handle hundreds of concurrent users. Performance depends on your server specifications.

**Q: How do I backup the database?**
A: For SQLite, copy the database file. For PostgreSQL, use `pg_dump`. We recommend automated daily backups.

### User Management

**Q: How do I create user accounts in bulk?**
A: Use the CSV import feature in the admin panel to create multiple users at once.

**Q: Can users have multiple roles?**
A: No, each user has one role: admin, teacher, or student. However, admins have access to all features.

**Q: How do I reset a user's password?**
A: Go to User Management, find the user, and click "Reset Password" to generate a new temporary password.

### Security

**Q: How secure is GoCBT?**
A: GoCBT implements industry-standard security measures including:
- JWT authentication
- Password hashing with bcrypt
- Input validation and sanitization
- SQL injection prevention
- XSS protection
- Rate limiting

**Q: Can I integrate GoCBT with our existing authentication system?**
A: Currently, GoCBT uses its own authentication system. LDAP/SSO integration may be added in future versions.

**Q: How long are user sessions valid?**
A: Sessions are valid for 24 hours by default. This can be configured in the system settings.

## 🛠️ Technical Questions

### Installation and Setup

**Q: Can I run GoCBT on Windows?**
A: Yes, GoCBT runs on Windows, macOS, and Linux. Docker is the recommended deployment method.

**Q: Do I need a dedicated server?**
A: For small deployments (< 50 users), a shared server is fine. For larger deployments, a dedicated server is recommended.

**Q: Can I use MySQL instead of PostgreSQL?**
A: Currently, GoCBT supports SQLite and PostgreSQL. MySQL support may be added in future versions.

### Performance

**Q: How many students can take a test simultaneously?**
A: This depends on your server specifications. A typical setup can handle 100-200 concurrent test sessions.

**Q: The system is running slowly. What can I do?**
A: Check these factors:
- Server resources (CPU, RAM, disk space)
- Database performance
- Network connectivity
- Number of concurrent users

**Q: How much storage space do I need?**
A: Plan for approximately 1MB per user per year, plus additional space for test data and backups.

### Troubleshooting

**Q: Students can't access tests. What should I check?**
A: Verify:
- Test is published and within the scheduled time window
- Students are assigned to the test
- Server is running and accessible
- No firewall blocking access

**Q: The database connection failed. How do I fix it?**
A: Check:
- Database server is running
- Connection credentials are correct
- Network connectivity between application and database
- Database permissions

**Q: How do I enable debug logging?**
A: Set the environment variable `LOG_LEVEL=debug` and restart the application.

## 📊 Data and Analytics

**Q: What reports are available?**
A: GoCBT provides:
- Individual student results
- Class performance summaries
- Question analysis reports
- Time usage statistics
- Grade distribution charts

**Q: Can I export data for external analysis?**
A: Yes, most data can be exported in CSV format for use in spreadsheet applications or statistical software.

**Q: How long is data retained?**
A: By default, data is retained indefinitely. You can configure automatic cleanup policies in the admin settings.

## 🔄 Updates and Maintenance

**Q: How do I update GoCBT?**
A: For Docker deployments, pull the latest image and restart containers. For manual installations, download the new version and follow the upgrade guide.

**Q: Will updates affect existing data?**
A: No, database migrations ensure existing data is preserved during updates.

**Q: How often should I update?**
A: We recommend updating monthly for security patches and quarterly for feature updates.

## 💡 Best Practices

**Q: What are the recommended test settings?**
A: 
- Time limit: 1-2 minutes per question
- Passing grade: 60-70%
- Question randomization: Enabled for large question pools
- Immediate feedback: Disabled during tests

**Q: How should I prepare students for online testing?**
A: 
- Provide practice tests
- Ensure students know the system
- Test technical setup beforehand
- Have backup plans for technical issues

**Q: What's the best way to prevent cheating?**
A: 
- Randomize question order
- Use large question pools
- Set appropriate time limits
- Monitor students during tests
- Use lockdown browser if available

## 🆘 Getting Help

**Q: Where can I find more documentation?**
A: Check the `/docs` folder in the project repository for comprehensive documentation.

**Q: How do I report a bug?**
A: Create an issue on the GitHub repository with:
- Steps to reproduce the problem
- Expected vs actual behavior
- System information
- Error messages or screenshots

**Q: Is there a user community?**
A: Join our GitHub Discussions for community support and feature requests.

**Q: Do you offer professional support?**
A: Community support is available through GitHub. Professional support options may be available upon request.

## 🔮 Future Features

**Q: What features are planned for future releases?**
A: Planned features include:
- Advanced question types (drag-and-drop, matching)
- Integration with learning management systems
- Mobile app for iOS and Android
- Advanced analytics and reporting
- Plagiarism detection
- Video proctoring integration

**Q: Can I request new features?**
A: Yes! Submit feature requests through GitHub Issues or Discussions. Community feedback helps prioritize development.

**Q: Is GoCBT open source?**
A: Yes, GoCBT is open source under the MIT license. Contributions are welcome!

---

**Still have questions?** Check our [GitHub Discussions](https://github.com/your-username/gocbt/discussions) or create an [issue](https://github.com/your-username/gocbt/issues) for technical problems.
