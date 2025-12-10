# Test Suite Organization

Maestro supports hierarchical test organization using `suite` and `test` commands for meaningful groupings in reports.

> **Note:** `describe` and `it` are aliases for `suite` and `test` respectively.

## Structure Overview

```
auth-suite.yaml              # Parent suite file (uses runFlow)
auth-suite/
  ├── login-valid-primary.yaml      # suite: "Valid Login"
  ├── login-valid-alternate.yaml    # suite: "Valid Login"
  ├── login-invalid-password.yaml   # suite: "Invalid Login"
  └── login-invalid-username.yaml   # suite: "Invalid Login"
```

## Parent Suite File

```yaml
# auth-suite.yaml
appId: com.testhiveapp
name: "Authentication Test Suite"
tags:
  - auth
  - suite
---
- runFlow: auth-suite/login-valid-primary.yaml
- runFlow: auth-suite/login-valid-alternate.yaml
- runFlow: auth-suite/login-invalid-password.yaml
- runFlow: auth-suite/login-invalid-username.yaml
```

## Individual Flow File

Each flow declares its `suite` and `test`:

```yaml
# auth-suite/login-valid-primary.yaml
appId: com.testhiveapp
name: "Login - Valid Primary Credentials"
---
- launchApp:
    clearState: true

- suite: "Valid Login"

- test: "Should login with primary credentials"
  steps:
    - assertVisible: "Welcome Back"
    - runFlow:
        file: ../common/login.yaml
        env:
          username: devicelab
          password: robustest
    - assertVisible: "TestHive"
    - assertVisible:
        id: "products-screen"
```

## How It Maps to Reports

| Flow YAML | JSON Report |
|-----------|-------------|
| `suite: "Valid Login"` | `suite.name` |
| `test: "Should login..."` | `flow.name` |
| filename | `flow.fileName` |

## Grouping Tests Across Files

Flows with the same `suite` name are merged in reports:

| File | suite | test |
|------|-------|------|
| login-valid-primary.yaml | Valid Login | Should login with primary credentials |
| login-valid-alternate.yaml | Valid Login | Should login with alternate credentials |
| login-invalid-password.yaml | Invalid Login | Should show error for invalid password |
| login-invalid-username.yaml | Invalid Login | Should show error for invalid username |

**Result:** Two suites in reports - "Valid Login" (2 tests) and "Invalid Login" (2 tests).

## Running the Suite

```bash
maestro test auth-suite.yaml
```

## Without suite/test Commands

When flows don't use `suite`/`test` commands:
- All tests are grouped under a single suite
- The suite name defaults to "Test Suite" (or use `--test-suite-name`)
- Test names come from the flow's `name` field or filename

## Multiple Tests Per Flow

A single flow file can contain multiple tests:

```yaml
# login-tests.yaml
appId: com.example.app
---
- launchApp:
    clearState: true

- suite: "Login Validation"

- test: "Should show error for empty username"
  steps:
    - tapOn:
        id: "login-button"
    - assertVisible: "Username is required"

- test: "Should show error for empty password"
  steps:
    - tapOn:
        id: "username-input"
    - inputText: "testuser"
    - tapOn:
        id: "login-button"
    - assertVisible: "Password is required"

- test: "Should login successfully"
  steps:
    - tapOn:
        id: "username-input"
    - inputText: "testuser"
    - tapOn:
        id: "password-input"
    - inputText: "password123"
    - tapOn:
        id: "login-button"
    - assertVisible:
        id: "home-screen"
```

## Shared Flows

Use `runFlow` to reuse common steps:

```yaml
# common/login.yaml
- tapOn:
    id: "username-input"
- inputText: ${username}
- tapOn:
    id: "password-input"
- inputText: ${password}
- tapOn:
    id: "login-button"
```

```yaml
# auth-suite/login-valid-primary.yaml
- suite: "Valid Login"
- test: "Should login with primary credentials"
  steps:
    - runFlow:
        file: ../common/login.yaml
        env:
          username: testuser
          password: password123
    - assertVisible:
        id: "home-screen"
```

## Running Tests

```bash
# Run a suite file
maestro test auth-suite.yaml

# Run a folder of flows
maestro test auth-suite/

# Run with custom report directory
maestro test auth-suite.yaml --report-dir ./reports

# Run with flat output (for CI)
maestro test auth-suite.yaml --report-dir ./reports --flatten-report-output
```

## Best Practices

1. **One suite per feature** - Group related tests under the same `suite` name
2. **Descriptive test names** - Use "Should..." format for clarity
3. **Shared flows for common steps** - Avoid duplicating login, navigation, etc.
4. **Keep tests independent** - Each test should work standalone with `clearState: true`
5. **Use tags** - Add tags for filtering (`smoke`, `regression`, `auth`, etc.)
