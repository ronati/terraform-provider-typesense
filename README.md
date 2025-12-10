<div align="center">
  <h1>Terraform Provider for Typesense</h1>
  <strong>This is a Terraform provider for Typesense</strong>

  <br><br>

  [![Tests](https://github.com/ronati/terraform-provider-typesense/actions/workflows/build-and-test.yml/badge.svg)](https://github.com/ronati/terraform-provider-typesense/actions/workflows/build-and-test.yml)
  [![codecov](https://codecov.io/gh/ronati/terraform-provider-typesense/branch/master/graph/badge.svg)](https://codecov.io/gh/ronati/terraform-provider-typesense)
  [![Go Report Card](https://goreportcard.com/badge/github.com/ronati/terraform-provider-typesense)](https://goreportcard.com/report/github.com/ronati/terraform-provider-typesense)
</div>

<hr>

## Support

- Supports v28.0+ version of Typesense.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= v0.12.0 (v0.11.x may work but not supported actively)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/ronati/terraform-provider-typesense`

```console
$ mkdir -p $GOPATH/src/github.com/ronati; cd $GOPATH/src/github.com/ronati
$ git clone git@github.com:ronati/terraform-provider-typesense
Enter the provider directory and build the provider

$ cd $GOPATH/src/github.com/ronati/terraform-provider-typesense
$ make build
```

## Testing

### Running Tests Locally

#### Unit Tests
```bash
go test -v -short ./...
```

#### Acceptance Tests

Acceptance tests require a running Typesense instance:

```bash
# Start Typesense
docker run -d --name typesense-test \
  -p 8108:8108 \
  -e TYPESENSE_DATA_DIR=/tmp \
  -e TYPESENSE_API_KEY=test-api-key \
  typesense/typesense:29.0

# Wait for it to start
sleep 5

# Run tests (will use localhost:8108 by default)
make testacc

# Clean up
docker stop typesense-test && docker rm typesense-test
```

**Note:** Tests will automatically connect to `http://localhost:8108` with API key `test-api-key` if environment variables are not set.

## Contributing

**All contributions are welcome!**

This project uses [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for automated semantic versioning and releases. Please read our [Contributing Guide](CONTRIBUTING.md) for details on:

- Setting up your development environment
- Commit message format and validation
- Testing requirements
- Pull request process

### Quick Start for Contributors

```bash
# Clone and setup
git clone https://github.com/ronati/terraform-provider-typesense.git
cd terraform-provider-typesense
npm install

# Setup git hooks for commit validation (optional)
./scripts/setup-git-hooks.sh

# Make changes and commit following conventional commits format
git commit -m "feat: add new feature"
```

### Commit Message Format

All commits must follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
type(scope): subject

Examples:
  feat: add support for nested fields
  fix: resolve document update issue
  docs: update README
  test: add tests for synonym resource
```

**Note**: Commit messages are automatically validated in CI. PRs with invalid commit messages will fail the build.

## CI/CD

This project uses GitHub Actions for continuous integration and deployment:

- **Pull Requests**: Validates commit messages and runs all tests
- **Master/Beta Branch**: Automatically creates releases using semantic versioning
- **Version Tags**: Publishes provider to Terraform Registry

See [GitHub Workflows Documentation](.github/workflows/README.md) for more details.

## Notes for Maintainers

When you merge a PR from `beta` into `master` and it successfully publishes a new version on the `latest` channel, **don't forget to create a PR from `master` to `beta`**. This is mandatory for `semantic-release` to take it into account for next `beta` version.
