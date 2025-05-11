# Kanban Backend

A robust, scalable Go-based backend for a modern Kanban board application, designed with Hexagonal (Ports & Adapters) Architecture for clean separation of concerns and high maintainability.

## Tech Stack

-  **Language:** Go (Gin framework)
-  **Architecture:** Hexagonal (Ports & Adapters)
-  **Database:** PostgreSQL
-  **Real-time:** WebSocket (for instant updates)
-  **Auth:** JWT authentication & real-time authorization
-  **Deployment:** Docker & Docker Compose

## Features

-  RESTful API for all Kanban operations
-  Real-time updates via WebSocket
-  JWT-based authentication and real-time authorization
-  Team and project-based access control
-  Invitation and project sharing system
-  Domain-driven design with clear separation between business logic and infrastructure

## Main Domain Entities

-  **User:** Authentication, roles, and profile
-  **Project:** Project management and ownership
-  **Team:** Team-based access and roles (owner, admin, write, read)
-  **ProjectMember:** User's role and membership in projects/teams
-  **Column:** Kanban columns (customizable)
-  **Task:** Tasks within columns, with rich content
-  **Invitation:** Project invitations and status tracking

## Hexagonal Architecture Overview

-  **Core Domain:** Business logic and domain models (see `internal/core/domain`)
-  **Ports:** Interfaces for inbound (driver) and outbound (driven) operations (see `internal/core/ports`)
-  **Adapters:**
   -  **Driver Adapters:** HTTP handlers (REST API), WebSocket handlers (see `internal/adapters/driver`)
   -  **Driven Adapters:** PostgreSQL repositories (see `internal/adapters/driven`)

This structure ensures the core logic is isolated from frameworks and infrastructure, making the codebase highly testable and extensible.

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
4. Set up PostgreSQL locally (or use Docker)
5. Update `config/config.yaml` if needed
6. Run with Air (hot reloading):
   ```
   air
   ```

## API Endpoints

The API provides endpoints for:

-  User authentication & authorization
-  Project management
-  Kanban columns and tasks
-  Team management
-  Invitations and project sharing

---

For more details, see the code in the `internal/` directory and configuration in `config/`.
