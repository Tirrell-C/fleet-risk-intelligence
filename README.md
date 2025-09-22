# Fleet Risk Intelligence MVP

A production-ready MVP showcasing AI-driven fleet risk detection and management, built with Go microservices, React/TypeScript frontend, and modern containerized architecture.

> **🎯 Current Status**: Core implementation complete with working dashboard, authentication service, and demo mode enabled for immediate exploration.

## 🎯 Project Overview

This project demonstrates enterprise-level software engineering practices for IoT fleet management and risk intelligence:

- **Real-time telemetry processing** for vehicle safety monitoring
- **AI-powered risk detection** with pattern analysis and scoring
- **Scalable microservices architecture** using Go and GraphQL
- **Modern containerized deployment** with Docker Compose
- **Professional React/TypeScript frontend** with real-time capabilities

## 🚦 Implementation Status

### ✅ Fully Working
- **React Dashboard**: Complete fleet management UI with demo data
- **Authentication Service**: JWT-based auth with user management (Port 8084)
- **Database Layer**: MySQL with all tables and relationships
- **Risk Engine**: Go microservice for risk analysis (Port 8082)
- **Infrastructure**: Docker containerization with health checks

### 🟡 Partially Implemented
- **GraphQL API**: Schema complete, version compatibility issue (Port 8080)
- **Telemetry Ingest**: Built, minor constraint naming conflict (Port 8081)
- **WebSocket Service**: Ready for real-time updates (Port 8083)

### 🎯 Demo Features
- **3 Sample Fleets**: Metro Logistics, Cross Country Transport, City Services Inc
- **5 Demo Vehicles**: Ford Transit, Chevy Express, Mercedes Sprinter, etc.
- **5 Demo Drivers**: Complete with risk scores and performance metrics
- **Risk Analytics**: Color-coded dashboards with safety indicators
- **Professional UI**: Modern design with full navigation between all sections

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │  Risk Engine    │
│ React/TypeScript│────│   GraphQL/Go    │────│      Go         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                       ┌─────────────────┐    ┌─────────────────┐
                       │  Telemetry      │    │   WebSocket     │
                       │  Ingest (Go)    │    │   Service (Go)  │
                       └─────────────────┘    └─────────────────┘
                                │                       │
                       ┌─────────────────┐    ┌─────────────────┐
                       │     MySQL       │    │     Redis       │
                       │   Database      │    │     Cache       │
                       └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start (Demo Mode)

Experience the full fleet management dashboard immediately:

```bash
# 1. Clone the repository
git clone https://github.com/Tirrell-C/fleet-risk-intelligence.git
cd fleet-risk-intelligence

# 2. Start core services
docker compose -f docker-compose.dev.yml up mysql redis auth frontend -d

# 3. Open the dashboard
open http://localhost:3000
```

**What you'll see**: Professional fleet management interface with demo data including 3 fleets, 5 vehicles, 5 drivers, risk analytics, and full navigation.

## 🛠️ Full Development Setup

For complete backend development with all microservices:

```bash
# Start all services
docker compose -f docker-compose.dev.yml up --build

# Services available:
# - Frontend: http://localhost:3000
# - Auth API: http://localhost:8084
# - Main API: http://localhost:8080 (partial)
# - Risk Engine: http://localhost:8082
```

## 🛠️ Tech Stack

### Backend Services (Go)
- **GraphQL API**: Unified data access with `gqlgen`
- **Telemetry Ingest**: High-throughput data processing
- **Risk Engine**: AI-powered safety scoring and alerts
- **WebSocket Service**: Real-time dashboard updates

### Frontend (React/TypeScript)
- **Vite** for fast development and building
- **Apollo Client** for GraphQL integration
- **Tailwind CSS** for styling
- **Recharts** for data visualization

### Infrastructure (Pulumi + AWS)
- **ECS Fargate** for containerized services
- **RDS MySQL** for persistent storage
- **ElastiCache Redis** for real-time data
- **Application Load Balancer** for traffic routing

## 📋 Available Commands

```bash
# Setup and development
make setup          # Install dependencies
make dev            # Start development environment
make dev-down       # Stop development environment

# Building and testing
make build          # Build all services
make test           # Run all tests
make lint           # Run linters
make fmt            # Format code

# Docker and deployment
make docker-build   # Build Docker images
make deploy-infra   # Deploy infrastructure
make deploy         # Full deployment

# Code generation
make generate       # Generate GraphQL types
```

## 🔧 Development Workflow

### 1. Environment Setup
```bash
# Copy environment variables
cp .env.example .env

# Start infrastructure
make dev

# Verify services are running
docker ps
```

### 2. Code Development
```bash
# Make changes to Go services or frontend
# Tests run automatically on save in devcontainer

# Generate GraphQL schemas
make generate

# Run specific tests
go test ./services/api/...
cd frontend && npm test
```

### 3. Quality Assurance
```bash
# Full QA pipeline
make lint           # Code quality
make test           # All tests
make docker-build   # Container builds
```

## 🏭 Production Deployment

### Prerequisites
- AWS account with appropriate permissions
- Pulumi account and access token
- Docker registry access

### Deploy to AWS
```bash
# Configure AWS credentials
aws configure

# Deploy infrastructure
cd infrastructure
pulumi up

# Build and deploy services
make deploy
```

## 📊 Key Features Demonstrated

### Enterprise Architecture Patterns
- **Microservices** with domain separation
- **Event-driven architecture** with Redis pub/sub
- **API Gateway pattern** with GraphQL
- **CQRS** for read/write separation

### Scalability & Performance
- **Horizontal scaling** with ECS Fargate
- **Caching strategy** with Redis
- **Connection pooling** for database efficiency
- **Real-time processing** with WebSockets

### DevOps & SDLC
- **Infrastructure as Code** with Pulumi
- **CI/CD pipeline** with GitHub Actions
- **Automated testing** at all levels
- **Security scanning** with Trivy
- **Code quality gates** with golangci-lint

### AI/ML Integration
- **Risk scoring algorithms** for driver behavior
- **Pattern detection** for fleet safety
- **Real-time alerting** for critical events
- **Predictive analytics** for maintenance

## 🔒 Security Features

- **JWT authentication** for API access
- **Role-based access control** (RBAC)
- **Input validation** and sanitization
- **SQL injection prevention** with GORM
- **Container security** scanning
- **Secrets management** with AWS Systems Manager

## 📈 Monitoring & Observability

- **Structured logging** with logrus
- **Metrics collection** with Prometheus
- **Distributed tracing** capabilities
- **Health checks** for all services
- **Error tracking** and alerting

## 🧪 Testing Strategy

- **Unit tests** for business logic
- **Integration tests** for API endpoints
- **End-to-end tests** for critical workflows
- **Load testing** for performance validation
- **Security testing** in CI pipeline

## 📝 API Documentation

GraphQL schema and API documentation available at:
- **Local**: http://localhost:8080/graphql
- **Production**: https://your-domain.com/graphql

## 🤝 Contributing

This project follows enterprise development standards:

1. **Fork** the repository
2. **Create** a feature branch
3. **Commit** with conventional commit messages
4. **Test** thoroughly (tests must pass)
5. **Submit** a pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

Built with ❤️ to showcase production-ready Go microservices and modern web development practices.