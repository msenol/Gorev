# Interactive Demos and Examples

**Version:** v0.16.3 | **Last Updated:** October 6, 2025

Practical, copy-paste-ready examples for common Gorev workflows.

---

## Quick Start Demo

### 1. Install and Initialize

```bash
# Install Gorev MCP server
npm install -g @mehmetsenol/gorev-mcp-server

# Start daemon in background
npx gorev daemon --detach

# Verify daemon is running
npx gorev health
# Output: ✓ Daemon is healthy (PID: 12345, URL: http://localhost:5082)

# Initialize workspace database
npx gorev init
# Output: ✓ Workspace initialized at /home/user/project/.gorev/gorev.db
```

### 2. Create Your First Task

```bash
# Using bug template
npx gorev create --template bug \
  --hata-aciklama "Login button doesn't respond on mobile Safari" \
  --oncelik yuksek

# Using feature template
npx gorev create --template feature \
  --ozellik-adi "Dark mode toggle" \
  --aciklama "Add dark/light theme switcher to settings page" \
  --oncelik orta
```

### 3. List and Search Tasks

```bash
# List all active tasks
npx gorev list --durum devam_ediyor

# Search by keyword
npx gorev search "login"

# Advanced search with filters
npx gorev search --action advanced \
  --query "API" \
  --durum devam_ediyor \
  --oncelik yuksek
```

---

## MCP Tool Usage Examples

### Creating Tasks with Templates

**Scenario:** Bug report workflow

```python
# Python MCP Client
async def report_bug(title, description, priority="high"):
    result = await session.call_tool("templateden_gorev_olustur", {
        "template_id": "bug",
        "degerler": {
            "hata_aciklama": description,
            "oncelik": priority
        }
    })
    return result

# Usage
task = await report_bug(
    "API timeout on /users endpoint",
    "GET /api/users returns 504 after 30 seconds",
    "yuksek"
)
```

**Scenario:** Feature request workflow

```typescript
// TypeScript MCP Client
async function createFeature(name: string, description: string) {
  const result = await client.callTool("templateden_gorev_olustur", {
    template_id: "feature",
    degerler: {
      ozellik_adi: name,
      aciklama: description,
      oncelik: "orta"
    }
  });
  return JSON.parse(result.content[0].text);
}

// Usage
const task = await createFeature(
  "Email notifications",
  "Send email when task status changes"
);
console.log(`Created task: ${task.id}`);
```

---

### Task Hierarchies and Dependencies

**Scenario:** Break down a large feature into subtasks

```bash
# 1. Create parent task
PARENT_ID=$(npx gorev create --template feature \
  --ozellik-adi "User authentication system" \
  --aciklama "Complete auth implementation" \
  --oncelik yuksek \
  | grep -oP 'ID: \K[a-z0-9]+')

# 2. Create subtasks
SUB1=$(npx gorev create --template feature \
  --ozellik-adi "Login API endpoint" \
  --parent-id $PARENT_ID \
  | grep -oP 'ID: \K[a-z0-9]+')

SUB2=$(npx gorev create --template feature \
  --ozellik-adi "Password hashing" \
  --parent-id $PARENT_ID \
  | grep -oP 'ID: \K[a-z0-9]+')

SUB3=$(npx gorev create --template feature \
  --ozellik-adi "JWT token generation" \
  --parent-id $PARENT_ID \
  | grep -oP 'ID: \K[a-z0-9]+')

# 3. Add dependencies
npx gorev add-dependency $SUB3 --depends-on $SUB2
# JWT generation depends on password hashing being done first

# 4. View hierarchy
npx gorev hierarchy $PARENT_ID
```

**MCP equivalent:**

```python
# Create task hierarchy via MCP
async def create_auth_system():
    # Parent task
    parent = await session.call_tool("templateden_gorev_olustur", {
        "template_id": "feature",
        "degerler": {
            "ozellik_adi": "User authentication system",
            "aciklama": "Complete auth implementation",
            "oncelik": "yuksek"
        }
    })
    parent_id = json.loads(parent.content[0].text)["id"]

    # Subtasks
    subtasks = []
    for name in ["Login API endpoint", "Password hashing", "JWT token generation"]:
        sub = await session.call_tool("gorev_altgorev_olustur", {
            "ust_gorev_id": parent_id,
            "baslik": name,
            "oncelik": "orta"
        })
        subtasks.append(json.loads(sub.content[0].text)["id"])

    # Add dependency: JWT depends on password hashing
    await session.call_tool("gorev_bagimlilik_ekle", {
        "gorev_id": subtasks[2],
        "bagli_gorev_id": subtasks[1]
    })

    return parent_id, subtasks
```

