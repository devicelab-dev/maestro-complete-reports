# Test Suite Organization

Maestro supports hierarchical test organization using `suite` and `test` commands for meaningful groupings in reports.

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

## Complete Examples

See [samples/flows/](../samples/flows/) for working examples.
