# Contributing to GoCBT

Thank you for your interest in contributing to GoCBT! This document provides guidelines and information for contributors.

## 🤝 How to Contribute

### Types of Contributions

We welcome various types of contributions:

- **Bug Reports**: Help us identify and fix issues
- **Feature Requests**: Suggest new functionality
- **Code Contributions**: Submit bug fixes or new features
- **Documentation**: Improve or add documentation
- **Testing**: Help test new features and report issues
- **Design**: Improve UI/UX design
- **Translations**: Help make GoCBT available in more languages

## 🐛 Reporting Bugs

### Before Submitting a Bug Report

1. **Check existing issues** to avoid duplicates
2. **Update to the latest version** to see if the issue persists
3. **Test in a clean environment** to rule out local configuration issues

### How to Submit a Bug Report

Create an issue with the following information:

**Title**: Brief, descriptive summary of the bug

**Description**:
- **Expected behavior**: What should happen
- **Actual behavior**: What actually happens
- **Steps to reproduce**: Detailed steps to recreate the issue
- **Environment**: OS, browser, GoCBT version
- **Screenshots**: If applicable
- **Error messages**: Full error text or logs

**Template**:
```markdown
## Bug Description
Brief description of the bug.

## Expected Behavior
What you expected to happen.

## Actual Behavior
What actually happened.

## Steps to Reproduce
1. Go to '...'
2. Click on '...'
3. Scroll down to '...'
4. See error

## Environment
- OS: [e.g., Windows 10, macOS 12, Ubuntu 20.04]
- Browser: [e.g., Chrome 95, Firefox 94]
- GoCBT Version: [e.g., 1.2.0]

## Additional Context
Any other context about the problem.
```

## 💡 Suggesting Features

### Before Submitting a Feature Request

1. **Check existing feature requests** to avoid duplicates
2. **Consider the scope** - is this a core feature or edge case?
3. **Think about implementation** - how would this work?

### How to Submit a Feature Request

Create an issue with:

**Title**: Clear, concise feature description

**Description**:
- **Problem statement**: What problem does this solve?
- **Proposed solution**: How should it work?
- **Alternatives considered**: Other approaches you've thought about
- **Use cases**: Who would benefit and how?
- **Implementation notes**: Technical considerations (optional)

## 🔧 Code Contributions

### Development Setup

1. **Fork the repository**
```bash
git clone https://github.com/your-username/gocbt.git
cd gocbt
```

2. **Set up development environment**
```bash
# Backend
go mod tidy
cp .env.example .env

# Frontend
cd frontend
npm install
```

3. **Create a feature branch**
```bash
git checkout -b feature/your-feature-name
```

### Coding Standards

#### Go Code Standards

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for formatting
- Run `golangci-lint` for linting
- Write tests for new functionality
- Use meaningful variable and function names
- Add comments for exported functions

**Example**:
```go
// CreateUser creates a new user with the provided information.
// It validates the input, checks for duplicates, and returns the created user.
func (s *UserService) CreateUser(req CreateUserRequest) (*User, error) {
    if err := s.validateUser(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // Check for existing user
    if exists, err := s.repo.UserExists(req.Username, req.Email); err != nil {
        return nil, fmt.Errorf("failed to check user existence: %w", err)
    } else if exists {
        return nil, ErrUserAlreadyExists
    }
    
    // Create user
    user := &User{
        Username:  req.Username,
        Email:     req.Email,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Role:      req.Role,
    }
    
    return s.repo.Create(user)
}
```

#### TypeScript/React Standards

- Use TypeScript for type safety
- Follow React best practices
- Use functional components with hooks
- Write meaningful component and prop names
- Add JSDoc comments for complex functions

**Example**:
```typescript
interface UserCardProps {
  user: User;
  onEdit?: (user: User) => void;
  onDelete?: (userId: number) => void;
}

/**
 * UserCard component displays user information with optional edit/delete actions.
 */
export const UserCard: React.FC<UserCardProps> = ({ 
  user, 
  onEdit, 
  onDelete 
}) => {
  const handleEdit = () => {
    onEdit?.(user);
  };

  const handleDelete = () => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      onDelete?.(user.id);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h3 className="text-lg font-semibold">
        {user.firstName} {user.lastName}
      </h3>
      <p className="text-gray-600">{user.email}</p>
      <div className="mt-4 space-x-2">
        {onEdit && (
          <button onClick={handleEdit} className="btn btn-primary">
            Edit
          </button>
        )}
        {onDelete && (
          <button onClick={handleDelete} className="btn btn-danger">
            Delete
          </button>
        )}
      </div>
    </div>
  );
};
```

### Testing Requirements

#### Backend Tests