---

### Bulk Operations

**Scenario:** Update multiple tasks at once

```bash
# Change priority for all "API" related tasks
TASK_IDS=$(npx gorev search "API" --format ids)

npx gorev bulk update \
  --ids $TASK_IDS \
  --data '{"oncelik": "yuksek"}'
```

**MCP equivalent:**

```typescript
// Bulk update via MCP
async function escalateApiTasks() {
  // Find all API-related tasks
  const search = await client.callTool("gorev_search_advanced", {
    query: "API",
    durum: "devam_ediyor"
  });

  const tasks = JSON.parse(search.content[0].text);
  const ids = tasks.map(t => t.id);

  // Bulk update to high priority
  const result = await client.callTool("gorev_bulk", {
    operation: "update",
    ids: ids,
    data: {
      oncelik: "yuksek"
    }
  });

  return JSON.parse(result.content[0].text);
}
```

---

## Real-World Workflows

### Workflow 1: Sprint Planning

```bash
#!/bin/bash
# sprint-planning.sh - Create tasks for new sprint

SPRINT_NAME="Sprint 23 - Q4 2025"
SPRINT_TAG="sprint-23"

# Create sprint project
npx gorev project create --name "$SPRINT_NAME"
npx gorev project set-active "$SPRINT_NAME"

# High priority features
npx gorev create --template feature \
  --ozellik-adi "Payment gateway integration" \
  --oncelik yuksek \
  --etiket $SPRINT_TAG

npx gorev create --template feature \
  --ozellik-adi "User profile settings page" \
  --oncelik orta \
  --etiket $SPRINT_TAG

# Bug fixes from backlog
npx gorev create --template bug \
  --hata-aciklama "Memory leak in WebSocket connection" \
  --oncelik yuksek \
  --etiket $SPRINT_TAG

# List sprint tasks
echo "Sprint tasks:"
npx gorev list --etiket $SPRINT_TAG
```

### Workflow 2: Daily Standup Report

```bash
#!/bin/bash
# standup-report.sh - Generate daily standup report

echo "🌅 Daily Standup Report - $(date +%Y-%m-%d)"
echo ""

echo "✅ Completed Yesterday:"
npx gorev list --durum tamamlandi \
  --updated-after "yesterday" \
  --format summary

echo ""
echo "🏃 In Progress:"
npx gorev list --durum devam_ediyor \
  --format summary

echo ""
echo "🔴 Blocked Tasks:"
npx gorev list --etiket blocked \
  --format summary

echo ""
echo "📊 Summary:"
npx gorev summary --period today
```

### Workflow 3: Code Review Task Creation

```bash
#!/bin/bash
# create-review-task.sh - Create task from GitHub PR

PR_NUMBER=$1
PR_TITLE=$(gh pr view $PR_NUMBER --json title -q .title)
PR_URL=$(gh pr view $PR_NUMBER --json url -q .url)
PR_AUTHOR=$(gh pr view $PR_NUMBER --json author -q .author.login)

# Create review task
npx gorev create --template feature \
  --ozellik-adi "Review: $PR_TITLE" \
  --aciklama "PR #$PR_NUMBER by @$PR_AUTHOR
URL: $PR_URL

Review checklist:
- [ ] Code quality and style
- [ ] Test coverage
- [ ] Documentation updated
- [ ] Breaking changes documented
- [ ] Security considerations" \
  --oncelik orta \
  --etiket "code-review,pr-$PR_NUMBER"

echo "✓ Review task created for PR #$PR_NUMBER"
```

---

## Claude Desktop Integration Examples

### Example 1: AI-Assisted Task Breakdown

**User prompt to Claude:**
> "I need to implement OAuth2 authentication. Break this down into subtasks."

**Claude's workflow:**

