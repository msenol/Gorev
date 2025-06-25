# Start Task - Systematic Task Implementation

## Initial Setup
First read CLAUDE.md to understand the project architecture and conventions. Then check the TodoRead tool for the current task list and pick the next pending task. If no tasks exist, check if there are any issues or feature requests to work on.

## Pre-Implementation Analysis

### 1. Deep Context Understanding
Before implementing:
- **THINK DEEPLY** about the task in context of existing Go architecture
- **ANALYZE** the relevant code base and package structure
- Check MCP SDK documentation for best practices
- Review related packages and their interactions
- Identify potential side effects or dependencies

### 2. Go/MCP Specific Requirements
If the task involves:
- **MCP Tools**: Follow the established handler pattern in `internal/mcp/handlers.go`
- **Database**: Use prepared statements and proper transaction handling
- **Error Handling**: Return explicit errors, never panic in library code
- **Testing**: Write table-driven tests following Go conventions
- **Concurrency**: Use goroutines and channels appropriately

### 3. Critical Questions to Ask
**ASK questions** if you find:
- Multiple valid approaches (present pros/cons of each)
- Unclear requirements or acceptance criteria
- Potential architectural impacts or breaking changes
- Missing dependencies or prerequisites
- Conflicts with existing patterns or standards
- Performance or scalability concerns
- Security implications

### 4. Proposal Presentation
Present your understanding including:
- Task summary and objectives
- Current state analysis
- Proposed implementation approach
- Step-by-step execution plan
- Testing strategy
- Documentation updates needed
- Estimated complexity and risks

### 5. Wait for Agreement
- Present complete analysis and approach
- Wait for explicit agreement before proceeding
- Be ready to adjust based on feedback
- Confirm understanding of any clarifications

## Implementation Rules

### Absolute Rules
- **NO shortcuts** ("for now", "temporarily", "to save time")
- **STOP** if tempted to change code outside current task scope
- **One task at a time** - complete focus on selected task
- **Test after each step** - all tests must pass
- **Update TodoWrite** when task status changes

### Gorev Project Specific Rules
1. **Go Standards**: Run `gofmt` and `go vet` before committing
2. **Package Structure**: Maintain clean separation between internal packages
3. **MCP Compliance**: All tools must follow MCP protocol specifications
4. **Turkish Naming**: Keep Turkish terms for domain concepts (gorev, proje, etc.)
5. **Error Messages**: Provide clear, actionable error messages in Turkish
6. **SQLite**: Use transactions for multi-statement operations

### Code Quality Standards
- Follow standards defined in CLAUDE.md "Code Style" section
- Use Turkish terms for domain concepts (gorev, proje, durum, oncelik)
- Keep English for technical terms and comments
- Handle all errors explicitly with context
- Add comments for exported functions
- Use context.Context for cancellation

## Task Execution Framework

### 1. Implementation Steps
```
1. Create/modify files following project structure
2. Implement core logic with proper error handling
3. Add necessary types/interfaces
4. Write unit tests (table-driven)
5. Implement integration tests if needed
6. Run go mod tidy if dependencies changed
7. Update relevant documentation
```

### 2. Testing Checklist
- [ ] Unit tests pass (`make test`)
- [ ] Integration tests pass
- [ ] No race conditions (`go test -race ./...`)
- [ ] Code formatted (`make fmt`)
- [ ] Vet checks pass (`go vet ./...`)
- [ ] Lint checks pass (`make lint`)
- [ ] Coverage acceptable (`make test-coverage`)

### 3. Documentation Updates
- [ ] Update docs/ if features changed
- [ ] Update README.md if setup/usage changed
- [ ] Update mcp-araclari.md if tools added/modified
- [ ] Add/update code comments
- [ ] Update CHANGELOG if needed

## Output Format

Start by providing:

```
📋 Selected Task: [From TodoRead or identified work]

🎯 Task Objective:
[Clear description of what needs to be achieved]

🔍 Current State Analysis:
- [What exists currently]
- [What's missing or needs change]
- [Related packages/modules]

❓ Initial Questions:
1. [Specific question about requirement/approach]
2. [Technical clarification needed]
3. [Architecture/design decision needed]

💡 Proposed Approach:
1. [Step 1 with rationale]
2. [Step 2 with rationale]
3. [Continue...]

⚠️ Potential Risks/Concerns:
- [Risk 1 and mitigation strategy]
- [Risk 2 and mitigation strategy]

🔗 Dependencies:
- [Go packages needed]
- [Other modules affected]
- [Prerequisites]

Ready to proceed after your confirmation and answers to my questions.
```

## Remember
- Quality over speed - no temporary solutions
- Ask when uncertain rather than assume
- Follow established Go patterns
- Keep changes focused and minimal
- Document decisions and rationale
- Use TodoWrite to track progress