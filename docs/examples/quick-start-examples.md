# Quick Start Examples

**Version:** v0.16.3 | **Last Updated:** October 6, 2025

5-minute quick start examples for common Gorev use cases.

---

## Example 1: Personal Task Management (5 min)

**Goal:** Set up personal to-do list with Gorev

```bash
# Install
npm install -g @mehmetsenol/gorev-mcp-server

# Start daemon
npx gorev daemon --detach

# Initialize workspace
cd ~/my-project
npx gorev init

# Create personal project
npx gorev project create --name "Personal Tasks"
npx gorev project set-active "Personal Tasks"

# Add tasks
npx gorev create --template feature \
  --ozellik-adi "Learn Gorev" \
  --oncelik yuksek

npx gorev create --template feature \
  --ozellik-adi "Write blog post" \
  --oncelik orta

npx gorev create --template feature \
  --ozellik-adi "Review pull requests" \
  --oncelik yuksek

# View tasks
npx gorev list

# Mark one complete
npx gorev update <task-id> --durum tamamlandi

# View summary
npx gorev summary
```

**Result:** Personal task tracker with 3 tasks, ready to use.

---

## Example 2: Claude Desktop Integration (3 min)

**Goal:** Use Gorev with Claude for AI-assisted task management

### 1. Configure Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or
`%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server@latest", "mcp-proxy"],
      "env": {
        "GOREV_LANG": "en"
      }
    }
  }
}
```

### 2. Restart Claude Desktop

### 3. Start Using

**Talk to Claude:**

> "Create a new task for implementing user authentication"

**Claude will:**
1. Use `templateden_gorev_olustur` tool
2. Select appropriate template (`feature`)
3. Create task with proper structure
4. Show you the task details

**Advanced usage:**

> "Show me all high-priority tasks in progress, then create subtasks for the authentication feature"

**Result:** AI-powered task management with natural language.

---

## Example 3: Team Project Setup (10 min)

**Goal:** Set up Gorev for a team development project

### 1. Project Structure

```bash
cd ~/projects/team-project

# Initialize Gorev
npx gorev init

# Create project
npx gorev project create --name "E-commerce Platform v2.0"
npx gorev project set-active "E-commerce Platform v2.0"
```

### 2. Create Sprint Structure

```bash
# Create epic tasks
EPIC_AUTH=$(npx gorev create --template feature \
  --ozellik-adi "User Authentication System" \
  --oncelik yuksek \
  | grep -oP 'ID: \K[a-z0-9]+')

EPIC_PAYMENT=$(npx gorev create --template feature \
  --ozellik-adi "Payment Gateway Integration" \
  --oncelik yuksek \
  | grep -oP 'ID: \K[a-z0-9]+')

EPIC_ADMIN=$(npx gorev create --template feature \
  --ozellik-adi "Admin Dashboard" \
  --oncelik orta \
  | grep -oP 'ID: \K[a-z0-9]+')

# Create subtasks for Auth epic
npx gorev create --template feature \
  --ozellik-adi "Login API endpoint" \
  --parent-id $EPIC_AUTH

npx gorev create --template feature \
  --ozellik-adi "JWT token management" \
  --parent-id $EPIC_AUTH

npx gorev create --template feature \
  --ozellik-adi "Password reset flow" \
  --parent-id $EPIC_AUTH

# View hierarchy
npx gorev hierarchy $EPIC_AUTH
```

### 3. Share Database with Team

```bash
# Option 1: Git (for small teams)
# Add .gorev/ to git
git add .gorev/
git commit -m "chore: add Gorev task database"
git push

# Team members:
git pull
npx gorev daemon --detach  # Each person runs their own daemon

# Option 2: Shared network drive
# Configure workspace path
export GOREV_WORKSPACE_PATH="/shared/team-project/.gorev"
npx gorev daemon --detach
```

**Result:** Team-wide task tracking with hierarchical structure.

---

## Example 4: Bug Tracking Workflow (5 min)

**Goal:** Track and manage bug reports efficiently

### 1. Create Bug Template Tasks

```bash
# High priority production bug
npx gorev create --template bug \
  --hata-aciklama "Payment processing fails for amounts > $1000" \
  --oncelik yuksek \
  --etiket "production,critical,payment"

# UI bug
npx gorev create --template bug \
  --hata-aciklama "Mobile menu doesn't close after navigation" \
  --oncelik orta \
  --etiket "ui,mobile"

# Performance issue
npx gorev create --template bug \
  --hata-aciklama "Dashboard loads slow with 1000+ users" \
  --oncelik orta \
  --etiket "performance,backend"
```

### 2. Track Bug Lifecycle

