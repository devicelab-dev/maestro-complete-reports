# Analyzing Reports

## Using jq

### Extract Pass Rate

```bash
jq '.summary.passedFlows / .summary.totalFlows * 100' report.json
```

### Find Failed Flows

```bash
jq '.flows[] | select(.status == "failed") | .name' report.json
```

### Find Slow Commands (>5s)

```bash
jq '.flows[].commands[] | select(.durationMs > 5000) | {type, description, durationMs}' report.json
```

### Get Summary

```bash
jq '.summary' report.json
```

### List All Flow Names and Status

```bash
jq '.flows[] | {name, status, durationMs}' report.json
```

## Allure Report Generation

```bash
# Install Allure CLI
brew install allure  # macOS

# Generate report
allure generate ./reports/*/allure-results/ -o ./allure-report/

# Open in browser
allure open ./allure-report/
```

## JSON Report Schema

```json
{
  "schemaVersion": "1.0.0",
  "summary": {
    "totalFlows": 5,
    "passedFlows": 4,
    "failedFlows": 1,
    "totalCommands": 45,
    "totalDurationMs": 125340,
    "startTime": "2024-12-10T14:30:22+00:00",
    "endTime": "2024-12-10T14:32:27+00:00"
  },
  "device": {
    "id": "emulator-5554",
    "name": "Pixel 6",
    "platform": "android"
  },
  "flows": [...]
}
```

## JUnit XML Structure

```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="5" failures="1" skipped="0">
  <testsuite name="Login Tests" tests="3" failures="1">
    <testcase name="Valid Login" classname="login.yaml" time="23.0">
      <properties>
        <property name="file" value="login.yaml"/>
        <property name="device.name" value="Pixel 6"/>
        <property name="device.id" value="emulator-5554"/>
        <property name="device.platform" value="android"/>
      </properties>
    </testcase>
  </testsuite>
</testsuites>
```

## Sample Reports

- [report.json](../samples/report.json) - JSON report example
- [junit-report.xml](../samples/junit-report.xml) - JUnit XML example
