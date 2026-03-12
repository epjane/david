# PHASED FIX PLAN FOR DAVID PROJECT

## Overview
Comprehensive phased plan to fix critical issues in the david WebDAV server project.

---

## PHASE 1: CRITICAL BUG FIXES

### 1.1 security.go - Authentication & Authorization
**File:** `app/security.go`

**Issues:**
- Lines 151-156: Duplicate error handling code
- Line 139: Premature Read permission check before successful authentication
- Lines 192-196: Duplicate "DELETE" in allowedMethods array
- Line 251: Wrong return value (!ok) for OPTIONS method
- Lines 266-288: Flawed PROPFIND non-existent file logic
- Line 344: Generic 501 error for unhandled methods

**Status:** COMPLETED
**Fixes Applied:**
- Separated Authenticated and Read permission checks
- Consolidated error handling
- Removed duplicate DELETE from allowedMethods
- Fixed OPTIONS return to true
- Simplified PROPFIND logic
- Fixed default case to return 501

### 1.2 config.go - User Permission Processing
**File:** `app/config.go`

**Issues:**
- Line 65: Commented out setDefaults() - dead code
- Lines 92-107: Iterates viper keys but accesses cfg.Users before population
- Lines 190-197: CRUD update logic doesn't properly preserve existing objects

**Status:** COMPLETED
**Fixes Applied:**
- Removed commented setDefaults() call
- Initialize cfg.Users[user] before accessing
- Fixed CRUD comparison to use Permissions field

### 1.3 crud.go - Context Handling
**File:** `app/crud.go`

**Issues:**
- Lines 67-71: Checks if ctx is nil after WithValue (ctx will never be nil)
- Lines 18-20: Unused contextKey and crudContextKey variables

**Status:** COMPLETED
**Fixes Applied:**
- Removed context update logic (not needed)
- Deleted unused contextKey declarations

### 1.4 fs.go - Permission Logic Bugs
**File:** `app/fs.go`

**Issues:**
- Line 152: Checks !Read for write operations, should check !Create
- Lines 161-171: Contradictory logic with confusing comments
- Line 145: Uses Warn log level, should be Debug
- Lines 46-71: Commented-out duplicate resolve function

**Status:** COMPLETED
**Fixes Applied:**
- Removed commented-out resolve function
- Changed line 152: Read to Create permission check
- Simplified write permission logic (lines 161-171)
- Changed line 145 log level to Debug

### 1.5 main.go - Typos
**File:** `cmd/david/main.go`

**Issues:**
- Lines 33, 37: "prdouction" typo

**Status:** COMPLETED
**Fixes Applied:**
- Fixed "prdouction" to "production"

### 1.5 main.go - Typos
**File:** `cmd/david/main.go`

**Issues:**
- Lines 33, 37: "prdouction" typo

**Status:** COMPLETED
**Fixes Applied:**
- Fixed "prdouction" to "production"

### 1.6 magefile.go - Wrong Paths
**File:** `magefile.go`

**Issues:**
- Lines 32, 50, 120: References "dave"/"davecli" → should be "david"/"bcpt"
- Lines 64-78: Variables use dave/daveCli → should be david/bcpt
- Line 73: Archive uses "dave" → should be "david"
- Lines 148-172: Source paths and executable names incorrect

**Status:** COMPLETED
**Fixes Applied:**
- Updated all comments to use david/bcpt
- Updated variable names from dave/daveCli to david/bcpt
- Fixed source paths to cmd/david and cmd/bcpt
- Fixed archive name to use david

---

## PHASE 2: CODE QUALITY

### 2.1 Remove Dead Code
- security.go: Remove unused testCrudType (line 47)
- crud.go: Remove unused contextKey variables (lines 18-20)

**Status:** PENDING

### 2.2 Standardize Error Handling
- Use consistent error wrapping (fmt.Errorf with %w or errors.Wrap)
- Add error context to all errors

**Status:** PENDING

### 2.3 Improve Logging Consistency
- Standardize log levels (Debug/Info/Warn/Error)
- Use consistent field names ("user" not "username")

