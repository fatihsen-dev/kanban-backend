# Kanban Backend

A Go-based backend for a Kanban board application.

## Features

-  RESTful API using Gin
-  PostgreSQL database
-  WebSocket for real-time updates
-  JWT authentication

## Running with Docker

### Prerequisites

-  Docker
-  Docker Compose

### Steps to Run

1. Clone the repository:

   ```
   git clone https://github.com/fatihsen-dev/kanban-backend.git
   cd kanban-backend
   ```

2. Start the services:

   ```
   docker-compose up -d
   ```

   This will start both the API and PostgreSQL database.

3. Access the API at `http://localhost:5000`

## Development Setup

1. Install Go (1.23 or later)
2. Clone the repository
3. Install dependencies:
   ```
   go mod download
   ```
4. Set up PostgreSQL locally
5. Update `config/config.yaml` if needed
6. Run with Air (hot reloading):
   ```
   air
   ```

## API Endpoints

The API provides endpoints for:

-  User authentication
-  Project management
-  Kanban columns and tasks
-  Team management
-  Invitations
