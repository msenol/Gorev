# Multi-Workspace Support

**Version**: v0.16.0
**Last Updated**: October 5, 2025
**Feature Status**: Production Ready ✅

---

## Overview

Gorev v0.16.0 introduces **multi-workspace support**, allowing you to manage separate task databases for different projects. Each workspace is completely isolated with its own `.gorev/gorev.db` SQLite database, enabling context-specific task management.

### Key Benefits

- ✅ **Workspace Isolation**: Each project has independent task database
- ✅ **Automatic Detection**: Gorev finds `.gorev/` directory in current folder
- ✅ **SHA256 Workspace IDs**: Secure, unique identification based on workspace path
- ✅ **Web UI Switcher**: Seamlessly switch between workspaces in browser
- ✅ **VS Code Integration**: Extension auto-registers workspace on activation
- ✅ **Context Propagation**: HTTP headers carry workspace context to all APIs

---

## Concepts

### Workspace

A **workspace** is a project directory containing:

- `.gorev/` directory (created by `gorev init`)
- `gorev.db` SQLite database file
- Optional MCP configuration (`.kilocode/mcp.json`, `.cursor/mcp.json`, etc.)

**Example**:

```
/projects/my-app/
├── .gorev/
│   └── gorev.db        # Workspace-specific database
├── .kilocode/
│   └── mcp.json        # Optional MCP config
├── src/
└── package.json
```

### Workspace ID

Each workspace is identified by a **SHA256 hash** of its absolute path:

```
Workspace Path:  /home/user/projects/my-app
SHA256 Hash:     4a5d7c9b8f3e2a1d6c4b3a9f8e7d6c5b4a3d2c1b9f8e7d6c5b4a3d2c1b9f8e7d6
Workspace ID:    4a5d7c9b (first 8 characters used for display)
```

**Benefits**:

- Deterministic (same path → same ID)
- Collision-resistant
- Platform-independent
- No manual configuration needed

### Workspace Registration

Workspaces must be **registered** before use. Registration stores metadata in a separate `workspace_registry.db`:

| Field | Description |
|-------|-------------|
| `workspace_id` | SHA256 hash (first 8 chars) |
| `name` | Human-readable name (from directory) |
| `path` | Absolute path to workspace |
| `db_path` | Full path to `.gorev/gorev.db` |
| `registered_at` | Registration timestamp |
| `last_accessed` | Last access timestamp |

**Registration Methods**:

1. Automatic: VS Code extension registers on activation
2. Manual: Web UI detects and registers on first access
3. CLI: `gorev init` creates and registers workspace

---

## Usage

### CLI Workflow

#### 1. Initialize Workspace

Navigate to your project directory and initialize:

```bash
cd /path/to/your/project
gorev init
```

**Output**:

```
✓ Created workspace directory: /path/to/your/project/.gorev
✓ Initialized database: gorev.db
✓ Workspace registered: my-project (ID: 4a5d7c9b)
```

#### 2. Verify Registration

List all registered workspaces:

```bash
gorev workspace list
```

**Output**:

```
Registered Workspaces:
  1. my-project (ID: 4a5d7c9b)
     Path: /path/to/your/project
     Tasks: 15 total (8 pending, 5 in progress, 2 completed)
     Last Accessed: 2025-10-05 10:30:00

  2. other-project (ID: 7f3a8c2d)
     Path: /path/to/other/project
     Tasks: 23 total (10 pending, 10 in progress, 3 completed)
     Last Accessed: 2025-10-04 15:20:00
```

#### 3. Switch Workspace

Change current working directory to switch workspace context:

```bash
cd /path/to/other/project
gorev task list  # Lists tasks from other-project's database
```

Gorev automatically detects the nearest `.gorev/` directory by walking up the directory tree.

### VS Code Extension Workflow

#### 1. Open Project

Open workspace in VS Code:

```bash
code /path/to/your/project
```

#### 2. Automatic Registration

