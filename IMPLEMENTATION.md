# Fleet Risk Intelligence System - Implementation Guide

## 📋 Project Overview

**Project Name**: Fleet Risk Intelligence MVP
**Purpose**: Real-time fleet management system for monitoring vehicle telemetry, analyzing driver behavior, and generating risk assessments
**Tech Stack**: Go 1.24, React 18, TypeScript, GraphQL, MySQL, Redis, Docker
**Architecture**: Microservices with event-driven communication

## 🏗️ System Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  React Frontend │────▶│   API Gateway    │────▶│  Auth Service   │
│   (Port 3000)   │     │   (Port 8080)    │     │  (Port 8084)    │
└─────────────────┘     └──────────────────┘     └─────────────────┘
                               │
                               ▼
        ┌──────────────────────┴───────────────────────┐
        │                                               │
        ▼                                               ▼
┌─────────────────┐                            ┌─────────────────┐
│Telemetry Ingest │                            │  Risk Engine    │
│  (Port 8081)    │                            │  (Port 8082)    │
└─────────────────┘                            └─────────────────┘
        │                                               │
        └──────────────┬───────────────────────────────┘
                       ▼
                ┌─────────────────┐
                │   WebSocket     │
                │  (Port 8083)    │
                └─────────────────┘
                       │
        ┌──────────────┴───────────────────┐
        ▼                                   ▼
   ┌─────────┐                        ┌─────────┐
   │  MySQL  │                        │  Redis  │
   └─────────┘                        └─────────┘
