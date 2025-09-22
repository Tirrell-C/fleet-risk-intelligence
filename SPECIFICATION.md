# Fleet Risk Intelligence System - Technical Specification

## 📋 Executive Summary

The Fleet Risk Intelligence System is a comprehensive platform for monitoring vehicle fleets, analyzing driver behavior, and providing real-time risk assessments. Built as a microservices architecture using Go and React, it provides scalable, secure, and real-time fleet management capabilities.

## 🎯 Business Requirements

### Primary Objectives
1. **Real-time Fleet Monitoring**: Track vehicle locations, status, and performance metrics
2. **Risk Assessment**: Analyze driving behavior and identify safety risks
3. **Alert Management**: Generate and manage safety alerts and notifications
4. **Driver Performance**: Monitor and score driver behavior and safety
5. **Fleet Optimization**: Provide insights for operational efficiency

### User Roles
- **Super Admin**: System-wide access and configuration
- **Fleet Admin**: Full access to assigned fleets
- **Fleet Manager**: Operational access to fleet data
- **Driver**: Limited access to personal performance data

## 🏗️ System Architecture

### Architecture Style
- **Pattern**: Microservices with API Gateway
- **Communication**: HTTP/GraphQL for sync, Redis pub/sub for async
- **Data Storage**: MySQL for persistent data, Redis for caching
- **Authentication**: JWT-based stateless authentication
- **Real-time**: WebSocket connections for live updates

### Service Boundaries

#### 1. Authentication Service (`services/auth/`)
**Responsibility**: User authentication, authorization, session management
- User registration and login
- JWT token generation and validation
- Role-based access control
- Session tracking and logout

#### 2. API Gateway Service (`services/api/`)
**Responsibility**: External API interface, request routing, data aggregation
- GraphQL endpoint for complex queries
- REST endpoints for simple operations
- Request authentication and authorization
- Data aggregation from multiple services

#### 3. Telemetry Ingest Service (`services/telemetry-ingest/`)
**Responsibility**: Vehicle data collection and validation
- Real-time telemetry data ingestion
- Data validation and sanitization
- Event publishing to message queue
- Batch processing capabilities

#### 4. Risk Engine Service (`services/risk-engine/`)
**Responsibility**: Risk analysis and alert generation
- Real-time risk event analysis
- Driver behavior scoring
- Alert generation and notification
- Risk pattern detection

#### 5. WebSocket Service (`services/websocket/`)
**Responsibility**: Real-time client notifications
- WebSocket connection management
- Real-time event broadcasting
- Client session management
- Message filtering by permissions

## 🗄️ Data Models

### Core Entities

#### User
```go
type User struct {
    ID        uint      `json:"id"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`        // Hidden, bcrypt hashed
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Role      string    `json:"role"`     // super_admin, fleet_admin, fleet_manager, driver
    Status    string    `json:"status"`   // active, inactive, suspended
    FleetIDs  string    `json:"-"`        // JSON array of accessible fleet IDs
    LastLogin *time.Time `json:"last_login"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### Fleet