```bash
# Developer picks up bug
BUG_ID="abc123"
npx gorev update $BUG_ID --durum devam_ediyor

# Add investigation notes
npx gorev update $BUG_ID --aciklama "Found issue in payment processor timeout setting. Default is 10s, needs to be 30s for large transactions."

# Mark as fixed
npx gorev update $BUG_ID --durum tamamlandi

# Generate bug report
npx gorev list --etiket production --durum tamamlandi
```

**Result:** Structured bug tracking with clear workflow.

---

## Example 5: VS Code Extension Workflow (3 min)

**Goal:** Visual task management in VS Code

### 1. Install Extension

```bash
# In VS Code
# Extensions → Search "Gorev" → Install
```

### 2. Configure Auto-Start

Extension automatically starts daemon on activation. No manual setup needed.

### 3. Use Tree View

- **View → Gorev Tasks** - Opens task tree
- **Right-click task** → Context menu:
  - Mark as In Progress
  - Mark as Complete
  - Create Subtask
  - Delete Task
  - Show Details

### 4. Create Task from File

```typescript
// Select code in editor
function calculateTotal(items: Item[]): number {
  // TODO: Add tax calculation
  return items.reduce((sum, item) => sum + item.price, 0);
}

// Right-click selection → "Create Gorev Task from Selection"
```

**Result:** Visual task management integrated into VS Code.

---

## Example 6: CI/CD Integration (10 min)

**Goal:** Automate task creation from CI/CD pipeline

### GitHub Actions Example

```yaml
# .github/workflows/create-task-on-failure.yml
name: Create Task on Test Failure

on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run tests
        id: tests
        run: npm test
        continue-on-error: true

      - name: Install Gorev
        if: steps.tests.outcome == 'failure'
        run: npm install -g @mehmetsenol/gorev-mcp-server

      - name: Start Gorev daemon
        if: steps.tests.outcome == 'failure'
        run: npx gorev daemon --detach

      - name: Create bug task
        if: steps.tests.outcome == 'failure'
        run: |
          npx gorev create --template bug \
            --hata-aciklama "Test failure in ${{ github.workflow }}
          Branch: ${{ github.ref }}
          Commit: ${{ github.sha }}
          Triggered by: ${{ github.actor }}

          See: ${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}" \
            --oncelik yuksek \
            --etiket "ci-failure,automated"

      - name: Fail if tests failed
        if: steps.tests.outcome == 'failure'
        run: exit 1
```

### GitLab CI Example

```yaml
# .gitlab-ci.yml
test:
  script:
    - npm test || true
  after_script:
    - |
      if [ $? -ne 0 ]; then
        npm install -g @mehmetsenol/gorev-mcp-server
        npx gorev daemon --detach
        npx gorev create --template bug \
          --hata-aciklama "Pipeline failure: $CI_PIPELINE_URL" \
          --oncelik yuksek \
          --etiket "ci-failure"
      fi
```

**Result:** Automatic bug creation on test failures.

---

## Example 7: Data Export/Import (5 min)

**Goal:** Backup and migrate tasks between projects

### Export Tasks

```bash
# Export all tasks
npx gorev export --output backup.json

# Export specific project
npx gorev export \
  --project "E-commerce Platform" \
  --output ecommerce-tasks.json

# Export by filter
npx gorev export \
  --durum tamamlandi \
  --created-after "2025-01-01" \
  --output completed-2025.json
```

### Import Tasks

```bash
# Import to new project
cd ~/new-project
npx gorev init
npx gorev import backup.json

# Import with merge
npx gorev import tasks.json --merge
```

### Convert Format

```bash
# Export as CSV
npx gorev export --output tasks.csv --format csv

# Export as Markdown
npx gorev export --output TASKS.md --format markdown
```

**Result:** Portable task data across projects and systems.

---

## Example 8: Search and Filtering (3 min)

**Goal:** Find tasks quickly with advanced search

### Basic Search

```bash
# Keyword search
npx gorev search "authentication"

# Search in specific fields
npx gorev search --field baslik "API"
npx gorev search --field aciklama "database"
```

### Advanced Filters

```bash
# Combine filters
npx gorev list \
  --durum devam_ediyor \
  --oncelik yuksek \
  --etiket backend \
  --created-after "2025-10-01"

# Date range search
npx gorev list \
  --created-after "2025-10-01" \
  --created-before "2025-10-07"

# Full-text search with FTS5
npx gorev search --action advanced \
  --query "payment gateway integration" \
  --oncelik yuksek
```

### Saved Filters

```bash
# Create filter profile
npx gorev filter create \
  --name "High Priority Active" \
  --filters '{"durum": "devam_ediyor", "oncelik": "yuksek"}'

# Use saved filter
npx gorev filter apply "High Priority Active"
```

**Result:** Fast task discovery with powerful search.

---

## Example 9: Python Scripting (5 min)

