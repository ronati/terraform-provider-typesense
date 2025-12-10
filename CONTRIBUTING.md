# Contributing to Terraform Typesense Provider

Thank you for your interest in contributing to the Terraform Typesense Provider! This guide will help you get started.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/terraform-provider-typesense.git
   cd terraform-provider-typesense
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/ronati/terraform-provider-typesense.git
   ```

4. **Install dependencies**:
   ```bash
   # Install Go dependencies
   go mod download

   # Install Node.js dependencies for commit validation
   npm install
   ```

5. **Setup git hooks** (optional but recommended):
   ```bash
   ./scripts/setup-git-hooks.sh
   ```
   This will install a commit-msg hook that validates your commit messages locally.

## Development Setup

### Prerequisites

- **Go 1.22+**: [Download Go](https://golang.org/dl/)
- **Node.js 18+**: Required for commit message validation (optional for development)
- **Docker**: Required for running acceptance tests locally
- **Terraform CLI**: [Install Terraform](https://developer.hashicorp.com/terraform/downloads)

### Building the Provider

```bash
make build
```

The compiled binary will be created in the project root.

### Running Typesense Locally

For running acceptance tests, you need a Typesense server:

```bash
docker run -d --name typesense \
  -p 8108:8108 \
  -e TYPESENSE_DATA_DIR=/tmp \
  -e TYPESENSE_API_KEY=test-api-key \
  typesense/typesense:29.0
```

## Commit Message Guidelines

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for semantic versioning and automated releases. **All commit messages must follow this format.**

### Format

```
type(scope): subject

body (optional)

footer (optional)
```

### Types

- **feat**: New feature (triggers MINOR version bump)
- **fix**: Bug fix (triggers PATCH version bump)
- **docs**: Documentation changes
- **style**: Code style changes (formatting, missing semicolons, etc.)
- **refactor**: Code refactoring without changing functionality
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Maintenance tasks, dependency updates
- **ci**: CI/CD configuration changes
- **build**: Build system changes

### Breaking Changes

For breaking changes that require a MAJOR version bump:

**Option 1**: Add `!` after the type:
```
feat!: change API key schema

BREAKING CHANGE: expires_at field is now required
```

**Option 2**: Include `BREAKING CHANGE:` in the footer:
```
feat: update collection schema

BREAKING CHANGE: Removed support for auto field type
```

### Examples

#### Good Commit Messages ✅

```
feat: add support for nested object fields

Add support for object and object[] field types with nested field configuration.
This enables more complex document structures in Typesense collections.

Closes #123
```

```
fix: resolve document update status code issue

The document update was incorrectly treating 201 as an error.
Now both 200 and 201 are accepted as successful responses.

Fixes #456
```

```
docs: add comprehensive testing guide

Added examples for running unit and acceptance tests locally,
including Docker setup for Typesense.
```

```
test: add tests for synonym resource

Added comprehensive tests covering:
- Multi-way synonyms
- One-way synonyms with root
- Multiple synonyms per collection
```

#### Bad Commit Messages ❌

```
# Too vague
fix: bug fix

# Missing type
add support for new fields

# Wrong type for breaking change
feat: change API (should be feat!)

# Not descriptive
updated code
```

### Validation

Your commit messages are automatically validated:

- **Locally**: If you ran `./scripts/setup-git-hooks.sh`, validation happens on every commit
- **CI**: All PR commits are validated in the GitHub Actions workflow

To manually validate your commits:
```bash
# Check last 3 commits
npx commitlint --from HEAD~3 --to HEAD --verbose