```

## 📁 Project Structure

```
samsara_healthcare/
├── pkg/                      # Shared packages
│   ├── auth/                # JWT & middleware
│   ├── config/             # Configuration management
│   ├── database/           # Database connection
│   ├── errors/             # Error handling
│   ├── models/             # GORM models
│   ├── server/             # Base server setup
│   └── validation/         # Request validation
│
├── services/               # Microservices
│   ├── api/               # GraphQL & REST API
│   ├── auth/              # Authentication service
│   ├── risk-engine/       # Risk analysis
│   ├── telemetry-ingest/  # Data ingestion
│   └── websocket/         # Real-time updates
│
├── frontend/              # React application
│   ├── src/
│   │   ├── components/    # Reusable components
│   │   ├── contexts/      # React contexts
│   │   ├── graphql/       # Queries & mutations
│   │   ├── pages/         # Route pages
│   │   └── types/         # TypeScript types
│   └── package.json
│
├── docker/                # Docker configurations
├── scripts/               # Utility scripts
├── .github/workflows/     # CI/CD pipelines
├── docker-compose.dev.yml # Development environment
├── go.mod                # Go dependencies
└── go.work               # Go workspace
```

## 🔑 Key Components

### 1. **Authentication System**
- **Location**: `pkg/auth/`, `services/auth/`
- **Features**: JWT tokens, bcrypt passwords, role-based access
- **Roles**: super_admin, fleet_admin, fleet_manager, driver
- **Implementation**:
  - JWT with 24-hour expiration
  - Session tracking in database
  - Fleet-level access control

### 2. **Data Models** (`pkg/models/models.go`)
```go
- User: Authentication and authorization
- Fleet: Organization/company entity
- Vehicle: Fleet vehicles with telemetry
- Driver: Vehicle operators
- TelemetryEvent: Raw vehicle data
- RiskEvent: Analyzed risk incidents
- Alert: System notifications
- DriverScore: Performance metrics
- Session: Login sessions
```

### 3. **API Service** (`services/api/`)
- GraphQL endpoint at `/graphql`
- REST endpoints at `/api/v1/*`
- Protected by JWT authentication
- GraphQL playground in dev mode

### 4. **Frontend Application** (`frontend/`)
- React 18 with TypeScript
- Apollo GraphQL client
- Protected routes with auth context
- Tailwind CSS for styling
- Vite for development/build

### 5. **Database Schema**
- MySQL 8.0 with GORM auto-migration
- Foreign key relationships
- JSON fields for flexible data
- Indexes on frequently queried fields

## 🚀 Development Setup

### Prerequisites
```bash
# Required software
- Go 1.24+
- Node.js 20+
- MySQL 8.0
- Redis 7+
- Docker & Docker Compose
```

### Environment Variables
```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=fleet
DB_PASSWORD=devpass
DB_NAME=fleet_dev

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Server
API_PORT=8080
LOG_LEVEL=debug
ENV=development
```

### Local Development (Without Docker)

1. **Start MySQL and Redis**
```bash
# MySQL
mysql -u root -p
CREATE DATABASE fleet_dev;
CREATE USER 'fleet'@'localhost' IDENTIFIED BY 'devpass';
GRANT ALL PRIVILEGES ON fleet_dev.* TO 'fleet'@'localhost';

# Redis
redis-server
```

2. **Run Backend Services**
```bash
# Terminal 1: Auth Service
cd services/auth
go run main.go

# Terminal 2: API Service
cd services/api
go run main.go

# Terminal 3: Telemetry Ingest
cd services/telemetry-ingest
go run main.go

# Terminal 4: Risk Engine
cd services/risk-engine
go run main.go

# Terminal 5: WebSocket Service
cd services/websocket
go run main.go
```

3. **Run Frontend**
```bash
cd frontend
npm install
npm run dev
```

### Docker Development
```bash
# Start all services
docker-compose -f docker-compose.dev.yml up

# Rebuild after changes
docker-compose -f docker-compose.dev.yml up --build

# View logs
docker-compose -f docker-compose.dev.yml logs -f [service_name]

# Stop services
docker-compose -f docker-compose.dev.yml down
```

## 🔧 Service Configuration

### Port Mapping
| Service | Internal Port | External Port | Description |
|---------|--------------|---------------|-------------|
| Frontend | 3000 | 3000 | React dev server |
| API | 8080 | 8080 | GraphQL & REST |
| Telemetry | 8080 | 8081 | Data ingestion |
| Risk Engine | 8080 | 8082 | Risk analysis |
| WebSocket | 8080 | 8083 | Real-time updates |
| Auth | 8080 | 8084 | Authentication |
| MySQL | 3306 | 3306 | Database |
| Redis | 6379 | 6379 | Cache & pub/sub |

### API Endpoints

**Authentication** (`/api/v1/auth/`)
- POST `/login` - User login
- POST `/register` - User registration
- GET `/me` - Get profile (protected)
- PUT `/me` - Update profile (protected)
- POST `/logout` - Logout (protected)

**GraphQL** (`/graphql`)
- Requires Bearer token
- Full CRUD for all entities
- Real-time subscriptions

**REST API** (`/api/v1/`)
- GET `/fleets` - List fleets
- GET `/vehicles` - List vehicles
- GET `/drivers` - List drivers
- GET `/risk-events` - List risk events
- GET `/alerts` - List alerts

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/models
go test ./pkg/validation
go test ./pkg/errors

# Frontend tests
cd frontend
npm test
```

### Test Coverage Areas
- Model validations
- JWT token generation/validation
- Error handling
- API endpoints
- React components

## 📊 Data Flow

### Telemetry Processing
1. Vehicle sends telemetry → Telemetry Ingest
2. Telemetry Ingest validates & stores → MySQL
3. Telemetry Ingest publishes event → Redis
4. Risk Engine subscribes to events → Analyzes risk
5. Risk Engine creates alerts → MySQL
6. WebSocket broadcasts updates → Connected clients

### Authentication Flow
1. User login → Auth Service
2. Auth Service validates credentials → MySQL
3. Auth Service generates JWT → Returns to client
4. Client includes JWT in requests → API Service
5. API Service validates JWT → Processes request

## 🐛 Common Issues & Solutions

### Issue: Services can't connect to MySQL
**Solution**: Ensure MySQL is running and credentials are correct
```bash
mysql -h localhost -u fleet -pdevpass fleet_dev
```

### Issue: JWT authentication fails
**Solution**: Ensure JWT_SECRET is same across all services
```bash
export JWT_SECRET="your-super-secret-jwt-key-change-in-production"
```

### Issue: Frontend can't reach backend
**Solution**: Check proxy configuration in `vite.config.ts`
```typescript
proxy: {
  '/api/v1/auth': { target: 'http://localhost:8084' },
  '/api': { target: 'http://localhost:8080' },
  '/graphql': { target: 'http://localhost:8080' },
}
```

### Issue: Docker build fails
**Solution**: Clear Docker cache and rebuild
```bash
docker system prune -a
docker-compose -f docker-compose.dev.yml build --no-cache
```

## 🔄 Git Workflow

### Branch Strategy
- `main` - Production ready code
- `develop` - Integration branch
- `feature/*` - New features
- `fix/*` - Bug fixes
- `hotfix/*` - Production fixes

### Commit Convention
```
type(scope): description

- feat: New feature
- fix: Bug fix
- docs: Documentation
- test: Testing
- refactor: Code refactoring
- style: Formatting
- chore: Maintenance
```

## 📈 Monitoring & Logging

### Logging
- All services use structured logging (logrus)
- Log levels: debug, info, warn, error, fatal
- Logs include correlation IDs for tracing

### Health Checks
- All services expose `/health` endpoint
- Docker health checks configured
- Response includes service status and metadata

## 🚢 Deployment

### Development
```bash
docker-compose -f docker-compose.dev.yml up
```

### Production Considerations
1. Use environment-specific configs
2. Enable TLS/SSL
3. Set strong JWT secret
4. Configure database backups
5. Set up monitoring (Prometheus/Grafana)
6. Implement rate limiting
7. Configure CORS properly
8. Use production database passwords

## 📚 Additional Resources

### Key Files to Review
1. `pkg/models/models.go` - Data models
2. `pkg/auth/jwt.go` - Authentication logic
3. `services/api/graph/schema.graphqls` - GraphQL schema
4. `frontend/src/App.tsx` - Frontend routing
5. `docker-compose.dev.yml` - Service orchestration

### Dependencies
- Backend: gin, gorm, gqlgen, jwt, logrus, redis
- Frontend: react, apollo, tailwind, vite, typescript

### Next Steps
1. Add comprehensive logging
2. Implement metrics collection
3. Add integration tests
4. Set up staging environment
5. Configure production deployment
6. Add API documentation (Swagger/OpenAPI)
7. Implement data retention policies
8. Add backup strategies

---

## Quick Commands Reference

```bash
# Start everything
docker-compose -f docker-compose.dev.yml up

# Rebuild specific service
docker-compose -f docker-compose.dev.yml up --build api

# View logs
docker-compose logs -f api

# Access MySQL
docker exec -it samsara_healthcare_mysql_1 mysql -u fleet -pdevpass fleet_dev

# Access Redis
docker exec -it samsara_healthcare_redis_1 redis-cli

# Run Go tests
go test ./...

# Run frontend dev
cd frontend && npm run dev

# Build frontend
cd frontend && npm run build

# Format Go code
go fmt ./...

# Lint Go code
golangci-lint run

# Update Go dependencies
go mod tidy
```