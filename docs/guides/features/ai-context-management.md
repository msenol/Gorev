# AI Context Management Guide

**Version**: v0.16.0
**Last Updated**: October 5, 2025
**Feature Status**: Production Ready ✅

---

## Overview

Gorev's **AI Context Management** system helps AI assistants (Claude, Copilot, etc.) maintain conversation context when working with tasks. It tracks active tasks, recent interactions, and provides intelligent suggestions based on your workflow patterns.

### Key Features

- ✅ **Active Task Tracking**: Set and track the currently active task
- ✅ **Automatic Status Updates**: Tasks move to "in_progress" when set as active
- ✅ **Recent Task History**: Track last 20 accessed tasks
- ✅ **Context Summary**: Session overview for AI assistants
- ✅ **Batch Operations**: Update multiple tasks in one command
- ✅ **Natural Language Queries**: Search tasks using conversational language
- ✅ **Session Persistence**: Context survives server restarts

---

## Concepts

### Active Task

The **active task** is the task you're currently working on. Setting a task as active:

- Automatically updates status to `in_progress` (if currently `pending`)
- Adds task to recent history
- Provides AI assistant with focused context
- Enables quick task switching

**Example Workflow**:

```
User: "Make task #42 active"
→ Task #42 status: pending → in_progress
→ Active task set to #42
→ Added to recent history

User: "What am I working on?"
→ Shows active task #42 details
```

### Recent History

**Recent history** tracks the last 20 tasks you've viewed or modified:

- Tasks are added when viewed, edited, or set as active
- Most recent tasks appear first
- Helps AI assistant understand your workflow
- Used for context-aware suggestions

### AI Context Session

An **AI context session** represents a conversation thread with an AI assistant:

- Started when first MCP tool is called
- Tracks all tool invocations
- Stores query patterns and preferences
- Persists across server restarts (stored in database)

---

## MCP Tools

### gorev_set_active

Set a task as the active task and automatically update its status.

**Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `task_id` | string | ✅ | UUID of task to activate |

**Behavior**:

1. Validates task exists
2. If task status is `pending`, changes to `in_progress`
3. Sets as active task
4. Adds to recent history
5. Updates last accessed timestamp

**Example**:

```json
{
  "name": "gorev_set_active",
  "arguments": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Response**:

```
✓ Aktif görev ayarlandı: API authentication implementasyonu
  Durum güncellendi: beklemede → devam_ediyor
```

**AI Assistant Usage**:

```
User: "I'm going to work on the authentication task now"
AI: [Calls gorev_set_active with task_id]
     "Great! I've set the authentication task as active and moved it to in-progress.
     What would you like to work on first?"
```

### gorev_get_active

Retrieve the currently active task.

**Parameters**: None

**Example**:

```json
{
  "name": "gorev_get_active",
  "arguments": {}
}
```

**Response**:

```markdown
## Aktif Görev

**Görev:** API authentication implementasyonu
**ID:** 550e8400-e29b-41d4-a716-446655440000
**Durum:** devam_ediyor
**Öncelik:** yuksek
**Proje:** E-ticaret Sitesi

### Açıklama
JWT tabanlı authentication sistemi kur. Refresh token desteği olmalı.

### Yapılacaklar
- [ ] JWT library entegrasyonu
- [ ] User authentication endpoint
- [ ] Token refresh mekanizması
- [ ] Rate limiting
```

**No Active Task**:

```
Henüz aktif görev ayarlanmamış.
```

### gorev_recent

List recently accessed tasks.

**Parameters**:

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `limit` | number | ❌ | 5 | Maximum number of tasks to return |

**Example**:

```json
{
  "name": "gorev_recent",
  "arguments": {
    "limit": 10
  }
}
```

**Response**:

```markdown
## Son Görüntülenen Görevler

1. **API authentication implementasyonu** (devam_ediyor)
   ID: 550e8400-e29b-41d4-a716-446655440000
   Son Erişim: 2 dakika önce

2. **Veritabanı şeması tasarımı** (tamamlandi)
   ID: 7c9e6679-7425-40de-944b-e07fc1f90ae7
   Son Erişim: 15 dakika önce