The Gorev VS Code extension automatically:

- Detects workspace root
- Runs `gorev init` if `.gorev/` doesn't exist
- Registers workspace with metadata

**Status Bar** shows:

```
Gorev: my-project (15 tasks)
```

#### 3. Workspace-Specific Views

All TreeView panels (Tasks, Projects, Templates) filter to current workspace:

- Tasks panel shows only tasks from current workspace database
- Projects panel shows projects from current workspace
- Creating tasks automatically uses current workspace

#### 4. Multi-Workspace Setup

Open multiple workspace folders in VS Code:

```
File → Add Folder to Workspace
```

Extension handles each workspace independently with separate database connections.

### Web UI Workflow

#### 1. Start Server

```bash
gorev serve --api-port 5082
```

#### 2. Open Web UI

Navigate to http://localhost:5082

#### 3. Workspace Switcher

**Top-right corner** shows workspace dropdown:

```
┌────────────────────────────┐
│  my-project (15 tasks)   ▼ │
├────────────────────────────┤
│ ✓ my-project (15 tasks)    │
│   other-project (23 tasks) │
│   research (7 tasks)       │
└────────────────────────────┘
```

Click to switch. UI automatically:

- Invalidates all React Query caches
- Refetches tasks/projects from selected workspace
- Updates HTTP headers (`X-Workspace-Id`, `X-Workspace-Path`, `X-Workspace-Name`)
- Saves preference to localStorage

---

## Configuration

### Environment Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `GOREV_DB_PATH` | Override database path | Auto-detected | `/custom/path/.gorev/gorev.db` |
| `GOREV_WORKSPACE_NAME` | Override workspace name | Directory name | `My Project` |

**Example**:

```bash
GOREV_DB_PATH=/custom/workspace/.gorev/gorev.db gorev serve
```

### MCP Configuration

#### Auto-Detection (Recommended)

Don't specify `GOREV_DB_PATH`. Gorev will:

1. Look for `.gorev/gorev.db` in current directory
2. Walk up parent directories to find workspace root
3. Fall back to global database (`~/.gorev/gorev.db`) if not found

**Example** (`.kilocode/mcp.json`):

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": [
        "-y",
        "@mehmetsenol/gorev-mcp-server@latest"
      ],
      "env": {
        "GOREV_LANG": "en"
      }
    }
  }
}
```

#### Explicit Path Per Workspace

Specify `GOREV_DB_PATH` for each workspace:

**Project A** (`.kilocode/mcp.json`):

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": [
        "-y",
        "@mehmetsenol/gorev-mcp-server@latest"
      ],
      "env": {
        "GOREV_DB_PATH": "/projects/project-a/.gorev/gorev.db"
      }
    }
  }
}
```

