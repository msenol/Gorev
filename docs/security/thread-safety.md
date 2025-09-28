# Thread-Safety Guidelines for Gorev Development

This document provides comprehensive guidance for implementing and maintaining thread-safe code in the Gorev project, with specific focus on concurrent access patterns and race condition prevention.

## Overview

As of v0.11.1, Gorev implements comprehensive thread-safety measures in critical components, particularly the AI Context Manager (`AIContextYonetici`). This document outlines the patterns, best practices, and testing strategies for concurrent programming in the Gorev codebase.

## AI Context Manager Thread-Safety Pattern

### Implementation Details

The `AIContextYonetici` struct in `internal/gorev/ai_context_yonetici.go` demonstrates the canonical thread-safety pattern for Gorev:

```go
type AIContextYonetici struct {
    veriYonetici     VeriYoneticiInterface
    autoStateManager *AutoStateManager
    mu               sync.RWMutex // Protects concurrent access to context operations
}
```

### Read vs Write Lock Usage

#### Read Operations (sync.RWMutex.RLock)

Use read locks for operations that only read data without modification:

```go
func (acy *AIContextYonetici) GetActiveTask() (*Gorev, error) {
    acy.mu.RLock()         // Multiple concurrent readers allowed
    defer acy.mu.RUnlock()
    
    context, err := acy.getContextUnsafe()
    // ... read-only operations
}
```

**Read Lock Operations:**

- `GetActiveTask()` - Reading active task ID
- `GetRecentTasks()` - Reading recent task list
- `GetContext()` - Reading AI context data

#### Write Operations (sync.RWMutex.Lock)

Use exclusive locks for operations that modify shared data:

```go
func (acy *AIContextYonetici) SetActiveTask(taskID string) error {
    acy.mu.Lock()         // Exclusive access required
    defer acy.mu.Unlock()
    
    // ... modify context data
    context.ActiveTaskID = taskID
    context.LastUpdated = time.Now()
    
    return acy.saveContextUnsafe(context)
}
```

**Write Lock Operations:**

- `SetActiveTask()` - Modifying active task and recent tasks
- `saveContext()` - Persisting context changes

### Internal Unsafe Method Pattern

To avoid deadlocks when calling methods within locked sections, use the internal unsafe pattern:

```go
// Public methods with mutex protection
func (acy *AIContextYonetici) GetContext() (*AIContext, error) {
    acy.mu.RLock()
    defer acy.mu.RUnlock()
    return acy.getContextUnsafe()
}

func (acy *AIContextYonetici) saveContext(context *AIContext) error {
    acy.mu.Lock()
    defer acy.mu.Unlock()
    return acy.saveContextUnsafe(context)
}

// Internal methods without mutex protection (unsafe outside of locks)
func (acy *AIContextYonetici) getContextUnsafe() (*AIContext, error) {
    return acy.veriYonetici.AIContextGetir()
}

func (acy *AIContextYonetici) saveContextUnsafe(context *AIContext) error {
    return acy.veriYonetici.AIContextKaydet(context)
}
```

**Benefits:**

- Prevents recursive locking (deadlocks)
- Allows internal method reuse within locked sections
- Maintains clear separation between protected and unprotected code paths

## Best Practices

### 1. Mutex Placement

- **Struct-level mutexes**: For protecting entire object state
- **Fine-grained locking**: Only when absolutely necessary for performance
- **Avoid global mutexes**: Prefer composition and dependency injection

### 2. Lock Ordering

- Always acquire locks in the same order to prevent deadlocks
- Document lock ordering requirements in complex scenarios
- Consider using a single mutex for related data structures

### 3. Critical Section Minimization

- Keep locked sections as small as possible
- Perform expensive operations outside of locks when safe
- Use defer for unlock to ensure proper cleanup

### 4. Error Handling in Concurrent Code

- Always use defer for unlocking, even in error paths
- Consider returning errors rather than panicking in concurrent code
- Log concurrent operation failures with sufficient context

### 5. Data Structure Design

- Prefer immutable data structures where possible
- Use channels for communication between goroutines
- Consider sync.Map for concurrent map operations

## When to Use Mutexes vs Channels

### Use Mutexes When

- Protecting shared memory access
- Simple critical sections
- Performance is critical (lower overhead than channels)
- Working with existing synchronous APIs

### Use Channels When

- Coordinating between goroutines
- Implementing producer-consumer patterns
- Need to pass ownership of data
- Building pipeline architectures

## Race Condition Prevention

### Common Race Condition Patterns in Gorev

#### 1. AI Context State Management

**Problem**: Multiple MCP tools accessing AI context simultaneously
**Solution**: RWMutex protection in `AIContextYonetici`

#### 2. Task Status Updates

**Problem**: Concurrent status updates to the same task
**Solution**: Database-level locking and atomic operations

#### 3. File System Operations

**Problem**: Multiple goroutines reading/writing project files
**Solution**: Coordinate through channels or file-level locking