3. **README dosyasını güncelle** (beklemede)
   ID: 9f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b5c
   Son Erişim: 1 saat önce
```

### gorev_context_summary

Get a summary of the current AI session context.

**Parameters**: None

**Example**:

```json
{
  "name": "gorev_context_summary",
  "arguments": {}
}
```

**Response**:

```markdown
## AI Oturum Özeti

### Aktif Görev
- **API authentication implementasyonu** (devam_ediyor)
- Öncelik: yuksek
- Proje: E-ticaret Sitesi

### Oturum İstatistikleri
- Toplam Sorgu: 42
- Oluşturulan Görev: 8
- Güncellenen Görev: 15
- Son Aktivite: 2 dakika önce

### Son İşlemler
1. Görev durumu güncellendi: #550e8400 → devam_ediyor
2. Yeni görev oluşturuldu: "README dosyasını güncelle"
3. Görev detayı görüntülendi: #7c9e6679
4. Görev arama yapıldı: "authentication"
5. Aktif görev değiştirildi

### Çalışma Alanı
- Aktif Proje: E-ticaret Sitesi
- Toplam Görev: 23 (8 beklemede, 12 devam ediyor, 3 tamamlandı)
```

### gorev_batch_update

Update multiple tasks in a single operation.

**Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `updates` | array | ✅ | List of task updates |

**Update Object**:

```typescript
{
  id: string,           // Task UUID
  durum?: string,       // Status: beklemede, devam_ediyor, tamamlandi
  oncelik?: string,     // Priority: dusuk, orta, yuksek
  baslik?: string,      // Title
  aciklama?: string     // Description
}
```

**Example**:

```json
{
  "name": "gorev_batch_update",
  "arguments": {
    "updates": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "durum": "tamamlandi"
      },
      {
        "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
        "oncelik": "yuksek"
      },
      {
        "id": "9f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b5c",
        "durum": "devam_ediyor",
        "baslik": "API documentation (updated)"
      }
    ]
  }
}
```

**Response**:

```
✓ Toplu güncelleme tamamlandı:
  - 3 görev başarıyla güncellendi
  - 0 görev hata aldı

Güncellenen Görevler:
  1. 550e8400... → durum: tamamlandi
  2. 7c9e6679... → oncelik: yuksek
  3. 9f1a2b3c... → durum: devam_ediyor, baslik: API documentation (updated)
```

### gorev_nlp_query

Search and filter tasks using natural language queries.

**Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | ✅ | Natural language query |

**Example Queries**:

```json
{
  "name": "gorev_nlp_query",
  "arguments": {
    "query": "show me all high priority bugs in production"
  }
}
```

```json
{
  "name": "gorev_nlp_query",
  "arguments": {
    "query": "tasks I worked on this week"
  }
}
```

```json
{
  "name": "gorev_nlp_query",
  "arguments": {
    "query": "overdue tasks in the backend project"
  }
}
```

**NLP Processing**:

The system extracts intent from natural language:

| Query | Extracted Filters |
|-------|------------------|
| "high priority bugs" | priority=yuksek, tags=bug |
| "this week" | date_range=last_7_days |
| "in production" | environment=production |
| "backend project" | project_name="backend" |
| "overdue" | due_date < today |

**Response**:

```markdown
## Arama Sonuçları: "high priority bugs in production"

**Bulunan:** 3 görev

1. 🐛 Login butonu çalışmıyor (yuksek öncelik)
   Proje: E-ticaret Sitesi
   Durum: devam_ediyor
   Ortam: production

2. 🐛 API timeout hatası (yuksek öncelik)
   Proje: Backend API
   Durum: beklemede
   Ortam: production

3. 🐛 Memory leak in server (yuksek öncelik)
   Proje: Infrastructure
   Durum: devam_ediyor
   Ortam: production
```

---

## AI Assistant Workflows

### Workflow 1: Starting Work

**User**: "I'm going to work on the authentication feature"

**AI Assistant Actions**:

1. Search for tasks matching "authentication"
2. If multiple matches, ask user to clarify
3. Call `gorev_set_active` with selected task
4. Provide task summary and next steps

**Example Conversation**:

```
User: "I'm going to work on the authentication feature"

