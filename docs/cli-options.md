# CLI Report Options

## Report Directory Configuration

| Option | Description | Default |
|--------|-------------|---------|
| `--report-dir <path>` | Base directory for reports | `./reports` |
| `--flatten-report-output` | Disable timestamp subfolders | `false` |

**Examples:**
```bash
# Default: creates ./reports/YYYY-MM-DD_HHmmss/
maestro test flows/

# Custom directory with timestamp subfolder
maestro test flows/ --report-dir /var/test-reports

# Flat output (no timestamp subfolder) - useful for CI
maestro test flows/ --report-dir ./reports --flatten-report-output
```

## Output Format Configuration

| Option | Description | Default |
|--------|-------------|---------|
| `--format <format>` | Report format (JUNIT, HTML, NOOP) | `NOOP` |
| `--output <path>` | File path for formatted report | - |
| `--test-suite-name <name>` | Custom suite name for JUnit | `"Test Suite"` |

## Debug Output Configuration

| Option | Description | Default |
|--------|-------------|---------|
| `--debug-output <path>` | Directory for debug files | `~/.maestro/tests/` |
| `--flatten-debug-output` | Disable timestamp subfolders for debug | `false` |

## Generated Report Files

```
reports/
└── 2024-12-10_143022/
    ├── report.json              # Complete JSON report
    ├── index.html               # Interactive HTML report
    ├── junit-report.xml         # JUnit XML for CI/CD
    ├── maestro.log              # Application logs
    ├── screenshots/             # Test screenshots
    └── allure-results/          # Allure report data
```

| File | Purpose | Usage |
|------|---------|-------|
| `report.json` | Complete test data in JSON | Programmatic analysis, custom dashboards |
| `index.html` | Interactive web report | Manual review, stakeholder communication |
| `junit-report.xml` | JUnit XML format | CI/CD integration (Jenkins, GitLab, etc.) |
| `maestro.log` | Real-time execution logs | Debugging, audit trails |
| `screenshots/` | Captured screenshots | Visual verification |
| `allure-results/` | Allure-compatible data | Trend analysis, historical comparison |
