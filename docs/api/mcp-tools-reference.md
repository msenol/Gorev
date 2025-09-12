# 🛠️ MCP Tools & API Reference - Complete Guide

**Version**: v0.14.0  
**Last Updated**: September 12, 2025  
**Status**: Production Ready  
**Total Tools**: 29 MCP Tools

---

## 📋 Quick Navigation

- [Task Management Tools](#task-management-tools) (7 tools)
- [Project Management Tools](#project-management-tools) (6 tools)  
- [Template System Tools](#template-system-tools) (4 tools)
- [AI Context Tools](#ai-context-tools) (5 tools)
- [File Watcher Tools](#file-watcher-tools) (4 tools)
- [Advanced Tools](#advanced-tools) (3 tools)

---

## 🎯 Task Management Tools

### 1. `gorev_listele` - List Tasks

**Description**: Retrieve tasks with filtering and pagination support

**Parameters**:
```json
{
  "durum": "string (optional)",
  "sirala": "string (optional)", 
  "filtre": "string (optional)",
  "etiket": "string (optional)",
  "tum_projeler": "boolean (optional)",
  "limit": "number (optional)",
  "offset": "number (optional)"
}
```

**Example Usage**:
```javascript
// List urgent tasks
{
  "filtre": "acil",
  "durum": "beklemede",
  "limit": 10
}

// List completed tasks with pagination
{
  "durum": "tamamlandi",
  "limit": 20,
  "offset": 40
}
```

### 2. `gorev_detay` - Get Task Details

**Description**: Retrieve detailed information about a specific task

**Parameters**:
```json
{
  "id": "string (required)"
}
```

**Response Format**:
```
# Task Details
**Title**: API Integration
**Status**: In Progress  
**Priority**: High
**Project**: Mobile App v2
**Due Date**: 2025-09-15
**Tags**: api, integration, urgent
**Dependencies**: 2 pending
**Description**: Implement REST API integration for user authentication
```

### 3. `gorev_guncelle` - Update Task Status

**Description**: Update task status with automatic validation

**Parameters**:
```json
{
  "id": "string (required)",
  "durum": "string (required)"
}
```

**Valid Status Values**:
- `beklemede` (pending)
- `devam_ediyor` (in progress)  
- `tamamlandi` (completed)

### 4. `gorev_duzenle` - Edit Task Properties

**Description**: Modify task properties with comprehensive validation

**Parameters**:
```json
{
  "id": "string (required)",
  "baslik": "string (optional)",
  "aciklama": "string (optional)", 
  "oncelik": "string (optional)",
  "son_tarih": "string (optional)"
}
```

**Priority Values**: `dusuk`, `orta`, `yuksek`

### 5. `gorev_sil` - Delete Task

**Description**: Permanently delete a task with confirmation

**Parameters**:
```json
{
  "id": "string (required)",
  "onay": "boolean (required)"
}
```

⚠️ **Warning**: This action cannot be undone. Always set `onay: true` for confirmation.

### 6. `gorev_bagimlilik_ekle` - Add Task Dependencies

**Description**: Create dependencies between tasks

**Parameters**:
```json
{
  "kaynak_id": "string (required)",
  "hedef_id": "string (required)",
  "baglanti_tipi": "string (required)"
}
```

**Dependency Types**:
- `tamamla_oncebi` - Must complete before target
- `baslat_oncebi` - Must start before target

---

## 📁 Project Management Tools

### 7. `proje_olustur` - Create Project

**Description**: Create new project with metadata

**Parameters**:
```json
{
  "isim": "string (required)",
  "tanim": "string (required)"
}
```

### 8. `proje_listele` - List Projects

**Description**: Retrieve all projects with status information

**Response Format**:
```
# Projects
1. **Mobile App v2** (15 tasks) - Active
2. **Website Redesign** (8 tasks) - Completed  
3. **API Documentation** (3 tasks) - In Progress
```

### 9. `proje_gorevleri` - Get Project Tasks

**Description**: List all tasks belonging to a specific project

**Parameters**:
```json
{
  "proje_id": "string (required)"
}
```

### 10. `aktif_proje_ayarla` - Set Active Project

**Description**: Set the default active project for task operations

**Parameters**:
```json
{
  "proje_id": "string (required)"
}
```

### 11. `aktif_proje_goster` - Show Active Project

**Description**: Display current active project information

### 12. `aktif_proje_kaldir` - Remove Active Project

**Description**: Clear active project setting

---

## 📋 Template System Tools

### 13. `template_listele` - List Templates

**Description**: Show available task templates

**Response includes**:
- Template ID and name
- Required fields
- Template category
- Usage examples

### 14. `templateden_gorev_olustur` - Create Task from Template

**Description**: Create tasks using predefined templates (replaces deprecated `gorev_olustur`)

**Parameters**:
```json
{
  "template_id": "string (required)",
  "degerler": "object (required)"
}
```

**Example - Bug Report Template**:
```javascript
{
  "template_id": "bug_report_v2",
  "degerler": {
    "baslik": "Login Bug",
    "aciklama": "Users cannot log in with valid credentials", 
    "modul": "authentication",
    "severity": "high",
    "steps": "1. Go to login page\n2. Enter valid credentials\n3. Click login\n4. Error appears",
    "environment": "production"
  }
}
```

**Available Templates**:
- `bug_report_v2` - Enhanced bug reports
- `feature_request` - New feature requests
- `spike_research` - Time-boxed research tasks
- `performance_issue` - Performance optimization tasks
- `security_fix` - Security-related fixes

---

## 🤖 AI Context Tools

### 15. `ai_context_ekle` - Add AI Context

**Description**: Add contextual information for AI processing

**Parameters**:
```json
{
  "context_type": "string (required)",
  "content": "string (required)", 
  "metadata": "object (optional)"
}
```

### 16. `ai_context_listele` - List AI Contexts

**Description**: Show all available AI contexts

### 17. `ai_context_sil` - Delete AI Context

**Description**: Remove AI context entry

**Parameters**:
```json
{
  "context_id": "string (required)"
}
```

### 18. `ai_context_guncelle` - Update AI Context

**Description**: Modify existing AI context

**Parameters**:
```json
{
  "context_id": "string (required)",
  "content": "string (optional)",
  "metadata": "object (optional)"
}
```

### 19. `ai_context_temizle` - Clear AI Context

**Description**: Remove all AI context data

---

## 👀 File Watcher Tools

### 20. `file_watch_ekle` - Add File Watcher

**Description**: Monitor files for changes and auto-update task status

**Parameters**:
```json
{
  "file_path": "string (required)",
  "gorev_id": "string (required)",
  "watch_type": "string (optional)"
}
```

**Watch Types**:
- `modify` - File modification
- `create` - File creation  
- `delete` - File deletion
- `all` - All events (default)

### 21. `file_watch_listele` - List File Watchers

**Description**: Show all active file watchers

### 22. `file_watch_sil` - Remove File Watcher

**Description**: Stop monitoring a file

**Parameters**:
```json
{
  "watcher_id": "string (required)"
}
```

### 23. `file_watch_durdur` - Stop All Watchers

**Description**: Temporarily stop all file monitoring

---

## 📊 Advanced Tools

### 24. `ozet_goster` - Show Summary

**Description**: Display comprehensive project and task statistics

**Response includes**:
- Total tasks by status
- Project completion rates
- Upcoming deadlines
- Performance metrics

### 25. `batch_gorev_guncelle` - Batch Update Tasks

**Description**: Update multiple tasks simultaneously

**Parameters**:
```json
{
  "gorev_ids": "array<string> (required)",
  "updates": "object (required)"
}
```

### 26. `export_data` - Export Data

**Description**: Export tasks and projects in various formats

**Parameters**:
```json
{
  "format": "string (required)", 
  "filter": "object (optional)"
}
```

**Supported Formats**: `json`, `csv`, `markdown`

### 27. `import_data` - Import Data

**Description**: Import tasks from external sources

**Parameters**:
```json
{
  "source": "string (required)",
  "data": "string (required)",
  "options": "object (optional)"
}
```

### 28. `system_status` - System Status

**Description**: Get system health and performance metrics

### 29. `debug_info` - Debug Information

**Description**: Get detailed system debugging information

---

## 🔧 Error Handling Patterns

All MCP tools follow **Rule 15 compliance** with comprehensive error handling:

### Standard Error Response Format
```json
{
  "error": true,
  "code": "ERROR_CODE",
  "message": "Human-readable error message",
  "details": {
    "field": "validation_error_details",
    "context": "additional_context"
  }
}
```

### Common Error Codes
- `VALIDATION_ERROR` - Invalid input parameters
- `NOT_FOUND` - Resource not found
- `PERMISSION_DENIED` - Access not allowed
- `DEPENDENCY_ERROR` - Dependency constraint violation
- `SYSTEM_ERROR` - Internal system error

### Error Handling Example
```javascript
// ✅ GOOD: Proper error handling
try {
  const result = await mcpClient.callTool('gorev_detay', { id: 'invalid-id' });
  if (result.isError) {
    console.error('Task not found:', result.content);
    return;
  }
  // Process result...
} catch (error) {
  console.error('MCP call failed:', error);
}
```

---

## 🚀 Performance & Rate Limiting

### Rate Limiting
- **Default Rate**: 100 requests/minute per client
- **Burst Limit**: 20 concurrent requests  
- **Retry Policy**: Exponential backoff (1s, 2s, 4s, 8s)

### Optimization Tips

1. **Use Pagination**: For large datasets, use `limit` and `offset`
2. **Batch Operations**: Use batch tools for multiple updates
3. **Cache Results**: Cache frequently accessed data
4. **Filter Early**: Apply filters to reduce response size

### Performance Monitoring
```javascript
// Monitor tool performance
const startTime = Date.now();
const result = await mcpClient.callTool('gorev_listele', params);
const duration = Date.now() - startTime;
console.log(`Tool execution time: ${duration}ms`);
```

---

## 🔒 Authentication & Security

### MCP Protocol Security
- **Local Communication**: stdio-based secure communication
- **Input Validation**: All parameters validated before processing
- **SQL Injection Protection**: 100% parameterized queries  
- **Path Traversal Prevention**: Secure file path handling

### Best Practices
```javascript
// ✅ GOOD: Validate inputs
function validateTaskId(id) {
  if (!id || typeof id !== 'string' || id.trim().length === 0) {
    throw new Error('Invalid task ID');
  }
  return id.trim();
}

// ✅ GOOD: Handle sensitive data
const taskData = {
  title: sanitizeInput(userInput.title),
  description: sanitizeInput(userInput.description)
};
```

---

## 📚 Integration Examples

### Claude Desktop Integration
```json
{
  "mcpServers": {
    "gorev": {
      "command": "/path/to/gorev",
      "args": ["serve"],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

### VS Code Extension Integration
```typescript
import { GorevMCPClient } from './gorev-mcp-client';

class TaskManager {
  private mcpClient: GorevMCPClient;

  async createTask(template: string, values: any) {
    const result = await this.mcpClient.callTool('templateden_gorev_olustur', {
      template_id: template,
      degerler: values
    });
    
    if (result.isError) {
      throw new Error(result.content);
    }
    
    return result.content;
  }
}
```

### cURL Examples
```bash
# Not applicable - MCP uses stdio protocol
# Use MCP-compatible clients like Claude Desktop or VS Code
```

---

## 🎯 Migration Guide

### From v0.13.x to v0.14.0

#### Deprecated Tool Replacements
```javascript
// ❌ OLD (deprecated since v0.10.0, removed in v0.11.1)
mcpClient.callTool('gorev_olustur', {
  baslik: 'Bug Fix',
  aciklama: 'Fix login issue'
});

// ✅ NEW (required since v0.10.0)
mcpClient.callTool('templateden_gorev_olustur', {
  template_id: 'bug_report_v2',
  degerler: {
    baslik: 'Bug Fix',
    aciklama: 'Fix login issue',
    severity: 'high',
    modul: 'authentication'
  }
});
```

#### Enhanced Error Responses
- More detailed error information
- Structured error codes
- Better validation messages

#### New Features in v0.14.0
- AI Context management tools
- File watcher integration
- Batch processing capabilities
- Enhanced template system

---

<div align="center">

**[⬆ Back to Top](#-mcp-tools--api-reference---complete-guide)**

Made with ❤️ by the Gorev Team | Enhanced by Claude (Anthropic)

*Following Rule 15 & DRY Principles for Reliable API Design*

</div>