AI: I found 2 tasks related to authentication:
    1. API authentication implementasyonu (pending)
    2. Frontend login form (in progress)

    Which one would you like to work on?

User: "The first one"

AI: [Calls gorev_set_active with task #1]

    Great! I've set "API authentication implementasyonu" as your active task.

    Here's what needs to be done:
    - JWT library entegrasyonu
    - User authentication endpoint
    - Token refresh mekanizması
    - Rate limiting

    Would you like to start with the JWT integration?
```

### Workflow 2: Context Switching

**User**: "What am I working on?"

**AI Assistant Actions**:

1. Call `gorev_get_active`
2. Display active task details
3. Show progress and next steps

**Example**:

```
User: "What am I working on?"

AI: [Calls gorev_get_active]

    You're currently working on "API authentication implementasyonu" (high priority).

    Status: In Progress
    Project: E-ticaret Sitesi

    Todo items:
    ✓ JWT library entegrasyonu (completed)
    → User authentication endpoint (next up)
    ○ Token refresh mekanizması
    ○ Rate limiting

    You've been working on this for 45 minutes.
```

### Workflow 3: Bulk Status Update

**User**: "Mark all completed tasks from yesterday as done and archive them"

**AI Assistant Actions**:

1. Query tasks from yesterday with status "in_progress"
2. Call `gorev_batch_update` to set all as "completed"
3. Optionally export to archive

**Example**:

```
User: "Mark all tasks I finished yesterday as completed"

AI: [Queries recent tasks]

    I found 5 tasks you worked on yesterday:
    1. Database schema design
    2. User model implementation
    3. API endpoint testing
    4. Code review for PR #42
    5. Documentation update

    Should I mark all of these as completed?

User: "Yes"

AI: [Calls gorev_batch_update with 5 task IDs]

    ✓ Updated 5 tasks to completed status.

    Great work! You completed 5 tasks yesterday.
    Your completion rate this week is 85%.
```

### Workflow 4: Smart Suggestions

**User**: "What should I work on next?"

**AI Assistant Actions**:

1. Call `gorev_context_summary` for session overview
2. Analyze recent history and priorities
3. Suggest next task based on:
   - High priority tasks
   - Blocked tasks becoming unblocked
   - Tasks in same project as active task
   - User's work patterns

**Example**:

```
User: "What should I work on next?"

AI: [Calls gorev_context_summary and analyzes]

    Based on your recent work on authentication, I recommend:

    1. 🔥 Token refresh mekanizması (high priority, same project)
       This continues your authentication work and unblocks
       the frontend team.

    2. ⚡ API endpoint testing (medium priority, dependency resolved)
       The database schema is done, so you can now test
       the endpoints you created.

    3. 📝 API documentation (medium priority, quick win)
       Estimated 30 minutes to document the new auth endpoints.

    Which one would you like to tackle?
```

---

## Database Schema

### AI Interactions Table

```sql
CREATE TABLE ai_interactions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id TEXT NOT NULL,           -- Session UUID
  tool_name TEXT NOT NULL,             -- MCP tool invoked
  arguments TEXT,                      -- JSON arguments
  result TEXT,                         -- Tool result
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  duration_ms INTEGER                  -- Execution time
);

CREATE INDEX idx_ai_session ON ai_interactions(session_id);
CREATE INDEX idx_ai_timestamp ON ai_interactions(created_at);
```

### AI Context Table

```sql
CREATE TABLE ai_context (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id TEXT UNIQUE NOT NULL,
  active_task_id TEXT,                 -- Current active task UUID
  last_query TEXT,                     -- Last NLP query
  preferences TEXT,                    -- JSON user preferences
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (active_task_id) REFERENCES gorevler(id)
);
```

### Recent Tasks Tracking

```sql
-- Stored in ai_interactions table as tool calls
-- Reconstructed from:
SELECT DISTINCT
  json_extract(arguments, '$.task_id') as task_id,
  MAX(created_at) as last_accessed
