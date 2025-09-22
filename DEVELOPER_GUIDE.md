# Developer Guide - Fleet Risk Intelligence System

## 🚀 Quick Start

### Prerequisites Checklist
```bash
# Check if you have required tools
go version          # Should be 1.24+
node --version      # Should be 20+
npm --version       # Should be 10+
docker --version    # Should be 24+
mysql --version     # Should be 8.0+
redis-cli --version # Should be 7.0+
```

### 5-Minute Setup
```bash
# 1. Clone and setup
git clone <repository-url>
cd samsara_healthcare

# 2. Start infrastructure
docker-compose -f docker-compose.dev.yml up mysql redis -d

# 3. Install dependencies
go mod download
cd frontend && npm install && cd ..

# 4. Run services (in separate terminals)
cd services/auth && go run main.go        # Terminal 1
cd services/api && go run main.go         # Terminal 2
cd frontend && npm run dev                 # Terminal 3

# 5. Open browser
open http://localhost:3000
```

## 🎯 Demo Mode (Recommended for First Look)

The system includes a demo mode for immediate exploration without full backend setup:

### Quick Demo Setup
```bash
# 1. Start core services only
docker compose -f docker-compose.dev.yml up mysql redis auth frontend -d

# 2. Open dashboard
open http://localhost:3000
```

### What You'll See in Demo Mode
- **Professional Dashboard**: Complete fleet management interface
- **Sample Data**: 3 fleets, 5 vehicles, 5 drivers with realistic metrics
- **Risk Analytics**: Color-coded risk scores and safety indicators
- **Full Navigation**: Access all pages (Dashboard, Vehicles, Drivers, Fleets, etc.)
- **Modern UI**: Tailwind CSS styling with responsive design

### Demo Features
- **Bypass Authentication**: No login required for exploration
- **Fallback Data**: Shows demo data when API services aren't available
- **All Pages Functional**: Navigate between all sections without backend dependency
- **Realistic Metrics**: Proper risk scores, status indicators, company data

### Switching Back to Full Mode
To restore full authentication and backend integration:
```typescript
// In frontend/src/contexts/AuthContext.tsx, revert the initialState:
const initialState: AuthState = {
  user: null,
  token: localStorage.getItem('token'),
  isAuthenticated: false,
  isLoading: true,
}
```

## 🛠️ Development Workflow

### Day-to-Day Development

#### Starting Work Session
```bash
# 1. Pull latest changes
git pull origin main

# 2. Start infrastructure
docker-compose -f docker-compose.dev.yml up mysql redis -d

# 3. Start your service
cd services/[service-name]
go run main.go

# 4. Start frontend (if needed)
cd frontend
npm run dev
```

#### Making Changes

**Backend Changes:**
```bash
# 1. Make your changes
# 2. Run tests
go test ./...

# 3. Check linting
go fmt ./...

# 4. Restart service to see changes
# (Go services don't have hot reload)
```

**Frontend Changes:**
```bash
# 1. Make your changes
# 2. Vite will auto-reload
# 3. Run tests if needed
npm test

# 4. Type check
npm run type-check
```

#### Database Changes
```bash
# 1. Update models in pkg/models/models.go
# 2. Restart any service (auto-migration will run)
# 3. Verify changes
docker exec -it samsara_healthcare_mysql_1 mysql -u fleet -pdevpass fleet_dev
SHOW TABLES;
DESCRIBE [table_name];
```

### Adding New Features

#### 1. Backend Service Changes

**Adding New API Endpoint:**
```go
// 1. Add to services/api/main.go
api.GET("/new-endpoint", getNewData(server))

// 2. Implement handler
func getNewData(server *server.BaseServer) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Implementation
        c.JSON(http.StatusOK, data)
    }
}
```

**Adding New GraphQL Field:**
```go
// 1. Update services/api/graph/schema.graphqls
type Query {
    newField(input: String!): String!
}

// 2. Regenerate GraphQL code
cd services/api
go generate ./...

// 3. Implement resolver in schema.resolvers.go
func (r *queryResolver) NewField(ctx context.Context, input string) (string, error) {
    // Implementation
    return result, nil
}
```

**Adding New Model:**
```go
// 1. Add to pkg/models/models.go
type NewModel struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// 2. Update Migrate function
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        // ... existing models
        &NewModel{},
    )
}
```

#### 2. Frontend Changes

**Adding New Page:**
```typescript
// 1. Create component in src/pages/NewPage.tsx
export default function NewPage() {
    return <div>New Page</div>
}

// 2. Add route in src/App.tsx
<Route path="/new-page" element={<NewPage />} />

// 3. Add navigation in src/components/Layout.tsx
const navigation = [
    // ... existing items
    { name: 'New Page', href: '/new-page', icon: SomeIcon },
]
```

**Adding GraphQL Query:**
```typescript
// 1. Add to src/graphql/queries.ts
export const GET_NEW_DATA = gql`
  query GetNewData($input: String!) {
    newField(input: $input)
  }
`

// 2. Use in component
const { data, loading, error } = useQuery(GET_NEW_DATA, {
    variables: { input: "value" }
})
```

