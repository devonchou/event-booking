# Go Application - CI/CD with Jenkins and Kubernetes

This branch integrates CI/CD pipelines with Kubernetes deployment, automating the entire process from building, testing, and code analysis to deploying the application in a Kubernetes cluster.

## CI/CD Pipeline Features

- **Build**: Compiles the Go application.
- **Test**: Runs tests.
- **SonarQube Analysis**: Scans code for quality and security issues.
- **Quality Gate**: Ensures the application meets code quality standards.
- **Docker Build**: Builds a Docker image for the Go application.
- **Docker Push**: Pushes the Docker image to a registry (e.g., DockerHub).
- **Kubernetes Deployment**: Deploys or updates the application on a Kubernetes cluster using Helm charts.

Refer to the `Jenkinsfile` in this branch for more details.