FROM ai_interactions
WHERE tool_name IN (
  'gorev_detay',
  'gorev_set_active',
  'gorev_guncelle'
)
GROUP BY task_id
ORDER BY last_accessed DESC
LIMIT 20;
```

---

## Advanced Usage

### Session Management

**Start New Session**:

```json
{
  "name": "gorev_context_summary",
  "arguments": {}
}
```

Creates new session if none exists.

**Continue Existing Session**:
Sessions are automatically continued when you restart the MCP server. Session ID is stored in database and restored on startup.

**Clear Session**:

```bash
# CLI (planned feature)
gorev session clear
```

### Context Persistence

**Automatic Persistence**:

- Active task: Saved to database immediately
- Recent history: Updated on every task access
- Session preferences: Stored in JSON blob

**Manual Export**:

```bash
gorev export --include-ai-context --output backup.json
```

**Import Context**:

```bash
gorev import --input backup.json --restore-ai-context
```

### Multi-User Context (Future)

**User-Specific Context** (planned for v0.17.0):

```json
{
  "name": "gorev_set_active",
  "arguments": {
    "task_id": "550e8400...",
    "user_id": "john@example.com"
  }
}
```

Allows multiple users to have separate active tasks and context.

---

## Best Practices

### 1. Use Active Task Consistently

**Good Workflow**:

```
Morning:
  - Review task list
  - Set active task: "Implement user authentication"
  - Work on task
  - Complete or pause
  - Set next active task

Evening:
  - Review completed tasks
  - Update statuses
  - Clear active task if done for the day
```

**Avoid**:

- Setting multiple tasks as active
- Forgetting to clear active task when switching
- Not updating active task when context switching

### 2. Leverage Recent History

Use recent history for:

- Quick access to recently viewed tasks
- Understanding your work patterns
- Resuming interrupted work
- Context-aware AI suggestions

**AI Prompt Examples**:

```
"Show me what I worked on yesterday"
"Continue where I left off"
"What tasks did I touch this week?"
```

### 3. Batch Updates for Efficiency

Instead of updating tasks one by one:

**Inefficient**:

```
gorev task update --id task1 --status completed
gorev task update --id task2 --status completed
gorev task update --id task3 --status completed
```

**Efficient**:

```json
{
  "name": "gorev_batch_update",
  "arguments": {
    "updates": [
      {"id": "task1", "durum": "tamamlandi"},
      {"id": "task2", "durum": "tamamlandi"},
      {"id": "task3", "durum": "tamamlandi"}
    ]
  }
}
```

### 4. Descriptive NLP Queries

**Good Queries**:

```
"Show me high priority bugs in production"
"Tasks I worked on this week"
"Overdue tasks in the backend project"
"Features ready for testing"
```

**Vague Queries**:

```
"Show tasks"
"What's there?"
"My stuff"
```

### 5. Regular Context Summaries

Check context summary periodically:

```
"What's my context summary?"
"How many tasks did I complete today?"
"What am I working on across all projects?"
```

Helps maintain awareness of overall progress and priorities.

---

## Integration Examples

### Claude Desktop

**Typical Workflow**:

```
User: "Good morning! What should I work on today?"

Claude: [Calls gorev_context_summary]
         Good morning! Here's your task overview:

         Active Task: None (you finished yesterday's task)

         High Priority Tasks (3):
         1. API authentication (backend project)
         2. Login UI (frontend project)
         3. Database optimization (infrastructure)

         Based on your recent work on backend features,
         I recommend starting with "API authentication".

         Shall I set it as your active task?

User: "Yes, let's do that"

Claude: [Calls gorev_set_active]
        ✓ Set "API authentication" as active task.

        This task involves:
        - JWT library integration
        - User authentication endpoint
        - Token refresh mechanism

        Would you like to start with JWT integration?
```

### VS Code with Copilot

**Code Context Integration**:

```typescript
// Copilot sees active task context
// and provides relevant suggestions

// Active Task: "API authentication implementasyonu"

// User types: "Create auth middleware"
// Copilot suggests:

import jwt from 'jsonwebtoken';
import { Request, Response, NextFunction } from 'express';

export const authMiddleware = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const token = req.headers.authorization?.split(' ')[1];

  if (!token) {
    return res.status(401).json({ error: 'No token provided' });
  }

  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET!);
    req.user = decoded;
    next();
  } catch (error) {
    return res.status(401).json({ error: 'Invalid token' });
  }
};
```

### Cursor with Gorev

**Task-Aware Code Generation**:

```
User: "Generate the user authentication endpoint"

