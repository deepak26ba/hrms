# HRMS (Human Resource Management System)

## Overview
HRMS is a Go-based Human Resource Management System built with:

- Go
- Fiber v3
- GORM
- PostgreSQL
- JWT Authentication

The application follows a layered architecture with:
- Routes
- Handlers
- Services
- Repositories
- Models

## Features

### Authentication & Authorization
- User authentication
- JWT-based authorization
- Role management

### Employee Management
- Employee records
- Employee onboarding status
- Position management
- Department management

### Attendance
- Attendance tracking
- Attendance audit logs

### Leave Management
- Leave requests
- Leave balance tracking

### Payroll
- Payroll management

### Probation Management
- Employee probation tracking

## Project Structure

```text
hrms/
├── cmd/server/           # Application entry point
├── api/
│   ├── handlers/         # HTTP handlers
│   ├── routes/           # Route definitions
│   ├── service/          # Business logic
│   └── repository/       # Database operations
├── internals/
│   ├── config/           # Environment configuration
│   └── connection/       # Database connection
├── pkg/models/           # Data models
└── tests/                # Test files
```

## Prerequisites

- Go 1.25+
- PostgreSQL
- Git

## Environment Variables

Create a `.env` file in the project root.

```env
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=hrms
DB_PORT=5432
DB_SSLMODE=disable

JWT_SECRET_KEY=your-secret-key

SERVER_PORT=:8080

CLIENT_ID=google-client-id
CLIENT_SECRET=google-client-secret
```

## Installation

```bash
git clone <repository-url>
cd hrms

go mod tidy
```

## Database Setup

1. Create a PostgreSQL database.
2. Update the `.env` file.
3. Execute any required SQL scripts such as:

```bash
adminInsert.sql
```

## Run the Application

```bash
go run ./cmd/server
```

Expected output:

```text
Server is running
```

## API Modules

- Authentication
- Users
- Roles
- Employees
- Departments
- Positions
- Attendance
- Attendance Audit
- Leave Requests
- Leave Balances
- Payroll
- Probation
- Employee Boarding Status

## Testing

```bash
go test ./...
```

## Dependencies

Major dependencies:

- Fiber v3
- GORM
- PostgreSQL Driver
- JWT v5
- Validator
- Godotenv
- OAuth2

## License

Internal project.
