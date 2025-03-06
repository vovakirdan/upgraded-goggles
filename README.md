# user-post-service: User & Post Management Service
> [Документация на Русском](README.ru.md)
## Table of Contents
- [Introduction](#introduction)
- [Project Structure](#project-structure)
- [Installation & Setup](#installation--setup)
- [Generate .proto files](#generate-proto-files)
- [Configuration](#configuration)
- [Services Overview](#services-overview)
- [API Gateway](#api-gateway)
- [gRPC Services](#grpc-services)
- [Database Schema](#database-schema)
- [Logging](#logging)
- [Authentication & Security](#authentication--security)
- [Swagger Documentation](#swagger-documentation)
- [Testing](#testing)
- [Deployment](#deployment)

---

## Introduction
The **user-post-service** project is a multi-functional service designed to manage users and posts. The project integrates authentication, post management, logging, and gRPC-based service communication. An API Gateway provides an HTTP interface for external access, while Swagger is used for API documentation.

### Key Features:
1. **User & Post Management**
   - User registration and authentication (JWT-based authentication)
   - CRUD operations for posts
   - Data validation
2. **Logging System**
   - Logs all incoming and outgoing requests
   - Captures errors and exceptions
   - Stores logs in files
3. **gRPC & API Gateway**
   - Implements gRPC for inter-service communication
   - API Gateway to route HTTP requests to gRPC services
   - Authorization middleware in API Gateway
4. **Swagger API Documentation**
   - Automatic documentation generation based on API handlers
   - User-friendly interface for API testing
5. **Comprehensive Testing & Documentation**
   - Integration and unit tests for all components
   - Clear documentation on setup and usage

---

## Project Structure
```
upgraded-goggles/
├── README.md
├── api/
│   ├── gateway/
│   │   ├── config.go
│   │   ├── main.go
│   │   ├── middleware.go
│   │   ├── routes.go
│   │   └── swagger/
│   │       ├── index.html
│   │       └── swagger.yaml
│   ├── grpc/
│   │   ├── post_service.go
│   │   ├── server.go
│   │   └── user_service.go
│   └── proto/
│       ├── common.proto
│       ├── gateway.proto
│       ├── post.proto
│       └── user.proto
├── cmd/
│   ├── post-service/
│   │   └── main.go
│   └── user-service/
│       └── main.go
├── internal/
│   ├── logger/
│   │   └── logger.go
│   ├── post/
│   │   ├── models.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── validator.go
│   └── user/
│       ├── models.go
│       ├── repository.go
│       ├── service.go
│       └── validator.go
├── pkg/
│   └── database/
│       └── db.go
├── test/
│   ├── integration_test.go
│   ├── middleware_test.go
│   ├── post_test.go
│   └── user_test.go
├── go.mod
├── go.sum
├── Dockerfile
└── docker-compose.yml
```

---

## Installation & Setup
### Prerequisites
- Go 1.23+
- Docker & Docker Compose
- PostgreSQL

### Steps
1. Clone the repository:
   ```sh
   git clone https://github.com/your-repo/upgraded-goggles.git
   cd upgraded-goggles
   ```
2. Set up environment variables:
   ```sh
   export HTTP_PORT=":8080"
   export USER_SERVICE_ADDRESS="localhost:50051"
   export POST_SERVICE_ADDRESS="localhost:50052"
   ```
   Or create a `.env` file in the root directory with the following variables:
   ```
   # Database configuration
   DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable

   # API Gateway configuration
   HTTP_PORT=:8080
   USER_SERVICE_ADDRESS=localhost:50051
   POST_SERVICE_ADDRESS=localhost:50052
   GATEWAY_ADDRESS=localhost:50053
   ```
3. Build the project:
   ```sh
   go mod download
   ```
4. Start services with Docker:
   ```sh
   docker-compose up --build
   ```
5. Verify services:
   ```sh
   curl http://localhost:8080/health
   ```

--- 

## Generate .proto files
To generate .proto files, run the following command:
1. Install `protoc` and `protoc-gen-go` plugins:
   ```sh
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```
2. Generate .proto files:
   * Common
   ```bash
   protoc -I api/proto \
  --go_out=paths=source_relative:api/proto/common \
   --go-grpc_out=paths=source_relative:api/proto/common \
   api/proto/common/common.proto
   ```
   * Gateway
   ```bash
   protoc -I api/proto \
   --go_out=paths=source_relative:api/proto/gateway \
   --go-grpc_out=paths=source_relative:api/proto/gateway \
   api/proto/gateway/gateway.proto
   ```
   * Post
   ```bash
   protoc -I api/proto \
   --go_out=paths=source_relative:api/proto/post \
   --go-grpc_out=paths=source_relative:api/proto/post \
   api/proto/post/post.proto
   ```
   * User
   ```bash
   protoc -I api/proto \
   --go_out=paths=source_relative:api/proto/user \
   --go-grpc_out=paths=source_relative:api/proto/user \
   api/proto/user/user.proto
   ```

---

## Configuration
Configurations are loaded from environment variables:
| Variable               | Description                             | Default Value |
|------------------------|-----------------------------------------|---------------|
| `HTTP_PORT`           | Port for API Gateway                   | `:8080`       |
| `USER_SERVICE_ADDRESS`| gRPC User service address               | `localhost:50051` |
| `POST_SERVICE_ADDRESS`| gRPC Post service address               | `localhost:50052` |
| `DATABASE_URL`        | PostgreSQL connection string            | -             |

---

## Services Overview
### API Gateway
Handles HTTP requests and forwards them to gRPC services.
- Uses JWT authentication middleware.
- Serves static Swagger documentation.
- Provides endpoints for users and posts.

### User Service
Manages user registration, authentication, and retrieval.
- Uses bcrypt for password hashing.
- Returns JWT token upon successful login.

### Post Service
Manages posts CRUD operations.
- Validates post data before saving.
- Supports updating and deleting posts.

---

## API Gateway
Routes API requests to gRPC services.
### Endpoints:
| Method | Path                 | Description           |
|--------|----------------------|-----------------------|
| POST   | `/v1/users/register` | Register new user    |
| POST   | `/v1/users/login`    | Login user           |
| GET    | `/v1/users/{id}`     | Get user details     |
| POST   | `/v1/posts`          | Create post          |
| GET    | `/v1/posts/{id}`     | Get post details     |
| PUT    | `/v1/posts/{id}`     | Update post          |
| DELETE | `/v1/posts/{id}`     | Delete post          |

---

## gRPC Services
### UserService
- `RegisterUser(RegisterRequest) -> RegisterResponse`
- `LoginUser(LoginRequest) -> LoginResponse`
- `GetUser(UserRequest) -> User`

### PostService
- `CreatePost(CreatePostRequest) -> CreatePostResponse`
- `GetPost(GetPostRequest) -> GetPostResponse`
- `UpdatePost(UpdatePostRequest) -> UpdatePostResponse`
- `DeletePost(DeletePostRequest) -> DeletePostResponse`

---

## Logging
- Uses `log.Logger` to log requests and errors.
- Logs are stored in `logs/app.log`.

---

## Authentication & Security
- Uses JWT authentication middleware.
- Passwords are hashed using bcrypt.
- Authorization checks on protected routes.

---

## Swagger Documentation
Swagger UI is available at:
```
http://localhost:8080/swagger/index.html
```

---

## Testing
Unit tests are available in the `test/` directory.
Run tests with:
```sh
go test ./test/...
```

---

## Deployment
To deploy with Docker:
```sh
docker-compose up --build -d
```