## 🧪 Testing Guide

### Running Tests

**All Tests:**
```bash
# Backend tests
go test ./...

# Frontend tests
cd frontend
npm test

# With coverage
go test -cover ./...
npm run test:coverage
```

**Specific Tests:**
```bash
# Test specific package
go test ./pkg/models
go test ./pkg/auth

# Test specific file
go test -run TestUserModel ./pkg/models

# Frontend component tests
npm test -- --testNamePattern="Dashboard"
```

### Writing Tests

**Backend Unit Test Example:**
```go
// pkg/models/models_test.go
func TestUserModel(t *testing.T) {
    // Setup
    db := setupTestDB(t)

    // Test data
    user := User{
        Email:     "test@example.com",
        Password:  "password123",
        FirstName: "John",
        LastName:  "Doe",
    }

    // Execute
    err := db.Create(&user).Error

    // Assert
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    assert.NotEqual(t, "password123", user.Password) // Should be hashed
}
```

**Frontend Component Test Example:**
```typescript
// src/components/Dashboard.test.tsx
import { render, screen } from '@testing-library/react'
import { MockedProvider } from '@apollo/client/testing'
import Dashboard from './Dashboard'

test('renders dashboard with data', () => {
    const mocks = [
        {
            request: { query: GET_FLEETS },
            result: { data: { fleets: [] } }
        }
    ]

    render(
        <MockedProvider mocks={mocks}>
            <Dashboard />
        </MockedProvider>
    )

    expect(screen.getByText('Dashboard')).toBeInTheDocument()
})
```

## 🐛 Debugging Guide

### Common Issues

#### "Service won't start"
```bash
# Check if port is in use
lsof -i :8080

# Check environment variables
env | grep -E "(DB_|REDIS_|JWT_)"

# Check database connection
mysql -h localhost -u fleet -pdevpass fleet_dev

# Check logs
docker-compose logs mysql
docker-compose logs redis
```

#### "Authentication not working"
```bash
# Verify JWT secret is consistent
echo $JWT_SECRET

# Check token in browser
# Open DevTools → Application → Local Storage → token

# Test auth endpoint
curl -X POST http://localhost:8084/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
```

#### "Frontend can't connect to backend"
```bash
# Check Vite proxy config
cat frontend/vite.config.ts

# Verify services are running
curl http://localhost:8080/health
curl http://localhost:8084/health

# Check network requests in DevTools
```

#### "Database migration issues"
```bash
# Connect to database
docker exec -it samsara_healthcare_mysql_1 mysql -u fleet -pdevpass fleet_dev

# Check tables
SHOW TABLES;

# Drop and recreate (CAUTION: loses data)
DROP DATABASE fleet_dev;
CREATE DATABASE fleet_dev;

# Restart service to re-migrate
```

### Logging and Monitoring

#### Service Logs
```bash
# View live logs
docker-compose logs -f api
docker-compose logs -f auth

# Search logs
docker-compose logs api | grep ERROR
docker-compose logs | grep "correlation_id=123"
```

#### Database Queries
```bash
# Enable query logging (development only)
# Add to mysql config: general_log = 1

# View slow queries
docker exec -it samsara_healthcare_mysql_1 mysql -u root -p
SHOW PROCESSLIST;
```

#### Performance Monitoring
```bash
# Check Go service metrics
curl http://localhost:8080/debug/pprof/
curl http://localhost:8080/debug/vars

# Monitor resource usage
docker stats

# Database performance
docker exec -it samsara_healthcare_mysql_1 mysql -u root -p
SHOW STATUS LIKE 'Slow_queries';
SHOW STATUS LIKE 'Threads_connected';
```

## 📁 Project Structure Deep Dive

### Shared Packages (`pkg/`)

#### `pkg/auth/`
- `jwt.go`: JWT token management
- `middleware.go`: Authentication middleware
- Usage: Import in services that need auth

#### `pkg/config/`
- `config.go`: Configuration management
- Loads from environment variables
- Usage: `config := config.Load()`

#### `pkg/database/`
- `database.go`: Database connection setup
- GORM configuration
- Usage: Shared across all services

#### `pkg/models/`
- `models.go`: All database models
- GORM hooks for password hashing
- Migration function

#### `pkg/server/`
- `server.go`: Base server setup
- Common middleware and configuration
- Usage: Foundation for all HTTP services

### Service Structure

Each service follows this pattern:
```
services/[service-name]/
├── main.go              # Service entry point
├── handlers/            # HTTP handlers (optional)
├── config/              # Service-specific config (optional)
└── [service-name].go    # Core business logic (optional)
```

### Frontend Structure

```
frontend/src/
├── components/          # Reusable UI components
│   ├── Layout.tsx      # Main layout wrapper
│   ├── ProtectedRoute.tsx # Route protection
│   └── [Component].tsx # Individual components
├── contexts/           # React contexts
│   └── AuthContext.tsx # Authentication state
├── graphql/           # GraphQL operations
│   ├── queries.ts     # GraphQL queries
│   ├── mutations.ts   # GraphQL mutations
│   └── subscriptions.ts # GraphQL subscriptions
├── pages/             # Route pages
│   ├── Dashboard.tsx  # Dashboard page
│   ├── Login.tsx      # Login page
│   └── [Page].tsx     # Other pages
├── types/             # TypeScript type definitions
└── main.tsx           # App entry point
```