### Detection Strategies

#### 1. Go Race Detector

Always test with the race detector enabled:

```bash
go test -race ./...
go run -race ./cmd/gorev
```

#### 2. Stress Testing

Create tests that exercise concurrent scenarios:

```go
func TestAIContextRaceCondition(t *testing.T) {
    const numGoroutines = 50
    const operationsPerGoroutine = 10
    
    // Launch concurrent operations
    for i := 0; i < numGoroutines; i++ {
        go func() {
            for j := 0; j < operationsPerGoroutine; j++ {
                // Mix of read and write operations
                acy.SetActiveTask("test-task")
                acy.GetActiveTask()
                acy.GetContext()
            }
        }()
    }
    
    // Verify no race conditions occurred
}
```

#### 3. Load Testing

Test under realistic concurrent load using MCP tools:

```bash
# Multiple concurrent MCP clients
for i in {1..10}; do
    gorev mcp call gorev_set_active '{"task_id": "task-'$i'"}' &
done
wait
```

## Testing Concurrent Code

### Test Pattern: Table-Driven Concurrent Tests

```go
func TestConcurrentOperations(t *testing.T) {
    tests := []struct {
        name          string
        numGoroutines int
        operations    []func()
        expectError   bool
    }{
        {
            name:          "concurrent reads",
            numGoroutines: 10,
            operations:    []func(){acy.GetActiveTask, acy.GetContext},
            expectError:   false,
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Test Requirements

- **Race Detector**: All concurrent tests must pass with `-race`
- **Error Collection**: Use channels to collect errors from goroutines
- **Deterministic Verification**: Verify final state consistency
- **Timeout Protection**: Use context.WithTimeout for test timeouts

## Performance Considerations

### RWMutex Performance Characteristics

- **Read locks**: Allow concurrent readers, good for read-heavy workloads
- **Write locks**: Exclusive access, blocks all other operations
- **Lock contention**: Monitor for performance degradation under high concurrency

### Monitoring Concurrent Performance

- Use Go's built-in profiler to identify lock contention
- Monitor goroutine counts and blocking operations
- Consider alternative approaches if lock contention becomes problematic

## Integration with MCP Protocol

### MCP Handler Thread-Safety

All MCP handlers in `internal/mcp/handlers.go` must be thread-safe since multiple clients can invoke tools simultaneously:

```go
func (h *GorevHandlers) GorevSetActive(args map[string]interface{}) (*mcp.CallToolResult, error) {
    // This handler is called concurrently by multiple MCP clients
    // Must ensure thread-safe access to shared resources
    return h.aiContextYonetici.SetActiveTask(taskID)
}
```

### Database Connection Thread-Safety

SQLite connections in Gorev are configured for concurrent access:

- Connection pooling handled by `database/sql`
- Transactions used for atomic operations
- Row-level locking for concurrent updates

## Debugging Race Conditions

### Tools and Techniques

1. **Go Race Detector**: Primary tool for detecting race conditions
2. **Logging**: Add detailed logging with goroutine IDs
3. **Stress Testing**: Reproduce issues under high concurrency
4. **Code Review**: Focus on shared data access patterns

### Common Debugging Commands

```bash
# Run with race detector
go test -race -v ./internal/gorev/

# Run specific concurrent test with verbose output
go test -race -v -run TestAIContextRaceCondition ./internal/gorev/

# Build and run server with race detector
go build -race ./cmd/gorev
./gorev serve --debug
```

## Migration Guide

### Converting Existing Code to Thread-Safe

When making existing code thread-safe:

1. **Identify Shared State**: Find all shared data structures
2. **Choose Protection Mechanism**: Mutex vs channels vs sync packages
3. **Add Synchronization**: Implement chosen protection mechanism
4. **Create Tests**: Add comprehensive concurrent tests
5. **Verify with Race Detector**: Ensure no race conditions remain

### Example Migration: Simple to Thread-Safe

```go
// Before: Not thread-safe
type SimpleManager struct {
    activeTask string
}

func (sm *SimpleManager) SetActive(taskID string) {
    sm.activeTask = taskID  // Race condition!
}

// After: Thread-safe
type SafeManager struct {
    activeTask string
    mu         sync.RWMutex
}

func (sm *SafeManager) SetActive(taskID string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.activeTask = taskID  // Protected by mutex
}
```

## Conclusion

Thread-safety is critical for the reliability and correctness of the Gorev system, especially as it serves multiple concurrent MCP clients. The patterns established in v0.11.1 provide a solid foundation for concurrent programming in Go.

Key takeaways:

- Use `sync.RWMutex` for read-heavy workloads with occasional writes
- Implement internal unsafe methods to avoid deadlocks
- Always test concurrent code with the Go race detector
- Follow established patterns for consistency across the codebase
- Prioritize correctness over premature optimization

For specific implementation questions or complex concurrent scenarios, refer to the `AIContextYonetici` implementation as the canonical example of thread-safe design in the Gorev project.