**Project B** (`.kilocode/mcp.json`):

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": [
        "-y",
        "@mehmetsenol/gorev-mcp-server@latest"
      ],
      "env": {
        "GOREV_DB_PATH": "/projects/project-b/.gorev/gorev.db"
      }
    }
  }
}
```

---

## Architecture

### Database Structure

#### Global Registry (`~/.gorev/workspace_registry.db`)

Stores metadata for all workspaces:

```sql
CREATE TABLE workspaces (
  workspace_id TEXT PRIMARY KEY,  -- SHA256 hash (first 8 chars)
  name TEXT NOT NULL,             -- Human-readable name
  path TEXT UNIQUE NOT NULL,      -- Absolute workspace path
  db_path TEXT NOT NULL,          -- Path to workspace's gorev.db
  registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Workspace Database (`.gorev/gorev.db`)

Each workspace has isolated database with full schema:

```sql
-- Task management tables
CREATE TABLE gorevler (...);
CREATE TABLE projeler (...);
CREATE TABLE baglantilar (...);
CREATE TABLE etiketler (...);
CREATE TABLE gorev_etiketleri (...);

-- Template system
CREATE TABLE gorev_templateleri (...);

-- AI context management
CREATE TABLE ai_interactions (...);
CREATE TABLE ai_context (...);

-- Search and filtering
CREATE TABLE gorevler_fts (...);
CREATE TABLE filter_profiles (...);
CREATE TABLE search_history (...);
```

### Context Propagation

#### HTTP Headers

All Web UI → REST API requests include:

```
X-Workspace-Id: 4a5d7c9b8f3e2a1d
X-Workspace-Path: /home/user/projects/my-app
X-Workspace-Name: my-app
```

#### Middleware Layer

Fiber middleware intercepts requests:

```go
func WorkspaceMiddleware(c *fiber.Ctx) error {
    workspaceID := c.Get("X-Workspace-Id")
    workspacePath := c.Get("X-Workspace-Path")

    if workspaceID == "" || workspacePath == "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "Missing workspace context",
        })
    }

    // Attach workspace context to request
    c.Locals("workspaceID", workspaceID)
    c.Locals("workspacePath", workspacePath)

    return c.Next()
}
```

#### Database Connection

Server opens workspace-specific database based on headers:

```go
func GetWorkspaceDB(c *fiber.Ctx) (*sql.DB, error) {
    workspacePath := c.Locals("workspacePath").(string)
    dbPath := filepath.Join(workspacePath, ".gorev", "gorev.db")

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    return db, nil
}
```

---

## Migration from Single Workspace

### Pre-v0.16.0 Behavior

Before v0.16.0, Gorev used a single global database:

```
~/.gorev/gorev.db  # All tasks for all projects
```

### v0.16.0 Migration

Automatic migration on first run:

1. **Detect Legacy Database**:
   - Check if `~/.gorev/gorev.db` exists
   - Check if it contains tasks

2. **Create Workspaces**:
   - For each unique `proje_id` in legacy database:
     - Create workspace directory
     - Copy relevant tasks to workspace database
     - Register workspace

3. **Preserve Legacy**:
   - Keep original `~/.gorev/gorev.db` as backup
   - Rename to `~/.gorev/gorev.db.backup-v0.15`

**Manual Migration**:

```bash
# Export legacy tasks
gorev export --output ~/backup/legacy-tasks.json

# Initialize new workspace
cd /path/to/project
gorev init

# Import tasks
gorev import --input ~/backup/legacy-tasks.json --conflict-resolution merge
```

---

## Advanced Usage

### Workspace Commands

#### List Workspaces

```bash
gorev workspace list
```

**Options**:

- `--active-only`: Show only recently accessed workspaces
- `--sort-by-tasks`: Sort by task count (descending)
- `--json`: Output as JSON

#### Register Workspace

```bash
gorev workspace register --path /path/to/project --name "My Project"
```

#### Unregister Workspace

```bash
gorev workspace unregister --id 4a5d7c9b
```

**Note**: This removes workspace from registry but doesn't delete `.gorev/` directory or database.

#### Workspace Info

```bash
gorev workspace info --id 4a5d7c9b
```

**Output**:

```yaml
Workspace: my-project
ID: 4a5d7c9b8f3e2a1d
Path: /home/user/projects/my-app
Database: /home/user/projects/my-app/.gorev/gorev.db
Registered: 2025-10-01 09:00:00
Last Accessed: 2025-10-05 10:30:00

Statistics:
  Total Tasks: 15
  Pending: 8
  In Progress: 5
  Completed: 2

  Total Projects: 3
  Active Project: Backend API

  Templates: 6 default templates
```

### Workspace Cleanup

#### Remove Orphaned Workspaces

```bash
gorev workspace cleanup
```

Removes workspace registrations where:

- Workspace path no longer exists
- Database file is missing or corrupted
- Last accessed > 90 days ago

#### Archive Old Workspaces

```bash
gorev workspace archive --id 4a5d7c9b
```

Exports workspace data to JSON and marks as archived:

```
~/.gorev/archives/my-project-2025-10-05.json
```

---

## Best Practices

### 1. One Workspace Per Project

**Recommended**:

```
/projects/
├── project-a/
│   └── .gorev/gorev.db    # Workspace A
├── project-b/
│   └── .gorev/gorev.db    # Workspace B
└── project-c/
    └── .gorev/gorev.db    # Workspace C
```

**Avoid**:

```
/projects/
└── .gorev/gorev.db        # Single database for all projects
```

### 2. Consistent Naming

Use descriptive workspace names:

```bash
gorev init --name "E-commerce Backend"
gorev init --name "Mobile App v2.0"
```

Avoid generic names:

```bash
gorev init --name "project"
gorev init --name "test"
```

### 3. Regular Backups

Export workspace data weekly:

```bash
# Automated backup script
#!/bin/bash
WORKSPACE_ID="4a5d7c9b"
DATE=$(date +%Y-%m-%d)
gorev export \
  --workspace-id $WORKSPACE_ID \
  --output ~/backups/gorev-$DATE.json \
  --include-completed \
  --include-dependencies
```

### 4. Share Workspace Config

Commit MCP config to version control:

```
project-root/
├── .kilocode/
│   └── mcp.json           # Commit this
├── .gorev/
│   └── gorev.db           # DO NOT commit this
└── .gitignore             # Add .gorev/gorev.db
```

**.gitignore**:

```
.gorev/gorev.db
.gorev/gorev.db-shm
.gorev/gorev.db-wal
```

---

## Troubleshooting

### Issue: Workspace Not Found

**Symptoms**:

- CLI shows "No workspace detected"
- VS Code extension shows empty task list

**Solutions**:

```bash
# Verify .gorev/ directory exists
ls -la .gorev/

# Check if database file exists
ls -la .gorev/gorev.db

# Re-initialize if missing
gorev init

# Verify registration
gorev workspace list
```

### Issue: Wrong Workspace Selected

**Symptoms**:

- Web UI shows tasks from different project

**Solutions**:

```bash
# Check workspace switcher (top-right corner)
# Select correct workspace from dropdown

# Or clear localStorage and reload
# Browser DevTools → Application → Local Storage → Clear

# Or specify workspace explicitly in URL
http://localhost:5082/?workspace=4a5d7c9b
```

### Issue: Database Locked

**Symptoms**:

- "database is locked" error in VS Code or Web UI

**Solutions**:

```bash
# Check for running Gorev processes
ps aux | grep gorev

# Kill duplicate server instances
killall gorev

# Disable WAL mode if on network drive
cd .gorev/
sqlite3 gorev.db "PRAGMA journal_mode=DELETE;"

# Use local directory instead of network share
# (VirtualBox shared folders don't support WAL)
```

### Issue: Workspace ID Collision

**Symptoms**:

- Two workspaces with same ID (extremely rare)

**Solutions**:

```bash
# Unregister conflicting workspace
gorev workspace unregister --id 4a5d7c9b

# Re-register with explicit path
gorev workspace register \
  --path /full/absolute/path/to/workspace \
  --name "Unique Workspace Name"

# Verify unique ID generated
gorev workspace list
```

---

## Performance

### Optimization Tips

1. **Limit Workspace Count**: Keep < 20 active workspaces
2. **Archive Old Projects**: Move completed projects to archives
3. **Regular Cleanup**: Run `gorev workspace cleanup` monthly
4. **Database Maintenance**: Vacuum databases quarterly

```bash
# Vacuum all workspace databases
gorev workspace list --json | jq -r '.[] | .db_path' | while read db; do
  sqlite3 "$db" "VACUUM;"
done
```

### Benchmarks

| Operation | Single Workspace | 10 Workspaces | 50 Workspaces |
|-----------|-----------------|---------------|---------------|
| List Tasks (100 tasks) | 5ms | 5ms | 5ms |
| Create Task | 2ms | 2ms | 2ms |
| Workspace Switch (Web UI) | N/A | 50ms | 80ms |
| Database Size (1000 tasks) | 2.5 MB | 2.5 MB | 2.5 MB |

**Conclusion**: Performance is workspace-isolated. Number of workspaces doesn't affect individual workspace operations.

---

## Security

### Considerations

1. **Database Access**: SQLite databases are unencrypted
2. **Workspace Isolation**: No authentication between workspaces
3. **File Permissions**: Rely on OS-level file permissions
4. **Shared Environments**: Not recommended for multi-user systems

### Hardening

**File Permissions** (Linux/macOS):

```bash
chmod 700 ~/.gorev
chmod 600 ~/.gorev/gorev.db
```

**Encrypted Workspace** (Future Enhancement):

```bash
# Create encrypted container for workspace
# (Future feature - not yet implemented)
gorev workspace create --encrypted --passphrase "strong-password"
```

---

## API Reference

### Workspace Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/workspaces` | GET | List all registered workspaces |
| `/api/workspaces/register` | POST | Register new workspace |
| `/api/workspaces/:id` | GET | Get workspace details |
| `/api/workspaces/:id` | DELETE | Unregister workspace |
| `/api/workspaces/:id/stats` | GET | Get workspace statistics |

**Example Request**:

```bash
curl -X GET http://localhost:5082/api/workspaces
```

**Response**:

```json
{
  "workspaces": [
    {
      "workspaceId": "4a5d7c9b8f3e2a1d",
      "name": "my-project",
      "path": "/home/user/projects/my-app",
      "dbPath": "/home/user/projects/my-app/.gorev/gorev.db",
      "registeredAt": "2025-10-01T09:00:00Z",
      "lastAccessed": "2025-10-05T10:30:00Z",
      "taskCount": 15,
      "projectCount": 3
    }
  ]
}
```

---

## Comparison: Single vs. Multi-Workspace

| Aspect | Single Workspace (v0.15) | Multi-Workspace (v0.16) |
|--------|--------------------------|-------------------------|
| **Database** | `~/.gorev/gorev.db` | `.gorev/gorev.db` per project |
| **Task Isolation** | ❌ All tasks mixed | ✅ Project-specific tasks |
| **Context Switching** | Manual filtering | Automatic workspace selection |
| **Backup** | One file | Per-project exports |
| **Collaboration** | Shared database | Independent databases |
| **IDE Integration** | Global task list | Workspace-aware TreeViews |
| **Web UI** | Not available | Workspace switcher |
| **Performance** | Degrades with > 1000 tasks | Scales per workspace |

---

## Migration Checklist

### From v0.15 to v0.16

- [ ] **Backup existing database**:

  ```bash
  cp ~/.gorev/gorev.db ~/gorev.db.backup-$(date +%Y%m%d)
  ```

- [ ] **Update to v0.16.0**:

  ```bash
  npm install -g @mehmetsenol/gorev-mcp-server@latest
  ```

- [ ] **Initialize workspaces** for each project:

  ```bash
  cd /path/to/project-a && gorev init
  cd /path/to/project-b && gorev init
  ```

- [ ] **Export legacy tasks**:

  ```bash
  gorev export --output ~/legacy-tasks.json
  ```

- [ ] **Import into workspaces**:

  ```bash
  cd /path/to/project-a
  gorev import --input ~/legacy-tasks.json --project-filter "Project A"
  ```

- [ ] **Update MCP configs** in each project

- [ ] **Test workspace switching** in Web UI

- [ ] **Verify VS Code extension** detects workspaces

---

## Additional Resources

- **Web UI Guide**: [Web UI Documentation](web-ui.md)
- **API Reference**: [REST API Guide](../developer/api-reference.md)
- **Template System**: [Template Guide](template-system.md)
- **GitHub Issues**: https://github.com/msenol/gorev/issues

---

**Need Help?** Open an issue at [GitHub Issues](https://github.com/msenol/gorev/issues) with the `multi-workspace` label.