```python
# Claude internally uses MCP tools

# 1. Create parent task
parent = await call_tool("templateden_gorev_olustur", {
    "template_id": "feature",
    "degerler": {
        "ozellik_adi": "OAuth2 Authentication",
        "aciklama": "Implement complete OAuth2 flow",
        "oncelik": "yuksek"
    }
})

# 2. Create logical subtasks
subtasks = [
    "OAuth2 provider configuration (Google, GitHub)",
    "Authorization endpoint implementation",
    "Token exchange endpoint",
    "User session management",
    "Token refresh mechanism",
    "Logout and token revocation"
]

for task in subtasks:
    await call_tool("gorev_altgorev_olustur", {
        "ust_gorev_id": parent["id"],
        "baslik": task,
        "oncelik": "orta"
    })

# 3. Show hierarchy to user
hierarchy = await call_tool("gorev_hiyerarsi_goster", {
    "gorev_id": parent["id"]
})
```

### Example 2: Context-Aware Task Creation

**User prompt to Claude:**
> "I'm working on the login page. Create a task for adding 'forgot password' functionality."

**Claude's workflow:**

```python
# Claude analyzes current context and creates appropriate task

# 1. Check active project
active = await call_tool("aktif_proje_goster")

# 2. Search for related tasks
related = await call_tool("gorev_search_advanced", {
    "query": "login",
    "durum": "devam_ediyor"
})

# 3. Create new task with context
task = await call_tool("templateden_gorev_olustur", {
    "template_id": "feature",
    "degerler": {
        "ozellik_adi": "Forgot password functionality",
        "aciklama": """Add password reset flow:
        1. Email input form on login page
        2. Send reset token via email
        3. Token validation endpoint
        4. Password update form
        5. Success confirmation""",
        "oncelik": "orta"
    }
})

# 4. Suggest linking to related tasks
if related["results"]:
    print(f"Related task found: {related['results'][0]['baslik']}")
    print(f"Consider adding dependency or parent-child relationship")
```

---

## VS Code Extension Examples

### Example 1: File Watcher Auto-Task

```typescript
// Configure file watcher in VS Code settings
{
  "gorev.fileWatcher.patterns": [
    "**/*.todo.md"
  ],
  "gorev.fileWatcher.autoCreate": true
}
```

**Create task from file:**

```markdown
<!-- tasks/implement-search.todo.md -->
# Implement Advanced Search

Priority: High
Tags: search, backend

## Description
Add full-text search with filters for:
- Task title and description
- Tags
- Priority levels
- Date ranges

## Acceptance Criteria
- [ ] FTS5 index created
- [ ] REST API endpoint /api/search
- [ ] MCP tool gorev_search_advanced
- [ ] Unit tests with 80%+ coverage
```

When saved, VS Code extension automatically creates a task using the template matching the file structure.

### Example 2: Quick Task Creation from Selection

```typescript
// Select code in editor, right-click → "Create Gorev Task"
// Extension creates task with code context

const selectedText = editor.document.getText(editor.selection);
const fileName = editor.document.fileName;
const lineNumber = editor.selection.start.line + 1;

await gorevClient.createTask({
  template_id: "bug",
  degerler: {
    hata_aciklama: `Issue in ${fileName}:${lineNumber}\n\n${selectedText}`,
    oncelik: "orta"
  }
});
```

---

## REST API Integration Examples

### cURL Examples

**Create task:**

```bash
curl -X POST http://localhost:5082/api/gorev \
  -H "Content-Type: application/json" \
  -d '{
    "baslik": "Fix API timeout issue",
    "aciklama": "Users endpoint times out after 30s",
    "oncelik": "yuksek",
    "durum": "devam_ediyor"
  }'
```

**List tasks with filters:**

```bash
curl -X GET "http://localhost:5082/api/gorevler?durum=devam_ediyor&oncelik=yuksek"
```

**Update task:**

```bash
curl -X PUT http://localhost:5082/api/gorev/abc123 \
  -H "Content-Type: application/json" \
  -d '{
    "durum": "tamamlandi",
    "tamamlanma_zamani": "2025-10-06T10:30:00Z"
  }'
```

### Postman Collection

Import this JSON into Postman for a complete API test suite:

