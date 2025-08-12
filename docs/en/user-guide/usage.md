# Usage Guide

> **Version**: This documentation is valid for v0.11.0  
> **Last Updated**: August 13, 2025

Master task management with Gorev - from basic concepts to advanced workflows.

## 🎯 Core Concepts

### Task
- The fundamental unit of work in Gorev
- Each task has a unique ID for precise tracking
- **Status Options**: `pending`, `in_progress`, `completed`
- **Priority Levels**: `low`, `medium`, `high`
- **Lifecycle**: Created → Started → Completed

### Project
- Organizational container for grouping related tasks
- Optional but recommended for structured workflow
- Supports hierarchical task organization
- Enables project-level reporting and tracking

### Template System (v0.10.0+)
- **Mandatory** for task creation since v0.10.0
- Ensures consistency across similar task types
- Pre-defined fields and validation
- Streamlines repetitive task creation

## 🚀 Getting Started

### Server Management

**Start the server:**
```bash
# Standard mode
gorev serve

# Debug mode with verbose logging
gorev serve --debug

# Custom data directory
gorev serve --data-dir /path/to/data

# Specify port
gorev serve --port 8080
```

**Version check:**
```bash
gorev version
```

**Health check:**
```bash
curl http://localhost:3000/health
```

## 💬 AI Assistant Integration

### Task Creation (Template-Based)

> ⚠️ **BREAKING CHANGE (v0.10.0)**: Direct task creation via `gorev_olustur` is no longer available. All tasks must be created using templates.

**View available templates:**
```
Show me all available task templates
List templates by category
```

**Create from template:**
```
Create a new bug report task:
- Title: Login form validation error
- Module: Authentication
- Environment: production
- Priority: high
```

**Feature request template:**
```
Create a feature request task using the template:
- Title: Dark mode support
- Description: Add dark/light theme toggle
- Priority: medium
- Estimated hours: 8
```

### Task Management

**List tasks with filters:**
```
Show all pending tasks
List high priority tasks
Show tasks in progress
Display completed tasks from this week
```

**Update task status:**
```
Start working on task [task-id]
Complete task [task-id]
Set task [task-id] back to pending
```

**Search and filter:**
```
Find tasks containing "API"
Show tasks assigned to frontend project
List tasks with "bug" tag
```

### Project Operations

**Create and manage projects:**
```
Create a new project called "Mobile App Redesign"
Add tasks to the "Mobile App" project
Show project summary for "Backend API"
```

**Project reporting:**
```
Generate project status report
Show completion statistics by project
List projects with overdue tasks
```

## 📊 Advanced Workflows

### Sprint Planning Workflow

1. **Setup Sprint Project:**
   ```
   Create project "Sprint 15 - Q4 2024"
   Set sprint duration: 2 weeks
   Set sprint goal: "User authentication improvements"
   ```

2. **Backlog Planning:**
   ```
   Create tasks from user story template:
   - User login with OAuth
   - Password reset functionality  
   - Two-factor authentication
   - User profile management
   ```

3. **Task Prioritization:**
   ```
   Set OAuth task as high priority
   Set profile management as low priority
   Show tasks sorted by priority
   ```

4. **Sprint Execution:**
   ```
   Start first task: OAuth implementation
   Update daily progress on active tasks
   Move completed tasks to done
   ```

### Bug Tracking Workflow

**Bug Report Creation:**
```
Create bug report using template:
- Title: "Search results not loading"
- Module: Search Engine
- Environment: production
- Steps to reproduce: [detailed steps]
- Expected behavior: [description]
- Actual behavior: [description]
- Priority: high
```

**Bug Lifecycle Management:**
```
Assign bug [bug-id] to developer
Start investigation on bug [bug-id]
Add fix details to bug [bug-id]
Complete bug [bug-id] and mark as resolved
```

### Feature Development Workflow

**Feature Planning:**
```
Create feature request: "Real-time notifications"
Break down into subtasks:
- Backend WebSocket implementation
- Frontend notification UI
- Push notification service
- User preference settings
```

**Development Tracking:**
```
Start backend WebSocket task
Update progress: "Authentication layer completed"
Block task on external dependency
Complete task and update documentation
```

## 🎨 Best Practices

### 1. Effective Task Titles
- ❌ "Fix bug"
- ✅ "Fix email validation error in user registration form"

- ❌ "Add feature"
- ✅ "Implement dark mode toggle with user preference persistence"

### 2. Comprehensive Descriptions
Include essential context:
- **Background**: Why this task is needed
- **Acceptance Criteria**: Definition of done
- **Resources**: Links to designs, specs, or related tasks
- **Constraints**: Technical limitations or requirements