Cursor: [Reads active task context from Gorev]

        Based on your task "API authentication implementasyonu",
        here's the authentication endpoint:

[Generates code with JWT, refresh tokens, rate limiting]
```

---

## Performance

### Benchmarks

| Operation | Time | Notes |
|-----------|------|-------|
| `gorev_set_active` | 8ms | Includes status update |
| `gorev_get_active` | 3ms | Cached in memory |
| `gorev_recent` (5 tasks) | 12ms | Database query |
| `gorev_context_summary` | 25ms | Multiple queries aggregated |
| `gorev_batch_update` (10 tasks) | 45ms | Transaction-based |
| `gorev_nlp_query` | 180ms | NLP processing + search |

### Optimization

**Caching Strategy**:

- Active task: Cached in memory, invalidated on update
- Recent history: Cached for 60 seconds
- Context summary: Cached for 30 seconds

**Database Optimization**:

- Indexes on `session_id`, `task_id`, `created_at`
- Vacuum database monthly
- Archive old sessions (> 90 days)

---

## Troubleshooting

### Issue: Active Task Not Updating Status

**Symptoms**:

- Task set as active but status stays "pending"

**Solutions**:

```bash
# Check task current status
gorev task show --id 550e8400...

# Manually update if stuck
gorev task update --id 550e8400... --status in_progress

# Verify active task
gorev context active
```

### Issue: Recent History Empty

**Symptoms**:

- `gorev_recent` returns no tasks

**Solutions**:

```bash
# Verify AI interactions are being logged
sqlite3 .gorev/gorev.db "SELECT COUNT(*) FROM ai_interactions;"

# Check table exists
sqlite3 .gorev/gorev.db ".schema ai_interactions"

# Rebuild history by viewing tasks
gorev task list
gorev task show --id <task-id>
```

### Issue: NLP Query Returns Unexpected Results

**Symptoms**:

- Query like "high priority bugs" returns unrelated tasks

**Solutions**:

- Use more specific queries: "high priority bugs tagged with 'security'"
- Check tag spelling and priority values
- Fall back to manual filters if NLP fails

---

## API Reference

### Context Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/context/active` | GET | Get active task |
| `/api/context/active` | POST | Set active task |
| `/api/context/recent` | GET | List recent tasks |
| `/api/context/summary` | GET | Get session summary |
| `/api/context/batch-update` | POST | Batch update tasks |
| `/api/context/nlp-query` | POST | Natural language search |

**Example**:

```bash
curl -X POST http://localhost:5082/api/context/active \
  -H "Content-Type: application/json" \
  -H "X-Workspace-Id: 4a5d7c9b" \
  -d '{"taskId": "550e8400-e29b-41d4-a716-446655440000"}'
```

---

## Future Enhancements

### Planned Features (v0.17.0+)

- [ ] **Multi-User Context**: Per-user active tasks and history
- [ ] **Context Templates**: Predefined context patterns for workflows
- [ ] **AI Learning**: ML-based task suggestions
- [ ] **Voice Commands**: "Alexa, what's my active task?"
- [ ] **Calendar Integration**: Sync with Google Calendar/Outlook
- [ ] **Pomodoro Timer**: Time tracking for active task
- [ ] **Context Switching Costs**: Measure productivity impact
- [ ] **Team Context**: See what teammates are working on

---

## Additional Resources

- **MCP Tools Reference**: [MCP Tools Guide](../../legacy/tr/mcp-araclari.md)
- **Template System**: [Template Guide](template-system.md)
- **Web UI Guide**: [Web UI Documentation](web-ui.md)
- **GitHub Issues**: https://github.com/msenol/gorev/issues

---

**Need Help?** Open an issue at [GitHub Issues](https://github.com/msenol/gorev/issues) with the `ai-context` label.