```go
type Fleet struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`
    CompanyName string    `json:"company_name"`
    ContactEmail string   `json:"contact_email"`
    Status      string    `json:"status"`    // active, inactive
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### Vehicle
```go
type Vehicle struct {
    ID          uint      `json:"id"`
    VIN         string    `json:"vin"`         // Unique vehicle identifier
    Make        string    `json:"make"`
    Model       string    `json:"model"`
    Year        int       `json:"year"`
    LicensePlate string   `json:"license_plate"`
    FleetID     uint      `json:"fleet_id"`
    Fleet       Fleet     `json:"fleet"`
    DriverID    *uint     `json:"driver_id"`   // Current assigned driver
    Driver      *Driver   `json:"driver,omitempty"`
    Status      string    `json:"status"`      // active, maintenance, inactive
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### Driver
```go
type Driver struct {
    ID          uint      `json:"id"`
    EmployeeID  string    `json:"employee_id"`
    FirstName   string    `json:"first_name"`
    LastName    string    `json:"last_name"`
    Email       string    `json:"email"`
    Phone       string    `json:"phone"`
    LicenseNum  string    `json:"license_number"`
    FleetID     uint      `json:"fleet_id"`
    Fleet       Fleet     `json:"fleet"`
    Status      string    `json:"status"`      // active, suspended, inactive
    RiskScore   float64   `json:"risk_score"`  // 0-10 calculated risk score
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### TelemetryEvent
```go
type TelemetryEvent struct {
    ID          uint      `json:"id"`
    VehicleID   uint      `json:"vehicle_id"`
    Vehicle     Vehicle   `json:"vehicle"`
    EventType   string    `json:"event_type"`    // location, speed, acceleration
    Timestamp   time.Time `json:"timestamp"`
    Latitude    *float64  `json:"latitude"`
    Longitude   *float64  `json:"longitude"`
    Speed       *float64  `json:"speed"`         // mph
    Acceleration *float64 `json:"acceleration"`  // m/s²
    Data        string    `json:"data"`          // JSON for additional data
    ProcessedAt *time.Time `json:"processed_at"`
    CreatedAt   time.Time `json:"created_at"`
}
```

#### RiskEvent
```go
type RiskEvent struct {
    ID          uint      `json:"id"`
    VehicleID   uint      `json:"vehicle_id"`
    Vehicle     Vehicle   `json:"vehicle"`
    DriverID    *uint     `json:"driver_id"`
    Driver      *Driver   `json:"driver,omitempty"`
    EventType   string    `json:"event_type"`    // speeding, harsh_braking, etc.
    Severity    string    `json:"severity"`      // low, medium, high, critical
    RiskScore   float64   `json:"risk_score"`    // 0-100
    Timestamp   time.Time `json:"timestamp"`
    Latitude    *float64  `json:"latitude"`
    Longitude   *float64  `json:"longitude"`
    Description string    `json:"description"`
    Data        string    `json:"data"`          // JSON for event details
    Status      string    `json:"status"`        // open, acknowledged, resolved
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Data Relationships
```
Fleet (1) → (N) Vehicles
Fleet (1) → (N) Drivers
Fleet (1) → (N) Alerts
Vehicle (1) → (N) TelemetryEvents
Vehicle (1) → (N) RiskEvents
Driver (1) → (N) RiskEvents
Driver (1) → (1) DriverScore
```

## 🔐 Security Specification

### Authentication
- **Method**: JWT (JSON Web Tokens)
- **Expiration**: 24 hours (configurable)
- **Storage**: Client-side localStorage (frontend), memory (backend)
- **Refresh**: Manual re-login (can be extended to refresh tokens)

### Authorization
- **Model**: Role-Based Access Control (RBAC) with resource-level permissions
- **Roles**:
  - `super_admin`: Full system access
  - `fleet_admin`: Full access to assigned fleets
  - `fleet_manager`: Operational access to assigned fleets
  - `driver`: Read-only access to personal data

### Password Security
- **Hashing**: bcrypt with default cost (10)
- **Requirements**: Minimum 8 characters (can be enhanced)
- **Storage**: Hashed in database, never transmitted in responses

### API Security
- **Authentication**: Bearer token in Authorization header
- **CORS**: Configured for development, should be restricted in production
- **Rate Limiting**: Not implemented (recommend adding for production)
- **Input Validation**: GORM validation and custom validators

## 📡 API Specification

### GraphQL Schema

#### Types
```graphql
type Fleet {
  id: ID!
  name: String!
  companyName: String!
  contactEmail: String!
  status: String!
  vehicles: [Vehicle!]!
  drivers: [Driver!]!
  createdAt: String!
  updatedAt: String!
}

type Vehicle {
  id: ID!
  vin: String!
  make: String!
  model: String!
  year: Int!
  licensePlate: String!
  fleetId: ID!
  fleet: Fleet!
  driverId: ID
  driver: Driver
  status: VehicleStatus!
  riskScore: Float!
  currentLocation: Location
  lastTelemetry: TelemetryEvent
  createdAt: String!
  updatedAt: String!
}

type Driver {
  id: ID!
  employeeId: String!
  firstName: String!
  lastName: String!
  email: String!
  phone: String!
  licenseNumber: String!
  fleetId: ID!
  fleet: Fleet!
  status: DriverStatus!
  riskScore: Float!
  currentVehicle: Vehicle
  driverScore: DriverScore
  createdAt: String!
  updatedAt: String!
}
```

#### Queries
```graphql
type Query {
  # Fleet operations
  fleets: [Fleet!]!
  fleet(id: ID!): Fleet

  # Vehicle operations
  vehicles(fleetId: ID): [Vehicle!]!
  vehicle(id: ID!): Vehicle
  liveVehicleData(vehicleId: ID!): VehicleData

  # Driver operations
  drivers(fleetId: ID): [Driver!]!
  driver(id: ID!): Driver
  driverScores(fleetId: ID!): [DriverScore!]!

  # Risk & Alert operations
  riskEvents(vehicleId: ID, driverId: ID, limit: Int): [RiskEvent!]!
  alerts(fleetId: ID!, status: AlertStatus): [Alert!]!
}
```

#### Mutations
```graphql
type Mutation {
  # Fleet management
  createFleet(input: CreateFleetInput!): Fleet!
  updateFleet(id: ID!, input: UpdateFleetInput!): Fleet!

  # Vehicle management
  createVehicle(input: CreateVehicleInput!): Vehicle!
  updateVehicle(id: ID!, input: UpdateVehicleInput!): Vehicle!
  assignDriver(vehicleId: ID!, driverId: ID!): Vehicle!

  # Driver management
  createDriver(input: CreateDriverInput!): Driver!
  updateDriver(id: ID!, input: UpdateDriverInput!): Driver!

  # Alert management
  acknowledgeAlert(id: ID!): Alert!
  dismissAlert(id: ID!): Alert!
}
```

#### Subscriptions
```graphql
type Subscription {
  # Real-time updates
  vehicleUpdates(vehicleId: ID!): VehicleData!
  riskEventNotifications(fleetId: ID!): RiskEvent!
  alertNotifications(fleetId: ID!): Alert!
}
```

### REST API Endpoints

#### Authentication (`/api/v1/auth/`)
```
POST /login
POST /register
GET  /me (protected)
PUT  /me (protected)
POST /logout (protected)

# Admin endpoints (admin roles only)
GET    /admin/users
POST   /admin/users
PUT    /admin/users/:id
DELETE /admin/users/:id
```

#### Fleet Management (`/api/v1/`)
```
GET /fleets
GET /fleets/:id
GET /vehicles
GET /vehicles/:id
GET /drivers
GET /drivers/:id
GET /risk-events
GET /vehicles/:id/risk-events
GET /alerts
```

## 🔄 Event Flow Specification

### Telemetry Processing Flow
```
1. Vehicle → Telemetry Data → Telemetry Ingest Service
2. Telemetry Ingest → Validate & Store → MySQL
3. Telemetry Ingest → Publish Event → Redis (channel: telemetry_events)
4. Risk Engine → Subscribe → Redis
5. Risk Engine → Analyze Risk → Generate Risk Events/Alerts
6. Risk Engine → Store Results → MySQL
7. Risk Engine → Publish → Redis (channels: risk_events, alerts)
8. WebSocket Service → Subscribe → Redis
9. WebSocket Service → Broadcast → Connected Clients
```

### Authentication Flow
```
1. Client → Login Request → Auth Service
2. Auth Service → Validate Credentials → MySQL
3. Auth Service → Generate JWT → Return to Client
4. Client → API Request + JWT → API Service
5. API Service → Validate JWT → Process Request
6. API Service → Return Response → Client
```

### Real-time Notification Flow
```
1. Risk Engine → Detect Risk Event → Generate Alert
2. Risk Engine → Publish → Redis (channel: alerts)
3. WebSocket Service → Receive → Filter by User Permissions
4. WebSocket Service → Broadcast → Relevant Connected Clients
5. Frontend → Receive → Update UI
```

## 📊 Performance Specifications

### Response Time Requirements
- **Authentication**: < 200ms
- **API Queries**: < 500ms
- **GraphQL Queries**: < 1s
- **Real-time Updates**: < 100ms
- **Dashboard Load**: < 2s

### Throughput Requirements
- **Telemetry Ingestion**: 1000 events/second per vehicle
- **Concurrent Users**: 500 simultaneous WebSocket connections
- **API Requests**: 10,000 requests/minute
- **Database Queries**: < 50ms average

### Scalability Targets
- **Fleets**: 1000+ organizations
- **Vehicles**: 10,000+ vehicles per system
- **Users**: 5,000+ concurrent users
- **Data Retention**: 2 years of telemetry data

## 🛠️ Technology Specifications

### Backend Stack
- **Language**: Go 1.24+
- **Framework**: Gin (HTTP), gqlgen (GraphQL)
- **Database**: MySQL 8.0
- **Cache/Queue**: Redis 7+
- **Authentication**: JWT with golang-jwt/jwt/v5
- **Password Hashing**: bcrypt
- **Logging**: logrus with structured logging

### Frontend Stack
- **Framework**: React 18
- **Language**: TypeScript 5.2+
- **Build Tool**: Vite 4.5+
- **Styling**: Tailwind CSS 3.3+
- **HTTP Client**: Apollo Client 3.8+
- **State Management**: React Context + Hooks
- **Testing**: Vitest + React Testing Library

### Infrastructure
- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Docker Compose (dev), Kubernetes (prod)
- **Web Server**: Nginx (production reverse proxy)
- **CI/CD**: GitHub Actions
- **Monitoring**: Health checks + structured logging

## 🧪 Testing Strategy

### Unit Testing
- **Backend**: Go testing framework with testify
- **Frontend**: Vitest with React Testing Library
- **Coverage Target**: > 80%
- **Mock Strategy**: Interfaces for external dependencies

### Integration Testing
- **Database**: SQLite in-memory for tests
- **API**: HTTP test requests
- **Authentication**: JWT token validation
- **Real-time**: WebSocket connection testing

### End-to-End Testing
- **User Flows**: Complete authentication → data access → real-time updates
- **Cross-service**: Full telemetry processing pipeline
- **Performance**: Load testing with realistic data volumes

## 🚀 Deployment Specification

### Environment Requirements

#### Development
```yaml
Resources:
  - CPU: 4 cores
  - RAM: 8GB
  - Storage: 20GB SSD
Services:
  - MySQL: 1 instance
  - Redis: 1 instance
  - All microservices: 1 instance each
```

#### Production
```yaml
Resources:
  - CPU: 16+ cores
  - RAM: 32+ GB
  - Storage: 100+ GB SSD
Load Balancing:
  - Frontend: CDN + multiple instances
  - API Gateway: Load balanced
  - Services: Auto-scaling based on metrics
Database:
  - MySQL: Master-slave replication
  - Redis: Cluster mode
```

### Configuration Management
- **Environment Variables**: All configuration externalized
- **Secrets**: Kubernetes secrets or equivalent
- **Feature Flags**: Environment-based feature toggling
- **Monitoring**: Prometheus + Grafana (recommended)

## 📋 Compliance & Governance

### Data Privacy
- **Personal Data**: Driver information requires consent
- **Data Retention**: Configurable retention policies
- **Data Access**: Audit logs for all data access
- **Data Export**: User data export capabilities

### Security Standards
- **Authentication**: Industry-standard JWT implementation
- **Encryption**: TLS 1.3 for data in transit
- **Password Policy**: Configurable complexity requirements
- **Session Management**: Secure session handling

### Monitoring & Observability
- **Logging**: Structured JSON logs with correlation IDs
- **Metrics**: Service-level metrics collection
- **Health Checks**: Comprehensive health monitoring
- **Alerting**: Critical error notification system

## 🔮 Future Enhancements

### Phase 2 Features
1. **Advanced Analytics**: ML-based risk prediction
2. **Mobile Applications**: Native iOS/Android apps
3. **IoT Integration**: Direct vehicle hardware integration
4. **Reporting**: Advanced analytics and reporting dashboard
5. **Third-party Integrations**: Fleet management system APIs

### Technical Improvements
1. **Caching**: Redis-based response caching
2. **Rate Limiting**: API rate limiting implementation
3. **Message Queue**: Separate message queue for high-volume events
4. **Database Sharding**: Horizontal database scaling
5. **Microservice Mesh**: Service mesh for advanced networking

---

This specification serves as the definitive technical reference for the Fleet Risk Intelligence System. All development decisions should align with these specifications, and any deviations should be documented and approved through the change management process.