# Go Application - Main Branch

This project is a Go application built using the Gin web framework. The code structure follows Dependency Injection (DI) principles to ensure flexibility and testability.

## Branch Overview

- **docker**: Containerizes the Go application using Docker.
- **ci-jenkins**: Implements a CI pipeline with Jenkins for building, testing, SonarQube analysis, Nexus artifact storage, and Slack notifications.
- **kubernetes**: Adds Kubernetes definition files and deploys the application to a Kubernetes cluster.
- **cicd-jenkins-kubernetes**: Adds a full CI/CD pipeline with Jenkins for building, testing, analyzing code, building Docker image, pushing image to a registry, and deploying/updating with Helm charts on Kubernetes.

## API Endpoints

All API routes start with `/api`.

### User Endpoints

- **POST /users**: Create a new user.
- **POST /users/login**: Login and verify user credentials.
- **GET /users**: Retrieve all user data (admin access only).
- **GET /users/:userId**: Retrieve user data by user ID.
- **PUT /users/:userId**: Update user data by user ID.
- **DELETE /users/:userId**: Delete user by user ID.

> Note: All user-related endpoints except `POST /users` and `POST /users/login` require JWT authentication.

### Event Endpoints

- **GET /events**: Get all events.
- **GET /events/:eventId**: Get event data by event ID.
- **POST /events**: Create a new event.
- **PUT /events/:eventId**: Update event data by event ID (only the event owner can modify).
- **DELETE /events/:eventId**: Delete event by event ID (only the event owner can delete).
- **POST /events/:eventId/register**: Register a user for an event.
- **DELETE /events/:eventId/register**: Cancel user registration for an event.
- **GET /events/:eventId/attendees**: Get a list of attendee emails (event owner access only).

> Note: All event-related endpoints except `GET /events` and `GET /events/:eventId` require JWT authentication.
