# Codecov PR Comments Setup

Codecov will now automatically post coverage reports as comments on every pull request!

## What You'll See on PRs

When you open a PR, Codecov will post a comment showing:

```
## Codecov Report

Attention: Patch coverage is 85.71% with 2 lines in your changes missing coverage.

| Comparison | Base | Head | Diff |
|------------|------|------|------|
| Coverage   | 78.23% | 79.15% | +0.92% ⬆️ |
| Files      | 12   | 12   | - |
| Lines      | 1234 | 1256 | +22 |

### Changes Missing Coverage

| File | Lines |
|------|-------|
| internal/provider/resource_new.go | L45-46 |

📊 View full report in Codecov by Sentry
```

## Features Enabled

### 1. PR Comments ✅
- **Automatic comments** on every PR
- **Coverage diff** showing what changed
- **File-by-file breakdown** of coverage changes
- **Missing coverage highlights** showing uncovered lines
- **Updates existing comment** instead of spamming

### 2. Coverage Status Checks ✅
Two checks will appear on your PR:

**Project Coverage:**
- Compares overall project coverage to base branch
- ✅ Passes if coverage doesn't drop more than 0.5%
- ❌ Fails if coverage drops significantly

**Patch Coverage:**
- Checks coverage of NEW code only
- ✅ Passes if new code is 80%+ covered
- ❌ Fails if new code is poorly tested

### 3. Configuration (`codecov.yml`)

```yaml
coverage:
  status:
    project:
      target: auto          # Compare to base branch
      threshold: 0.5%       # Allow small drops
    patch:
      target: 80%           # New code must be 80% covered
```

## Setup Required

### For Public Repositories (Recommended)
No token needed! Codecov works automatically with:
- ✅ Public GitHub repos
- ✅ No configuration required
- ✅ Free forever

Just push your code and open a PR. Codecov will work automatically.

### For Private Repositories
You need to add a Codecov token:

1. **Get Token:**
   - Go to https://codecov.io/
   - Sign in with GitHub
   - Select your repository
   - Copy the token from Settings → General

2. **Add to GitHub Secrets:**
   - Go to your repo → Settings → Secrets → Actions
   - Click "New repository secret"
   - Name: `CODECOV_TOKEN`
   - Value: (paste your token)
   - Click "Add secret"

3. **Already Configured!** The workflow already has:
   ```yaml
   - name: Upload coverage to Codecov
     uses: codecov/codecov-action@v4
     env:
       CODECOV_TOKEN: ${{ secrets.CODECOV_TOKEN }}
   ```

## What the Comment Includes

### Header
- Overall coverage change (e.g., "+0.92% ⬆️")
- Comparison table (Base vs Head)

### Diff
- Files that changed coverage
- Line-by-line coverage changes

### Changes
- New files added
- Removed files

### Suggestions
- Tips to improve coverage
- Links to uncovered code

### Files
- Full list of all files with coverage %
- Sortable by coverage percentage

### Footer
- Link to full report on Codecov
- Build information

## PR Workflow Example

```bash
# 1. Create a branch
git checkout -b feature/new-resource

# 2. Make changes
# ... edit code ...

# 3. Commit
git add .
git commit -m "feat: add new resource"

# 4. Push
git push origin feature/new-resource

# 5. Open PR on GitHub
# → CI runs tests with coverage
# → Codecov posts comment automatically
# → You see coverage diff in PR
# → Review coverage before merging
```

## Configuration Options

### Strict Mode (Current)
```yaml
coverage:
  status:
    project:
      informational: false  # Fail PR if coverage drops
    patch:
      informational: false  # Fail PR if new code poorly tested
```

**Result:** PR checks will fail if coverage is poor

### Informational Mode
```yaml
coverage:
  status:
    project:
      informational: true   # Just report, don't fail
    patch:
      informational: true   # Just report, don't fail
```

**Result:** PR checks always pass, coverage is FYI only

### Adjust Thresholds

**More Strict:**
```yaml
coverage:
  status:
    project:
      threshold: 0%         # No coverage drops allowed
    patch:
      target: 90%           # 90% coverage required for new code
```

**More Lenient:**
```yaml
coverage:
  status:
    project:
      threshold: 2%         # Allow 2% drop
    patch:
      target: 70%           # 70% coverage for new code
```

## Troubleshooting

### Comment Not Appearing?

**Check Permissions:**
- Workflow has `pull-requests: write` permission ✅ (already set)

**Check Token (private repos only):**
```bash
# Verify secret exists
gh secret list | grep CODECOV_TOKEN
```

**Check Codecov Status:**
- Go to https://app.codecov.io/gh/ronati/terraform-provider-typesense
- Check if builds are uploading

### Coverage Seems Wrong?

**Check Ignored Files:**
Edit `codecov.yml`:
```yaml
ignore:
  - "**/*_test.go"  # Test files
  - ".old/**"       # Old code
  - "examples/**"   # Examples
```

**View Raw Coverage:**
```bash
# Generate locally
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Multiple Comments Posted?

Should not happen with current config:
```yaml
comment:
  behavior: default  # Updates existing comment
```

If it happens, Codecov may be running multiple times. Check workflow.

## Benefits

### For Contributors
- 👀 **See coverage impact** before merging
- 🎯 **Know what to test** - highlighted uncovered lines
- 📊 **Track progress** - see coverage trends

### For Maintainers
- ✅ **Enforce standards** - fail PRs with poor coverage
- 📈 **Improve quality** - coverage visible to all
- 🔍 **Catch gaps** - see exactly what's untested

### For Project
- 🛡️ **Prevent regressions** - coverage can't drop
- 📚 **Documentation** - coverage as quality metric
- 🏆 **Professional** - shows commitment to quality

## Example PR Comment

Here's what a real comment looks like:

> ## [Codecov](https://app.codecov.io/gh/ronati/terraform-provider-typesense/pull/123) Report
> 
> **Attention**: Patch coverage is `75.00%` with `3 lines` in your changes missing coverage. Please review.
> 
> | Comparison | Base | Head | +/- |
> |------------|------|------|-----|
> | **Coverage** | 78.45% | 79.23% | +0.78% ⬆️ |
> | Files | 12 | 13 | +1 |
> | Lines | 1234 | 1268 | +34 |
> 
> ### Files with missing lines
> 
> | File | Coverage Δ | Complexity Δ | Missing Lines |
> |------|------------|--------------|---------------|
> | internal/provider/resource_new.go | 72.50% | 0.00 | L45-47 |
> 
> 📢 Thoughts on this report? [Let us know!](https://about.codecov.io/codecov-free-trial/)

Clean, professional, and actionable!

## Summary

✅ **PR comments enabled** - automatic coverage reports
✅ **Status checks enabled** - fail if coverage drops
✅ **Configured for success** - sensible defaults
✅ **Public repo ready** - no token needed
✅ **Private repo ready** - just add token

Your PRs will now have beautiful coverage reports! 🎉

---

## Quick Reference

**Current Settings:**
- Project coverage threshold: 0.5% drop allowed
- Patch coverage requirement: 80% for new code
- Comment behavior: Update existing (no spam)
- Ignored: test files, .old/, examples/

**Files:**
- `codecov.yml` - Configuration
- `.github/workflows/build-and-test.yml` - CI workflow
- `CODECOV_SETUP.md` - This documentation

**Links:**
- Codecov Dashboard: https://app.codecov.io/gh/ronati/terraform-provider-typesense
- Documentation: https://docs.codecov.com/
- Support: https://about.codecov.io/support/