# Check specific commit
npx commitlint --edit <commit-hash>
```

## Testing

### Unit Tests

Run all unit tests:
```bash
go test -v -short ./...
```

### Acceptance Tests

Acceptance tests require a running Typesense instance.

#### Easy Way: Use the Test Runner Script

```bash
# Automatically starts Typesense, runs tests, and cleans up
./scripts/run-tests.sh
```

This script will:
- Check if Typesense is already running
- Start Typesense in Docker if needed
- Wait for it to be ready
- Run all acceptance tests
- Clean up the container when done

#### Manual Way: Start Typesense Yourself

```bash
# Start Typesense using Docker
docker run -d --name typesense-test \
  -p 8108:8108 \
  -e TYPESENSE_DATA_DIR=/tmp \
  -e TYPESENSE_API_KEY=test-api-key \
  typesense/typesense:29.0

# Wait a few seconds for Typesense to start
sleep 5

# Run acceptance tests (environment variables will default to localhost:8108)
make testacc

# Optional: Set custom environment variables
export TYPESENSE_API_KEY=test-api-key
export TYPESENSE_API_ADDRESS=http://localhost:8108
make testacc

# Clean up after testing
docker stop typesense-test && docker rm typesense-test
```

**Note:** If you don't set the environment variables, the tests will automatically use:
- `TYPESENSE_API_KEY=test-api-key`
- `TYPESENSE_API_ADDRESS=http://localhost:8108`

### Running Specific Tests

```bash
# Run specific test file
TF_ACC=1 go test -v ./internal/provider/ -run TestAccCollectionResource

# Run specific test case
TF_ACC=1 go test -v ./internal/provider/ -run TestAccCollectionResource_WithNestedFields
```

### Writing Tests

When adding new features or fixing bugs, please include tests:

1. **Unit tests**: For utility functions and non-resource logic
2. **Acceptance tests**: For resource CRUD operations

Test files should be named `*_test.go` and follow the existing patterns in the codebase.

Example acceptance test structure:
```go
func TestAccMyResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Create and Read
            {
                Config: testAccMyResourceConfig("test"),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("typesense_my_resource.test", "name", "test"),
                ),
            },
            // ImportState
            {
                ResourceName:      "typesense_my_resource.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
            // Update and Read
            {
                Config: testAccMyResourceConfig("test_updated"),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("typesense_my_resource.test", "name", "test_updated"),
                ),
            },
        },
    })
}
```

## Pull Request Process

1. **Create a feature branch**:
   ```bash
   git checkout -b feat/my-new-feature
   # or
   git checkout -b fix/issue-123
   ```

2. **Make your changes** following the code style and conventions

3. **Write tests** for your changes

4. **Commit your changes** using conventional commit messages:
   ```bash
   git add .
   git commit -m "feat: add support for X"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feat/my-new-feature
   ```

6. **Open a Pull Request** on GitHub

7. **Address review feedback** - maintainers will review your PR and may request changes

### PR Requirements

Your PR must pass all checks:

- ✅ All commits follow Conventional Commits format
- ✅ All tests pass (unit and acceptance)
- ✅ Code builds successfully
- ✅ No merge conflicts with base branch

### PR Description

Please include in your PR description:

- What does this PR do?
- Why is this change needed?
- How was this tested?
- Any breaking changes?
- Related issues (use `Fixes #123` or `Closes #456`)

## Release Process

Releases are automated using semantic versioning:

1. **Commits are merged to `master`** via approved PRs
2. **CI validates all commits** follow conventional commits
3. **Semantic Release analyzes commits** and determines version bump:
   - `feat:` → minor version bump (0.X.0)
   - `fix:` → patch version bump (0.0.X)
   - `feat!:` or `BREAKING CHANGE:` → major version bump (X.0.0)
4. **Release is created automatically** with:
   - Updated version number
   - Generated CHANGELOG
   - GitHub release notes
   - Git tag
5. **Provider is published** to Terraform Registry

You don't need to manually update version numbers or create releases - it's all automated!

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the code, not the person
- Help others learn and grow

## Questions or Issues?

- **Bug reports**: Open an issue with details about the bug and steps to reproduce
- **Feature requests**: Open an issue describing the feature and its use case
- **Questions**: Open a discussion or issue for clarification

## License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project (check the LICENSE file).

---

Thank you for contributing! 🎉