```json
{
  "info": {
    "name": "Gorev API Collection",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "header": [],
        "url": "{{base_url}}/api/health"
      }
    },
    {
      "name": "List Tasks",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/gorevler?limit=10",
          "host": ["{{base_url}}"],
          "path": ["api", "gorevler"],
          "query": [
            {"key": "limit", "value": "10"},
            {"key": "durum", "value": "devam_ediyor", "disabled": true},
            {"key": "oncelik", "value": "yuksek", "disabled": true}
          ]
        }
      }
    },
    {
      "name": "Create Task",
      "request": {
        "method": "POST",
        "header": [{"key": "Content-Type", "value": "application/json"}],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"baslik\": \"New task\",\n  \"aciklama\": \"Task description\",\n  \"oncelik\": \"orta\"\n}"
        },
        "url": "{{base_url}}/api/gorev"
      }
    }
  ],
  "variable": [
    {"key": "base_url", "value": "http://localhost:5082"}
  ]
}
```

---

## WebSocket Real-Time Updates Example

```javascript
// Connect to WebSocket for real-time task updates
const ws = new WebSocket('ws://localhost:5082/ws');

ws.onopen = () => {
  console.log('✓ Connected to Gorev WebSocket');
};

ws.onmessage = (event) => {
  const update = JSON.parse(event.data);

  switch (update.type) {
    case 'task_created':
      console.log('New task:', update.data.baslik);
      notifyUser(`New task: ${update.data.baslik}`);
      refreshTaskList();
      break;

    case 'task_updated':
      console.log('Task updated:', update.data.id);
      updateTaskInUI(update.data);
      break;

    case 'task_deleted':
      console.log('Task deleted:', update.data.id);
      removeTaskFromUI(update.data.id);
      break;

    default:
      console.log('Unknown event:', update.type);
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket connection closed');
  // Attempt reconnection
  setTimeout(() => connectWebSocket(), 5000);
};
```

---

## Performance Testing Examples

### Load Testing with k6

```javascript
// load-test.js - Test Gorev API performance
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 10 },   // Ramp up to 10 users
    { duration: '1m', target: 50 },    // Ramp up to 50 users
    { duration: '30s', target: 0 },    // Ramp down
  ],
};

export default function () {
  // List tasks
  let listRes = http.get('http://localhost:5082/api/gorevler?limit=50');
  check(listRes, {
    'list tasks status 200': (r) => r.status === 200,
    'list tasks duration < 50ms': (r) => r.timings.duration < 50,
  });

  // Create task
  let createRes = http.post(
    'http://localhost:5082/api/gorev',
    JSON.stringify({
      baslik: `Load test task ${Date.now()}`,
      aciklama: 'Generated by k6 load test',
      oncelik: 'orta',
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(createRes, {
    'create task status 201': (r) => r.status === 201,
    'create task duration < 20ms': (r) => r.timings.duration < 20,
  });

  sleep(1);
}
```

**Run load test:**

```bash
k6 run load-test.js
```

---

## Troubleshooting Examples

### Example 1: Daemon Not Starting

```bash
# Check if daemon is already running
npx gorev health

# If stale lock file exists
rm -f ~/.gorev-daemon/.lock

# Start with debug logging
npx gorev daemon --debug --log-file /tmp/gorev-daemon.log

# Check logs
tail -f /tmp/gorev-daemon.log
```

### Example 2: Database Corruption Recovery

```bash
# Backup existing database
cp .gorev/gorev.db .gorev/gorev.db.backup

# Check database integrity
sqlite3 .gorev/gorev.db "PRAGMA integrity_check;"

# If corrupted, export and reimport
npx gorev export --output tasks-backup.json
npx gorev init --force
npx gorev import tasks-backup.json
```

### Example 3: MCP Connection Issues

```python
# Debug MCP connection with verbose logging
import logging
logging.basicConfig(level=logging.DEBUG)

async with stdio_client(server_params) as (read, write):
    # Check connection
    print("✓ stdio connection established")

    async with ClientSession(read, write) as session:
        # Initialize with timeout
        await asyncio.wait_for(session.initialize(), timeout=10.0)
        print("✓ MCP session initialized")

        # List tools to verify
        tools = await session.list_tools()
        print(f"✓ Available tools: {len(tools.tools)}")
```

---

## See Also

- [API Integration Guide](../api/integration-examples.md) - Complete API reference
- [MCP Tools Reference](../MCP_TOOLS_REFERENCE.md) - All 24 MCP tools
- [Architecture Diagrams](../architecture/sequence-diagrams.md) - System flow diagrams
- [CLI Reference](../cli-reference.md) - Command-line usage

---

**Last Updated:** October 6, 2025 | **Version:** v0.16.3
