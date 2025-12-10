# Test Fixes Applied

## Problem
Acceptance tests were failing with error: `unsupported protocol scheme ""` because the environment variables were set to invalid values.

## Root Cause
The `testAccPreCheck` function in `internal/provider/provider_test.go` was setting:
```go
os.Setenv("TYPESENSE_API_KEY", "1")
os.Setenv("TYPESENSE_API_ADDRESS", "1")
```

The value "1" is not a valid URL, causing the Typesense client to fail.

## Solution
Updated `testAccPreCheck` to:
1. Check if environment variables are already set (for CI)
2. If not set, use sensible defaults for local testing:
   - `TYPESENSE_API_KEY=test-api-key`
   - `TYPESENSE_API_ADDRESS=http://localhost:8108`
3. Validate that variables are properly set

## Changes Made

### 1. Fixed `provider_test.go`
- Updated `testAccPreCheck` to use proper default values
- Added validation to ensure variables are set

### 2. Created Test Helper Script
- `scripts/run-tests.sh` - Automatically manages Typesense container
- Starts Typesense if not running
- Runs tests
- Cleans up afterward

### 3. Updated Documentation
- `README.md` - Added clear testing instructions
- `CONTRIBUTING.md` - Detailed testing guide with script usage

## How to Run Tests Now

### Option 1: Use the helper script (easiest)
```bash
./scripts/run-tests.sh
```

### Option 2: Manual
```bash
# Start Typesense
docker run -d --name typesense-test \
  -p 8108:8108 \
  -e TYPESENSE_DATA_DIR=/tmp \
  -e TYPESENSE_API_KEY=test-api-key \
  typesense/typesense:29.0

# Wait for startup
sleep 5

# Run tests
make testacc

# Cleanup
docker stop typesense-test && docker rm typesense-test
```

### Option 3: CI (GitHub Actions)
Tests run automatically with Typesense service container on every PR.

## Testing
The tests should now pass when:
- Running in CI (uses service container with env vars)
- Running locally (uses defaults or custom env vars)
- Running with the helper script (fully automated)
