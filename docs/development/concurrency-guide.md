# Concurrency Development Guide

This guide provides practical guidance for developers working with concurrent programming patterns in the Gorev project. It covers established patterns, development workflows, and testing strategies specific to the Gorev codebase.

## Table of Contents

- [Gorev Concurrency Patterns](#gorev-concurrency-patterns)
- [Development Workflow](#development-workflow)
- [Testing Strategies](#testing-strategies)
- [Performance Guidelines](#performance-guidelines)
- [Common Pitfalls and Solutions](#common-pitfalls-and-solutions)

## Gorev Concurrency Patterns

### 1. AI Context Manager Pattern (v0.11.1+)

The `AIContextYonetici` implements the canonical concurrent access pattern for Gorev:

#### Structure

```go
type AIContextYonetici struct {
    veriYonetici     VeriYoneticiInterface
    autoStateManager *AutoStateManager
    mu               sync.RWMutex
}
```

#### Implementation Pattern

```go
// Public API with synchronization
func (acy *AIContextYonetici) SetActiveTask(taskID string) error {
    acy.mu.Lock()
    defer acy.mu.Unlock()
    
    context, err := acy.getContextUnsafe()
    if err != nil {
        context = &AIContext{
            RecentTasks: []string{},
            SessionData: make(map[string]interface{}),
        }
    }
    
    context.ActiveTaskID = taskID
    context.LastUpdated = time.Now()
    
    return acy.saveContextUnsafe(context)
}

// Internal unsafe methods for use within locks
func (acy *AIContextYonetici) getContextUnsafe() (*AIContext, error) {
    return acy.veriYonetici.AIContextGetir()
}
```

#### Key Benefits

- **No Deadlocks**: Internal unsafe methods prevent recursive locking
- **Read Optimization**: RWMutex allows concurrent reads
- **Clear API**: Public methods always safe, unsafe methods clearly marked

### 2. MCP Handler Concurrent Access

All MCP handlers must be thread-safe since multiple clients can invoke tools simultaneously:

```go
func (h *GorevHandlers) GorevSetActive(args map[string]interface{}) (*mcp.CallToolResult, error) {
    taskID, exists := args["task_id"].(string)
    if !exists {
        return mcp.NewToolResultError("INVALID_ARGS", "task_id gerekli"), nil
    }
    
    // Thread-safe call to AI context manager
    if err := h.aiContextYonetici.SetActiveTask(taskID); err != nil {
        return mcp.NewToolResultError("SET_ACTIVE_FAILED", err.Error()), nil
    }
    
    return mcp.NewToolResultText(i18n.T("success.activeTaskSet", map[string]interface{}{
        "TaskID": taskID,
    })), nil
}
```

### 3. Database Connection Safety

Gorev uses SQLite with proper concurrent access configuration:

```go
// Connection configuration for concurrent access
db, err := sql.Open("sqlite3", "gorev.db?_journal_mode=WAL&_sync=NORMAL&_cache_size=1000")
if err != nil {
    return err
}

// Set connection pool parameters
db.SetMaxOpenConns(10)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(time.Hour)
```

#### Transaction Pattern

```go
func (vy *VeriYonetici) AtomicUpdate(updates []func(*sql.Tx) error) error {
    tx, err := vy.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // Safe to call even after successful commit
    
    for _, update := range updates {
        if err := update(tx); err != nil {
            return err // Rollback will be called automatically
        }
    }
    
    return tx.Commit()
}
```

## Development Workflow

### 1. Adding Concurrent Components

When adding new components that may be accessed concurrently:

#### Step 1: Design Phase

```markdown
1. Identify shared state that needs protection
2. Choose appropriate synchronization mechanism:
   - sync.RWMutex: Read-heavy workloads
   - sync.Mutex: Simple exclusive access
   - Channels: Communication between goroutines
   - sync.Map: Concurrent map operations
3. Design internal/external API separation
```

#### Step 2: Implementation

```go
// Example: Adding a new concurrent manager
type NewConcurrentManager struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

// Public methods with synchronization
func (ncm *NewConcurrentManager) GetData(key string) (interface{}, bool) {
    ncm.mu.RLock()
    defer ncm.mu.RUnlock()
    
    value, exists := ncm.data[key]
    return value, exists
}

func (ncm *NewConcurrentManager) SetData(key string, value interface{}) {
    ncm.mu.Lock()
    defer ncm.mu.Unlock()
    
    ncm.data[key] = value
}
```

#### Step 3: Testing

```go
func TestNewConcurrentManager_RaceCondition(t *testing.T) {
    ncm := &NewConcurrentManager{
        data: make(map[string]interface{}),
    }
    
    const numGoroutines = 50
    const operationsPerGoroutine = 100
    
    // Launch concurrent operations
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            for j := 0; j < operationsPerGoroutine; j++ {
                key := fmt.Sprintf("key-%d-%d", id, j)
                ncm.SetData(key, fmt.Sprintf("value-%d-%d", id, j))
                _, _ = ncm.GetData(key)
            }
        }(i)
    }
    
    time.Sleep(100 * time.Millisecond)
    
    // Verify no race conditions (test passes if no panic/race detected)
}
```

### 2. Code Review Checklist

When reviewing concurrent code:

```markdown
- [ ] Shared data properly protected by synchronization primitives
- [ ] Consistent lock ordering to prevent deadlocks
- [ ] Minimal critical sections (small locked areas)
- [ ] Proper use of defer for cleanup
- [ ] Error handling doesn't bypass unlock operations
- [ ] Tests include concurrent scenarios
- [ ] Race detector passes (`go test -race`)
```

### 3. Integration with Existing Systems

When integrating concurrent components with existing Gorev systems:

#### Database Integration

```go
// Ensure database operations are atomic
func (manager *ConcurrentManager) UpdateWithDatabase(key string, value interface{}) error {
    manager.mu.Lock()
    defer manager.mu.Unlock()
    
    // Update in-memory state
    manager.data[key] = value
    
    // Persist to database
    return manager.veriYonetici.SaveData(key, value)
}
```

#### MCP Handler Integration

```go
func (h *GorevHandlers) NewConcurrentOperation(args map[string]interface{}) (*mcp.CallToolResult, error) {
    // Validate arguments first (outside of locks)
    key, exists := args["key"].(string)
    if !exists {
        return mcp.NewToolResultError("INVALID_ARGS", "key required"), nil
    }
    
    // Thread-safe operation
    result, err := h.concurrentManager.ProcessKey(key)
    if err != nil {
        return mcp.NewToolResultError("OPERATION_FAILED", err.Error()), nil
    }
    
    return mcp.NewToolResultText(result), nil
}
```

## Testing Strategies

### 1. Race Condition Detection

#### Standard Race Detector Usage

```bash
# Run all tests with race detector
go test -race ./...

# Run specific concurrent tests
go test -race -v -run TestConcurrent ./internal/gorev/

# Build and run with race detector
go build -race ./cmd/gorev
./gorev serve --debug
```

#### Custom Race Detection Tests

```go
func TestAIContextManager_ConcurrentAccess(t *testing.T) {
    // Test demonstrates comprehensive concurrent testing pattern
    manager := setupTestManager()
    
    // Error collection channel
    errors := make(chan error, 100)
    
    // Define concurrent operations
    operations := []func(){
        func() { manager.SetActiveTask("task-1") },
        func() { manager.GetActiveTask() },
        func() { manager.GetContext() },
        func() { manager.GetRecentTasks(5) },
    }
    
    // Launch concurrent goroutines
    const numGoroutines = 50
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            defer func() {
                if r := recover(); r != nil {
                    errors <- fmt.Errorf("goroutine %d panicked: %v", id, r)
                }
            }()
            
            for j := 0; j < 10; j++ {
                op := operations[j%len(operations)]
                if err := op(); err != nil {
                    errors <- fmt.Errorf("goroutine %d operation failed: %w", id, err)
                    return
                }
            }
        }(i)
    }
    
    // Wait and collect results
    time.Sleep(200 * time.Millisecond)
    close(errors)
    
    var collectedErrors []error
    for err := range errors {
        collectedErrors = append(collectedErrors, err)
    }
    
    assert.Empty(t, collectedErrors, "No race conditions should occur")
}
```

### 2. Load Testing Patterns

#### MCP Tool Load Testing

```bash
#!/bin/bash
# Script: test_concurrent_mcp.sh

echo "Testing concurrent MCP tool access..."

# Function to call MCP tool
call_mcp_tool() {
    local task_id="task-$1"
    ./gorev mcp call gorev_set_active "{\"task_id\": \"$task_id\"}" 2>&1
}

# Launch concurrent calls
for i in {1..20}; do
    call_mcp_tool $i &
done

# Wait for all background jobs
wait

echo "Concurrent MCP test completed"
```

#### Performance Measurement

```go
func BenchmarkConcurrentAccess(b *testing.B) {
    manager := setupBenchmarkManager()
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // Mix of operations
            switch rand.Intn(3) {
            case 0:
                manager.SetActiveTask("benchmark-task")
            case 1:
                manager.GetActiveTask()
            case 2:
                manager.GetContext()
            }
        }
    })
}
```

### 3. Integration Testing

#### End-to-End Concurrent Scenarios

```go
func TestE2E_ConcurrentClientUsage(t *testing.T) {
    // Start test server
    server := startTestServer()
    defer server.Stop()
    
    // Simulate multiple MCP clients
    const numClients = 10
    results := make(chan error, numClients)
    
    for i := 0; i < numClients; i++ {
        go func(clientID int) {
            client := createMCPClient()
            defer client.Close()
            
            // Perform sequence of operations
            taskID := fmt.Sprintf("client-%d-task", clientID)
            
            // Create task
            err := client.Call("templateden_gorev_olustur", map[string]interface{}{
                "template_id": "default_task",
                "degerler": map[string]string{
                    "baslik": fmt.Sprintf("Client %d Task", clientID),
                },
            })
            if err != nil {
                results <- err
                return
            }
            
            // Set as active
            err = client.Call("gorev_set_active", map[string]interface{}{
                "task_id": taskID,
            })
            if err != nil {
                results <- err
                return
            }
            
            results <- nil
        }(i)
    }
    
    // Collect results
    for i := 0; i < numClients; i++ {
        err := <-results
        assert.NoError(t, err, "Client %d should complete successfully", i)
    }
}
```

## Performance Guidelines

### 1. Lock Granularity

#### Coarse-Grained Locking (Current Pattern)

```go
// Good: Single mutex for related data
type AIContextYonetici struct {
    mu           sync.RWMutex
    activeTask   string
    recentTasks  []string
    sessionData  map[string]interface{}
}
```

#### Fine-Grained Locking (When Needed)

```go
// Use sparingly: Only when contention is proven problematic
type FineLockingExample struct {
    activeMu    sync.RWMutex
    activeTask  string
    
    recentMu    sync.RWMutex
    recentTasks []string
}
```

### 2. Performance Monitoring

#### Identifying Lock Contention

```go
func init() {
    // Enable lock profiling in development
    if os.Getenv("GOREV_DEBUG") == "1" {
        runtime.SetMutexProfileFraction(1)
    }
}
```

#### Profiling Commands

```bash
# Run with profiling
go test -mutexprofile=mutex.prof ./internal/gorev/

# Analyze mutex contention
go tool pprof mutex.prof
(pprof) top10
(pprof) list AIContextYonetici
```

### 3. Alternative Approaches

#### Channel-Based Coordination

```go
// For complex coordination, consider channels
type ChannelBasedManager struct {
    requests chan Request
    done     chan struct{}
}

func (cbm *ChannelBasedManager) Start() {
    go func() {
        for {
            select {
            case req := <-cbm.requests:
                cbm.handleRequest(req)
            case <-cbm.done:
                return
            }
        }
    }()
}
```

#### Atomic Operations

```go
// For simple counters and flags
type AtomicCounter struct {
    count int64
}

func (ac *AtomicCounter) Increment() {
    atomic.AddInt64(&ac.count, 1)
}

func (ac *AtomicCounter) Get() int64 {
    return atomic.LoadInt64(&ac.count)
}
```

## Common Pitfalls and Solutions

### 1. Deadlock Prevention

#### Problem: Lock Ordering

```go
// BAD: Inconsistent lock ordering can cause deadlocks
func badTransfer(from, to *Account, amount int) {
    from.mu.Lock()
    to.mu.Lock()   // Deadlock if another goroutine locks to, then from
    // ... transfer logic
    to.mu.Unlock()
    from.mu.Unlock()
}
```

#### Solution: Consistent Ordering

```go
// GOOD: Always acquire locks in the same order
func goodTransfer(from, to *Account, amount int) {
    first, second := from, to
    if from.ID > to.ID {
        first, second = to, from
    }
    
    first.mu.Lock()
    second.mu.Lock()
    // ... transfer logic
    second.mu.Unlock()
    first.mu.Unlock()
}
```

### 2. Resource Leaks

#### Problem: Missing Unlock

```go
// BAD: Missing unlock in error path
func badFunction() error {
    mu.Lock()
    
    if err := doSomething(); err != nil {
        return err  // BUG: Lock never released!
    }
    
    mu.Unlock()
    return nil
}
```

#### Solution: Always Use Defer

```go
// GOOD: Defer ensures unlock happens
func goodFunction() error {
    mu.Lock()
    defer mu.Unlock()
    
    if err := doSomething(); err != nil {
        return err  // Lock properly released
    }
    
    return nil
}
```

### 3. Race Conditions in Tests

#### Problem: Unreliable Test Timing

```go
// BAD: Relying on timing for synchronization
func badTest(t *testing.T) {
    go doAsyncWork()
    time.Sleep(100 * time.Millisecond)  // Unreliable!
    checkResult()
}
```

#### Solution: Proper Synchronization

```go
// GOOD: Use channels or WaitGroup for coordination
func goodTest(t *testing.T) {
    done := make(chan struct{})
    
    go func() {
        doAsyncWork()
        close(done)
    }()
    
    select {
    case <-done:
        checkResult()
    case <-time.After(5 * time.Second):
        t.Fatal("Test timeout")
    }
}
```

## Integration with Rule 15

Concurrency patterns in Gorev must follow Rule 15 principles:

### Comprehensive Solutions Only

- **NO Quick Fixes**: Don't add locks "just in case" without understanding the problem
- **Root Cause Analysis**: Identify actual race conditions, not just potential ones
- **Proper Testing**: Every concurrent component must have race condition tests

### No Technical Debt

- **Complete Implementation**: Concurrent safety from the start, not added later
- **No Workarounds**: Fix synchronization properly, don't disable race detector
- **Clean Abstractions**: Clear separation between safe and unsafe operations

### Example: v0.11.1 Implementation

The AI Context Manager thread-safety implementation exemplifies Rule 15 compliance:

- ✅ **Comprehensive**: All methods properly synchronized
- ✅ **No Workarounds**: Clean mutex implementation, no "temporary" solutions
- ✅ **Fully Tested**: 50-goroutine stress test with race detector
- ✅ **Zero Breaking Changes**: Backward compatibility maintained
- ✅ **Production Ready**: Handles real concurrent MCP client scenarios

## Conclusion

Concurrent programming in Gorev follows established patterns that prioritize correctness, maintainability, and performance. The AI Context Manager implementation in v0.11.1 serves as the reference implementation for thread-safety patterns in the project.

Key principles:

1. **Use established patterns** from the AI Context Manager
2. **Test thoroughly** with race detector and stress tests
3. **Follow Rule 15** - no shortcuts or workarounds
4. **Document concurrent behavior** clearly
5. **Monitor performance** and optimize when necessary

For specific questions about implementing concurrent features, refer to the `AIContextYonetici` implementation and this guide's patterns.
