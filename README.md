<p align="center">
  <h1 align="center">Maestro Complete Reports</h1>
  <p align="center">
    <strong>Enhanced reporting for Maestro mobile UI testing</strong>
  </p>
  <p align="center">
    <a href="https://github.com/devicelab-dev/maestro-complete-reports/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="License"></a>
    <a href="https://github.com/devicelab-dev/maestro-complete-reports/releases"><img src="https://img.shields.io/badge/maestro-2.0.9%20%7C%202.0.10-green.svg" alt="Maestro Version"></a>
    <a href="https://devicelab.dev"><img src="https://img.shields.io/badge/by-DeviceLab.dev-orange.svg" alt="DeviceLab.dev"></a>
  </p>
</p>

---

> **Community Build** - Not affiliated with Maestro or mobile.dev
>
> Built by [DeviceLab.dev](https://devicelab.dev) - Turn Your Devices Into a Distributed Device Lab

## Features

| Feature | Description |
|---------|-------------|
| **JSON Reports** | Complete test data for programmatic analysis |
| **HTML Reports** | Interactive web reports for stakeholders |
| **JUnit XML** | CI/CD integration (Jenkins, GitLab, GitHub Actions) |
| **Allure Support** | Trend analysis and historical comparison |
| **Test Suites** | Hierarchical suite/test organization |
| **Screenshots** | Automatic capture with linking |

## Quick Start

### Option 1: Using the Binary (Automated)

```bash
# Setup - detects Maestro, backs up JARs, downloads and replaces with patched JARs
./maestro-complete-reports setup

# Restore original JARs
./maestro-complete-reports restore
```

### Option 2: Manual Installation

```bash
# Copy JARs to Maestro lib directory
cp jars/2.0.10/*.jar ~/.maestro/lib/

# To restore
cp ~/.maestro/backup/*.jar ~/.maestro/lib/
```

## Usage

```bash
maestro test flows/ --report-dir report_folder
```

## Sample Reports

| Format | Sample |
|--------|--------|
| JSON | [report.json](samples/report.json) |
| JUnit XML | [junit-report.xml](samples/junit-report.xml) |
| Sample Flows | [samples/flows/](samples/flows/) |

### Console Output

```
===== Test Summary =====

✅ Invalid Login
  ✅ Should show error for invalid password
  ✅ Should show error for invalid username
❌ Valid Login
  ✅ Should login with alternate credentials
  ❌ Should login with alternate credentials. fail as password is wrong
  ✅ Should login with primary credentials

========================
Suites: 1 passed, 1 failed
Tests:  4 passed, 1 failed

  4/5 Flows Passed

==========================================================================================
  Flow                                   Status  Steps   Pass   Fail   Skip   Duration
------------------------------------------------------------------------------------------
  Login - Invalid Password                    ✓     10     10      0      0      23.0s
  Login - Valid Alternate Credentials         ✓     13     13      0      0      22.0s
  Login - Valid Alternate Credentials         ✕     11     10      1      0      30.0s
  Login - Valid Primary Credentials           ✓     13     13      0      0      23.0s
  Login - Invalid Username                    ✓     10     10      0      0      22.0s
------------------------------------------------------------------------------------------
  Total                                             57     56      1      0      2m 0s
==========================================================================================
```

### HTML Report

![HTML Report](assets/html-report.png)

## Documentation

| Doc | Description |
|-----|-------------|
| [CLI Options](docs/cli-options.md) | All report configuration options |
| [Test Suites](docs/test-suites.md) | Organizing tests with suite/test commands |
| [CI/CD Integration](docs/ci-cd.md) | GitHub, GitLab, Jenkins, CircleCI examples |

## Backup Location

Original JARs are backed up to `~/.maestro/backup/`

---

## Contributing

Issues and PRs welcome at [GitHub](https://github.com/devicelab-dev/maestro-complete-reports).

## License

Apache 2.0 (same as Maestro)

## Disclaimer

This project is not affiliated with, endorsed by, or connected to mobile.dev or the official Maestro project. This tool patches your existing Maestro installation to add reporting functionality not yet available in the official release.

**Use at your own risk.** We recommend switching to official Maestro once these reporting features are officially released.