**Status:** PENDING

### 2.4 Code Formatting
- Run: `gofmt -s -w .`
- Run: `go vet ./...`
- Fix all linting issues

**Status:** PENDING

---

## PHASE 3: TEST FIXES

### 3.1 security_test.go
- Test case 4 (line 132): Add missing CrudType to expected result
- Test case 7 (line 188): Change Authenticated from false to true
- TestHandle tests (lines 270-315): Add proper Crud permissions to users
- Update expectations based on new authorization logic

**Status:** PENDING

### 3.2 fs_test.go
- TestDirOpenFile (lines 284-290): Fix error expectations
- TestRename (lines 483-500): Ensure Update permissions in admin context
- TestDirStat (lines 567-571): Verify preconditions

**Status:** PENDING

### 3.3 Add Missing Tests
- config_test.go: Test config reload functionality
- security_test.go: Test different HTTP methods with permissions
- integration_test.go: End-to-end server tests

**Status:** PENDING

---

## PHASE 4: DEPENDENCY UPDATES

**File:** `go.mod`

**Updates Needed:**
- Go version: 1.21 → 1.22
- github.com/sirupsen/logrus: v1.6.0 → v1.9.3
- github.com/spf13/viper: v1.15.0 → v1.18.0
- golang.org/x/crypto: v0.14.0 → v0.17.0
- golang.org/x/net: v0.17.0 → v0.19.0
- golang.org/x/term: v0.13.0 → v0.15.0
- github.com/fsnotify/fsnotify: v1.6.0 → v1.7.0
- github.com/spf13/cobra: v1.6.1 → v1.8.0
- github.com/magefile/mage: v1.10.0 → v1.15.0

**Commands:**
```bash
go get -u ./...
go mod tidy
```

**Status:** PENDING

---

## PHASE 5: DOCUMENTATION

### 5.1 Readme.md
**Issues:**
- Line 44: Step numbering starts at 3 (missing steps 1-2)
- Lines 21, 34, 36: Inconsistent "dave"/"david" naming
- Lines 85-90: Password hashes are examples, should note this
- Lines 130-140: CORS origin "*" is insecure, should warn

**Status:** PENDING

### 5.2 config-sample.yaml
**Issues:**
- Comments reference "dave"
- Password hashes are examples
- CORS origin "*" for production

**Status:** PENDING

---

## EXECUTION ORDER

1. PHASE 1 → Critical bugs (security, logic)
2. PHASE 2 → Code quality
3. PHASE 3 → Test fixes
4. PHASE 4 → Dependencies
5. PHASE 5 → Documentation

---

## PHASE 1 SUMMARY

**Status:** COMPLETED

**Files Modified:**
- app/security.go - Fixed authentication flow, removed duplicate DELETE, simplified PROPFIND, fixed OPTIONS return
- app/config.go - Removed dead code, fixed user permission processing, fixed CRUD update logic
- app/crud.go - Removed buggy context handling and unused variables
- app/fs.go - Fixed permission checks, simplified logic, removed commented code
- cmd/david/main.go - Fixed typos
- magefile.go - Fixed binary names and paths

**Test Results:** 
- TestAuthenticate: PASS (7/7 tests)
- TestHandle: 2/3 pass (1 test needs permission updates)
- Other tests: Need verification after Phase 2-3

**Note:** Some tests fail because they don't set Crud permissions. This is expected - Phase 3 will fix the tests to match the improved authorization logic.

---

## SUCCESS CRITERIA

- [x] All critical bugs fixed
- [ ] Code passes go vet
- [ ] 90%+ test coverage
- [ ] Dependencies updated
- [ ] Documentation accurate
- [ ] mage Build succeeds
- [ ] go test -v -cover passes

---

## VERIFICATION COMMANDS

```bash
cd /home/steam/Documents/david

# Build
mage Clean && mage Build

# Test
cd app && go test -v -cover ./...

# Vet
go vet -v ./...
```

---

**Started: Thu Mar 12 2026**
**Branch: fix-issues**