# s4s Backend

## Description
s4s (Smart for Sales) backend is the server-side component of a SaaS platform for visual sales automation, targeted at SMB non-tech users (sales/marketing managers). It handles workflow orchestration, integrations, user auth, monitoring, and monetization. Inspired by n8n, but simplified for sales-focused use cases with drag-and-drop workflows, templates, and freemium model.

Key features:
- API for frontend (REST with JWT auth).
- Workflow engine with triggers, actions, logic nodes.
- Integrations with services (e.g., Gmail, Slack via OAuth).
- Logging and notifications (Elasticsearch, RabbitMQ).
- Monetization via Stripe (freemium limits).
- Scalable architecture (monolith with MSA-ready modules).

Built with performance in mind for 10k+ users (as of September 2025 benchmarks).

## Tech Stack
- **Language**: Go 1.23 (for concurrency with goroutines).
- **Database**: PostgreSQL (main storage) + Redis (caching/queues).
- **Queue**: RabbitMQ (async workflow executions).
- **Logging/Monitoring**: Elasticsearch (logs), Sentry (errors).
- **Auth**: JWT (with bcrypt hashing).
- **API**: REST (with GraphQL option for v2).
- **Deployment**: Docker, Kubernetes-ready.
- **Tools**: Gin (web framework), GORM (ORM), Viper (config), Zerolog (logging).

## Prerequisites
- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 16+
- Redis 7+
- RabbitMQ 3.12+
- Environment variables (see .env.example)

## Installation
1. Clone the repo:
   ```
   git clone https://github.com/your-org/s4s-backend.git
   cd s4s-backend
   ```

2. Set up environment:
    - Copy `.env.example` to `.env` and fill in values (e.g., DB_URL, JWT_SECRET, STRIPE_KEY).

3. Install dependencies:
   ```
   go mod tidy
   ```

4. Run with Docker (recommended for dev):
   ```
   docker-compose up -d
   ```

5. Or run locally:
   ```
   go run cmd/main.go
   ```

The server starts at `http://localhost:8080/api/v1`. Health check: `/health`.

## Configuration
`.env` example:
```
DB_URL=postgres://user:pass@localhost:5432/s4s?sslmode=disable
REDIS_ADDR=localhost:6379
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
JWT_SECRET=your-secret-key
STRIPE_KEY=sk_test_123
SENTRY_DSN=https://your-sentry-dsn
ELASTICSEARCH_URL=http://localhost:9200
PORT=8080
```

- For production: Use secure secrets (e.g., Vault or env vars in K8s).
- Migrations: Use Goose or GORM auto-migrate in main.go.

## API Endpoints
Base URL: `/api/v1`

- **Auth**:
    - POST `/auth/register`: Register user (input: fullName, email, password; output: JWT token).
    - POST `/auth/login`: Login (input: email, password; output: JWT).
    - POST `/auth/forgot-password`: Reset link (input: email).

- **Workflows**:
    - GET `/workflows`: List workflows (query: active, page, limit).
    - POST `/workflows`: Create (input: name, json).
    - GET/PUT/DELETE `/workflows/{id}`: Get/update/delete.
    - POST `/workflows/{id}/test`: Test (input: testData).
    - POST `/workflows/{id}/run`: Run (async, returns executionId).

- **Connections**:
    - GET/POST `/connections`: List/create (input: service, credentials).
    - GET/PUT/DELETE `/connections/{id}`: Manage.

- **Templates**:
    - GET `/templates`: List (query: category).
    - GET `/templates/{id}`: Get.

- **Executions**:
    - GET `/executions`: List (query: workflowId, status).
    - GET `/executions/{id}`: Get details.

- **Subscriptions**:
    - GET `/subscriptions`: Get status.
    - POST `/subscriptions/upgrade`: Upgrade (input: plan; output: Stripe session URL).

- **Notifications**:
    - GET/PUT `/notifications/settings`: Get/update (email, slack, channels).

- **Admin (RBAC-protected)**:
    - GET `/admin/users`: List users.
    - GET `/admin/workflows`: List all workflows.

Full docs: Use Swagger at `/swagger/index.html` (integrated with swaggo).

## Testing
- Unit/Integration: `go test ./... -v` (coverage 80%+).
- E2E: Use Postman or Cypress (frontend-integrated).
- Load test: k6 for 10k users (target: <2 sec latency).

## Deployment
- **Local**: `go run`.
- **Docker**: `docker build -t s4s-backend .` then `docker run -p 8080:8080 -env-file .env s4s-backend`.
- **Prod**: Kubernetes with Helm chart (included in repo). Scale with replicas for workers. Monitor with Prometheus/Grafana.

## Contributing
- Fork repo, create branch (`feature/xxx`).
- Commit with conventional commits.
- PR with tests, lint (golangci-lint run).
