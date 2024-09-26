# Go Application - Kubernetes Deployment

This branch includes Kubernetes definition files for deploying both the application (API) and database to a Kubernetes cluster.

## Features

- **API Deployment (`api-dep.yaml`)**: Defines the deployment for the Go API, specifying replicas, container image, etc.
- **API Service (`api-svc.yaml`)**: Exposes the Go API externally via a LoadBalancer service, making it accessible outside the cluster.
- **Database Deployment (`db-dep.yaml`)**: Defines the deployment for the database (e.g., MySQL), including container specifications.
- **Database Service (`db-svc.yaml`)**: Exposes the database internally within the cluster using a ClusterIP service.
- **Database Secrets (`db-secret.yaml`)**: Stores sensitive information like database credentials securely as Kubernetes secrets.
