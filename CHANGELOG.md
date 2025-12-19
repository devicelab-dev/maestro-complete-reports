# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- **Enhanced assertion output with variable values for debugging**

  Improved `assertTrue`/`assertCondition` output to show actual variable values when evaluating complex JavaScript expressions, making it easier to debug failing assertions.

  **Before:**
  ```
  ${output.actualText.includes("expected")} is false
  ```

  **After:**
  ```
  ${output.actualText.includes("expected")} is false (output.actualText = "actual value")
  → "actual value".includes("expected")
  ```

- Enhanced JUnit report with custom labels and failure types
- Device details added to shard log prefix
- Colorful banners and improved API endpoints
- API integration for JAR download with version check
- Makefile for cross-platform builds

### Fixed

- Nested zip extraction issues
