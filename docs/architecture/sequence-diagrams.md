# Sequence Diagrams - Gorev Daemon

**Version:** v0.16.3 | **Last Updated:** October 6, 2025

Detailed sequence diagrams showing the flow of operations in Gorev daemon architecture.

---

## 1. Daemon Startup Sequence

```mermaid
sequenceDiagram
    participant User
    participant CLI as gorev CLI
    participant Daemon
    participant LockFile as Lock File
    participant API as REST API
    participant DB as SQLite

    User->>CLI: gorev daemon --detach
    CLI->>LockFile: Check existing lock
    alt Lock exists and process alive
        LockFile-->>CLI: Daemon already running
        CLI-->>User: Error: Already running
    else No lock or stale
        CLI->>Daemon: Start daemon process
        Daemon->>DB: Initialize database
        DB-->>Daemon: Connection established
        Daemon->>API: Start Fiber server (port 5082)
        API-->>Daemon: Server ready
        Daemon->>LockFile: Write lock file
        Note over LockFile: PID, port, version,<br/>daemon URL, start time
        Daemon-->>CLI: Daemon started
        CLI-->>User: Success (PID: 12345)
        Daemon->>Daemon: Listen for connections
    end
```

---

## 2. VS Code Extension Auto-Start

```mermaid
sequenceDiagram
    participant VSCode as VS Code Extension
    participant Lock as Lock File
    participant Health as Health Endpoint
    participant Daemon
    participant TreeView as Task Tree View

    VSCode->>VSCode: activate()
    VSCode->>Lock: Read ~/.gorev-daemon/.lock
    alt Lock file exists
        Lock-->>VSCode: {PID, URL, port}
        VSCode->>Health: GET /api/health
        alt Daemon alive
            Health-->>VSCode: {status: "healthy"}
            VSCode->>VSCode: Connect to existing daemon
        else Daemon dead
            Health--xVSCode: Connection failed
            VSCode->>Daemon: Start new daemon
            Daemon-->>VSCode: Started (PID: 12345)
        end
    else No lock file
        VSCode->>Daemon: Start daemon
        Daemon->>Lock: Create lock file
        Daemon-->>VSCode: Started (PID: 12345)
    end
    VSCode->>TreeView: Initialize tree provider
    VSCode->>Daemon: GET /api/gorevler
    Daemon-->>VSCode: Task list
    VSCode->>TreeView: Render tasks
```

---

## 3. MCP Client Connection (Claude)

```mermaid
sequenceDiagram
    participant Claude as Claude Desktop
    participant MCPProxy as MCP Proxy
    participant Handler as MCP Handlers
    participant DB as SQLite

    Claude->>MCPProxy: stdio connect
    MCPProxy->>MCPProxy: Generate client ID
    Note over MCPProxy: Client: abc-123
    MCPProxy-->>Claude: Connection established

    Claude->>MCPProxy: list_tools()
    MCPProxy->>Handler: GetToolList()
    Handler-->>MCPProxy: 24 tools
    MCPProxy-->>Claude: Tool list

    Claude->>MCPProxy: call_tool("gorev_listele")
    MCPProxy->>Handler: Execute gorev_listele
    Handler->>DB: SELECT tasks...
    DB-->>Handler: Task rows
    Handler->>Handler: Build tree structure
    Handler-->>MCPProxy: JSON result
    MCPProxy-->>Claude: Task list (markdown)

    Note over Claude,DB: Client remains connected<br/>for subsequent requests
```

---

## 4. Task Creation Flow