**Goal:** Automate task management with Python

```python
#!/usr/bin/env python3
# daily-report.py - Generate daily task report

import asyncio
import json
from datetime import datetime
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

async def generate_daily_report():
    server_params = StdioServerParameters(
        command="npx",
        args=["-y", "@mehmetsenol/gorev-mcp-server@latest", "mcp-proxy"],
        env={"GOREV_LANG": "en"}
    )

    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            await session.initialize()

            # Get tasks by status
            in_progress = await session.call_tool("gorev_listele", {
                "durum": "devam_ediyor"
            })
            completed = await session.call_tool("gorev_listele", {
                "durum": "tamamlandi",
                "tamamlanma_zamani_baslangic": datetime.now().isoformat()
            })

            # Generate report
            print("=" * 50)
            print(f"📊 Daily Report - {datetime.now().strftime('%Y-%m-%d')}")
            print("=" * 50)

            ip_tasks = json.loads(in_progress.content[0].text)
            print(f"\n🏃 In Progress ({len(ip_tasks)}):")
            for task in ip_tasks[:5]:
                print(f"  • {task['baslik']}")

            comp_tasks = json.loads(completed.content[0].text)
            print(f"\n✅ Completed Today ({len(comp_tasks)}):")
            for task in comp_tasks[:5]:
                print(f"  • {task['baslik']}")

if __name__ == "__main__":
    asyncio.run(generate_daily_report())
```

**Run:**

```bash
python3 daily-report.py
```

**Result:** Automated daily reports with Python.

---

## Example 10: TypeScript Dashboard (10 min)

**Goal:** Build a real-time task dashboard

```typescript
// dashboard.ts - Real-time task dashboard
import axios from 'axios';
import WebSocket from 'ws';

const API_URL = 'http://localhost:5082/api';
const WS_URL = 'ws://localhost:5082/ws';

interface Task {
  id: string;
  baslik: string;
  durum: string;
  oncelik: string;
}

class TaskDashboard {
  private ws: WebSocket;
  private tasks: Task[] = [];

  async start() {
    // Load initial tasks
    await this.loadTasks();
    this.displayTasks();

    // Connect WebSocket for real-time updates
    this.ws = new WebSocket(WS_URL);

    this.ws.on('open', () => {
      console.log('✓ Connected to Gorev WebSocket');
    });

    this.ws.on('message', (data) => {
      const update = JSON.parse(data.toString());
      this.handleUpdate(update);
    });
  }

  async loadTasks() {
    const response = await axios.get(`${API_URL}/gorevler`, {
      params: { durum: 'devam_ediyor' }
    });
    this.tasks = response.data;
  }

  handleUpdate(update: any) {
    switch (update.type) {
      case 'task_created':
        this.tasks.push(update.data);
        break;
      case 'task_updated':
        const idx = this.tasks.findIndex(t => t.id === update.data.id);
        if (idx >= 0) this.tasks[idx] = update.data;
        break;
      case 'task_deleted':
        this.tasks = this.tasks.filter(t => t.id !== update.data.id);
        break;
    }
    this.displayTasks();
  }

  displayTasks() {
    console.clear();
    console.log('📊 Task Dashboard - Real-time');
    console.log('='.repeat(50));

    const byPriority = {
      yuksek: this.tasks.filter(t => t.oncelik === 'yuksek'),
      orta: this.tasks.filter(t => t.oncelik === 'orta'),
      dusuk: this.tasks.filter(t => t.oncelik === 'dusuk'),
    };

    console.log(`\n🔴 High Priority (${byPriority.yuksek.length}):`);
    byPriority.yuksek.forEach(t => console.log(`  • ${t.baslik}`));

    console.log(`\n🟡 Medium Priority (${byPriority.orta.length}):`);
    byPriority.orta.forEach(t => console.log(`  • ${t.baslik}`));

    console.log(`\n🟢 Low Priority (${byPriority.dusuk.length}):`);
    byPriority.dusuk.forEach(t => console.log(`  • ${t.baslik}`));
  }
}

// Run
const dashboard = new TaskDashboard();
dashboard.start();
```

**Compile and run:**

```bash
npm install axios ws
npx tsx dashboard.ts
```

**Result:** Live updating task dashboard with WebSocket.

---

## Next Steps

After these quick start examples, explore:

- [**Full Interactive Demos**](./interactive-demos.md) - More complex workflows
- [**API Integration Guide**](../api/integration-examples.md) - Production-ready clients
- [**MCP Tools Reference**](../MCP_TOOLS_REFERENCE.md) - All 24 tools documented
- [**Architecture Docs**](../architecture/daemon-architecture.md) - System internals

---

**Last Updated:** October 6, 2025 | **Version:** v0.16.3
