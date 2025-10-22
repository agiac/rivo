# Contributing to rivo

Thank you for your interest in contributing to rivo! This document provides guidelines for contributing to the project.

## Reporting Issues

When reporting issues, please include:

- A clear and descriptive title
- A detailed description of the problem
- Steps to reproduce the issue
- Expected behavior vs. actual behavior
- Go version and operating system information
- Any relevant code snippets or error messages

## Opening Pull Requests

We welcome pull requests! Please follow these guidelines:

1. **Fork the repository** and create your branch from `main`
2. **Branch naming**: Use descriptive names like:
   - `feature/add-new-operator`
   - `fix/handle-edge-case`
   - `docs/improve-readme`
3. **Make your changes** following the coding style of the project
4. **Add tests** for any new functionality
5. **Run tests** to ensure nothing is broken (see Testing section below)
6. **Write clear commit messages** (see Commit Messages section below)
7. **Open a pull request** with a clear description of your changes

## Commit Message Style

We follow conventional commit message format:

- `feat: add new feature` - for new features
- `fix: resolve bug in X` - for bug fixes
- `docs: update documentation` - for documentation changes
- `test: add tests for Y` - for test additions
- `refactor: improve code structure` - for code refactoring
- `chore: update dependencies` - for maintenance tasks

Keep the first line under 72 characters and provide additional details in the body if needed.

## Code Style

- Follow standard Go conventions and idioms
- Run `go fmt` on your code before committing
- Ensure all exported functions, types, and packages have proper documentation comments
- Keep functions focused and maintainable
- Use meaningful variable and function names

## Testing

Before submitting your pull request:

1. Run all tests:
   ```bash
   go test ./...
   ```

2. Run tests with race detection:
   ```bash
   go test -race ./...
   ```

3. Or use the Makefile:
   ```bash
   make test
   ```

All tests should pass before your PR can be merged.

## Documentation

- Add or update documentation for any new features or API changes
- Include examples in code comments where appropriate
- Update the README.md if adding significant new functionality
- Consider adding examples in the `examples/` directory for complex features

## Questions?

If you have questions about contributing, please open an issue with the `question` label, and we'll be happy to help!

## Code of Conduct

Please note that this project follows a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.