```mermaid
sequenceDiagram
    participant Client
    participant API as REST API
    participant Handler
    participant Business as Business Logic
    participant DB as SQLite
    participant WS as WebSocket Hub

    Client->>API: POST /api/gorev
    Note over Client,API: {baslik, aciklama, oncelik}

    API->>Handler: CreateTask(params)
    Handler->>Business: ValidateTask()
    Business-->>Handler: Validation OK

    Handler->>Business: GenerateTaskID()
    Business-->>Handler: ID: abc12345

    Handler->>DB: BEGIN TRANSACTION
    DB-->>Handler: TX started

    Handler->>DB: INSERT INTO gorevler...
    DB-->>Handler: Row inserted

    Handler->>DB: INSERT INTO etiketler...
    DB-->>Handler: Tags inserted

    Handler->>DB: COMMIT
    DB-->>Handler: TX committed

    Handler->>WS: Broadcast("task_created", task)
    WS->>WS: Notify all WS clients

    Handler-->>API: Task created
    API-->>Client: 201 Created {task}

    WS-->>Client: WebSocket: {type: "created", data: task}
```

---

## 5. Bulk Update Flow

```mermaid
sequenceDiagram
    participant Client
    participant Handler as gorev_bulk
    participant Transform as Parameter Transform
    participant DB as SQLite
    participant WS as WebSocket Hub

    Client->>Handler: gorev_bulk
    Note over Client,Handler: {operation: "update",<br/>ids: [a,b,c],<br/>data: {oncelik: "yuksek"}}

    Handler->>Transform: Transform params
    Note over Transform: {ids, data} →<br/>{updates: [{id:a,...},{id:b,...}]}
    Transform-->>Handler: Transformed

    Handler->>DB: BEGIN TRANSACTION

    loop For each task
        Handler->>DB: UPDATE gorevler SET...
        DB-->>Handler: Rows affected
        alt Update succeeded
            Handler->>WS: Broadcast("task_updated")
        else Update failed
            Handler->>Handler: Log error
        end
    end

    Handler->>DB: COMMIT
    DB-->>Handler: TX committed

    Handler-->>Client: {success: 3, failed: 0}
```

---

## 6. Advanced Search Flow (with FTS5)

```mermaid
sequenceDiagram
    participant Client
    participant Handler as gorev_search
    participant Parser as Query Parser
    participant FTS as FTS5 Index
    participant DB as SQLite

    Client->>Handler: gorev_search
    Note over Client,Handler: {action: "advanced",<br/>query: "durum:devam API"}

    Handler->>Parser: parseQueryFilters()
    Note over Parser: Extract key:value pairs
    Parser-->>Handler: {<br/>  query: "API",<br/>  filters: {durum: "devam"}<br/>}

    Handler->>FTS: Search FTS5
    Note over FTS: SELECT * FROM gorevler_fts<br/>WHERE gorevler_fts MATCH 'API'
    FTS-->>Handler: Task IDs with relevance

    Handler->>DB: Apply additional filters
    Note over DB: WHERE durum = 'devam_ediyor'<br/>AND id IN (...)
    DB-->>Handler: Filtered results

    Handler->>Handler: Sort by relevance + date
    Handler-->>Client: Ranked search results
```

---

## 7. WebSocket Real-Time Updates

```mermaid
sequenceDiagram
    participant VSCode1 as VS Code Client 1
    participant VSCode2 as VS Code Client 2
    participant Hub as WebSocket Hub
    participant API as REST API
    participant DB as Database

    VSCode1->>Hub: Connect /ws
    Hub->>Hub: Register client 1
    VSCode2->>Hub: Connect /ws
    Hub->>Hub: Register client 2

    Note over VSCode1,DB: Client 1 updates a task

    VSCode1->>API: PUT /api/gorev/abc123
    API->>DB: UPDATE gorevler...
    DB-->>API: Updated
    API->>Hub: Broadcast(task_updated)

    Hub->>VSCode1: {type: "updated", task_id: "abc123"}
    Hub->>VSCode2: {type: "updated", task_id: "abc123"}

    Note over VSCode1,VSCode2: Both clients receive<br/>real-time notification

    VSCode2->>VSCode2: Refresh task tree
```

---

## 8. Graceful Shutdown Sequence