## ⚡ Performance Tips

### Backend Optimization

**Database Performance:**
```go
// Use indexes
type User struct {
    Email string `gorm:"index"`
}

// Preload relationships
var users []User
db.Preload("Fleet").Find(&users)

// Limit results
db.Limit(100).Find(&users)

// Use specific fields
db.Select("id", "name").Find(&users)
```

**Memory Management:**
```go
// Close database connections
defer db.Close()

// Use context timeouts
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

### Frontend Optimization

**React Performance:**
```typescript
// Use React.memo for expensive components
const ExpensiveComponent = React.memo(({ data }) => {
    return <div>{data}</div>
})

// Use useMemo for expensive calculations
const expensiveValue = useMemo(() => {
    return heavyCalculation(data)
}, [data])

// Lazy load components
const LazyComponent = lazy(() => import('./LazyComponent'))
```

**GraphQL Performance:**
```typescript
// Use specific fields only
const GET_USERS = gql`
  query GetUsers {
    users {
      id
      name
      # Don't fetch unnecessary fields
    }
  }
`

// Use pagination
const GET_PAGINATED = gql`
  query GetPaginated($limit: Int!, $offset: Int!) {
    users(limit: $limit, offset: $offset) {
      id
      name
    }
  }
`
```

## 🔧 Configuration Guide

### Environment Variables

**Required for Development:**
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
JWT_SECRET=your-secret-key-here

# Server
API_PORT=8080
LOG_LEVEL=debug
ENV=development
```

**Production Additions:**
```env
# Security
JWT_SECRET=super-secure-production-secret-256-bits-long

# Database
DB_HOST=production-db-host
DB_PASSWORD=secure-production-password

# Performance
LOG_LEVEL=info
ENV=production

# Features
ENABLE_RATE_LIMITING=true
RATE_LIMIT_REQUESTS_PER_MINUTE=1000
```

### Service Configuration

**Per-Service Ports:**
```env
# Auth Service
AUTH_PORT=8084

# API Service
API_PORT=8080

# Telemetry Service
TELEMETRY_PORT=8081

# Risk Engine
RISK_ENGINE_PORT=8082

# WebSocket Service
WEBSOCKET_PORT=8083
```

## 📝 Code Standards

### Go Code Style
```go
// Package naming: lowercase, no underscores
package auth

// Function naming: CamelCase for exported, camelCase for internal
func GenerateToken() string        // Exported
func validatePassword() bool       // Internal

// Error handling: explicit error checking
result, err := someFunction()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Struct tags: consistent formatting
type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Email string `json:"email" gorm:"uniqueIndex;size:255"`
}
```

### TypeScript Code Style
```typescript
// Interface naming: PascalCase with descriptive names
interface UserProfile {
    id: string
    email: string
}

// Component naming: PascalCase
export default function UserDashboard() {
    return <div>Dashboard</div>
}

// Hook naming: use prefix
function useUserData() {
    // Hook implementation
}

// File naming: PascalCase for components, camelCase for utilities
// UserDashboard.tsx (component)
// userUtils.ts (utility)
```

## 🚀 Deployment Checklist

### Pre-deployment Verification
```bash
# 1. Run all tests
go test ./...
cd frontend && npm test

# 2. Build frontend
cd frontend && npm run build

# 3. Check Docker builds
docker-compose -f docker-compose.dev.yml build

# 4. Verify environment variables
env | grep -E "(DB_|REDIS_|JWT_)"

# 5. Run security scan
# Add security scanning tool

# 6. Check database migrations
# Ensure all models are in Migrate function
```

### Production Deployment
```bash
# 1. Set production environment variables
export ENV=production
export JWT_SECRET=production-secret
export DB_PASSWORD=production-password

# 2. Build production images
docker build -f docker/api.Dockerfile -t fleet-api:latest .
docker build -f docker/auth.Dockerfile -t fleet-auth:latest .

# 3. Deploy to production environment
# (Kubernetes, AWS ECS, etc.)

# 4. Run health checks
curl https://production-api/health

# 5. Monitor logs
kubectl logs -f deployment/fleet-api
```

---

## 🆘 Getting Help

### Internal Resources
1. **IMPLEMENTATION.md** - Architecture and setup details
2. **SPECIFICATION.md** - Technical specifications
3. **README.md** - Project overview

### External Resources
1. **Go Documentation**: https://golang.org/doc/
2. **React Documentation**: https://react.dev/
3. **GORM Documentation**: https://gorm.io/docs/
4. **Apollo Client**: https://www.apollographql.com/docs/react/

### Troubleshooting Steps
1. Check this guide for common issues
2. Review service logs
3. Verify environment variables
4. Test individual components
5. Check database/Redis connectivity
6. Validate authentication flow

Remember: When in doubt, start simple and build complexity gradually!