- Write unit tests for all new functions
- Include integration tests for API endpoints
- Aim for >80% code coverage
- Use table-driven tests where appropriate

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        req     CreateUserRequest
        wantErr bool
        errType error
    }{
        {
            name: "valid user",
            req: CreateUserRequest{
                Username:  "testuser",
                Email:     "test@example.com",
                FirstName: "Test",
                LastName:  "User",
                Role:      "student",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            req: CreateUserRequest{
                Username:  "testuser",
                Email:     "invalid-email",
                FirstName: "Test",
                LastName:  "User",
                Role:      "student",
            },
            wantErr: true,
            errType: ErrInvalidEmail,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := NewUserService(mockRepo)
            user, err := service.CreateUser(tt.req)
            
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errType != nil {
                    assert.ErrorIs(t, err, tt.errType)
                }
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, user)
                assert.Equal(t, tt.req.Username, user.Username)
            }
        })
    }
}
```

#### Frontend Tests

- Write tests for all new components
- Test user interactions and edge cases
- Use React Testing Library

```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { UserCard } from '../UserCard';

const mockUser = {
  id: 1,
  username: 'testuser',
  email: 'test@example.com',
  firstName: 'Test',
  lastName: 'User',
  role: 'student',
};

describe('UserCard', () => {
  it('renders user information', () => {
    render(<UserCard user={mockUser} />);
    
    expect(screen.getByText('Test User')).toBeInTheDocument();
    expect(screen.getByText('test@example.com')).toBeInTheDocument();
  });

  it('calls onEdit when edit button is clicked', () => {
    const mockOnEdit = jest.fn();
    render(<UserCard user={mockUser} onEdit={mockOnEdit} />);
    
    fireEvent.click(screen.getByText('Edit'));
    expect(mockOnEdit).toHaveBeenCalledWith(mockUser);
  });
});
```

### Commit Message Format

Use conventional commits format:

```
type(scope): description

[optional body]

[optional footer]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples**:
```
feat(auth): add password reset functionality

fix(api): handle null values in user response

docs(readme): update installation instructions

test(user): add unit tests for user service
```

### Pull Request Process

1. **Ensure your branch is up to date**
```bash
git checkout main
git pull origin main
git checkout your-feature-branch
git rebase main
```

2. **Run tests and linting**
```bash
# Backend
go test ./...
golangci-lint run

# Frontend
cd frontend
npm test
npm run lint
```

3. **Create a pull request** with:
   - Clear title and description
   - Reference to related issues
   - Screenshots for UI changes
   - Testing instructions

4. **Pull request template**:
```markdown
## Description
Brief description of changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Related Issues
Fixes #123

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Screenshots (if applicable)
[Add screenshots here]

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Tests added/updated
- [ ] Documentation updated
```

## 📚 Documentation Contributions

### Types of Documentation

- **API Documentation**: Update when adding/changing endpoints
- **User Guides**: Help users understand features
- **Developer Documentation**: Technical implementation details
- **README Updates**: Keep installation and setup current

### Documentation Standards

- Use clear, concise language
- Include code examples where helpful
- Keep documentation up to date with code changes
- Use proper markdown formatting

## 🎨 Design Contributions

### UI/UX Guidelines

- Follow existing design patterns
- Ensure accessibility (WCAG 2.1 AA)
- Test on different screen sizes
- Use consistent color scheme and typography
- Provide mockups or prototypes for major changes

### Design Tools

- Figma files available upon request
- Use Tailwind CSS classes for styling
- Follow component-based design principles

## 🌍 Translation Contributions

We welcome translations to make GoCBT accessible to more users:

1. Check existing translations in `/locales`
2. Create new language files following the existing structure
3. Translate all user-facing strings
4. Test the translation in the application
5. Submit a pull request

## 📋 Review Process

### What We Look For

- **Functionality**: Does it work as intended?
- **Code Quality**: Is it well-written and maintainable?
- **Testing**: Are there adequate tests?
- **Documentation**: Is it properly documented?
- **Security**: Are there any security concerns?
- **Performance**: Does it impact performance?

### Review Timeline

- Initial review: Within 3-5 business days
- Follow-up reviews: Within 1-2 business days
- Merge: After approval from at least one maintainer

## 🏆 Recognition

Contributors are recognized in:
- GitHub contributors list
- Release notes for significant contributions
- Special mentions in documentation

## 📞 Getting Help

### Communication Channels

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and community support
- **Email**: For sensitive security issues

### Mentorship

New contributors can request mentorship for:
- Understanding the codebase
- Learning development practices
- Getting started with first contributions

## 📜 Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive environment for all contributors, regardless of:
- Experience level
- Gender identity and expression
- Sexual orientation
- Disability
- Personal appearance
- Body size
- Race
- Ethnicity
- Age
- Religion
- Nationality

### Expected Behavior

- Use welcoming and inclusive language
- Be respectful of differing viewpoints
- Accept constructive criticism gracefully
- Focus on what's best for the community
- Show empathy towards other community members

### Unacceptable Behavior

- Harassment or discriminatory language
- Personal attacks or trolling
- Public or private harassment
- Publishing others' private information
- Other conduct inappropriate in a professional setting

### Enforcement

Violations can be reported to the project maintainers. All complaints will be reviewed and investigated promptly and fairly.

---

**Thank you for contributing to GoCBT!** Your efforts help make online testing better for everyone. 🎉
