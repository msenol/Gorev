# VS Code Extension: MCP to REST API Migration Guide

**Created:** 30 September 2025
**Status:** Planning Phase
**Target Version:** v0.17.0

## 📋 Overview

This document outlines the migration plan for transitioning the Gorev VS Code extension from MCP (Model Context Protocol) stdio communication to REST API HTTP communication.

## 🎯 Motivation

### Why Migrate

1. **Markdown Parsing Hell**
   - MCP responses are markdown strings
   - `markdownParser.ts` constantly breaks with format changes
   - Fragile string parsing logic
   - Difficult to debug

2. **No Type Safety**
   - Responses are untyped strings
   - Requires manual parsing and type casting
   - Runtime errors common
   - Loss of TypeScript benefits

3. **Stdio Complexity**
   - Process spawn/kill management
   - Buffer handling
   - Timeout handling
   - Platform-specific issues (Windows vs Linux)

4. **Web UI Consistency**
   - Web UI already uses REST API successfully
   - Shared endpoint definitions
   - Same response format
   - Easier maintenance

### Benefits of REST API

1. **Structured JSON Responses**

   ```typescript
   // Before (MCP - bad)
   const response = "## Task List\n- [pending] Task 1..."
   const tasks = parseMarkdown(response) // 🤮

   // After (REST API - good)
   const response: ApiResponse<Task[]> = await api.get('/tasks')
   const tasks = response.data // 😎
   ```

2. **TypeScript Type Safety**

   ```typescript
   interface ApiResponse<T> {
     success: boolean;
     data?: T;
     error?: string;
     total?: number;
   }

   const tasks: Task[] = await apiClient.getTasks({ durum: 'pending' })
   ```

3. **HTTP Standard Features**
   - Status codes (200, 404, 500)
   - Headers (Content-Type, CORS)
   - Query parameters
   - Request/Response interceptors
   - Built-in caching

4. **Easier Debugging**
   - Use browser DevTools for inspection
   - Network tab shows all requests
   - cURL-compatible for testing
   - Standard HTTP debugging tools

## 📊 Current State Analysis

### Existing Files

1. **MCP Client** (to be deprecated)
   - `src/mcp/client.ts` - Main MCP client (stdio)
   - `src/mcp/types.ts` - MCP type definitions
   - `src/utils/markdownParser.ts` - Fragile markdown parsing

2. **API Client** (already exists!)
   - `src/api/client.ts` - Basic REST API client
   - Currently only used for language synchronization
   - Needs expansion to cover all operations

3. **Unified Client** (hybrid approach)
   - `src/unified/client.ts` - Combines MCP + API
   - Experimental, not fully implemented

### Current Usage

```typescript
// Current MCP usage in providers
const mcpClient = new MCPClient();
await mcpClient.connect();
const response = await mcpClient.callTool('gorev_listele', { durum: 'pending' });
const tasks = parseMarkdownTaskList(response); // ❌ Fragile
```

## 🗺️ Migration Plan

### Phase 1: Expand API Client (Week 1)

**Goal:** Make `src/api/client.ts` feature-complete

**Tasks:**

1. Add all CRUD methods to API client

   ```typescript
   class GorevApiClient {
     async getTasks(filters?: TaskFilters): Promise<Task[]>
     async getTask(id: string): Promise<Task>
     async createTaskFromTemplate(templateId: string, values: Record<string, string>): Promise<Task>
     async updateTask(id: string, updates: Partial<Task>): Promise<Task>
     async deleteTask(id: string): Promise<void>

     async getProjects(): Promise<Project[]>
     async getProject(id: string): Promise<Project>
     async createProject(name: string, description?: string): Promise<Project>
     async activateProject(id: string): Promise<Project>

     async getTemplates(category?: string): Promise<Template[]>

     async createSubtask(parentId: string, data: SubtaskData): Promise<Task>
     async changeParent(taskId: string, newParentId: string): Promise<Task>
     async getHierarchy(taskId: string): Promise<TaskHierarchy>

     async addDependency(targetId: string, sourceId: string): Promise<void>
   }
   ```

2. Add proper error handling

   ```typescript
   class ApiError extends Error {
     constructor(
       public statusCode: number,
       public apiError: string,
       public endpoint: string
     ) {
       super(`API Error ${statusCode}: ${apiError}`);
     }
   }
   ```

