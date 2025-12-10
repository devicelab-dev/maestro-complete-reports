# CI/CD Integration

## GitHub Actions

```yaml
- name: Run Maestro Tests
  run: |
    maestro test flows/ \
      --report-dir ./reports \
      --flatten-report-output

- name: Upload Reports
  uses: actions/upload-artifact@v3
  with:
    name: test-reports
    path: reports/

- name: Publish Results
  uses: EnricoMi/publish-unit-test-result-action@v2
  with:
    files: './reports/junit-report.xml'
```

## GitLab CI

```yaml
test:
  script:
    - maestro test flows/
        --report-dir $CI_PROJECT_DIR/reports
        --flatten-report-output
  artifacts:
    reports:
      junit: reports/junit-report.xml
    paths:
      - reports/
    expire_in: 30 days
```

## Jenkins Pipeline

```groovy
pipeline {
    agent any
    stages {
        stage('Run Tests') {
            steps {
                sh '''
                    maestro test flows/ \
                      --report-dir "${WORKSPACE}/reports" \
                      --flatten-report-output
                '''
            }
        }
    }
    post {
        always {
            junit 'reports/junit-report.xml'
            archiveArtifacts artifacts: 'reports/**'
            publishHTML([
                reportDir: 'reports',
                reportFiles: 'index.html',
                reportName: 'Maestro Report'
            ])
        }
    }
}
```

## CircleCI

```yaml
- run:
    name: Run Tests
    command: |
      maestro test flows/ \
        --report-dir /tmp/reports \
        --flatten-report-output

- store_test_results:
    path: /tmp/reports

- store_artifacts:
    path: /tmp/reports
```

## Best Practices

1. **Use `--flatten-report-output`** in CI for predictable paths
2. **Archive reports as artifacts** for historical analysis
3. **Generate JUnit reports** for CI platform integration
4. **Keep maestro.log** for debugging failed tests
5. **Set artifact expiration** to manage storage (e.g., 30 days)
