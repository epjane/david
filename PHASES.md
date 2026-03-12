# PHASED FIX PLAN FOR DAVID PROJECT

## Overview
Comprehensive phased plan to fix critical issues in the david WebDAV server project.

---

## PHASE 1: CRITICAL BUG FIXES

### 1.1 security.go - Authentication & Authorization
**File:** `app/security.go`

**Issues Fixed:**
- Lines 151-156: Duplicate error handling code - FIXED
- Line 139: Premature Read permission check - FIXED
- Lines 192-196: Duplicate "DELETE" in allowedMethods - FIXED
- Line 251: Wrong return value for OPTIONS - FIXED
- Lines 266-288: Flawed PROPFIND logic - FIXED
- Line 344: Generic 501 for unhandled methods - FIXED

**Status:** ✓ COMPLETED

### 1.2 config.go - User Permission Processing
**File:** `app/config.go`

**Issues Fixed:**
- Line 65: Dead commented code - FIXED
- Lines 92-107: cfg.Users accessed before population - FIXED
- Lines 190-197: CRUD update logic - FIXED

**Status:** ✓ COMPLETED

### 1.3 crud.go - Context Handling
**File:** `app/crud.go`

**Issues Fixed:**
- Lines 67-71: Buggy context nil check - FIXED
- Lines 18-20: Unused contextKey variables - FIXED

**Status:** ✓ COMPLETED

### 1.4 fs.go - Permission Logic Bugs
**File:** `app/fs.go`

**Issues Fixed:**
- Line 152: Wrong permission check (Read vs Create) - FIXED
- Lines 161-171: Contradictory logic - FIXED
- Line 145: Wrong log level - FIXED
- Lines 46-71: Commented duplicate code - FIXED

**Status:** ✓ COMPLETED

### 1.5 main.go - Typos
**File:** `cmd/david/main.go`

**Issues Fixed:**
- Lines 33, 37: "prdouction" typo - FIXED

**Status:** ✓ COMPLETED

### 1.6 magefile.go - Wrong Paths
**File:** `magefile.go`

**Issues Fixed:**
- All references to dave/davecli updated to david/bcpt
- Source paths corrected
- Archive names fixed

**Status:** ✓ COMPLETED

### 1.7 Additional Fix - nil Safety
**File:** `app/security.go`

**Issue Fixed:** authInfo could be nil when user not found

**Status:** ✓ COMPLETED

---

## PHASE 1 SUMMARY

**Files Modified:**
- app/security.go
- app/config.go
- app/crud.go
- app/fs.go
- cmd/david/main.go
- magefile.go

**Test Results:**
- TestAuthenticate: PASS (7/7 tests) ✓
- TestHandle: 2/3 pass (1 test needs Crud permissions)
- go vet: PASS ✓

**Note:** TestHandle/ok fails because test doesn't set Crud permissions. Phase 3 will fix this.

---

## PHASE 2: CODE QUALITY

### 2.1 Remove Dead Code
- ✓ crud.go: Removed unused contextKey variables
- security.go: testCrudType is used, kept

### 2.2 Standardize Error Handling
- Error handling follows Go best practices
- No major issues found

### 2.3 Improve Logging Consistency
- Log levels appropriately set
- Field names consistent

### 2.4 Code Formatting
**Status:** ✓ COMPLETED
- Applied: gofmt -s -w .
- Verified: go vet ./... passes

---

## PHASE 3: TEST FIXES

**Status:** IN PROGRESS

### 3.1 security_test.go

**Issues:**
- Test case 4 (line 132): Missing CrudType in expected result
- Test case 7 (line 188): Authenticated should be true for successful auth
- TestHandle "ok" test (line 303): User lacks Crud permissions

**Fixes Needed:**
```go
// Test case 4 - add CrudType:
&AuthInfo{
    Username:      "",
    Authenticated: false,
    CrudType:      noAuthentication.CrudType,  // ADD THIS
},

// Test case 7 - change Authenticated to true:
&AuthInfo{
    Username:      "foo",
    Authenticated: true,  // CHANGE FROM false
    CrudType:      &CrudType{Crud: "", Create: false, Read: false, Update: false, Delete: false},
},

// TestHandle "ok" test - add Crud permissions:
&App{
    Config: &Config{
        Users: map[string]*UserInfo{
            "foo": {
                Password:    GenHash([]byte("password")),
                Permissions: "crud",
                Crud:        &CrudType{Crud: "crud", Create: true, Read: true, Update: true, Delete: true},
            },
        },
    },
    Handler: &webdav.Handler{...},
},
```

**Status:** PENDING

### 3.2 fs_test.go

**Issues:**
- Test expectations may need updates due to improved permission logic

**Status:** PENDING

### 3.3 Add Missing Tests

**Status:** PENDING

**Target:** 90%+ code coverage

---

## PHASE 4: DEPENDENCY UPDATES

**Status:** PENDING

**Updates Needed:**
```
Go version: 1.21 → 1.22
github.com/sirupsen/logrus: v1.6.0 → v1.9.3
github.com/spf13/viper: v1.15.0 → v1.18.0
golang.org/x/crypto: v0.14.0 → v0.17.0
golang.org/x/net: v0.17.0 → v0.19.0
golang.org/x/term: v0.13.0 → v0.15.0
github.com/fsnotify/fsnotify: v1.6.0 → v1.7.0
github.com/spf13/cobra: v1.6.1 → v1.8.0
github.com/magefile/mage: v1.10.0 → v1.15.0
```

**Commands:**
```bash
go get -u ./...
go mod tidy
```

---

## PHASE 5: DOCUMENTATION

**Status:** PENDING

### 5.1 Readme.md
- Add missing installation steps 1-2
- Fix dave/david naming inconsistencies
- Add password hash notes
- Add CORS security warning

### 5.2 config-sample.yaml
- Add password hash generation note
- Set production to false
- Add CORS warning comment

---

## EXECUTION ORDER

1. ✓ PHASE 1 → Critical bugs (security, logic)
2. ✓ PHASE 2 → Code quality
3. ⏳ PHASE 3 → Test fixes (IN PROGRESS)
4. ⏳ PHASE 4 → Dependencies
5. ⏳ PHASE 5 → Documentation

---

## SUCCESS CRITERIA

- [x] All critical bugs fixed
- [x] Code passes go vet
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

**Current Status:**
- Phase 1: ✓ COMPLETED
- Phase 2: ✓ COMPLETED
- Phase 3: 30% Complete
- Phase 4: PENDING
- Phase 5: PENDING