3. Add request/response interceptors

   ```typescript
   private async request<T>(options: RequestOptions): Promise<ApiResponse<T>> {
     try {
       const response = await fetch(this.buildUrl(options.path), {
         method: options.method,
         headers: this.getHeaders(),
         body: options.body ? JSON.stringify(options.body) : undefined,
       });

       if (!response.ok) {
         throw new ApiError(response.status, await response.text(), options.path);
       }

       return await response.json();
     } catch (error) {
       this.handleError(error, options);
       throw error;
     }
   }
   ```

### Phase 2: Update TreeView Providers (Week 2)

**Goal:** Replace MCP calls with API calls in all providers

**Files to Update:**

1. `src/providers/enhancedGorevTreeProvider.ts`
2. `src/providers/projeTreeProvider.ts`
3. `src/providers/templateTreeProvider.ts`

**Before:**

```typescript
// src/providers/enhancedGorevTreeProvider.ts (old)
const response = await this.mcpClient.callTool('gorev_listele', { durum: 'pending' });
const tasks = parseMarkdownTaskList(response);
```

**After:**

```typescript
// src/providers/enhancedGorevTreeProvider.ts (new)
const tasks = await this.apiClient.getTasks({ durum: 'pending' });
// No parsing needed! Direct TypeScript types
```

### Phase 3: Update Commands (Week 2)

**Goal:** Replace MCP calls with API calls in all commands

**Files to Update:**

- `src/commands/index.ts`
- All command handlers

**Example Migration:**

**Before:**

```typescript
// commands/createTask.ts (old)
async function createTask(templateId: string) {
  const response = await mcpClient.callTool('templateden_gorev_olustur', {
    template_id: templateId,
    degerler: formData
  });

  const taskId = extractTaskIdFromMarkdown(response); // ❌ Fragile
  vscode.window.showInformationMessage('Task created!');
}
```

**After:**

```typescript
// commands/createTask.ts (new)
async function createTask(templateId: string) {
  try {
    const task = await apiClient.createTaskFromTemplate(templateId, formData);
    vscode.window.showInformationMessage(`Task created: ${task.baslik}`);
    return task; // ✅ Return typed object
  } catch (error) {
    if (error instanceof ApiError) {
      vscode.window.showErrorMessage(`Failed: ${error.apiError}`);
    }
    throw error;
  }
}
```

### Phase 4: Remove MCP Client (Week 3)

**Goal:** Deprecate and remove MCP client code

**Steps:**

1. Mark MCP client as deprecated

   ```typescript
   /**
    * @deprecated Use GorevApiClient instead
    */
   export class MCPClient {
     constructor() {
       throw new Error('MCPClient is deprecated. Use GorevApiClient from src/api/client.ts');
     }
   }
   ```

2. Remove markdown parser

   ```bash
   rm src/utils/markdownParser.ts
   ```

3. Update package.json dependencies
   - Remove any MCP-specific dependencies
   - Keep only HTTP/fetch dependencies

4. Clean up types

   ```bash
   rm src/mcp/types.ts
   ```

### Phase 5: Testing & Validation (Week 3)

**Goal:** Ensure all functionality works with REST API

**Test Plan:**

1. **Unit Tests**

   ```typescript
   // test/unit/apiClient.test.ts
   describe('GorevApiClient', () => {
     it('should fetch tasks with filters', async () => {
       const client = new GorevApiClient('http://localhost:5082');
       const tasks = await client.getTasks({ durum: 'pending' });
       expect(tasks).toBeArrayOf(Task);
     });
   });
   ```

2. **Integration Tests**

   ```typescript
   // test/integration/commands.test.ts
   describe('Task Commands', () => {
     it('should create task from template via API', async () => {
       const task = await vscode.commands.executeCommand('gorev.createTask');
       expect(task).toHaveProperty('id');
       expect(task).toHaveProperty('baslik');
     });
   });
   ```

3. **Manual Testing Checklist**
   - [ ] List tasks in TreeView
   - [ ] Create task from template
   - [ ] Update task status
   - [ ] Delete task
   - [ ] Project operations
   - [ ] Template operations
   - [ ] Subtask operations
   - [ ] Dependency operations
   - [ ] Error handling

## 📝 Implementation Checklist

### API Client Enhancement

- [ ] Add all CRUD methods
- [ ] Implement proper error handling
- [ ] Add request/response interceptors
- [ ] Add TypeScript interfaces for all API types
- [ ] Add JSDoc documentation
- [ ] Add unit tests for API client

