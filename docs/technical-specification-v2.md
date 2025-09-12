# 📋 Gorev Technical Specification v2.0

**Version**: v2.0.0  
**Project**: Gorev Task Management System  
**Release**: v0.14.0  
**Date**: September 12, 2025  
**Status**: Production Ready  

---

## 📑 Table of Contents

1. [Executive Summary](#executive-summary)
2. [System Requirements](#system-requirements)
3. [Architecture Overview](#architecture-overview)
4. [Component Specifications](#component-specifications)
5. [Configuration Management](#configuration-management)
6. [Database Schema](#database-schema)
7. [API Specifications](#api-specifications)
8. [Security & Compliance](#security--compliance)
9. [Performance Requirements](#performance-requirements)
10. [Monitoring & Observability](#monitoring--observability)
11. [Deployment Architecture](#deployment-architecture)
12. [Quality Assurance](#quality-assurance)
13. [Migration & Compatibility](#migration--compatibility)

---

## 🎯 Executive Summary

Gorev v0.14.0 is an enterprise-grade task management system built with Go, featuring advanced AI integration, MCP (Model Context Protocol) support, and comprehensive VS Code extension capabilities. The system maintains 100% Rule 15 compliance with zero error/warning tolerance and follows strict DRY principles throughout the codebase.

### Key Features
- **Multi-Platform Support**: Windows, macOS, Linux
- **AI-Powered NLP**: Advanced natural language processing for task analysis
- **MCP Integration**: 29 standardized tools for AI assistants
- **VS Code Extension**: Native IDE integration
- **Real-time File Watching**: Automated task status updates
- **Bilingual Support**: Turkish and English interfaces
- **Enterprise Security**: 100% parameterized queries, zero SQL injection risk

---

## 🔧 System Requirements

### Minimum Requirements

#### Runtime Environment
```yaml
Go Runtime:
  version: ">= 1.21.0"
  architecture: ["amd64", "arm64"]
  
Operating System:
  windows: ">= Windows 10 (1809)"
  macOS: ">= 10.15 Catalina"
  linux: ">= Ubuntu 20.04 LTS, RHEL 8, or equivalent"

Resources:
  RAM: "256 MB minimum, 512 MB recommended"
  Disk: "100 MB for binary + 500 MB for data"
  CPU: "1 core minimum, 2+ cores recommended"
```

#### Database Requirements
```yaml
SQLite:
  version: ">= 3.35.0"
  features: ["foreign_keys", "json1", "fts5"]
  
PostgreSQL (Optional):
  version: ">= 13.0"
  extensions: ["uuid-ossp", "pg_trgm"]
```

### Recommended Production Environment

```yaml
Hardware:
  CPU: "4+ cores, 2.5+ GHz"
  RAM: "2 GB minimum, 8 GB recommended"
  Disk: "SSD with 10 GB available space"
  Network: "100 Mbps for cloud deployments"

Software:
  OS: "Linux (Ubuntu 22.04 LTS recommended)"
  Database: "PostgreSQL 15+ for production"
  Monitoring: "Prometheus + Grafana stack"
  Reverse Proxy: "Nginx or Traefik"
```

### Development Environment

```yaml
Required Tools:
  - Go 1.21+
  - Node.js 18+ (for VS Code extension)
  - Git 2.30+
  - Make or Task runner

IDE Support:
  - VS Code with Go extension
  - GoLand
  - Vim/Neovim with gopls

Testing Tools:
  - gotestsum for test execution
  - golangci-lint for static analysis
  - goreleaser for builds
```

---

## 🏗️ Architecture Overview

### High-Level Architecture

```mermaid
graph TB
    A[VS Code Extension] --> B[MCP Protocol Layer]
    B --> C[Core Task Engine]
    C --> D[Database Layer]
    C --> E[NLP Processor]
    C --> F[File Watcher]
    C --> G[Auto State Manager]
    
    H[Claude Desktop] --> B
    I[CLI Interface] --> C
    
    E --> J[AI Context Store]
    F --> K[File System Events]
    
    subgraph "Data Layer"
        D --> L[SQLite/PostgreSQL]
        J --> L
    end
```

### System Components Interaction

```mermaid
sequenceDiagram
    participant U as User
    participant V as VS Code Extension
    participant M as MCP Server
    participant C as Core Engine
    participant N as NLP Processor
    participant D as Database
    
    U->>V: Create task request
    V->>M: MCP tool call
    M->>C: Process request
    C->>N: Extract entities/intent
    N-->>C: Structured data
    C->>D: Store task
    D-->>C: Confirmation
    C-->>M: Success response
    M-->>V: Task created
    V-->>U: Visual feedback
```

---

## 🧩 Component Specifications

### 1. Core Task Engine (`core/`)

```go
// Task engine specifications
type TaskEngine struct {
    db           Database
    nlp          *NLPProcessor
    fileWatcher  *FileWatcher
    stateManager *AutoStateManager
    config       *Config
}

// Performance requirements
const (
    MaxConcurrentTasks = 1000
    TaskProcessingTime = 50 * time.Millisecond  // P95 target
    DatabaseTimeout    = 5 * time.Second
    MemoryLimit        = 512 * 1024 * 1024      // 512 MB
)
```

#### Responsibilities
- Task lifecycle management
- Business logic execution
- Resource coordination
- Error handling and recovery
- Performance monitoring

#### Error Handling Strategy
```go
// Rule 15 compliant error handling
type TaskError struct {
    Code      string                 `json:"code"`
    Message   string                 `json:"message"`
    Details   map[string]interface{} `json:"details"`
    Timestamp time.Time              `json:"timestamp"`
    Context   map[string]string      `json:"context"`
}

func (e *TaskError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}
```

### 2. MCP Protocol Server (`mcp/`)

```go
// MCP server specifications
type MCPServer struct {
    handlers     map[string]MCPHandler
    registry     *ToolRegistry
    rateLimiter  *RateLimiter
    validator    *InputValidator
    logger       *StructuredLogger
}

// Tool handler interface
type MCPHandler interface {
    Name() string
    Description() string
    Parameters() map[string]interface{}
    Execute(ctx context.Context, params map[string]interface{}) (*MCPResponse, error)
    Validate(params map[string]interface{}) error
}
```

#### Tool Registry Architecture
```yaml
Tool Categories:
  Task Management: ["gorev_listele", "gorev_detay", "gorev_guncelle", "gorev_duzenle", "gorev_sil"]
  Project Management: ["proje_olustur", "proje_listele", "proje_gorevleri"]
  AI Context: ["ai_context_ekle", "ai_context_listele", "ai_context_sil"]
  File Watching: ["file_watch_ekle", "file_watch_listele", "file_watch_sil"]
  Advanced Operations: ["batch_gorev_guncelle", "export_data", "import_data"]

Rate Limiting:
  per_client: 100 requests/minute
  burst_limit: 20 concurrent
  retry_policy: exponential backoff
```

### 3. NLP Processor (`nlp/`)

```go
// NLP processor specifications
type NLPProcessor struct {
    tokenizer    *Tokenizer
    entityExtractor *EntityExtractor
    intentClassifier *IntentClassifier
    contextStore *AIContextStore
    models       map[string]*Model
}

// Processing pipeline
type ProcessingPipeline struct {
    stages []ProcessingStage
    config *NLPConfig
}

const (
    // Performance targets
    ProcessingLatencyP50 = 100 * time.Millisecond
    ProcessingLatencyP95 = 500 * time.Millisecond
    ModelLoadTime        = 2 * time.Second
    MemoryUsageLimit     = 256 * 1024 * 1024  // 256 MB
)
```

#### NLP Capabilities
- Turkish and English language support
- Task intent classification (create, update, delete, query)
- Entity extraction (dates, priorities, projects)
- Context awareness and memory
- Semantic similarity matching

### 4. Auto State Manager (`state/`)

```go
// State manager specifications  
type AutoStateManager struct {
    watchers    map[string]*StateWatcher
    triggers    []StateTrigger
    persistence *StatePersistence
    eventBus    *EventBus
}

type StateTransition struct {
    From       string            `json:"from"`
    To         string            `json:"to"`
    Trigger    string            `json:"trigger"`
    Conditions map[string]string `json:"conditions"`
    Actions    []string          `json:"actions"`
}
```

### 5. File Watcher (`watcher/`)

```go
// File watcher specifications
type FileWatcher struct {
    watchers    map[string]*fsnotify.Watcher
    taskMap     map[string][]string  // file -> task IDs
    eventQueue  chan FileEvent
    processor   *EventProcessor
}

// Supported events
const (
    FileCreated  = "created"
    FileModified = "modified" 
    FileDeleted  = "deleted"
    FileRenamed  = "renamed"
)
```

### 6. VS Code Extension (`extension/`)

```typescript
// Extension architecture
interface ExtensionComponents {
  taskProvider: TaskProvider;
  mcpClient: MCPClient;
  uiManager: UIManager;
  configManager: ConfigManager;
  eventBus: EventBus;
}

// Performance requirements
const PERFORMANCE_TARGETS = {
  startupTime: 2000,        // 2 seconds max
  responseTime: 500,        // 500ms max for UI updates
  memoryUsage: 50 * 1024 * 1024,  // 50 MB limit
} as const;
```

---

## ⚙️ Configuration Management

### Configuration Schema

```yaml
# config.yaml schema
server:
  host: "localhost"
  port: 8080
  timeout: "30s"
  max_connections: 100
  
database:
  type: "sqlite"  # or "postgresql"
  path: "./gorev.db"
  connection_pool_size: 25
  max_idle_connections: 10
  connection_max_lifetime: "1h"
  
nlp:
  enabled: true
  language: "tr"  # "tr" or "en"
  model_path: "./models"
  cache_size: 1000
  processing_timeout: "5s"
  
mcp:
  enabled: true
  tools_enabled: ["all"]  # or specific tool names
  rate_limit: 100  # requests per minute
  max_concurrent: 20
  
logging:
  level: "info"  # debug, info, warn, error
  format: "json"  # json or text
  output: "stdout"  # stdout, stderr, or file path
  
file_watcher:
  enabled: true
  poll_interval: "1s"
  max_watchers: 1000
  buffer_size: 100

ai_context:
  max_entries: 10000
  cleanup_interval: "24h"
  context_ttl: "7d"
```

### Environment Variables

```bash
# Production deployment variables
GOREV_ENV=production
GOREV_DATABASE_URL=postgres://user:pass@host:5432/gorev
GOREV_SECRET_KEY=your-secret-key-here
GOREV_LOG_LEVEL=info
GOREV_METRICS_ENABLED=true
GOREV_SENTRY_DSN=https://your-sentry-dsn

# Development overrides  
GOREV_DEBUG=true
GOREV_DATABASE_URL=sqlite:///tmp/gorev_dev.db
GOREV_LOG_LEVEL=debug
```

### Configuration Validation

```go
// Configuration validation rules
type ConfigValidation struct {
    Required []string
    Ranges   map[string]Range
    Patterns map[string]*regexp.Regexp
}

var ValidationRules = ConfigValidation{
    Required: []string{
        "database.type",
        "database.path", 
        "server.port",
    },
    Ranges: map[string]Range{
        "server.port":                {1024, 65535},
        "database.connection_pool_size": {1, 100},
        "mcp.rate_limit":             {1, 10000},
    },
}
```

---

## 🗄️ Database Schema

### Entity Relationship Diagram

```mermaid
erDiagram
    PROJECTS {
        string id PK
        string name UK
        string description
        string status
        datetime created_at
        datetime updated_at
    }
    
    TASKS {
        string id PK
        string title
        string description
        string status
        string priority
        string project_id FK
        datetime due_date
        datetime created_at
        datetime updated_at
    }
    
    TASK_DEPENDENCIES {
        string id PK
        string source_task_id FK
        string target_task_id FK
        string dependency_type
        datetime created_at
    }
    
    AI_CONTEXTS {
        string id PK
        string context_type
        text content
        jsonb metadata
        datetime created_at
        datetime expires_at
    }
    
    FILE_WATCHERS {
        string id PK
        string file_path UK
        string task_id FK
        string watch_type
        boolean active
        datetime created_at
    }
    
    PROJECTS ||--o{ TASKS : contains
    TASKS ||--o{ TASK_DEPENDENCIES : "source"
    TASKS ||--o{ TASK_DEPENDENCIES : "target"
    TASKS ||--o{ FILE_WATCHERS : monitors
```

### Table Specifications

#### Projects Table
```sql
CREATE TABLE projects (
    id TEXT PRIMARY KEY DEFAULT (hex(randomblob(16))),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'archived')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_created_at ON projects(created_at);
```

#### Tasks Table
```sql
CREATE TABLE tasks (
    id TEXT PRIMARY KEY DEFAULT (hex(randomblob(16))),
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled')),
    priority TEXT NOT NULL DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
    project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
    due_date DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
```

### Database Performance Specifications

```yaml
Performance Requirements:
  Query Response Time (P95): < 100ms
  Connection Pool Size: 25 connections
  Transaction Timeout: 30 seconds
  Backup Retention: 30 days
  
Scaling Limits:
  Max Tasks: 1,000,000 per database
  Max Projects: 10,000 per database  
  Max Concurrent Connections: 100
  Max Query Complexity: 10 JOINs

Backup Strategy:
  Frequency: Every 6 hours
  Retention: 30 days online, 365 days archive
  Recovery Time Objective (RTO): < 1 hour
  Recovery Point Objective (RPO): < 6 hours
```

---

## 🔌 API Specifications

### MCP Protocol Implementation

#### Tool Call Request Format
```json
{
  "jsonrpc": "2.0",
  "id": "unique-request-id",
  "method": "tools/call",
  "params": {
    "name": "gorev_listele",
    "arguments": {
      "durum": "pending",
      "limit": 10,
      "offset": 0
    }
  }
}
```

#### Tool Call Response Format
```json
{
  "jsonrpc": "2.0",
  "id": "unique-request-id",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "# Pending Tasks\n1. Complete API documentation\n2. Fix authentication bug"
      }
    ],
    "isError": false
  }
}
```

### REST API (Internal)

#### Authentication
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
Accept: application/json
X-API-Version: v2.0
```

#### Rate Limiting Headers
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1694520000
Retry-After: 60
```

### WebSocket Events (Real-time Updates)

```typescript
// WebSocket event types
interface TaskEvent {
  type: 'task.created' | 'task.updated' | 'task.deleted';
  task: Task;
  timestamp: string;
  user_id: string;
}

interface FileWatchEvent {
  type: 'file.changed';
  file_path: string;
  watch_type: string;
  task_ids: string[];
  timestamp: string;
}
```

---

## 🔒 Security & Compliance

### Security Architecture

```mermaid
graph TB
    A[Client Request] --> B[Input Validation]
    B --> C[Authentication]
    C --> D[Authorization]
    D --> E[Rate Limiting]
    E --> F[Business Logic]
    F --> G[Data Access Layer]
    G --> H[Audit Logging]
    H --> I[Response]
```

### Security Requirements

#### Input Validation
```go
// Rule 15 compliant validation
type InputValidator struct {
    rules map[string]ValidationRule
}

type ValidationRule struct {
    Required bool
    Type     string
    MinLen   int
    MaxLen   int
    Pattern  *regexp.Regexp
    Enum     []string
}

// All user inputs MUST be validated
func (v *InputValidator) Validate(data map[string]interface{}) error {
    for field, rule := range v.rules {
        if err := v.validateField(field, data[field], rule); err != nil {
            return fmt.Errorf("validation failed for field '%s': %w", field, err)
        }
    }
    return nil
}
```

#### SQL Injection Prevention
```go
// 100% parameterized queries - NO string concatenation
func (db *Database) CreateTask(task *Task) error {
    query := `INSERT INTO tasks (id, title, description, status, project_id) 
              VALUES (?, ?, ?, ?, ?)`
    _, err := db.Exec(query, task.ID, task.Title, task.Description, task.Status, task.ProjectID)
    return err
}
```

#### Authentication & Authorization
```yaml
Authentication Methods:
  - JWT tokens (default)
  - API keys for service accounts
  - OAuth2 for third-party integrations
  
Token Specifications:
  Algorithm: RS256
  Expiration: 24 hours
  Refresh Token TTL: 30 days
  
Authorization Model:
  - Role-based access control (RBAC)
  - Resource-level permissions
  - Project-based isolation
```

### Compliance Matrix

| Standard | Coverage | Status | Notes |
|----------|----------|--------|-------|
| OWASP Top 10 2023 | 100% | ✅ Compliant | All vulnerabilities mitigated |
| ISO 27001 | 95% | 🟡 Partial | Data classification pending |
| SOC 2 Type II | 90% | 🟡 Partial | Audit logging complete |
| GDPR | 100% | ✅ Compliant | Data deletion, export ready |
| PCI DSS | N/A | N/A | No payment data processed |

### Security Monitoring

```yaml
Security Events:
  - Failed authentication attempts
  - Privilege escalation attempts
  - Unusual query patterns
  - Rate limit violations
  - Data export activities

Alerting Thresholds:
  Failed Logins: 5 attempts/5 minutes
  Rate Limit: 80% of limit reached
  Database Errors: 10 errors/minute
  Response Time: P95 > 2 seconds
```

---

## ⚡ Performance Requirements

### Response Time Targets

| Operation Type | P50 Target | P95 Target | P99 Target |
|----------------|------------|------------|------------|
| Task List | 50ms | 100ms | 200ms |
| Task Create | 100ms | 200ms | 500ms |
| Task Update | 75ms | 150ms | 300ms |
| NLP Processing | 100ms | 500ms | 1000ms |
| File Watch Events | 10ms | 50ms | 100ms |

### Throughput Requirements

```yaml
MCP Operations:
  Peak: 1,000 requests/second
  Sustained: 500 requests/second
  Concurrent Users: 100

Database Operations:
  Reads: 5,000 queries/second
  Writes: 1,000 queries/second
  Connections: 25 concurrent

File System:
  Watched Files: 10,000 maximum
  Events/Second: 1,000 maximum
  Processing Latency: < 100ms
```

### Resource Utilization Limits

```yaml
Memory Usage:
  Core Engine: 512 MB maximum
  NLP Processor: 256 MB maximum
  VS Code Extension: 50 MB maximum
  Total System: 1 GB recommended

CPU Usage:
  Normal Operation: < 10% (single core)
  Peak Load: < 50% (single core)
  NLP Processing: < 80% (during batch operations)

Disk I/O:
  Database: < 100 IOPS normal, < 500 IOPS peak
  Log Files: < 10 MB/hour
  Temporary Files: < 100 MB total
```

### Performance Monitoring

```go
// Performance metrics collection
type PerformanceMetrics struct {
    RequestDuration  prometheus.HistogramVec
    RequestCount     prometheus.CounterVec
    ActiveConnections prometheus.Gauge
    MemoryUsage      prometheus.Gauge
    DatabaseLatency  prometheus.HistogramVec
}

// SLA monitoring
const (
    SLA_ResponseTime = 200 * time.Millisecond  // P95
    SLA_Availability = 0.999                   // 99.9% uptime
    SLA_Throughput   = 500                     // requests/second
)
```

---

## 📊 Monitoring & Observability

### Metrics Architecture

```mermaid
graph LR
    A[Application] --> B[Prometheus Metrics]
    A --> C[Structured Logs]
    A --> D[Distributed Tracing]
    
    B --> E[Prometheus Server]
    C --> F[Log Aggregator]
    D --> G[Jaeger/Zipkin]
    
    E --> H[Grafana Dashboard]
    F --> H
    G --> H
    
    H --> I[Alertmanager]
    I --> J[PagerDuty/Slack]
```

### Application Metrics

```yaml
Business Metrics:
  - tasks_created_total
  - tasks_completed_total  
  - projects_active_count
  - users_active_count
  - nlp_requests_total
  - file_watch_events_total

Technical Metrics:
  - http_requests_duration_seconds
  - http_requests_total
  - database_connections_active
  - database_query_duration_seconds
  - memory_usage_bytes
  - cpu_usage_percent

Error Metrics:
  - errors_total
  - database_errors_total
  - nlp_errors_total
  - validation_errors_total
```

### Logging Standards

```go
// Structured logging format
type LogEntry struct {
    Timestamp string                 `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Component string                 `json:"component"`
    TraceID   string                 `json:"trace_id"`
    UserID    string                 `json:"user_id"`
    Fields    map[string]interface{} `json:"fields"`
}

// Example log entry
{
  "timestamp": "2025-09-12T10:30:00Z",
  "level": "info",
  "message": "Task created successfully",
  "component": "task_service",
  "trace_id": "abc123def456",
  "user_id": "user_789",
  "fields": {
    "task_id": "task_123",
    "project_id": "proj_456",
    "priority": "high",
    "duration_ms": 45
  }
}
```

### Health Check Endpoints

```yaml
Health Check APIs:
  /health/live:
    description: "Liveness probe"
    response_time: < 100ms
    dependencies: []
    
  /health/ready:
    description: "Readiness probe" 
    response_time: < 500ms
    dependencies: ["database", "nlp_service"]
    
  /health/detailed:
    description: "Detailed health status"
    response_time: < 1000ms
    includes: ["version", "uptime", "dependencies", "metrics"]
```

### Distributed Tracing

```go
// OpenTelemetry tracing implementation
func (s *TaskService) CreateTask(ctx context.Context, req *CreateTaskRequest) error {
    ctx, span := otel.Tracer("task_service").Start(ctx, "CreateTask")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("task.title", req.Title),
        attribute.String("task.priority", req.Priority),
        attribute.String("project.id", req.ProjectID),
    )
    
    // Business logic with child spans
    if err := s.validateRequest(ctx, req); err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "validation failed")
        return err
    }
    
    return s.repository.Create(ctx, task)
}
```

---

## 🚀 Deployment Architecture

### Container Specifications

```dockerfile
# Multi-stage Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gorev .

FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/gorev .
COPY --from=builder /app/config.yaml .
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ./gorev health || exit 1
CMD ["./gorev", "serve"]
```

### Kubernetes Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gorev
  labels:
    app: gorev
    version: v0.14.0
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gorev
  template:
    metadata:
      labels:
        app: gorev
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      containers:
      - name: gorev
        image: gorev:v0.14.0
        ports:
        - containerPort: 8080
        env:
        - name: GOREV_DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: gorev-secrets
              key: database-url
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi" 
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

### Cloud Architecture (AWS Example)

```mermaid
graph TB
    subgraph "Internet"
        U[Users]
    end
    
    subgraph "AWS Cloud"
        subgraph "Public Subnet"
            ALB[Application Load Balancer]
        end
        
        subgraph "Private Subnet"
            ECS1[ECS Task 1]
            ECS2[ECS Task 2]  
            ECS3[ECS Task 3]
        end
        
        subgraph "Database Subnet"
            RDS[RDS PostgreSQL]
            ElastiCache[ElastiCache Redis]
        end
        
        subgraph "Monitoring"
            CW[CloudWatch]
            XR[X-Ray]
        end
    end
    
    U --> ALB
    ALB --> ECS1
    ALB --> ECS2
    ALB --> ECS3
    
    ECS1 --> RDS
    ECS2 --> RDS
    ECS3 --> RDS
    
    ECS1 --> ElastiCache
    ECS2 --> ElastiCache
    ECS3 --> ElastiCache
    
    ECS1 --> CW
    ECS2 --> CW
    ECS3 --> CW
    
    ECS1 --> XR
    ECS2 --> XR
    ECS3 --> XR
```

### Infrastructure as Code

```terraform
# AWS ECS deployment
resource "aws_ecs_cluster" "gorev" {
  name = "gorev-cluster"
  
  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

resource "aws_ecs_service" "gorev" {
  name            = "gorev-service"
  cluster         = aws_ecs_cluster.gorev.id
  task_definition = aws_ecs_task_definition.gorev.arn
  desired_count   = 3
  
  load_balancer {
    target_group_arn = aws_lb_target_group.gorev.arn
    container_name   = "gorev"
    container_port   = 8080
  }
  
  deployment_configuration {
    maximum_percent         = 200
    minimum_healthy_percent = 100
  }
}
```

---

## 🧪 Quality Assurance

### Test Strategy

```yaml
Test Pyramid:
  Unit Tests: 70% coverage minimum
  Integration Tests: 20% coverage
  End-to-End Tests: 10% coverage
  
Test Categories:
  - Business Logic Tests
  - API Contract Tests
  - Database Integration Tests
  - Performance Tests
  - Security Tests
  - Chaos Engineering Tests
```

### Test Implementation

```go
// Example test following Rule 15 principles
func TestTaskService_CreateTask(t *testing.T) {
    tests := []struct {
        name           string
        input          *CreateTaskRequest
        expectedError  string
        expectedStatus TaskStatus
        setupMocks     func(*MockRepository)
    }{
        {
            name: "successful_task_creation",
            input: &CreateTaskRequest{
                Title:       "Test Task",
                Description: "Test Description",
                Priority:    PriorityHigh,
                ProjectID:   "proj_123",
            },
            expectedStatus: StatusPending,
            setupMocks: func(m *MockRepository) {
                m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
            },
        },
        {
            name: "validation_error_empty_title",
            input: &CreateTaskRequest{
                Title:     "", // Invalid: empty title
                Priority:  PriorityHigh,
                ProjectID: "proj_123",
            },
            expectedError: "validation failed: title is required",
            setupMocks:    func(m *MockRepository) {}, // No calls expected
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            
            mockRepo := NewMockRepository(ctrl)
            tt.setupMocks(mockRepo)
            
            service := NewTaskService(mockRepo, &Config{})
            
            // Execute
            result, err := service.CreateTask(context.Background(), tt.input)
            
            // Assert
            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
                assert.Nil(t, result)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
                assert.Equal(t, tt.expectedStatus, result.Status)
            }
        })
    }
}
```

### Performance Testing

```yaml
Performance Test Requirements:
  Load Testing:
    - 1,000 concurrent users
    - 10,000 requests over 10 minutes
    - Response time P95 < 200ms
    
  Stress Testing:
    - Gradually increase load to failure point
    - Identify breaking point
    - Verify graceful degradation
    
  Endurance Testing:
    - Normal load for 24 hours
    - Monitor memory leaks
    - Verify stable performance
```

### Security Testing

```yaml
Security Test Coverage:
  SAST (Static Analysis):
    - gosec for Go code
    - eslint-plugin-security for TypeScript
    - Semgrep for general patterns
    
  DAST (Dynamic Analysis):
    - OWASP ZAP scanning
    - SQL injection testing
    - XSS vulnerability testing
    
  Dependency Scanning:
    - Go mod vulnerability check
    - npm audit for Node.js deps
    - Snyk container scanning
```

---

## 🔄 Migration & Compatibility

### Version Compatibility Matrix

| Version | Go Min | Node.js Min | SQLite Min | PostgreSQL Min | Breaking Changes |
|---------|--------|-------------|------------|----------------|------------------|
| v0.14.0 | 1.21   | 18.0        | 3.35       | 13.0           | Template system changes |
| v0.13.x | 1.20   | 16.0        | 3.32       | 12.0           | MCP protocol updates |
| v0.12.x | 1.19   | 16.0        | 3.32       | 12.0           | None |

### Migration Scripts

```go
// Database migration system
type Migration struct {
    Version     string    `json:"version"`
    Description string    `json:"description"`
    Timestamp   time.Time `json:"timestamp"`
    Up          func(tx *sql.Tx) error
    Down        func(tx *sql.Tx) error
}

var migrations = []Migration{
    {
        Version:     "v0.14.0",
        Description: "Add AI context and file watcher tables",
        Up: func(tx *sql.Tx) error {
            queries := []string{
                `CREATE TABLE ai_contexts (
                    id TEXT PRIMARY KEY,
                    context_type TEXT NOT NULL,
                    content TEXT NOT NULL,
                    metadata TEXT,
                    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
                )`,
                `CREATE TABLE file_watchers (
                    id TEXT PRIMARY KEY,
                    file_path TEXT NOT NULL UNIQUE,
                    task_id TEXT REFERENCES tasks(id),
                    watch_type TEXT NOT NULL,
                    active BOOLEAN DEFAULT true,
                    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
                )`,
            }
            
            for _, query := range queries {
                if _, err := tx.Exec(query); err != nil {
                    return fmt.Errorf("failed to execute migration query: %w", err)
                }
            }
            return nil
        },
    },
}
```

### Backward Compatibility

```yaml
API Compatibility:
  MCP Protocol: Backward compatible with v1.0
  REST API: Semantic versioning (Major.Minor.Patch)
  Database Schema: Forward-only migrations
  
Configuration Compatibility:
  Config File: Backward compatible with defaults
  Environment Variables: Prefix-based isolation
  CLI Arguments: Deprecated flags supported with warnings

Extension Compatibility:
  VS Code: Minimum VS Code 1.60
  Language Server: Protocol v3.17
  File Formats: JSON schema versioning
```

### Migration Checklist

```markdown
## Pre-Migration Checklist
- [ ] Backup current database
- [ ] Review configuration changes
- [ ] Test migration in staging environment
- [ ] Verify VS Code extension compatibility
- [ ] Check custom MCP tool integrations

## Migration Steps
- [ ] Stop current Gorev instance
- [ ] Update binary/container image
- [ ] Run database migrations
- [ ] Update configuration files
- [ ] Restart Gorev service
- [ ] Verify health checks pass
- [ ] Test core functionality

## Post-Migration Verification
- [ ] All tasks visible and accessible
- [ ] Projects preserved correctly
- [ ] MCP tools responding properly
- [ ] VS Code extension connected
- [ ] Performance within SLA limits
- [ ] Monitoring and alerts functional
```

---

## 📋 Appendices

### A. Glossary

| Term | Definition |
|------|------------|
| MCP | Model Context Protocol - standardized interface for AI tools |
| NLP | Natural Language Processing - AI text analysis capabilities |
| Rule 15 | Zero errors/warnings policy - proper fixes only, no suppressions |
| DRY | Don't Repeat Yourself - eliminate code duplication |
| P95 | 95th percentile - metric where 95% of values are below this threshold |
| SLA | Service Level Agreement - guaranteed service performance levels |

### B. References

1. [Go Language Specification](https://golang.org/ref/spec)
2. [MCP Protocol Specification](https://modelcontextprotocol.io/docs)
3. [VS Code Extension API](https://code.visualstudio.com/api)
4. [OWASP Top 10 2023](https://owasp.org/Top10/)
5. [PostgreSQL Documentation](https://www.postgresql.org/docs/)
6. [Kubernetes Documentation](https://kubernetes.io/docs/)

### C. Contact Information

```yaml
Development Team:
  Lead Developer: Mustafa SENOL
  Email: info@gorev.dev
  Repository: https://github.com/msenol/gorev
  
Support:
  Documentation: https://docs.gorev.dev
  Issues: https://github.com/msenol/gorev/issues
  Discussions: https://github.com/msenol/gorev/discussions
```

---

<div align="center">

**Gorev Technical Specification v2.0**  
*Following Rule 15 & DRY Principles for Reliable System Design*

**Document Version**: 2.0.0  
**Last Updated**: September 12, 2025  
**Status**: ✅ Production Ready

</div>