# API Integration Examples - Gorev MCP Server

**Version:** v0.16.3 | **Last Updated:** October 6, 2025

Complete integration examples for connecting to Gorev MCP Server from different programming languages and platforms.

---

## Table of Contents

1. [Python Integration](#python-integration)
2. [TypeScript/Node.js Integration](#typescriptnodejs-integration)
3. [Go Integration](#go-integration)
4. [REST API Direct Access](#rest-api-direct-access)
5. [WebSocket Real-Time Updates](#websocket-real-time-updates)

---

## Python Integration

### Using MCP Python SDK

```python
#!/usr/bin/env python3
"""
Gorev MCP Client - Python Example
Demonstrates task management using MCP protocol
"""

import asyncio
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

async def main():
    # Connect to Gorev MCP Server
    server_params = StdioServerParameters(
        command="npx",
        args=["-y", "@mehmetsenol/gorev-mcp-server", "serve"],
        env={"GOREV_LANG": "en"}
    )

    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            # Initialize connection
            await session.initialize()

            # List all tools
            tools_result = await session.list_tools()
            print(f"Available tools: {len(tools_result.tools)}")

            # List tasks
            tasks = await session.call_tool("gorev_listele", {
                "durum": "devam_ediyor",
                "limit": 10
            })
            print(f"Active tasks: {tasks.content}")

            # Create a new task from template
            new_task = await session.call_tool("templateden_gorev_olustur", {
                "template_id": "feature",
                "degerler": {
                    "ozellik_adi": "User Authentication",
                    "aciklama": "Implement JWT-based authentication",
                    "oncelik": "yuksek"
                }
            })
            print(f"Created task: {new_task.content}")

            # Update task status
            update = await session.call_tool("gorev_guncelle", {
                "id": new_task.content["gorev_id"],
                "durum": "devam_ediyor",
                "oncelik": "yuksek"
            })
            print(f"Task updated: {update.content}")

            # Get task details
            details = await session.call_tool("gorev_detay", {
                "id": new_task.content["gorev_id"]
            })
            print(f"Task details: {details.content}")

if __name__ == "__main__":
    asyncio.run(main())
```

### Using REST API (Python Requests)

```python
#!/usr/bin/env python3
"""
Gorev REST API Client - Python Example
Direct HTTP access to daemon REST API
"""

import requests
import json

class GorevClient:
    def __init__(self, base_url="http://localhost:5082"):
        self.base_url = base_url
        self.session = requests.Session()

    def health_check(self):
        """Check daemon health"""
        response = self.session.get(f"{self.base_url}/api/health")
        return response.json()

    def list_tasks(self, status=None, limit=50):
        """List tasks with optional filtering"""
        params = {"limit": limit}
        if status:
            params["durum"] = status

        response = self.session.get(
            f"{self.base_url}/api/gorevler",
            params=params
        )
        return response.json()

    def create_task(self, baslik, aciklama, oncelik="orta"):
        """Create a new task"""
        data = {
            "baslik": baslik,
            "aciklama": aciklama,
            "oncelik": oncelik
        }
        response = self.session.post(
            f"{self.base_url}/api/gorev",
            json=data
        )
        return response.json()

    def update_task(self, task_id, durum=None, oncelik=None):
        """Update task status or priority"""
        data = {}
        if durum:
            data["durum"] = durum
        if oncelik:
            data["oncelik"] = oncelik

        response = self.session.put(
            f"{self.base_url}/api/gorev/{task_id}",
            json=data
        )
        return response.json()

    def search_tasks(self, query, filters=None):
        """Advanced search with filters"""
        data = {"query": query}
        if filters:
            data["filters"] = filters

        response = self.session.post(
            f"{self.base_url}/api/search",
            json=data
        )
        return response.json()

# Usage example
if __name__ == "__main__":
    client = GorevClient()

    # Check health
    health = client.health_check()
    print(f"Daemon status: {health['status']}")

    # List active tasks
    tasks = client.list_tasks(status="devam_ediyor")
    print(f"Active tasks: {len(tasks)}")

    # Create task
    new_task = client.create_task(
        baslik="Implement caching layer",
        aciklama="Add Redis caching for API responses",
        oncelik="yuksek"
    )
    print(f"Created task ID: {new_task['id']}")

    # Search tasks
    results = client.search_tasks(
        query="API",
        filters={"oncelik": ["yuksek"]}
    )
    print(f"Search results: {len(results)}")
```

---

## TypeScript/Node.js Integration

### Using MCP TypeScript SDK

```typescript
/**
 * Gorev MCP Client - TypeScript Example
 * Demonstrates task management with type safety
 */

import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";

interface Task {
  id: string;
  baslik: string;
  durum: "beklemede" | "devam_ediyor" | "tamamlandi";
  oncelik: "dusuk" | "orta" | "yuksek";
}

class GorevMCPClient {
  private client: Client;
  private transport: StdioClientTransport;

  constructor() {
    this.transport = new StdioClientTransport({
      command: "npx",
      args: ["-y", "@mehmetsenol/gorev-mcp-server", "serve"],
      env: { GOREV_LANG: "en" },
    });

    this.client = new Client(
      {
        name: "gorev-typescript-client",
        version: "1.0.0",
      },
      {
        capabilities: {},
      }
    );
  }

  async connect(): Promise<void> {
    await this.client.connect(this.transport);
    console.log("Connected to Gorev MCP Server");
  }

  async listTools(): Promise<void> {
    const tools = await this.client.listTools();
    console.log(`Available tools: ${tools.tools.length}`);
    tools.tools.forEach((tool) => {
      console.log(`  - ${tool.name}: ${tool.description}`);
    });
  }

  async listTasks(status?: string): Promise<Task[]> {
    const result = await this.client.callTool("gorev_listele", {
      durum: status,
      limit: 50,
    });

    return JSON.parse(result.content[0].text);
  }

  async createTask(
    baslik: string,
    aciklama: string,
    oncelik: "dusuk" | "orta" | "yuksek" = "orta"
  ): Promise<Task> {
    const result = await this.client.callTool("templateden_gorev_olustur", {
      template_id: "feature",
      degerler: {
        ozellik_adi: baslik,
        aciklama: aciklama,
        oncelik: oncelik,
      },
    });

    return JSON.parse(result.content[0].text);
  }

  async updateTask(
    taskId: string,
    durum?: string,
    oncelik?: string
  ): Promise<void> {
    const params: Record<string, string> = { id: taskId };
    if (durum) params.durum = durum;
    if (oncelik) params.oncelik = oncelik;

    await this.client.callTool("gorev_guncelle", params);
  }

  async bulkUpdate(
    taskIds: string[],
    data: { durum?: string; oncelik?: string }
  ): Promise<void> {
    await this.client.callTool("gorev_bulk", {
      operation: "update",
      ids: taskIds,
      data: data,
    });
  }

  async disconnect(): Promise<void> {
    await this.client.close();
    console.log("Disconnected from Gorev MCP Server");
  }
}

// Usage example
async function main() {
  const client = new GorevMCPClient();

  try {
    await client.connect();
    await client.listTools();

    // List tasks
    const tasks = await client.listTasks("devam_ediyor");
    console.log(`Active tasks: ${tasks.length}`);

    // Create task
    const newTask = await client.createTask(
      "Implement WebSocket notifications",
      "Add real-time notifications via WebSocket",
      "yuksek"
    );
    console.log(`Created task: ${newTask.id}`);

    // Bulk update
    await client.bulkUpdate([newTask.id], {
      durum: "devam_ediyor",
      oncelik: "yuksek",
    });
    console.log("Bulk update completed");
  } finally {
    await client.disconnect();
  }
}

main().catch(console.error);
```

### Using Axios (REST API)

```typescript
/**
 * Gorev REST API Client - TypeScript/Axios Example
 */

import axios, { AxiosInstance } from "axios";

interface GorevTask {
  id: string;
  baslik: string;
  aciklama: string;
  durum: string;
  oncelik: string;
}

class GorevAPIClient {
  private client: AxiosInstance;

  constructor(baseURL: string = "http://localhost:5082") {
    this.client = axios.create({
      baseURL,
      headers: {
        "Content-Type": "application/json",
      },
    });
  }

  async healthCheck(): Promise<boolean> {
    const { data } = await this.client.get("/api/health");
    return data.status === "healthy";
  }

  async listTasks(filters?: {
    durum?: string;
    oncelik?: string;
    limit?: number;
  }): Promise<GorevTask[]> {
    const { data } = await this.client.get("/api/gorevler", {
      params: filters,
    });
    return data;
  }

  async createTask(task: {
    baslik: string;
    aciklama: string;
    oncelik?: string;
  }): Promise<GorevTask> {
    const { data } = await this.client.post("/api/gorev", task);
    return data;
  }

  async updateTask(
    id: string,
    updates: { durum?: string; oncelik?: string }
  ): Promise<GorevTask> {
    const { data } = await this.client.put(`/api/gorev/${id}`, updates);
    return data;
  }

  async deleteTask(id: string): Promise<void> {
    await this.client.delete(`/api/gorev/${id}`);
  }

  async advancedSearch(query: string, filters?: object): Promise<GorevTask[]> {
    const { data } = await this.client.post("/api/search", {
      query,
      filters,
    });
    return data;
  }
}

// Usage
const client = new GorevAPIClient();

async function demo() {
  // Check health
  const isHealthy = await client.healthCheck();
  console.log(`Daemon healthy: ${isHealthy}`);

  // Create task
  const task = await client.createTask({
    baslik: "Optimize database queries",
    aciklama: "Add indexes and query optimization",
    oncelik: "yuksek",
  });

  // Update task
  await client.updateTask(task.id, { durum: "devam_ediyor" });

  // Search
  const results = await client.advancedSearch("database", {
    oncelik: ["yuksek"],
  });
  console.log(`Found ${results.length} tasks`);
}

demo().catch(console.error);
```

---

## Go Integration

### Using Go HTTP Client

```go
package main

import (
 "bytes"
 "encoding/json"
 "fmt"
 "io"
 "net/http"
)

// GorevClient represents a client for Gorev REST API
type GorevClient struct {
 BaseURL string
 Client  *http.Client
}

// Task represents a Gorev task
type Task struct {
 ID        string `json:"id"`
 Baslik    string `json:"baslik"`
 Aciklama  string `json:"aciklama"`
 Durum     string `json:"durum"`
 Oncelik   string `json:"oncelik"`
 ProjeID   string `json:"proje_id,omitempty"`
 SonTarih  string `json:"son_tarih,omitempty"`
}

// NewGorevClient creates a new Gorev API client
func NewGorevClient(baseURL string) *GorevClient {
 return &GorevClient{
  BaseURL: baseURL,
  Client:  &http.Client{},
 }
}

// HealthCheck checks daemon health
func (c *GorevClient) HealthCheck() (bool, error) {
 resp, err := c.Client.Get(c.BaseURL + "/api/health")
 if err != nil {
  return false, err
 }
 defer resp.Body.Close()

 var health struct {
  Status string `json:"status"`
 }
 if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
  return false, err
 }

 return health.Status == "healthy", nil
}

// ListTasks retrieves tasks with optional filtering
func (c *GorevClient) ListTasks(filters map[string]string) ([]Task, error) {
 req, err := http.NewRequest("GET", c.BaseURL+"/api/gorevler", nil)
 if err != nil {
  return nil, err
 }

 // Add query parameters
 q := req.URL.Query()
 for key, value := range filters {
  q.Add(key, value)
 }
 req.URL.RawQuery = q.Encode()

 resp, err := c.Client.Do(req)
 if err != nil {
  return nil, err
 }
 defer resp.Body.Close()

 var tasks []Task
 if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
  return nil, err
 }

 return tasks, nil
}

// CreateTask creates a new task
func (c *GorevClient) CreateTask(task Task) (*Task, error) {
 body, err := json.Marshal(task)
 if err != nil {
  return nil, err
 }

 resp, err := c.Client.Post(
  c.BaseURL+"/api/gorev",
  "application/json",
  bytes.NewBuffer(body),
 )
 if err != nil {
  return nil, err
 }
 defer resp.Body.Close()

 var createdTask Task
 if err := json.NewDecoder(resp.Body).Decode(&createdTask); err != nil {
  return nil, err
 }

 return &createdTask, nil
}

// UpdateTask updates an existing task
func (c *GorevClient) UpdateTask(id string, updates map[string]string) error {
 body, err := json.Marshal(updates)
 if err != nil {
  return err
 }

 req, err := http.NewRequest(
  "PUT",
  c.BaseURL+"/api/gorev/"+id,
  bytes.NewBuffer(body),
 )
 if err != nil {
  return err
 }
 req.Header.Set("Content-Type", "application/json")

 resp, err := c.Client.Do(req)
 if err != nil {
  return err
 }
 defer resp.Body.Close()

 if resp.StatusCode != http.StatusOK {
  body, _ := io.ReadAll(resp.Body)
  return fmt.Errorf("update failed: %s", string(body))
 }

 return nil
}

// AdvancedSearch performs advanced search with filters
func (c *GorevClient) AdvancedSearch(query string, filters map[string]interface{}) ([]Task, error) {
 searchReq := map[string]interface{}{
  "query":   query,
  "filters": filters,
 }

 body, err := json.Marshal(searchReq)
 if err != nil {
  return nil, err
 }

 resp, err := c.Client.Post(
  c.BaseURL+"/api/search",
  "application/json",
  bytes.NewBuffer(body),
 )
 if err != nil {
  return nil, err
 }
 defer resp.Body.Close()

 var tasks []Task
 if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
  return nil, err
 }

 return tasks, nil
}

func main() {
 client := NewGorevClient("http://localhost:5082")

 // Health check
 healthy, err := client.HealthCheck()
 if err != nil {
  panic(err)
 }
 fmt.Printf("Daemon healthy: %v\n", healthy)

 // List tasks
 tasks, err := client.ListTasks(map[string]string{
  "durum": "devam_ediyor",
  "limit": "10",
 })
 if err != nil {
  panic(err)
 }
 fmt.Printf("Active tasks: %d\n", len(tasks))

 // Create task
 newTask, err := client.CreateTask(Task{
  Baslik:   "Implement health metrics",
  Aciklama: "Add Prometheus metrics endpoints",
  Oncelik:  "yuksek",
 })
 if err != nil {
  panic(err)
 }
 fmt.Printf("Created task ID: %s\n", newTask.ID)

 // Update task
 err = client.UpdateTask(newTask.ID, map[string]string{
  "durum":   "devam_ediyor",
  "oncelik": "yuksek",
 })
 if err != nil {
  panic(err)
 }
 fmt.Println("Task updated successfully")

 // Advanced search
 results, err := client.AdvancedSearch("metrics", map[string]interface{}{
  "oncelik": []string{"yuksek"},
 })
 if err != nil {
  panic(err)
 }
 fmt.Printf("Search results: %d tasks\n", len(results))
}
```

---

## REST API Direct Access

### cURL Examples

```bash
#!/bin/bash
# Gorev REST API - cURL Examples

BASE_URL="http://localhost:5082"

# Health check
curl -X GET "$BASE_URL/api/health" | jq

# List tasks
curl -X GET "$BASE_URL/api/gorevler?durum=devam_ediyor&limit=10" | jq

# Create task
curl -X POST "$BASE_URL/api/gorev" \
  -H "Content-Type: application/json" \
  -d '{
    "baslik": "Fix memory leak",
    "aciklama": "Profile and fix memory leak in API server",
    "oncelik": "yuksek"
  }' | jq

# Update task
TASK_ID="abc12345"
curl -X PUT "$BASE_URL/api/gorev/$TASK_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "durum": "devam_ediyor",
    "oncelik": "yuksek"
  }' | jq

# Advanced search
curl -X POST "$BASE_URL/api/search" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "API",
    "filters": {
      "durum": ["devam_ediyor"],
      "oncelik": ["yuksek"]
    }
  }' | jq

# Export tasks
curl -X POST "$BASE_URL/api/export" \
  -H "Content-Type: application/json" \
  -d '{
    "format": "json",
    "include_subtasks": true
  }' | jq

# Get workspace summary
curl -X GET "$BASE_URL/api/ozet" | jq
```

---

## WebSocket Real-Time Updates

### JavaScript/TypeScript WebSocket Client

```typescript
/**
 * Gorev WebSocket Client - Real-time updates
 */

class GorevWebSocketClient {
  private ws: WebSocket;
  private reconnectInterval: number = 5000;
  private reconnectAttempts: number = 0;
  private maxReconnectAttempts: number = 10;

  constructor(private url: string = "ws://localhost:5082/ws") {
    this.connect();
  }

  private connect(): void {
    console.log("Connecting to Gorev WebSocket...");
    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      console.log("WebSocket connected");
      this.reconnectAttempts = 0;
    };

    this.ws.onmessage = (event) => {
      try {
        const update = JSON.parse(event.data);
        this.handleUpdate(update);
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    this.ws.onclose = () => {
      console.log("WebSocket disconnected");
      this.attemptReconnect();
    };
  }

  private handleUpdate(update: {
    type: string;
    task_id: string;
    project_id: string;
    timestamp: string;
    data: any;
  }): void {
    console.log(`Task ${update.type}:`, update.task_id);

    switch (update.type) {
      case "created":
        console.log("New task created:", update.data);
        break;
      case "updated":
        console.log("Task updated:", update.data);
        break;
      case "deleted":
        console.log("Task deleted:", update.task_id);
        break;
      default:
        console.log("Unknown update type:", update.type);
    }
  }

  private attemptReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error("Max reconnection attempts reached");
      return;
    }

    this.reconnectAttempts++;
    console.log(
      `Reconnecting in ${this.reconnectInterval}ms (attempt ${this.reconnectAttempts})`
    );

    setTimeout(() => {
      this.connect();
    }, this.reconnectInterval);
  }

  public disconnect(): void {
    if (this.ws) {
      this.ws.close();
    }
  }
}

// Usage
const wsClient = new GorevWebSocketClient();

// Keep connection alive
process.on("SIGINT", () => {
  console.log("Disconnecting...");
  wsClient.disconnect();
  process.exit(0);
});
```

---

## Error Handling Best Practices

### Comprehensive Error Handling Example

```typescript
import { GorevAPIClient } from "./gorev-client";

async function robustTaskCreation() {
  const client = new GorevAPIClient();

  try {
    // Verify daemon is running
    const isHealthy = await client.healthCheck();
    if (!isHealthy) {
      throw new Error("Gorev daemon is not healthy");
    }

    // Create task with retry logic
    let task;
    let retries = 3;

    while (retries > 0) {
      try {
        task = await client.createTask({
          baslik: "Critical task",
          aciklama: "Important work",
          oncelik: "yuksek",
        });
        break;
      } catch (error) {
        retries--;
        if (retries === 0) throw error;
        console.log(`Retry attempt ${3 - retries}...`);
        await new Promise((resolve) => setTimeout(resolve, 1000));
      }
    }

    console.log(`Task created successfully: ${task.id}`);
    return task;
  } catch (error) {
    if (error.response) {
      // Server responded with error
      console.error(`API Error: ${error.response.status}`);
      console.error(error.response.data);
    } else if (error.request) {
      // No response received
      console.error("No response from server. Is daemon running?");
      console.error("Start daemon: gorev daemon --detach");
    } else {
      // Other errors
      console.error("Error:", error.message);
    }
    throw error;
  }
}
```

---

## Performance Optimization Tips

1. **Connection Pooling**: Reuse HTTP clients instead of creating new ones
2. **Batch Operations**: Use `gorev_bulk` for multiple updates
3. **Pagination**: Always use `limit` and `offset` for large datasets
4. **Caching**: Cache frequently accessed data locally
5. **WebSocket**: Use WebSocket for real-time updates instead of polling

---

## See Also

- [MCP Tools Reference](./MCP_TOOLS_REFERENCE.md) - Complete tool documentation
- [REST API Reference](./rest-api-reference.md) - Detailed API endpoints
- [Daemon Architecture](../architecture/daemon-architecture.md) - System architecture
- [NPM Package](https://www.npmjs.com/package/@mehmetsenol/gorev-mcp-server) - Official package

---

**Last Updated:** October 6, 2025 | **Version:** v0.16.3
