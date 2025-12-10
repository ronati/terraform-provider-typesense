# GitHub Actions Workflows

This directory contains the CI/CD workflows for the Terraform Typesense Provider.

## Workflows

### 1. Build and Test (`build-and-test.yml`)

**Trigger:** Pull requests (opened, synchronized, reopened)

**Purpose:** Validates commits and runs tests for every pull request.

**Jobs:**

#### Commit Lint
- Validates that all commits in the PR follow the [Conventional Commits](https://www.conventionalcommits.org/) specification
- Required for semantic versioning and automated releases
- Format: `type(scope): subject`
- Examples:
  - `feat: add support for new field types`
  - `fix: resolve collection update issue`
  - `docs: update README with examples`
  - `chore: update dependencies`

**Commit Types:**
- `feat`: New features (triggers minor version bump)
- `fix`: Bug fixes (triggers patch version bump)
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `ci`: CI/CD changes
- `build`: Build system changes

**Breaking Changes:**
- Add `BREAKING CHANGE:` in the commit body or use `!` after type: `feat!: change API`
- Triggers major version bump

#### Lint
- Runs golangci-lint to check code quality
- Enforces consistent code style
- Catches common bugs and issues
- Configuration: `.golangci.yml`
- Enabled linters:
  - errcheck, gosimple, govet, ineffassign, staticcheck, unused
  - gofmt, goimports, misspell, revive
  - gosec (security), gocritic (quality)

#### Build and Test
- Runs against Typesense v29.0 server in a Docker container
- Executes both unit tests and acceptance tests
- Generates code coverage report (uploaded to Codecov)
- **Posts coverage report as PR comment** showing coverage diff
- **Coverage status checks**: Fails PR if coverage drops or new code poorly tested
- Validates that the provider builds successfully
- Required environment variables:
  - `TYPESENSE_API_KEY`: API key for Typesense (set to `test-api-key` in CI)
  - `TYPESENSE_API_ADDRESS`: Typesense server address (set to `http://localhost:8108` in CI)
  - `TF_ACC`: Must be set to `1` to run acceptance tests

### 2. SemVer Release (`semver-release.yml`)

**Trigger:** Push to `master` or `beta` branches

**Purpose:** Automatically creates releases using semantic versioning based on commit messages.

**Jobs:**
- Validates all commits since the last tag
- Runs tests to ensure the release is stable
- Generates documentation
- Creates GitHub release with changelog
- Updates version numbers automatically

**Requirements:**
- All commits must follow Conventional Commits format
- Tests must pass
- Requires GitHub App token for authentication

### 3. Release (`release-go.yml`)

**Trigger:** Push of version tags (e.g., `v1.2.3`)

**Purpose:** Publishes the provider to the Terraform Registry using GoReleaser.

**Jobs:**
- Builds binaries for multiple platforms
- Signs binaries with GPG
- Publishes to GitHub releases
- Makes provider available in Terraform Registry

**Requirements:**
- Valid GPG key for signing
- Tag must match pattern `v[0-9]+.[0-9]+.[0-9]+`

## Running Tests Locally

### Unit Tests
```bash
go test -v -short ./...
```

### Acceptance Tests
```bash
# Start Typesense locally
docker run -d -p 8108:8108 \
  -e TYPESENSE_DATA_DIR=/tmp \
  -e TYPESENSE_API_KEY=test-api-key \
  typesense/typesense:29.0

# Run tests
export TYPESENSE_API_KEY=test-api-key
export TYPESENSE_API_ADDRESS=http://localhost:8108
make testacc
```

### Validate Commit Messages
```bash
# Install dependencies
npm install

# Check recent commits
npx commitlint --from HEAD~3 --to HEAD --verbose

# Check specific commit
npx commitlint --edit <commit-hash>
```

## Commit Message Examples

### Feature Addition
```
feat: add support for nested object fields

Add support for object and object[] field types with nested field configuration.
This enables more complex document structures in Typesense collections.

Closes #123
```

### Bug Fix
```
fix: resolve document update returning incorrect status

The document update was returning 201 instead of 200. Added handling
to treat 201 as successful update response.

Fixes #456
```

### Breaking Change
```
feat!: change API key resource schema

BREAKING CHANGE: The expires_at field is now required for all API keys.
Previous behavior allowed omitting this field, which resulted in keys
with very long default expiration times.

Migration: Add expires_at = 64723363199 to existing API key resources
to maintain previous behavior.
```

### Documentation
```
docs: add comprehensive testing guide

Added examples for running unit and acceptance tests locally,
including Docker setup for Typesense.
```

## Troubleshooting

### Commit validation fails
- Ensure your commit messages follow the Conventional Commits format
- Use `npx commitlint --edit HEAD` to validate your last commit
- Amend your commit message if needed: `git commit --amend`

### Tests fail in CI but pass locally
- Check Typesense version compatibility (CI uses v29.0)
- Verify environment variables are set correctly
- Ensure no test data conflicts between test cases

### Release workflow doesn't trigger
- Verify commit messages include proper types (feat/fix)
- Check that you're pushing to `master` or `beta` branch
- Ensure GitHub App credentials are configured correctly

## Additional Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Terraform Plugin Testing](https://developer.hashicorp.com/terraform/plugin/testing)
- [Typesense Documentation](https://typesense.org/docs/)