### 3. Priority Management Strategy
- **High**: Urgent + Important (production bugs, critical features, security issues)
- **Medium**: Important but not urgent (new features, performance improvements)
- **Low**: Neither urgent nor important (nice-to-have features, refactoring)

### 4. Status Management Guidelines
- Keep only 1-3 tasks in `in_progress` status simultaneously
- Break large tasks into smaller, manageable subtasks (2-8 hours each)
- Update status regularly to maintain visibility
- Review completed tasks weekly for lessons learned

## 🔧 Power User Features

### Template Management

**View template details:**
```
Show details for bug report template
List all fields for feature request template
Display template usage statistics
```

**Template-based task creation:**
```
Create task from template [template-id] with:
- field1: value1
- field2: value2
- field3: value3
```

### Advanced Filtering and Search

**Complex queries:**
```
Show high priority tasks created this month
Find tasks with "API" in title and "backend" tag
List overdue tasks by project
Display tasks modified in the last 7 days
```

**Bulk operations:**
```
Update all pending backend tasks to medium priority
Add "refactoring" tag to all code cleanup tasks
Archive all completed tasks older than 3 months
```

### Dependency Management

**Task dependencies:**
```
Make task B dependent on task A
Show dependency chain for task [task-id]
Find tasks blocked by dependencies
List tasks ready to start (no pending dependencies)
```

### Time Tracking and Estimation

**Time management:**
```
Set estimated time for task: 4 hours
Log 2 hours of work on task [task-id]
Show time tracking report for this week
Compare estimated vs actual time by project
```

## 📈 Reporting and Analytics

### Daily Workflows

**Morning standup preparation:**
```
Show my tasks completed yesterday
List my tasks for today
Display any blocked or overdue tasks
Generate daily progress summary
```

**End of day review:**
```
Update progress on active tasks
Log time spent on completed tasks
Plan tomorrow's priorities
Review any new task assignments
```

### Weekly Reviews

**Team performance:**
```
Generate weekly completion report
Show productivity trends by team member
List missed deadlines and their causes
Display project progress summaries
```

**Process improvement:**
```
Show most common task types this week
Identify bottlenecks in workflow
Review estimation accuracy
Analyze task cycle time trends
```

## 🔍 Common Patterns

### Quick Task Creation
```
Quick bug report for login issue:
- Environment: production
- Priority: high
- Description: Users can't login with valid credentials

Quick feature request:
- Add export functionality to reports
- Priority: medium
- Estimated: 6 hours
```

### Batch Operations
```
Create 5 testing tasks for the new feature:
- Unit tests for authentication
- Integration tests for API endpoints
- UI tests for user flows
- Performance tests for database queries
- Security tests for data validation
```

### Status Updates
```
Daily standup update:
- Completed: User authentication API
- In progress: OAuth integration (80% done)
- Blocked: Waiting for design approval
- Next: Start password reset functionality
```

## ❓ Troubleshooting

### Common Issues

**Finding task IDs:**
```
Search for tasks containing "authentication"
Show task ID for "OAuth implementation" task
List recent tasks with their IDs
```

**Status confusion:**
```
Show current status of all my tasks
List tasks that haven't been updated in 7 days
Find tasks in unclear or invalid states
```

**Template problems:**
```
Show available templates if task creation fails
Validate template field requirements
Check template permissions and access
```

### Data Management

**Backup and recovery:**
```bash
# Create backup
cp ~/.gorev/data/gorev.db ~/.gorev/data/gorev-backup-$(date +%Y%m%d).db

# Restore from backup
cp ~/.gorev/data/gorev-backup-20240813.db ~/.gorev/data/gorev.db
```

**Database maintenance:**
```bash
# Check database integrity
sqlite3 ~/.gorev/data/gorev.db "PRAGMA integrity_check;"

# Optimize database
sqlite3 ~/.gorev/data/gorev.db "VACUUM;"
```

## 🚀 Next Steps

- **[MCP Tools Reference](mcp-tools.md)** - Complete reference for all available tools and commands
- **[Installation Guide](../getting-started/installation.md)** - Platform-specific installation instructions
- **[API Reference](../api/reference.md)** - Technical API documentation for developers
- **[Troubleshooting Guide](troubleshooting.md)** - Solutions for common issues and problems

## 💡 Pro Tips

1. **Template Consistency**: Always use templates for similar task types to maintain data quality
2. **Regular Reviews**: Schedule weekly reviews to keep your task list current and relevant
3. **Dependency Planning**: Map out task dependencies before starting complex projects
4. **Time Boxing**: Set realistic time estimates and track actual time to improve planning
5. **Status Hygiene**: Keep task statuses current to maintain team visibility and coordination

---

*✨ This comprehensive usage guide helps you master Gorev's powerful task management capabilities - from basic task creation to advanced project workflows and team collaboration patterns.*