```mermaid
sequenceDiagram
    participant User
    participant Daemon
    participant API as REST API
    participant WS as WebSocket Hub
    participant Clients as Connected Clients
    participant Lock as Lock File
    participant DB as Database

    User->>Daemon: SIGTERM
    Daemon->>Daemon: Initiate shutdown

    Daemon->>API: Stop accepting new requests
    Daemon->>WS: Stop accepting new connections

    Daemon->>Clients: Broadcast shutdown warning
    Clients-->>Daemon: Disconnect gracefully

    Daemon->>WS: Close all connections
    WS-->>Daemon: All closed

    Daemon->>Daemon: Wait for in-flight requests
    Note over Daemon: Max 30 seconds timeout

    alt All requests completed
        Daemon->>DB: Close connections
        DB-->>Daemon: Closed
    else Timeout reached
        Daemon->>DB: Force close
    end

    Daemon->>Lock: Remove lock file
    Lock-->>Daemon: Deleted

    Daemon-->>User: Exit 0
```

---

## 9. Multi-Workspace Handling

```mermaid
sequenceDiagram
    participant Client
    participant Daemon
    participant Manager as Workspace Manager
    participant DB1 as Workspace 1 DB
    participant DB2 as Workspace 2 DB

    Client->>Daemon: Request with workspace path
    Note over Client,Daemon: /home/user/project-a

    Daemon->>Manager: GetWorkspaceID(path)
    Manager->>Manager: SHA256 hash
    Manager-->>Daemon: workspace_id: abc123

    Daemon->>Manager: GetDatabase(workspace_id)

    alt Database exists
        Manager->>DB1: Get connection pool
        DB1-->>Manager: Connection
    else Database not exists
        Manager->>DB1: Create + initialize
        DB1-->>Manager: New connection
    end

    Manager-->>Daemon: Database connection
    Daemon->>DB1: Execute query
    DB1-->>Daemon: Results
    Daemon-->>Client: Response

    Note over Daemon,DB2: Subsequent request<br/>for different workspace

    Client->>Daemon: Request (/home/user/project-b)
    Daemon->>Manager: GetWorkspaceID(path)
    Manager-->>Daemon: workspace_id: def456
    Daemon->>Manager: GetDatabase(def456)
    Manager->>DB2: Get connection pool
    Manager-->>Daemon: Connection
    Daemon->>DB2: Execute query
    DB2-->>Daemon: Results
```

---

## 10. Error Handling and Retry Flow

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant DB
    participant Logger

    Client->>Handler: call_tool("gorev_guncelle")

    Handler->>DB: UPDATE gorevler...

    alt Success
        DB-->>Handler: OK
        Handler-->>Client: Success
    else Database locked
        DB--xHandler: SQLITE_BUSY
        Handler->>Handler: Wait + retry (3x)
        Handler->>DB: UPDATE gorevler...
        DB-->>Handler: OK
        Handler-->>Client: Success (after retry)
    else Constraint violation
        DB--xHandler: FOREIGN_KEY_CONSTRAINT
        Handler->>Logger: Log error
        Handler-->>Client: Error: Invalid reference
    else Unknown error
        DB--xHandler: Unknown error
        Handler->>Logger: Log with context
        Handler->>Handler: Rollback transaction
        Handler-->>Client: Error: Internal error
    end
```

---

## Performance Metrics

Based on v0.16.3 benchmarks:

| Operation | Avg Time | Components Involved |
|-----------|----------|---------------------|
| Daemon Startup | 500-800ms | Lock, DB init, API start |
| Client Connection | 50-100ms | MCP proxy, client registration |
| Task List (50) | 5-15ms | Handler → DB → Tree build |
| Task Create | 3-8ms | Validation → DB insert → WS broadcast |
| Bulk Update (10) | 11-33ms | Transform → 10x DB update → Broadcast |
| Advanced Search | 6-67ms | FTS5 query → Filter → Rank |
| WebSocket Broadcast | 1-2ms | Hub → All clients |

---

## See Also

- [Daemon Architecture](./daemon-architecture.md) - Complete technical documentation
- [API Integration Examples](../api/integration-examples.md) - Client implementation examples
- [Performance Benchmarks](../reports/performance-benchmarks.md) - Detailed metrics

---

**Last Updated:** October 6, 2025 | **Version:** v0.16.3