### TreeView Migration

- [ ] enhancedGorevTreeProvider.ts
- [ ] projeTreeProvider.ts
- [ ] templateTreeProvider.ts
- [ ] Update refresh logic
- [ ] Test drag & drop functionality
- [ ] Test inline editing

### Command Migration

- [ ] Create task command
- [ ] Update task command
- [ ] Delete task command
- [ ] Project commands
- [ ] Template commands
- [ ] Subtask commands
- [ ] Dependency commands

### Cleanup

- [ ] Mark MCPClient as deprecated
- [ ] Remove markdown parser
- [ ] Remove MCP types
- [ ] Update package.json
- [ ] Update README.md
- [ ] Update CHANGELOG.md

### Testing

- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Manual testing
- [ ] Performance testing
- [ ] Cross-platform testing (Windows, Linux, macOS)

## 🚧 Breaking Changes

### For Users

**None** - The migration is transparent to end users. All VS Code commands and UI remain the same.

### For Developers

1. **Import Changes**

   ```typescript
   // Before
   import { MCPClient } from './mcp/client';

   // After
   import { GorevApiClient } from './api/client';
   ```

2. **Method Signature Changes**

   ```typescript
   // Before
   const response: string = await mcpClient.callTool('gorev_listele', {});
   const tasks = parseMarkdown(response);

   // After
   const tasks: Task[] = await apiClient.getTasks();
   ```

3. **Error Handling Changes**

   ```typescript
   // Before
   try {
     const response = await mcpClient.callTool('gorev_olustur', data);
   } catch (error) {
     // Generic error
   }

   // After
   try {
     const task = await apiClient.createTask(data);
   } catch (error) {
     if (error instanceof ApiError) {
       // Structured error with status code
     }
   }
   ```

## 📈 Performance Considerations

### Expected Improvements

1. **Response Parsing**
   - Before: 10-50ms (markdown parsing overhead)
   - After: <1ms (JSON.parse is native)

2. **Type Safety**
   - Before: Runtime type errors common
   - After: Compile-time type checking

3. **Debugging Time**
   - Before: Hours (markdown format debugging)
   - After: Minutes (HTTP inspection tools)

### Potential Issues

1. **Network Latency**
   - HTTP has slightly more overhead than stdio
   - Mitigation: Keep API server on localhost
   - Impact: Negligible (<5ms difference)

2. **Connection Management**
   - Need to ensure API server is running
   - Mitigation: Add health check on extension activation
   - Fallback: Show clear error message if server unavailable

## 🔧 Configuration Changes

### Extension Settings

Add new settings for API configuration:

```json
{
  "gorev.apiBaseUrl": {
    "type": "string",
    "default": "http://localhost:5082/api/v1",
    "description": "Base URL for Gorev REST API"
  },
  "gorev.apiTimeout": {
    "type": "number",
    "default": 5000,
    "description": "API request timeout in milliseconds"
  },
  "gorev.serverMode": {
    "type": "string",
    "enum": ["api", "mcp"],
    "default": "api",
    "description": "Communication mode with Gorev server (DEPRECATED: mcp)"
  }
}
```

## 📅 Timeline

| Week | Phase | Deliverables |
|------|-------|--------------|
| 1 | API Client | Complete GorevApiClient with all methods |
| 2 | Migration | Update all providers and commands |
| 3 | Cleanup & Testing | Remove MCP code, comprehensive testing |
| 4 | Release | v0.17.0 with REST API only |

## 🎯 Success Criteria

1. ✅ All features work with REST API
2. ✅ No markdown parsing code remains
3. ✅ 100% test coverage for API client
4. ✅ 90%+ test coverage for commands
5. ✅ Performance equal or better than MCP
6. ✅ No user-facing breaking changes
7. ✅ Documentation updated

## 🔗 Related Documents

- [REST API Reference](../api/rest-api-reference.md)
- [Web UI Development Guide](web-ui-development.md)
- [VS Code Extension Guide](../guides/user/vscode-extension.md)

## 📝 Notes

- Keep MCP support for Claude Desktop and other MCP-only tools
- REST API is only for VS Code extension and Web UI
- Maintain backward compatibility during transition period
- Consider gradual rollout with feature flags

---

**Status:** This migration is currently in **planning phase**. Implementation will begin in v0.17.0 development cycle.
