# Start Task with Parameters - Systematic Implementation

## Initial Setup
First read CLAUDE.md to understand the project architecture and conventions. Then check the TodoRead tool for any existing tasks. If the provided argument matches a todo item, work on that. Otherwise, treat the argument as a new task description.

Your assigned task is defined by the following argument: **$ARGUMENTS**

## Pre-Implementation Protocol

Before implementing this task, you must follow these steps:

### 1. Deep Analysis Phase
- **THINK DEEPLY** about the task in the context of the existing Go architecture
- **ANALYZE** the relevant codebase and package structure
- Check MCP SDK documentation for best practices
- Map out all affected packages and their dependencies
- Consider edge cases and error scenarios

### 2. Go/MCP Specific Compliance
If the task involves:
- **MCP Tools**: Follow handler patterns in `internal/mcp/handlers.go`
- **Database**: Use prepared statements and transactions
- **Concurrency**: Proper goroutine and channel usage
- **Error Handling**: Explicit error returns, no panics
- **Testing**: Table-driven tests following Go conventions

### 3. Critical Questions Assessment
**ASK questions** if you find:
- Multiple valid approaches
- Unclear requirements
- Potential architectural impacts
- Missing dependencies
- Security or performance concerns
- Breaking changes potential

### 4. Proposal and Approval
- Present your understanding and proposed approach
- Wait for agreement before proceeding
- Be prepared to iterate on the approach based on feedback

## Implementation Constraints

### Absolute Rules
- **NO shortcuts** ("for now", "temporarily", "to save time")
- **STOP** if tempted to change code outside the current task's scope
- **One task at a time**
- **Test after each significant step** - all tests must pass
- **Update TodoWrite** to track progress

### Gorev Project Specific Rules
1. **Architecture**: Follow clean architecture pattern as defined in CLAUDE.md
2. **Go Standards**: Use CLAUDE.md "Development Commands" for building/testing
3. **Package Structure**: Maintain internal package separation (mcp/, gorev/)
4. **MCP Protocol**: Follow patterns in internal/mcp/handlers.go
5. **Turkish Terms**: Keep domain terms in Turkish per CLAUDE.md "Code Style"
6. **Error Handling**: Follow CLAUDE.md "Error Handling" section

## Response Format

Start by confirming the task you were given as "**$ARGUMENTS**" and presenting your initial analysis:

```
📋 Assigned Task: $ARGUMENTS

🎯 Task Understanding:
[My interpretation of what needs to be accomplished]

📍 Current State Assessment:
- Existing implementation: [what exists]
- Gap analysis: [what's missing]
- Related packages: [affected areas]

🔍 Technical Analysis:
- Architecture impact: [how it fits/changes current architecture]
- Dependencies: [required Go packages/modules]
- Complexity assessment: [simple/medium/complex with reasons]

❓ Questions Before Implementation:
1. [Specific clarification about requirement]
2. [Technical decision that needs agreement]
3. [Potential approach validation]

💡 Proposed Implementation Approach:
Step 1: [Description with rationale]
Step 2: [Description with rationale]
Step 3: [Continue...]

⚠️ Risk Assessment:
- Risk: [Description] → Mitigation: [Strategy]
- Risk: [Description] → Mitigation: [Strategy]

🧪 Testing Strategy:
- Unit tests: [what will be tested]
- Integration tests: [scope]
- Benchmarks: [if performance critical]

📚 Documentation Impact:
- Files to update: [list of docs]
- New documentation needed: [if any]

⏱️ Estimated Effort:
[Time estimate with breakdown]

Awaiting your confirmation and answers to proceed with implementation.
```

## Additional Context for Common Task Types

### New MCP Tool Implementation
- Follow CLAUDE.md "Adding New MCP Tools" section:
  1. Add handler method to internal/mcp/handlers.go
  2. Register tool in RegisterTools() with proper schema
  3. Add integration tests in test/integration_test.go
  4. Update docs/mcp-araclari.md with tool documentation

### Database Schema Changes
- Update tablolariOlustur() in veri_yonetici.go
- Follow CLAUDE.md "Database Schema" section
- Test with both new and existing data
- Update CLAUDE.md schema documentation

### Performance Optimization
- Benchmark before and after
- Use pprof for profiling
- Consider goroutine usage
- Document improvements

### Bug Fix
- Write failing test first
- Fix with minimal changes
- Verify no regressions
- Update relevant tests

## Post-Implementation Checklist

After implementation, ensure:
- [ ] All tests pass (`make test`)
- [ ] No race conditions (`go test -race ./...`)
- [ ] Code formatted (`make fmt`)
- [ ] Lint checks pass (`make lint`)
- [ ] Vet checks pass (`go vet ./...`)
- [ ] CLAUDE.md updated if architecture changes
- [ ] Documentation updated per CLAUDE.md standards
- [ ] TodoWrite status updated
- [ ] Ready for commit

Remember: Quality and maintainability over speed. When in doubt, ask for clarification.