# Go Application - CI with Jenkins

This branch adds a Continuous Integration (CI) pipeline using Jenkins.

## Jenkins Pipeline Features

- **Build**: Compiles the Go application.
- **Test**: Runs automated tests.
- **SonarQube Analysis**: Scans code for quality and security issues.
- **Quality Gate**: Enforces quality standards before proceeding with the pipeline.
- **Nexus Artifact Management**: Uploads and stores build artifacts (e.g., binaries) in Nexus.
- **Slack Notifications**: Sends pipeline status updates to a Slack channel.

Refer to the `Jenkinsfile` in this branch for more details.
